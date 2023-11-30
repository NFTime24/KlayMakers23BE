package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nftime/config"
	"github.com/nftime/logger"
	"github.com/nftime/util"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nftime/db"
	"github.com/nftime/model"
)

// @Tags User
// @Summary PostUser
// @Description PostUser
// @Accept json
// @Produce json
// @Param like body model.UserCreateParam true "user data"
// @Router /users [post]
// @Success 200 {object} model.UserResult
func PostUser(c *gin.Context) {
	streamPlatformDb := db.StreamPlatformDbManager()
	var ticketCount int
	userBody := model.UserCreateParam{}
	err := c.Bind(&userBody)
	if err != nil {
		logger.Error.Printf("Failed processing Binding: %v\n", err)
	}

	var social model.SocialUser
	streamPlatformDb.Model(social).Select("*").
		Where("id=?", userBody.SocialUserId).Scan(&social)

	logger.Info.Printf("social: %v\n", social)
	if len(social.NickName) <= 0 {
		social.NickName = userBody.NickName
	}
	if len(social.ProfilePath) <= 0 {
		social.ProfilePath = "https://secure.nftime.gallery/assets/uploadimage/upload-773568470.png"
	}
	var exists bool
	var user model.User
	streamPlatformDb.Model(&user).
		Select("count(*) > 0").
		Where("user_address = ?", userBody.Address).
		Find(&exists)

	streamPlatformDb.Select("count(ticket_count)").Table("tickets").Where("user_address=?", userBody.Address).Scan(&ticketCount)

	if ticketCount == 0 {
		ticketInsert := model.Ticket{UserAddress: userBody.Address, TicketCount: 3}
		streamPlatformDb.Create(&ticketInsert)
	}

	if exists {
		streamPlatformDb.Model(&user).Where("user_address = ?", userBody.Address).Updates(map[string]interface{}{
			"nick_name":    social.NickName,
			"profile_path": social.ProfilePath,
		})
	} else {
		userInsert := model.User{UserAddress: userBody.Address, NickName: social.NickName, ProfilePath: social.ProfilePath}
		streamPlatformDb.Create(&userInsert)
	}
	//TODO: ProfilePath

	var result model.UserResult
	streamPlatformDb.Select(`u.*, t.ticket_count`).
		Table("users as u").Joins("left join tickets as t on u.user_address = t.user_address").
		Where("u.user_address=?", userBody.Address).Scan(&result)

	c.JSON(http.StatusOK, result)
}

// @Param file formData file true "profile_image"
// @Param userParamBody body model.UserProfileInfoParam true "user data"

// @Tags User
// @Summary UploadProfile
// @Description UploadProfile
// @Accept multipart/form-data
// @Produce json
// @Param id formData string true "user id"
// @Param user_address formData string true "user address"
// @Param nick_name formData string true "user nickname"
// @Param file formData file true "profile_image"
// @Router /users/upload-profile [post]
// @Success 200 {object} model.UserResult
func UploadProfile(c *gin.Context) {
	streamPlatformDb := db.StreamPlatformDbManager()

	userAddress := c.PostForm("user_address")
	userIdStr := c.PostForm("id")
	userNickName := c.PostForm("nick_name")
	userId, _ := strconv.ParseUint(userIdStr, 10, 64)
	var nickName string
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}
	files := form.File["file"]
	logger.Info.Println(files)
	var originalLocation string
	var user model.User
	streamPlatformDb.Select("nick_name").Table("users").Where("user_address=? and id=?", userAddress, userId).Scan(&nickName)
	if nickName != userNickName {
		streamPlatformDb.Model(&user).Where("user_address=? and id =?", userAddress, userId).Update("nick_name", userNickName)
	}
	//Where("u.user_address=? and l.work_id=?", userAddress, workId)
	for _, file := range files {
		filename := filepath.Base(file.Filename)
		uploadedFile, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer uploadedFile.Close()

		// Create a temporary file
		tempFile, err := os.CreateTemp("", "upload-*.tmp")
		if err != nil {
			log.Fatal(err)
		}
		defer tempFile.Close()
		fileNameParts := strings.Split(filename, ".")
		filename = fileNameParts[0]
		fileName := util.CreateUnixFilename(filename)

		config.StreamS3Info.BucketName = config.Cfg.AWS.Stream.ProfileBucket

		result := config.StreamS3Info.UploadFile(uploadedFile, fileName, "images/")

		originalLocation = result.Location
		fmt.Println(originalLocation)
	}
	streamPlatformDb.Model(&user).Where("user_address=? and id =?", userAddress, userId).Update("profile_path", originalLocation)

	var result model.UserResult
	streamPlatformDb.Select(`u.*, t.ticket_count`).
		Table("users as u").Joins("left join tickets as t on u.user_address = t.user_address").
		Where("u.user_address=? and u.id=?", userAddress, userId).Scan(&result)

	c.JSON(http.StatusOK, result)

	//db.Model(&user).Where("user_address =? and id = ?", userAddress, userId).Update("profile_path", originalLocation)
	//
	//rows, err := streamPlatformDb.Select("*").Table("users").Where("user_address =? and id = ?", userAddress, userId).Rows()
	//for rows.Next() {
	//	db.ScanRows(rows, &user)
	//}
	//defer rows.Close()
	////err = row.Scan(&user)
	//if err != nil {
	//	logger.Error.Printf("failed to get user info: %v\n", err)
	//}
	//c.JSON(http.StatusOK, user)
}

func UploadCompany(c *gin.Context) {
	timeStorageDb := db.TimeStorageDbManager()

	companyName := c.PostForm("company_name")
	companyDescription := c.PostForm("company_description")
	companyWebsite := c.PostForm("company_website")

	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}
	files := form.File["company_image"]
	logger.Info.Println(files)
	var originalLocation string

	for _, file := range files {
		filename := filepath.Base(file.Filename)
		uploadedFile, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer uploadedFile.Close()

		// Create a temporary file
		tempFile, err := os.CreateTemp("", "upload-*.tmp")
		if err != nil {
			log.Fatal(err)
		}
		defer tempFile.Close()
		fileNameParts := strings.Split(filename, ".")
		filename = fileNameParts[0]
		fileName := util.CreateUnixFilename(filename)

		config.StreamS3Info.BucketName = config.Cfg.AWS.Stream.ProfileBucket

		result := config.StreamS3Info.UploadFile(uploadedFile, fileName, "images/")

		originalLocation = result.Location
		fmt.Println(originalLocation)
	}
	companyInsert := model.Company{CompanyName: companyName, CompanyImage: originalLocation, CompanyDescription: companyDescription, CompanyWebsite: companyWebsite}
	timeStorageDb.Create(&companyInsert)

	var result model.CompanyResult
	timeStorageDb.Select(`cu.*`).
		Table("certi_users as cu").
		Where("cu.company_name=? and cu.company_description and cu.company_website =?", companyName, companyDescription, companyWebsite).Scan(&result)

	c.JSON(http.StatusOK, result)
}

// @Summary GetUserWithAddress
// @Description Get requested user info with address
// @Tags User
// @Accept json
// @Produce json
// @Param address path string true "address of user"
// @Router /users/addresses/{address} [get]
// @Success 200 {object} model.UserResult
func GetUserWithAddress(c *gin.Context) {
	address := c.Param("address")
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")

	streamPlatformDb := db.StreamPlatformDbManager()
	var result model.UserResult
	streamPlatformDb.Select(`u.*, t.ticket_count`).
		Table("users as u").Joins("left join tickets as t on u.user_address = t.user_address").
		Where("u.user_address=?", address).Scan(&result)
	logger.Info.Println(result)

	if result == (model.UserResult{}) {
		c.JSON(http.StatusOK, nil)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

// @Summary UpdateUserNickname
// @Description Update User Nickname with address and nickname
// @Tags User
// @Accept json
// @Produce json
// @Param like body model.UserCreateParam true "user data"
// @Router /users/nickname [patch]
func UpdateUserNickname(c *gin.Context) {
	streamPlatformDb := db.StreamPlatformDbManager()

	userBody := model.UserCreateParam{}
	err := c.Bind(&userBody)
	if err != nil {
		logger.Error.Printf("Failed processing Binding: %v\n", err)
	}

	var exists bool
	var User model.User
	streamPlatformDb.Model(&User).
		Select("count(*) > 0").
		Where("nick_name = ?", userBody.NickName).
		Find(&exists)

	if exists {
		c.String(http.StatusOK, "U1002") // 중복된 닉네임
	} else {
		streamPlatformDb.Model(&User).Where("user_address=?", userBody.Address).Update("nick_name", userBody.NickName)
		c.String(http.StatusOK, "U1001") // 정상 동작
	}
}

// @Summary UpdateUserNickname
// @Description Update User Nickname with address and nickname
// @Tags User
// @Accept json
// @Produce json
// @Param like body model.UserCreateParam true "user data"
// @Router /users/nickname [post]
func PostUserNickname(c *gin.Context) {
	streamPlatformDb := db.StreamPlatformDbManager()

	userBody := model.UserCreateParam{}
	err := c.Bind(&userBody)
	if err != nil {
		logger.Error.Printf("Failed processing Binding: %v\n", err)
	}

	var exists bool
	var User model.User
	streamPlatformDb.Model(&User).
		Select("count(*) > 0").
		Where("nick_name = ?", userBody.NickName).
		Find(&exists)

	if exists {
		c.String(http.StatusOK, "U1002") // 중복된 닉네임
	} else {
		streamPlatformDb.Model(&User).Where("user_address=?", userBody.Address).Update("nick_name", userBody.NickName)
		c.String(http.StatusOK, "U1001") // 정상 동작
	}
}

// @Summary delete user
// @Description delete user
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "user_id"
// @Param address path string true "user_address"
// @Router /users/ids/{id}/addresses/{address} [delete]
func DeleteUser(c *gin.Context) {
	streamPlatformDb := db.StreamPlatformDbManager()
	userId := c.Param("id")
	userAddress := c.Param("address")

	var likes model.Like
	var user model.User
	result := streamPlatformDb.Where("id = ? and user_address =?", userId, userAddress).Delete(&user)

	result2 := streamPlatformDb.Where("user_address =?", userId).Delete(&likes)

	if result.Error != nil {
		logger.Error.Printf("err: %v\n", result.Error)
	}
	if result2.Error != nil {
		logger.Error.Printf("err: %v\n", result2.Error)
	}
	deleteInfo := fmt.Sprintf("user address : %s deleted", userAddress)
	c.String(http.StatusOK, deleteInfo)
}

// @Summary SetTicketCountTo10
// @Description SetTicketCountTo10 of specific address
// @Tags User
// @Accept json
// @Produce json
// @Param address path string true "address of user"
// @Router /users/test/{address} [get]
// @Success 200 {object} model.UserResult
func TestUser(c *gin.Context) {
	streamPlatformDb := db.StreamPlatformDbManager()
	//UserAddress string `json:"user_address"`
	addr := c.Param("address")
	//var userData UserData

	streamPlatformDb.Model(&model.Ticket{}).Where("user_address=?", addr).Update("ticket_count", 10)

	//rows, err := streamPlatformDb.Select(`distinct user_address`).Table("users").Rows()
	//if err != nil {
	//	logger.Error.Printf("err: %v\n", err)
	//}
	responseString := fmt.Sprintf("address %v set ticket_count to %v", addr, 10)
	c.JSON(http.StatusOK, responseString)
	//db.Select("count(ticket_count)").Table("ticket").Where("user_address=?", userBody.Address).Scan(&ticketCount)

}

// @Summary AddUserLog
// @Description Add logger of a user
// @Tags User
// @Accept json
// @Produce json
// @Deprecated True
// @Param user_address path string true "address of user"
// @Param status path string true "status of user"
// @Router /users/logs/{address}/{status} [post]
func AddUserLog(c *gin.Context) {
	userAddress := c.Param("address")
	status := c.Param("status")

	streamPlatformDb := db.StreamPlatformDbManager()

	streamPlatformDb.Create(&model.Log{
		UserAddress: userAddress,
		Status:      status,
	})

	resultStr := "Add Log Successfully"

	c.String(http.StatusOK, resultStr)
}

func SendCrashReport(c *gin.Context) {
	jsonBody := make(map[string]interface{})
	err := json.NewDecoder(c.Request.Body).Decode(&jsonBody)
	if err != nil {

		logger.Error.Printf("err: %v\n", err)
		return
	}

	filepathString := fmt.Sprint("./assets/crashReports/cr_", time.Now())
	absPath, _ := filepath.Abs(filepathString)
	f, err := os.Create(absPath)
	check(err)
	defer f.Close()
	jsonBytes, err := json.MarshalIndent(jsonBody, "", " ")
	check(err)
	f.Write(jsonBytes)

	c.JSON(http.StatusOK, jsonBody)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

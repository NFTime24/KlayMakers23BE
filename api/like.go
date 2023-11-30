package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nftime/logger"
	"github.com/nftime/model"
	"net/http"
	"strconv"

	"github.com/nftime/db"
)

// @Summary update like
// @Description update like
// @Tags Like
// @Accept json
// @Produce json
// @Param like body model.LikeCreateParam true "like data"
// @Router /likes [post]
func UpdateLike(c *gin.Context) {
	logger.Info.Println("=======UpdateLike=======")

	like := model.LikeCreateParam{}
	err := c.Bind(&like)
	if err != nil {
		logger.Error.Printf("Failed processing update like request: %s\n", err)
		c.String(http.StatusBadRequest, "err occured")
	}

	streamPlatformDb := db.StreamPlatformDbManager()
	var likes model.Like
	likeIndex := -1
	// select w.name as work_name, a.name as artist_name, w.description from works w join artists a on w.artist_id = a.id;
	streamPlatformDb.Model(likes).Select("id").
		Where("user_address=? and work_id=?", like.UserAddress, like.WorkId).Scan(&likeIndex)

	logger.Info.Printf("like index :%v\n", likeIndex)
	if likeIndex == -1 {
		// 좋아요가 되어있지 않은 상태

		updateLike := model.Like{UserAddress: like.UserAddress, WorkID: like.WorkId}
		result := streamPlatformDb.Create(&updateLike)
		if result.Error != nil {
			logger.Error.Printf("result err :%v\n", result.Error)
			logger.Error.Println("=======UpdateLike End With Err=======")
			c.String(http.StatusBadRequest, "err occured")
		}
		logger.Info.Println("=======UpdateLike End=======")
		c.String(http.StatusOK, "added")
	} else {
		// 좋아요를 한 상태라면 삭제
		streamPlatformDb.Where("id =?", likeIndex).Delete(&likes)
		logger.Info.Println("=======UpdateLike End with unlike=======")
		c.String(http.StatusOK, "deleted")
	}
}

// @Summary CheckLike
// @Description CheckLike
// @Tags Like
// @Accept json
// @Produce json
// @Param address path string true "address of user"
// @Param id path string true "work id"
// @Router /likes/check/addresses/{address}/works/{id} [get]
func CheckLike(c *gin.Context) {
	userAddress := c.Param("address")
	workIdStr := c.Param("id")
	workId, _ := strconv.ParseUint(workIdStr, 10, 64)

	c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")

	streamPlatformDb := db.StreamPlatformDbManager()
	likeIndex := -1
	streamPlatformDb.Select("l.id").
		Table("likes as l").
		Joins("left join users as u on u.user_address = l.user_address").
		Where("u.user_address=? and l.work_id=?", userAddress, workId).Scan(&likeIndex)

	if likeIndex == -1 {
		c.String(http.StatusOK, "false")
	} else {
		c.String(http.StatusOK, "true")
	}
}

// @Summary GetLikeCount
// @Description GetLikeCount
// @Tags Like
// @Accept json
// @Produce json
// @Param address query string true "user address"
// @Router /likes/list [get]
func GetLikeList(c *gin.Context) {
	userAddress := c.Query("address")
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")

	streamPlatformDb := db.StreamPlatformDbManager()

	//var result []model.LikeListResult
	//
	//streamPlatformDb.Select("l.id").
	//	Table("likes as l").
	//	Where("user_address=?", userAddress).Scan(&result)

	//var result []model.LikeListResult

	var results []model.WorkResult
	rows, err := streamPlatformDb.Select(`w.id as id, 
	 w.name as work_name, w.price as price, w.description as description, w.file_path as file_path, w.thumbnail_path as thumbnail_path, 
	 w.category as work_category, a.profile_path as profile_path, a.name as artist_name, a.address as artist_address`).
		Table(`works as w`).
		Joins("left join artists as a on w.artist_name = a.name").
		Joins("left join likes as l on l.work_id = w.id").
		Where("w.display =? and l.user_address =?", true, userAddress).Rows()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		streamPlatformDb.ScanRows(rows, &results)
	}
	c.JSON(http.StatusOK, results)

}

// @Summary GetLikeList
// @Description GetLikeList
// @Tags Like
// @Accept json
// @Produce json
// @Param id query string true "user_address"
// @Router /likes/works/count [get]
func GetLikeCount(c *gin.Context) {
	workId := c.Query("id")
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")

	streamPlatformDb := db.StreamPlatformDbManager()
	likeCount := 0
	streamPlatformDb.Select("count(l.id)").
		Table("likes as l").
		Where("work_id=?", workId).Scan(&likeCount)

	c.String(http.StatusOK, strconv.Itoa(likeCount))
}

//func GetLikesWithUserId(c *gin.Context) error {
//	user_id_str := c.QueryParam("user_id")
//	user_id, _ := strconv.ParseUint(user_id_str, 10, 64)
//
//	type Result struct {
//		NftId         int    `json:"nft_id"`
//		WorkName      string `json:"work_name"`
//		Price         int    `json:"work_price"`
//		Description   string `json:"description"`
//		WorkCategory  string `json:"category"`
//		FileName      string `json:"filename"`
//		FileSize      int    `json:"filesize"`
//		FileType      string `json:"filetype"`
//		FilePath      string `json:"path"`
//		ThumbnailPath string `json:"thumbnail_path"`
//		ArtistName    string `json:"artist_name"`
//		ProfilePath   string `json:"artist_profile_path"`
//		ArtistAddress string `json:"artist_address"`
//		LikeId        int    `json:"like_id"`
//	}
//
//	streamPlatformDb := db.StreamPlatformDbManager()
//	var users model.User
//	var results Result
//	db.Model(users).Select(`n.nft_id as nft_id,
//	 w.name as work_name, w.price as price, w.description as description,
//	 w.category as work_category,f.filename as file_name, f.filesize as file_size,
//	 f.filetype as file_type, f.path as file_path, t.path as thumbnail_path, a.name as artist_name, p.path as profile_path,
//	 a.address as artist_address`).
//		Joins("left join nfts as n on users.id = n.owner_id").
//		Joins("left join works as w on n.works_id = w.work_id").
//		Joins("left join files as f on w.file_id = f.id").
//		Joins("left join files as t on f.thumbnail_id = t.id").
//		Joins("left join artists as a on w.artist_id = a.id").
//		Joins("left join files as p on a.profile_id = p.id").
//		Joins("left join likes as l on users.id = l.owner_id").
//		Where("l.owner_id=?", user_id).Scan(&results)
//
//	logger.Info.Printf("result: %v\n", results)
//	return c.JSON(http.StatusOK, results)
//}

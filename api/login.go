package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nftime/config"
	"github.com/nftime/db"
	"github.com/nftime/logger"
	"github.com/nftime/model"
	"io"
	"net/http"
	"strconv"
)

// @Summary social login
// @Description social login
// @Tags Login
// @Accept json
// @Produce json
// @Param like body model.LoginInfo true "login data"
// @Param x-nftime-token header string true "Login Token" // Add this line to specify a header parameter
// @Router /login/social [post]
func LoginWithSocial(c *gin.Context) {
	var social model.SocialUser
	accessTokenHeader := c.GetHeader("x-nftime-token")

	logger.Info.Println("access token: ", accessTokenHeader)
	loginBody := model.LoginInfo{}
	socialUserId := -1
	streamPlatformDb := db.StreamPlatformDbManager()

	err := c.Bind(&loginBody)
	if err != nil {
		logger.Error.Printf("Failed processing Binding: %v\n", err)
		c.String(http.StatusInternalServerError, "err occured while binding")
		return
	}
	if loginBody.Target == "kakao" {

		if accessTokenHeader != "" {
			client := &http.Client{}

			kakaoValidationUrl := "https://kapi.kakao.com/v2/user/me"
			req, _ := http.NewRequest("GET", kakaoValidationUrl, nil)
			req.Header.Set("Authorization", "Bearer "+accessTokenHeader)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
			res, err := client.Do(req)
			if err != nil {
				logger.Error.Printf("err: %v\n", err)
			}
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				logger.Error.Printf("err: %v\n", err)
			}
			logger.Info.Println(string(body))
			var kakaoUser model.KakaoUser
			err = json.Unmarshal(body, &kakaoUser)
			if err != nil {
				fmt.Println("Error decoding JSON:", err)
				c.String(http.StatusInternalServerError, "Unmarshal error")
				return
			}
			nickName := kakaoUser.Properties.Nickname
			profilePath := kakaoUser.Properties.ThumbnailImage
			email := kakaoUser.KakaoAccount.Email
			gender := kakaoUser.KakaoAccount.Gender
			ageRange := kakaoUser.KakaoAccount.AgeRange
			birthDay := kakaoUser.KakaoAccount.Birthday
			birthDayType := kakaoUser.KakaoAccount.BirthdayType

			if len(email) <= 0 {
				c.String(http.StatusBadRequest, "unable to get email")
				return
			}
			// TODO: 사용자 정보 요청
			streamPlatformDb.Model(social).Select("id").
				Where("email=?", email).Scan(&socialUserId)
			if socialUserId != -1 {
				c.String(http.StatusOK, strconv.Itoa(socialUserId))
				return

			} else {
				userInsert := model.SocialUser{NickName: nickName, ProfilePath: profilePath, Email: email, Gender: gender, AgeRange: ageRange, Birthday: birthDay, BirthdayType: birthDayType}

				streamPlatformDb.Create(&userInsert)

				streamPlatformDb.Model(social).Select("id").
					Where("email=?", email).Scan(&socialUserId)

				c.String(http.StatusOK, strconv.Itoa(socialUserId))
				return
			}

		}

	} else {
		c.String(http.StatusInternalServerError, "No login service other than Kakao allowed")
		return
	}
}

func LoginBackoffice(c *gin.Context) {

	loginBody := model.LoginParam{}
	err := c.Bind(&loginBody)
	if err != nil {
		logger.Error.Printf("Failed processing Binding: %v\n", err)
	}
	if loginBody.Id == config.Cfg.Login.Id && loginBody.Password == config.Cfg.Login.Pw {
		c.JSON(http.StatusOK, "login success")
		return
	} else {
		c.JSON(http.StatusUnauthorized, "login failed")
		return
	}
}

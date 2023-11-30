package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nftime/db"
	"github.com/nftime/model"
	"net/http"
)

//func PostFantalk(c *gin.Context) error {
//	streamPlatformDb := db.StreamPlatformDbManager()
//	params := make(map[string]string)
//	bindingParams := c.Bind(&params)
//
//	logger.Info.Printf("binding: %v\n", bindingParams)
//	logger.Info.Printf("artist id: %v, owner id: %v, post_text: %v\n", params["artist_id"], params["owner_id"], params["post_text"])
//
//	var id uint
//	var fantalk_id model.Fantalk
//	artist_id, _ := strconv.ParseUint(params["artist_id"], 10, 32)
//	owner_id, _ := strconv.ParseUint(params["owner_id"], 10, 32)
//	post_text := params["post_text"]
//
//	db.Model(&fantalk_id).Select("post_id").Last(&id)
//	id += 1
//	logger.Info.Printf("new id : %v\n", id)
//
//	creative_time := time.Now()
//	modify_time := time.Now()
//
//	fantalk_insert := model.Fantalk{
//		Post_id:    id,
//		ArtistID:   uint(artist_id),
//		OwnerID:    uint(owner_id),
//		PostText:   post_text,
//		LikeCount:  0,
//		CreateTime: &creative_time,
//		ModifyTime: &modify_time,
//	}
//
//	db.Create(&fantalk_insert)
//
//	return c.String(http.StatusOK, strconv.FormatUint(uint64(id), 10))
//}

// @Summary GetArtistFantalks
// @Description Get fantalks of requested artist
// @Tags Fantalk
// @Accept json
// @Produce json
// @Param artist_name path string true "name of artist"
// @Router /fantalks/artists/{name} [get]
// @Success 200 {object} model.FantalkResult
func GetArtistFantalks(c *gin.Context) {
	artistName := c.Param("name")

	streamPlatformDb := db.StreamPlatformDbManager()
	var results []model.FantalkResult
	rows, err := streamPlatformDb.Select(`ft.*, u.nick_name, u.profile_path`).
		Table("fantalks as ft").
		Joins("left join users as u on u.user_address = ft.user_address").
		Where("ft.artist_name=?", artistName).
		Order("ft.id desc").Rows()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		streamPlatformDb.ScanRows(rows, &results)
	}
	defer rows.Close()
	c.JSON(http.StatusOK, results)
}

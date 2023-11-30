package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/nftime/db"
	"github.com/nftime/logger"
	"github.com/nftime/model"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// @Summary GetUserPlayListWithId
// @Description GetUserPlayListWithId
// @Tags Playlist
// @Accept json
// @Produce json
// @Param id path string true "id of user"
// @Router /playlists/ids/{id} [get]
// @Success 200 {object} model.PlaylistResult
func GetUserPlayListWithId(c *gin.Context) {
	userIdStr := c.Param("id")
	userId, _ := strconv.ParseUint(userIdStr, 10, 64)

	logger.Info.Printf("userIdStr: %v\n", userIdStr)
	logger.Info.Printf("userId: %v\n", userId)
	streamPlatformDb := db.StreamPlatformDbManager()
	var playlistResults []model.PlaylistResult

	// Define the query and execute it
	rows, err := streamPlatformDb.Table("playlists p").
		Select("p.id, p.created_at, p.updated_at, p.name, p.index, ARRAY_AGG(pw.work_id) as work_id").
		Joins("LEFT JOIN playlist_works as pw ON p.id = pw.playlist_id").
		Where("user_id = ?", userId).
		Group("p.id, p.created_at, p.updated_at, p.name, p.index").
		Order("p.index asc").
		Rows()

	if err != nil {
		logger.Error.Printf("err: %v\n", err)
		c.String(http.StatusBadRequest, "err occurred")
		return
	}

	defer rows.Close()

	for rows.Next() {
		var playlistResult model.PlaylistResult
		var workIDs pq.Int64Array // Use pq.Int64Array to scan the PostgreSQL array as int64

		err := rows.Scan(&playlistResult.Id, &playlistResult.CreatedAt, &playlistResult.UpdatedAt, &playlistResult.Name, &playlistResult.Index, &workIDs)
		if err != nil {
			logger.Error.Printf("err: %v\n", err)
			c.String(http.StatusBadRequest, "err occurred")
			return
		}

		// Convert pq.Int64Array to []uint
		for _, id := range workIDs {
			playlistResult.WorkId = append(playlistResult.WorkId, uint(id))
		}

		playlistResults = append(playlistResults, playlistResult)
	}
	c.JSON(http.StatusOK, playlistResults)
}

// @Summary DeleteUserPlaylistWithIds
// @Description DeleteUserPlaylistWithIds
// @Tags Playlist
// @Accept json
// @Produce json
// @Param id path string true "id of user"
// @Param ids path string true "ids of playlist"
// @Router /playlists/ids/{ids}/user/ids/{id} [delete]
// @Success 200 {object} model.WorkResult
func DeleteUserPlaylistWithIds(c *gin.Context) {
	userIdStr := c.Param("id")
	userId, _ := strconv.ParseUint(userIdStr, 10, 64)

	playlistIdsStr := c.Param("ids")
	selectedIds := strings.Split(playlistIdsStr, ",")

	streamPlatformDb := db.StreamPlatformDbManager()
	var playlist model.Playlist

	result := streamPlatformDb.Where("id IN ? and user_id =?", selectedIds, userId).Delete(&playlist)

	if result.Error != nil {
		logger.Error.Printf("err: %v\n", result.Error)
		c.String(http.StatusBadRequest, "err occurred")
		return
	}
	var playlistWork model.PlaylistWork

	result = streamPlatformDb.Where("playlist_id IN ?", selectedIds).Delete(&playlistWork)

	if result.Error != nil {
		logger.Error.Printf("err: %v\n", result.Error)
		c.String(http.StatusBadRequest, "err occurred")
		return
	}
	deleteInfo := fmt.Sprintf("playlist ids [%s] for user id [%v] deleted", playlistIdsStr, userId)
	c.String(http.StatusOK, deleteInfo)
}

// @Summary CreatePlaylist
// @Description CreatePlaylist
// @Tags Playlist
// @Accept json
// @Produce json
// @Param like body model.PlaylistCreateParam true "playlist data"
// @Router /playlists [post]
// @Success 200
func CreatePlaylist(c *gin.Context) {
	playlistParam := model.PlaylistCreateParam{}
	var playlists model.Playlist

	err := c.Bind(&playlistParam)
	if err != nil {
		logger.Error.Printf("Failed processing create playlist request: %s\n", err)
		c.String(http.StatusBadRequest, "err occurred")
	}

	if len(playlistParam.WorkIds) < 1 {
		logger.Error.Println("WorkId length must be longer than 0")
		c.String(http.StatusBadRequest, "WorkId length must be longer than 0")
		return
	}
	streamPlatformDb := db.StreamPlatformDbManager()

	lastIndex := -1
	streamPlatformDb.Model(playlists).Select(`index`).
		Where("user_id=?", playlistParam.UserId).Order("index desc").Limit(1).Scan(&lastIndex)

	var createPlaylist model.Playlist
	var index int
	if lastIndex == -1 {
		index = 1
	} else {
		index = lastIndex + 1
	}
	logger.Info.Printf("last_id: %v, index: %v\n", lastIndex, index)

	createPlaylist = model.Playlist{
		UserId: playlistParam.UserId,
		Name:   playlistParam.PlaylistName,
		Index:  uint(index),
	}
	//WorkId: []int32(playlistParam.WorkIds),

	result := streamPlatformDb.Create(&createPlaylist)
	if result.Error != nil {
		logger.Error.Printf("result err :%v\n", result.Error)
		logger.Error.Println("=======Create Playlist End With Err=======")
		c.String(http.StatusBadRequest, "err occurred may be playlist name dup error")
		return
	}

	var lastId uint
	streamPlatformDb.Model(playlists).Select(`id`).
		Where("user_id=?", playlistParam.UserId).Order("id desc").Limit(1).Scan(&lastId)
	var createPlaylistWorkList []model.PlaylistWork

	for _, v := range playlistParam.WorkIds {
		createPlaylistWork := model.PlaylistWork{
			PlaylistId: lastId,
			WorkId:     v,
		}
		createPlaylistWorkList = append(createPlaylistWorkList, createPlaylistWork)
	}

	result = streamPlatformDb.Create(&createPlaylistWorkList)
	if result.Error != nil {
		logger.Error.Printf("result err :%v\n", result.Error)
		logger.Error.Println("=======Create Playlist End With Err=======")
		c.String(http.StatusBadRequest, "err occurred")
	}

	lastIdStr := strconv.Itoa(int(lastId))
	logger.Info.Println("=======Create Playlist End=======")
	c.String(http.StatusOK, lastIdStr)
}

// @Summary UpdatePlaylistOrder
// @Description UpdatePlaylistOrder
// @Tags Playlist
// @Accept json
// @Produce json
// @Param like body model.PlaylistOrderParam true "playlist data"
// @Router /playlists [patch]
// @Success 200
func UpdatePlaylistOrder(c *gin.Context) {

	streamPlatformDb := db.StreamPlatformDbManager()

	playlistParam := model.PlaylistOrderParam{}
	var playlists model.Playlist

	err := c.Bind(&playlistParam)
	if err != nil {
		logger.Error.Printf("Failed processing update playlist order request: %s\n", err)
		c.String(http.StatusBadRequest, "err occurred")
		return
	}
	var playlistIds []uint

	err = streamPlatformDb.Model(&playlists).
		Where("user_id = ?", playlistParam.UserId).Order("index asc").Pluck("id", &playlistIds).Error

	playlistIdsTmp := make([]uint, len(playlistIds))
	playlistParamIdsTmp := make([]uint, len(playlistParam.PlaylistIds))

	if err != nil {
		logger.Error.Printf("Failed to get playlistIds")
		c.String(http.StatusInternalServerError, "err occurred")
		return
	}

	var needsUpdate bool
	logger.Info.Printf("playlistIds: %v / playlistparam: %v", playlistIds, playlistParam.PlaylistIds)
	if len(playlistIds) != len(playlistParam.PlaylistIds) {
		logger.Error.Printf("playlist ids length different [playlistIds] : %v, [playlistParam]: %v\n", playlistIds, playlistParam.PlaylistIds)
		c.String(http.StatusBadRequest, "playlist ids length not matched")
		return
	} else {

		copy(playlistIdsTmp, playlistIds)
		copy(playlistParamIdsTmp, playlistParam.PlaylistIds)
		sort.Slice(playlistParamIdsTmp, func(i, j int) bool {
			return playlistParamIdsTmp[i] < playlistParamIdsTmp[j]
		})
		sort.Slice(playlistIdsTmp, func(i, j int) bool {
			return playlistIdsTmp[i] < playlistIdsTmp[j]
		})

		if reflect.DeepEqual(playlistParam.PlaylistIds, playlistIds) {
			logger.Info.Println("here1")
			needsUpdate = false
		} else {

			// Check if the sorted slices are the same
			if reflect.DeepEqual(playlistParamIdsTmp, playlistIdsTmp) {
				logger.Info.Printf("[Param]: %v, [DB]: %v\n", playlistParamIdsTmp, playlistIdsTmp)
				logger.Info.Println("here2")
				needsUpdate = true // The slices have the same values in different orders
			} else {
				logger.Info.Println("here3")
				logger.Error.Printf("playlist ids different [playlistIds] : %v, [playlistParam]: %v\n", playlistIdsTmp, playlistParamIdsTmp)
				c.String(http.StatusBadRequest, "playlist ids are not matched")
				return
			}
		}
	}

	if needsUpdate {
		// Start a new transaction
		tx := streamPlatformDb.Begin()

		for i := 0; i < len(playlistParam.PlaylistIds); i++ {
			sql := `
            UPDATE "playlists"
            SET "index" = ?
            WHERE "id" = ? AND "user_id" = ?`

			err := tx.Exec(sql, i+1, playlistParam.PlaylistIds[i], playlistParam.UserId).Error

			if err != nil {
				// Rollback the transaction on error
				tx.Rollback()
				logger.Error.Printf("Failed to update index for playlist ID %d", playlistParam.PlaylistIds[i])
				c.String(http.StatusInternalServerError, "err occurred")
				return
			}
		}

		// Commit the transaction if all updates are successful
		tx.Commit()

		c.String(http.StatusOK, "ordered")
	} else {
		c.String(http.StatusOK, "no changes needed")
	}

}

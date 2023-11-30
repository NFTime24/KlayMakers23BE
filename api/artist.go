package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nftime/db"
	"github.com/nftime/logger"
	"github.com/nftime/model"
	"net/http"
	"strconv"
)

//func PostArtist(c *gin.Context) error {
//	streamPlatformDb := db.StreamPlatformDbManager()
//	params := make(map[string]string)
//	bindingParams := c.Bind(&params)
//
//	logger.Info.Printf("binding: %v\n", bindingParams)
//
//	logger.Info.Printf("params: artist_name: %v, artist_address: %v, artist_profile_id: %v\n", params["artist_name"], params["artist_address"], params["artist_profile_id"])
//	var id uint
//	var artist_id model.Artist
//	artist_name := params["artist_name"]
//	artist_address := params["artist_address"]
//	artist_profile_str := params["artist_profile_id"]
//
//	profile, _ := strconv.ParseUint(artist_profile_str, 10, 32)
//	artist_profile_id := uint(profile)
//
//	db.Model(&artist_id).Pluck("ID", &id)
//	id += 1
//	logger.Info.Printf("current id: %v\n", id)
//
//	artist_insert := model.Artist{ID: id, Name: artist_name, Address: artist_address, ProfileID: artist_profile_id}
//
//	db.Create(&artist_insert)
//	return c.JSON(http.StatusOK, params["artist_name"])
//}

// @Summary artist info
// @Description Get All Artist Info
// @Tags Artist
// @Accept json
// @Produce json
// @Deprecated True
// @Router /artist [get]
//func ShowAllArtists(c *gin.Context) error {
//
//	type Result struct {
//		Id         uint   `json:"id"`
//		Name       string `json:"name"`
//		Address    string `json:"address"`
//		Profile_id uint   `json:"profile_id"`
//	}
//
//	streamPlatformDb := db.StreamPlatformDbManager()
//	var artists model.Artist
//	var results []Result
//	rows, err := db.Model(artists).Select(`artists.*`).Rows()
//
//	if err != nil {
//		panic(err)
//	}
//	for rows.Next() {
//		db.ScanRows(rows, &results)
//	}
//	return c.JSON(http.StatusOK, results)
//}

// @Summary GetActiveArtists
// @Description Get currently active artists
// @Tags Artist
// @Accept json
// @Produce json
// @Router /artists/active [get]
// @Success 200 {object} model.ArtistResult
func GetActiveArtists(c *gin.Context) {

	streamPlatformDb := db.StreamPlatformDbManager()
	var results []model.ArtistResult
	rows, err := streamPlatformDb.Select(`DISTINCT on (name) a.*`).
		Table(`artists as a`).
		Joins("left join works as w on w.artist_name = a.name").Rows()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var test string
		rows.Scan(&test)
		fmt.Println(test)
		streamPlatformDb.ScanRows(rows, &results)
	}
	defer rows.Close()
	c.JSON(http.StatusOK, results)
}

// @Summary GetTopArtists
// @Description Get currently top artists
// @Tags Artist
// @Accept json
// @Produce json
// @Router /artists/top [get]
// @Success 200 {object} model.ArtistResult
func GetTopArtists(c *gin.Context) {

	streamPlatformDb := db.StreamPlatformDbManager()
	var results []model.ArtistResult
	rows, err := streamPlatformDb.Select(`a.*`).
		Table(`(
			select a.id, count(n.id) as count
			from artists as a
			left join works as w on w.artist_name = a.name
			left join nfts as n on n.works_id = w.id
			where w.id is not null
			group by a.id
			order by count desc
		)as base`).
		Joins("left join artists as a on base.id = a.id").Rows()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		streamPlatformDb.ScanRows(rows, &results)
	}
	defer rows.Close()
	c.JSON(http.StatusOK, results)
}

// @Summary GetArtistWithName
// @Description GetArtistWithName
// @Tags Artist
// @Accept json
// @Produce json
// @Param name path string true "name"
// @Router /artists/names/{name} [get]
// @Success 200 {object} model.ArtistResult
func GetArtistWithName(c *gin.Context) {
	name := c.Param("name")

	streamPlatformDb := db.StreamPlatformDbManager()
	var result model.ArtistResult
	streamPlatformDb.Select(`a.*`).
		Table(`artists as a`).
		Where("a.name=?", name).Scan(&result)
	logger.Info.Println(result)
	c.JSON(http.StatusOK, result)
}

// @Summary GetArtistWithId
// @Description GetArtistWithId
// @Tags Artist
// @Accept json
// @Produce json
// @Param id path string true "artist id"
// @Router /artists/ids/{id} [get]
// @Success 200 {object} model.ArtistResult
func GetArtistWithId(c *gin.Context) {
	artistIdStr := c.Param("id")
	artistId, _ := strconv.ParseUint(artistIdStr, 10, 64)
	streamPlatformDb := db.StreamPlatformDbManager()
	var result model.ArtistResult
	streamPlatformDb.Select(`a.*`).
		Table(`artists as a`).
		Where("a.id=?", artistId).Scan(&result)
	logger.Info.Println(result)
	c.JSON(http.StatusOK, result)
}

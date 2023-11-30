package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nftime/db"
	"github.com/nftime/logger"
	"github.com/nftime/model"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// @Summary update work
// @Description update work
// @Tags Work
// @Accept json
// @Produce json
// @Deprecated True
// @Param like body model.WorkInfoCreateParam true "work info data"
// @Router /workInfo [post]
//func PostWork(c *gin.Context) error {
//	streamPlatformDb := db.StreamPlatformDbManager()
//	params := make(map[string]string)
//	bindingParams := c.Bind(&params)
//
//	logger.Info.Printf("binding: %v\n", bindingParams)
//	logger.Info.Println(params["work_name"], params["work_price"], params["work_description"], params["work_category"], params["file_id"], params["artist_id"])
//
//	var id uint
//	var work_id model.Work
//	work_name := params["work_name"]
//	work_price_str := params["work_price"]
//	work_description := params["work_description"]
//	work_category := params["work_category"]
//
//	file_id_str := params["file_id"]
//	artist_id_str := params["artist_id"]
//
//	price, _ := strconv.ParseUint(work_price_str, 10, 32)
//	work_price := uint(price)
//
//	file, _ := strconv.ParseUint(file_id_str, 10, 32)
//	file_id := uint(file)
//
//	artist, _ := strconv.ParseUint(artist_id_str, 10, 32)
//	artist_id := uint(artist)
//	logger.Info.Printf("artist id :%v\n", artist_id)
//	db.Model(&work_id).Pluck("work_id", &id)
//	id += 1
//	logger.Info.Println(id)
//
//	work_insert := model.Work{WorkID: id, Name: work_name, Price: work_price, Description: work_description, Category: work_category, FileID: file_id, ArtistID: artist_id}
//
//	db.Create(&work_insert)
//	return c.JSON(http.StatusOK, params["work_name"])
//}

// @Summary Get specific NFT
// @Description Get nft info
// @Tags Work
// @Accept json
// @Produce json
// @Deprecated True
// @Param ex_id query string true "ex_id"
// @Param user_id query string true "user_id"
// @Router /getWorksInfo [get]
//func GetWorksInfoInExhibition(c *gin.Context) error {
//
//	// nft_owner := c.QueryParam("owner_address")
//	ex_id_str := c.QueryParam("ex_id")
//	ex_id, _ := strconv.ParseUint(ex_id_str, 10, 64)
//
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
//		UserId        uint   `json:"user_id"`
//		ExhibitionId  uint   `json:"exhibition_id"`
//		IsOwned       bool   `json:"is_owned"`
//	}
//
//	streamPlatformDb := db.StreamPlatformDbManager()
//	var users model.User
//	var results []Result
//	var results_user []Result
//
//	rows_user, err := db.Model(users).Select(`n.nft_id as nft_id, w.exhibitions_id as exhibition_id,
//    w.name as work_name, w.price as price, w.description as description,
//    w.category as work_category,f.filename as file_name, f.filesize as file_size,
//    f.filetype as file_type, f.path as file_path, t.path as thumbnail_path, a.name as artist_name, p.path as profile_path,
//    a.address as artist_address, users.id as user_id`).
//		Joins("left join nfts as n on users.id = n.owner_id").
//		Joins("left join works as w on n.works_id = w.work_id").
//		Joins("left join files as f on w.file_id = f.id").
//		Joins("left join files as t on f.thumbnail_id = t.id").
//		Joins("left join artists as a on w.artist_id = a.id").
//		Joins("left join files as p on a.profile_id = p.id").
//		Joins("left join exhibitions as e on e.exhibition_id = w.exhibitions_id").
//		Where("w.exhibitions_id=? and users.id=?", ex_id, user_id).Rows()
//	if err != nil {
//		panic(err)
//	}
//	for rows_user.Next() {
//		db.ScanRows(rows_user, &results_user)
//	}
//
//	logger.Info.Println(results_user)
//	rows, err := db.Model(users).Select(`n.nft_id as nft_id, w.exhibitions_id as exhibition_id,
//	 w.name as work_name, w.price as price, w.description as description,
//	 w.category as work_category,f.filename as file_name, f.filesize as file_size,
//	 f.filetype as file_type, f.path as file_path, t.path as thumbnail_path, a.name as artist_name, p.path as profile_path,
//	 a.address as artist_address, users.id as user_id`).
//		Joins("left join nfts as n on users.id = n.owner_id").
//		Joins("left join works as w on n.works_id = w.work_id").
//		Joins("left join files as f on w.file_id = f.id").
//		Joins("left join files as t on f.thumbnail_id = t.id").
//		Joins("left join artists as a on w.artist_id = a.id").
//		Joins("left join files as p on a.profile_id = p.id").
//		Joins("left join exhibitions as e on e.exhibition_id = w.exhibitions_id").
//		Where("w.exhibitions_id=?", ex_id).Rows()
//	if err != nil {
//		panic(err)
//	}
//	for rows.Next() {
//		db.ScanRows(rows, &results)
//	}
//
//	for key := range results {
//		for key := range results_user {
//			if results[key].WorkName == results_user[key].WorkName {
//				results[key].IsOwned = true
//			}
//			logger.Info.Println(key)
//		}
//		logger.Info.Println(results[key].IsOwned)
//	}
//
//	return c.JSON(http.StatusOK, results)
//}

// @Summary GetWorkInfoWithId
// @Description Get info of requested work
// @Tags Work
// @Accept json
// @Produce json
// @Param id path string true "work id"
// @Router /works/info/{id} [get]
// @Success 200 {object} model.WorkResult
func GetWorkInfoWithID(c *gin.Context) {

	// nft_owner := c.QueryParam("owner_address")
	workIdStr := c.Param("id")
	workId, err := strconv.ParseUint(workIdStr, 10, 64)
	if err != nil {
		logger.Error.Printf("err: %v\n", err)
	}
	streamPlatformDb := db.StreamPlatformDbManager()
	var results model.WorkResult
	streamPlatformDb.Select(`w.id as id,w.name as work_name, w.price as price, w.description as description, 
	 w.category as work_category, w.file_path as file_path, w.thumbnail_path as thumbnail_path, a.name as artist_name, a.profile_path as profile_path, 
	 a.address as artist_address`).
		Table("works as w").
		Joins("left join artists as a on w.artist_name = a.name").
		Where("w.id=? and w.display =?", workId, true).Scan(&results)

	logger.Info.Println(results)
	c.JSON(http.StatusOK, results)
}

// deprecated
//func GetTopWorks(c *gin.Context) error {
//	type Result struct {
//		WorkId        int    `json:"work_id"`
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
//		ExhibitionId  uint   `json:"exhibition_id"`
//	}
//
//	streamPlatformDb := db.StreamPlatformDbManager()
//	var results []Result
//	rows, err := streamPlatformDb.Select(`w.work_id as work_id, w.exhibitions_id as exhibition_id,
//	 w.name as work_name, w.price as price, w.description as description,
//	 w.category as work_category, f.filename as file_name, f.filesize as file_size,
//	 f.filetype as file_type, f.path as file_path, t.path as thumbnail_path, a.name as artist_name, p.path as profile_path,
//	 a.address as artist_address`).
//		Table(`(
//			select n.works_id, count(n.works_id) as count
//			from test.nfts as n
//			group by n.works_id
//			order by 2 desc
//			limit 10 ) as base`).
//		Joins("join works as w on base.works_id = w.work_id").
//		Joins("left join files as f on w.file_id = f.id").
//		Joins("left join files as t on f.thumbnail_id = t.id").
//		Joins("left join artists as a on w.artist_id = a.id").
//		Joins("left join files as p on a.profile_id = p.id").Rows()
//	if err != nil {
//		panic(err)
//	}
//	for rows.Next() {
//		db.ScanRows(rows, &results)
//	}
//	return c.JSON(http.StatusOK, results)
//}

// @Summary GetTopWorksWithCategory
// @Description GetTopWorksWithCategory
// @Tags Work
// @Accept json
// @Produce json
// @Param category path string true "category"
// @Router /works/top/{category} [get]
// @Success 200 {object} []model.WorkResult
func GetTopWorksWithCategory(c *gin.Context) {
	category := c.Param("category")

	//table := `(
	//	select w.*, count(w.id) as count
	//	from nftime.works as w
	//	left join nftime.nfts as n on w.id = n.works_id
	//	group by w.id
	//	order by count desc
	//	limit 10 ) as w`
	//
	//if category != "all" {
	//	table = fmt.Sprintf(`(
	//		select w.*, count(w.id) as count
	//		from nftime.works as w
	//		left join nftime.nfts as n on w.id = n.works_id
	//		where w.category = '%s'
	//		group by w.id
	//		order by count desc
	//		limit 10 ) as w`, category)
	//}

	streamPlatformDb := db.StreamPlatformDbManager()
	var results []model.WorkResult
	rows, err := streamPlatformDb.Select(`w.id as id, 
	 w.name as work_name, w.price as price, w.description as description, w.file_path as file_path, w.thumbnail_path as thumbnail_path, 
	 w.category as work_category, a.profile_path as profile_path, a.name as artist_name, a.address as artist_address`).
		Table("works as w").
		Joins("left join artists as a on w.artist_name = a.name").Where("w.category= ? and w.display =?", category, true).Limit(10).Rows()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		streamPlatformDb.ScanRows(rows, &results)
	}
	rows.Close()
	c.JSON(http.StatusOK, results)
}

// @Summary GetTodayWorks
// @Description GetTodayWorks
// @Tags Work
// @Accept json
// @Produce json
// @Router /works/today [get]
// @Success 200 {object} model.WorkResult
func GetTodayWorks(c *gin.Context) {
	now := time.Now()

	nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	logger.Info.Println(nextMidnight)
	durationUntilMidnight := nextMidnight.Sub(now)
	logger.Info.Println(durationUntilMidnight)
	durationInt := int(durationUntilMidnight.Minutes())
	c.Header("Cache-Control", "public, max-age="+strconv.Itoa(durationInt))

	fmt.Printf("Duration until the second midnight: %v\n", durationInt)

	streamPlatformDb := db.StreamPlatformDbManager()
	var results []model.WorkResult
	rows, err := streamPlatformDb.Select(`w.id as id, 
	 w.name as work_name, w.price as price, w.description as description, w.file_path as file_path, w.thumbnail_path as thumbnail_path, 
	 w.category as work_category, a.profile_path as profile_path, a.name as artist_name, a.address as artist_address`).
		Table(`(
			select * from nftime.works
			where display = true
			order by random()
			limit 6
		) as w`).
		Joins("left join artists as a on w.artist_name = a.name").Rows()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		streamPlatformDb.ScanRows(rows, &results)
	}
	c.JSON(http.StatusOK, results)
}

// @Summary GetFreeWorks
// @Description GetFreeWorks
// @Tags Work
// @Accept json
// @Produce json
// @Router /works/free [get]
// @Success 200 {object} model.WorkResult
func GetFreeWorks(c *gin.Context) {

	streamPlatformDb := db.StreamPlatformDbManager()
	var results []model.WorkResult
	rows, err := streamPlatformDb.Select(`w.id as id, 
	 w.name as work_name, w.price as price, w.description as description, w.file_path as file_path, w.thumbnail_path as thumbnail_path, 
	 w.category as work_category, a.profile_path as profile_path, a.name as artist_name, a.address as artist_address`).
		Table(`(
			select * from nftime.works
			order by random()
			limit 4
		) as w`).
		Joins("left join artists as a on w.artist_name = a.name").Where("w.display =?", true).Rows()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		streamPlatformDb.ScanRows(rows, &results)
	}
	c.JSON(http.StatusOK, results)
}

// @Summary GetNewWorks
// @Description GetNewWorks
// @Tags Work
// @Accept json
// @Produce json
// @Router /works/new [get]
// @Success 200 {object} model.WorkResult
func GetNewWorks(c *gin.Context) {

	streamPlatformDb := db.StreamPlatformDbManager()
	var results []model.WorkResult
	rows, err := streamPlatformDb.Select(`w.id as id, 
	 w.name as work_name, w.price as price, w.description as description, w.file_path as file_path, w.thumbnail_path as thumbnail_path, 
	 w.category as work_category, a.profile_path as profile_path, a.name as artist_name, a.address as artist_address`).
		Table(`(
			select * from nftime.works
			order by 1 desc
			limit 10
		) as w`).
		Joins("left join artists as a on w.artist_name = a.name").Where("w.display =?", true).Rows()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		streamPlatformDb.ScanRows(rows, &results)
	}
	c.JSON(http.StatusOK, results)
}

// @Summary GetAllWorks
// @Description GetAllWorks
// @Tags Work
// @Accept json
// @Produce json
// @Router /works/all [get]
// @Success 200 {object} model.WorkResult
func GetAllWorks(c *gin.Context) {

	streamPlatformDb := db.StreamPlatformDbManager()
	var results []model.WorkResult
	rows, err := streamPlatformDb.Select(`w.id as id, 
	 w.name as work_name, w.price as price, w.description as description, w.file_path as file_path, w.thumbnail_path as thumbnail_path, 
	 w.category as work_category, a.profile_path as profile_path, a.name as artist_name, a.address as artist_address`).
		Table(`works as w`).
		Joins("left join artists as a on w.artist_name = a.name").Where("w.display =?", true).Rows()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		streamPlatformDb.ScanRows(rows, &results)
	}
	c.JSON(http.StatusOK, results)
}

// @Summary GetArtistWorksWithName
// @Description Get works of requested artist
// @Tags Work
// @Accept json
// @Produce json
// @Param name path string true "artist name"
// @Router /works/artists/names/{name} [get]
// @Success 200 {object} model.WorkResult
func GetArtistWorksWithName(c *gin.Context) {
	artistName := c.Param("name")
	//artist_id, _ := strconv.ParseUint(artist_id_str, 10, 64)

	streamPlatformDb := db.StreamPlatformDbManager()
	var results []model.WorkResult
	rows, err := streamPlatformDb.Select(`w.id as id, 
	 w.name as work_name, w.price as price, w.description as description, w.file_path as file_path, w.thumbnail_path as thumbnail_path, 
	 w.category as work_category, a.profile_path as profile_path, a.name as artist_name, a.address as artist_address`).
		Table(`works as w`).
		Joins("left join artists as a on w.artist_name = a.name").
		Where("a.name=? and w.display = ?", artistName, true).Rows()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		streamPlatformDb.ScanRows(rows, &results)
	}
	defer rows.Close()
	c.JSON(http.StatusOK, results)
}

// @Summary GetArtistWorksWithId
// @Description Get works of requested artist
// @Tags Work
// @Accept json
// @Produce json
// @Param id path string true "artist Id"
// @Router /works/artists/ids/{id} [get]
// @Success 200 {object} model.WorkResult
func GetArtistWorksWithId(c *gin.Context) {
	artistIdStr := c.Param("id")
	artistId, _ := strconv.ParseUint(artistIdStr, 10, 64)

	streamPlatformDb := db.StreamPlatformDbManager()
	var results []model.WorkResult
	rows, err := streamPlatformDb.Select(`w.id as id, 
	 w.name as work_name, w.price as price, w.description as description, w.file_path as file_path, w.thumbnail_path as thumbnail_path, 
	 w.category as work_category, a.profile_path as profile_path, a.name as artist_name, a.address as artist_address`).
		Table(`works as w`).
		Joins("left join artists as a on w.artist_name = a.name").
		Where("a.id=? and w.display =?", artistId, true).Rows()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		streamPlatformDb.ScanRows(rows, &results)
	}
	defer rows.Close()
	c.JSON(http.StatusOK, results)
}

// @Summary GetArtistWorksWithName
// @Description GetArtistWorksWithName
// @Tags Artist
// @Accept json
// @Produce json
// @Deprecated true
// @Param artist_name query string true "artist_name"
// @Router /getArtistWorksWithName [get]
//func GetArtistWorksWithName(c *gin.Context) {
//	artist_name := c.Param("artist_name")
//
//	streamPlatformDb := db.StreamPlatformDbManager()
//	var results []model.WorkResult
//	rows, err := streamPlatformDb.Select(`w.work_id as work_id, w.exhibitions_id as exhibition_id,
//	 w.name as work_name, w.price as price, w.description as description,
//	 w.category as work_category, f.filename as file_name, f.filesize as file_size,
//	 f.filetype as file_type, f.path as file_path, t.path as thumbnail_path, a.name as artist_name, p.path as profile_path,
//	 a.address as artist_address`).
//		Table(`works as w`).
//		Joins("left join files as f on w.file_id = f.id").
//		Joins("left join files as t on f.thumbnail_id = t.id").
//		Joins("left join artists as a on w.artist_id = a.id").
//		Joins("left join files as p on a.profile_id = p.id").
//		Where("a.name=?", artist_name).Rows()
//	if err != nil {
//		panic(err)
//	}
//	for rows.Next() {
//		db.ScanRows(rows, &results)
//	}
//	c.JSON(http.StatusOK, results)
//}

// @Summary get specific work
// @Description Get works
// @Tags Work
// @Accept json
// @Produce json
// @Deprecated True
// @Param name query string true "name"
// @Router /work/specific [get]
//func GetSpecificWorkWithName(c *gin.Context) error {
//	name := c.QueryParam("name")
//	// 구조체 멤버변수 이름과 DB에서 가져오는 컬럼명이 일치해야함
//	type Result struct {
//		WorkName        string
//		ArtistName      string
//		WorkDescription string
//	}
//	streamPlatformDb := db.StreamPlatformDbManager()
//
//	// var artists model.Artist
//	var works model.Work
//	var results Result
//
//	// select w.name as work_name, a.name as artist_name, w.description from works w join artists a on w.artist_id = a.id;
//	db.Model(works).Select("works.name as work_name, works.description as work_description, artists.name as artist_name").Joins("left join artists on works.work_id = artists.id").Where("works.name=?", name).Scan(&results)
//	logger.Info.Println(results)
//	return c.JSON(http.StatusOK, results)
//}

// @Summary get top 10 works
// @Description get top 10 works
// @Tags Work
// @Accept json
// @Deprecated True
// @Produce json
// @Router /work/top10 [get]
//func GetTop10Works(c *gin.Context) error {
//	// 구조체 멤버변수 이름과 DB에서 가져오는 컬럼명이 일치해야함
//	// filepath, workname, artistname
//	type Result struct {
//		WorkName   string
//		ArtistName string
//		FilePath   string
//	}
//	streamPlatformDb := db.StreamPlatformDbManager()
//
//	// var artists model.Artist
//	var works model.Work
//	var results []Result
//
//	// select w.name as work_name, a.name as artist_name, w.description from works w join artists a on w.artist_id = a.id;
//	rows, err := db.Model(works).Select("works.name as work_name, f.path as file_path, a.name as artist_name").
//		Joins("left join files as f on works.file_id = f.id").
//		Joins("left join artists as a on works.artist_id = a.id").Rows()
//	if err != nil {
//		panic(err)
//	}
//	logger.Info.Println(rows)
//	defer rows.Close()
//	for rows.Next() {
//		db.ScanRows(rows, &results)
//	}
//	logger.Info.Println(results)
//	return c.JSON(http.StatusOK, results)
//}

// @Summary Get specific Work
// @Description Get work info in Exibition
// @Tags NFT
// @Accept json
// @Produce json
// @Deprecated True
// @Param ex_id query string true "ex_id"
// @Router /getWorksInExhibition [get]
//func GetWorksInExhibition(c *gin.Context) error {
//
//	// nft_owner := c.QueryParam("owner_address")
//	ex_id_str := c.QueryParam("ex_id")
//	ex_id, _ := strconv.ParseUint(ex_id_str, 10, 64)
//
//	type Result struct {
//		WorkId        int    `json:"work_id"`
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
//		ExhibitionId  uint   `json:"exhibition_id"`
//	}
//
//	streamPlatformDb := db.StreamPlatformDbManager()
//	// var users model.User
//	var results []Result
//
//	rows, err := streamPlatformDb.Select(`w.work_id as work_id, w.exhibitions_id as exhibition_id,
//	 w.name as work_name, w.price as price, w.description as description,
//	 w.category as work_category,f.filename as file_name, f.filesize as file_size,
//	 f.filetype as file_type, f.path as file_path, t.path as thumbnail_path, a.name as artist_name, p.path as profile_path,
//	 a.address as artist_address`).
//		Table("works as w").
//		Joins("left join files as f on w.file_id = f.id").
//		Joins("left join files as t on f.thumbnail_id = t.id").
//		Joins("left join artists as a on w.artist_id = a.id").
//		Joins("left join files as p on a.profile_id = p.id").
//		Joins("left join exhibitions as e on e.exhibition_id = w.exhibitions_id").
//		Where("w.exhibitions_id=?", ex_id).Rows()
//	if err != nil {
//		panic(err)
//	}
//	for rows.Next() {
//		db.ScanRows(rows, &results)
//	}
//
//	return c.JSON(http.StatusOK, results)
//}

// @Summary GetSelectedWorksWithName
// @Description GetSelectedWorksWithName
// @Tags Work
// @Accept json
// @Produce json
// @Param name path string true "work_name"
// @Router /works/select/names/{name} [get]
// @Success 200 {object} model.WorkResult
func GetSelectedWorksWithName(c *gin.Context) {
	selectedWorkName := c.Param("name")
	//selected_ids := strings.Split(selected_ids_str, ",")

	streamPlatformDb := db.StreamPlatformDbManager()
	var results []model.WorkResult
	rows, err := streamPlatformDb.Select(`w.id as id, 
	 w.name as work_name, w.price as price, w.description as description, w.file_path as file_path, w.thumbnail_path as thumbnail_path, 
	 w.category as work_category, a.profile_path as profile_path, a.name as artist_name, a.address as artist_address`).
		Table(`works as w`).
		Joins("left join artists as a on w.artist_name = a.name").
		Where("w.name LIKE ? and w.display = ?", "%"+selectedWorkName, true).Rows()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		streamPlatformDb.ScanRows(rows, &results)
	}
	c.JSON(http.StatusOK, results)
}

// @Summary GetSelectedWorksWithId
// @Description GetSelectedWorksWithId
// @Tags Work
// @Accept json
// @Produce json
// @Param ids path string true "ids of work"
// @Router /works/select/ids/{ids} [get]
// @Success 200 {object} model.WorkResult
func GetSelectedWorksWithId(c *gin.Context) {
	selectedIdsStr := c.Param("ids")
	selectedIds := strings.Split(selectedIdsStr, ",")
	streamPlatformDb := db.StreamPlatformDbManager()
	var results []model.WorkResult
	rows, err := streamPlatformDb.Select(`w.id as id, 
	 w.name as work_name, w.price as price, w.description as description, w.file_path as file_path, w.thumbnail_path as thumbnail_path, 
	 w.category as work_category, a.profile_path as profile_path, a.name as artist_name, a.address as artist_address`).
		Table(`works as w`).
		Joins("left join artists as a on w.artist_name = a.name").
		Where("w.id IN ? and w.display =?", selectedIds, true).Rows()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		streamPlatformDb.ScanRows(rows, &results)
	}
	defer rows.Close()
	c.JSON(http.StatusOK, results)
}

// @Summary GetSelectedWorksByteCodeWithId
// @Description GetSelectedWorksByteCodeWithId
// @Tags Work
// @Accept json
// @Produce json
// @Param ids path string true "ids of work"
// @Router /works/stream/select/ids/{ids} [get]
// @Success 200 {object} model.WorkByteCodeResult
func GetSelectedWorksByteCodeWithId(c *gin.Context) {
	selectedIdsStr := c.Param("ids")
	selectedIds := strings.Split(selectedIdsStr, ",")
	streamPlatformDb := db.StreamPlatformDbManager()
	var results []model.WorkByteCodeResult
	rows, err := streamPlatformDb.Select(`w.id as id, 
	 w.name as work_name, w.price as price, w.description as description, w.file_path as file_path, w.thumbnail_path as thumbnail_path, 
	 w.category as work_category, a.profile_path as profile_path, a.name as artist_name, a.address as artist_address`).
		Table(`works as w`).
		Joins("left join artists as a on w.artist_name = a.name").
		Where("w.id IN ? and w.display =?", selectedIds, true).Rows()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		streamPlatformDb.ScanRows(rows, &results)
	}
	defer rows.Close()

	for i, v := range results {
		fileByteCode, err := FileWithPathToByte(v.FilePath)
		if err != nil {
			logger.Error.Printf("Failed to get file bytecode: %s\n", err)
			c.String(http.StatusBadRequest, "err occured")
			return
		}
		results[i].FileByteCode = fileByteCode
		thumbnailByteCode, err := FileWithPathToByte(v.ThumbnailPath)
		if err != nil {
			logger.Error.Printf("Failed to get thumbnail bytecode: %s\n", err)
			c.String(http.StatusBadRequest, "err occured")
			return
		}
		results[i].ThumbnailByteCode = thumbnailByteCode
	}
	c.JSON(http.StatusOK, results)
}

func FileWithPathToByte(filePath string) ([]byte, error) {

	// Make an HTTP GET request to the image URL
	response, err := http.Get(filePath)
	if err != nil {
		logger.Error.Printf("err to get response", err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		logger.Error.Printf("wrong statusCode", err)
		return nil, err
	}

	imageBytes, err := ioutil.ReadAll(response.Body)

	if err != nil {
		logger.Error.Printf("unable to read", err)
		// Handle the error
		return nil, err
	}

	return imageBytes, err
}

//func GetWorkCount(c *gin.Context) (err error) {
//	user_address := c.QueryParam("address")
//	works_id := c.QueryParam("work_id")
//
//	streamPlatformDb := db.StreamPlatformDbManager()
//
//	var user_id string
//	db.Select("id").Table("users").Where("address=?", user_address).Scan(&user_id)
//
//	nft_id := -1
//	db.Select("nft_id").Table("nfts").Where("owner_id=? and works_id=?", user_id, works_id).Scan(&nft_id)
//
//	work_count := -1
//	if nft_id == -1 {
//		db.Select("count(works_id)").Table("nfts").Where("works_id=?", works_id).Scan(&work_count)
//	} else {
//		db.Select("count(works_id)").Table("nfts").Where("works_id=? and nft_id<?", works_id, nft_id).Scan(&work_count)
//	}
//
//	return c.String(http.StatusOK, strconv.Itoa(work_count))
//}

func CacheTestImage(c *gin.Context) {
	maxAgeStr := c.Query("max_age")
	defaultMaxAge := 60 // 1 hour

	maxAge, err := strconv.Atoi(maxAgeStr)
	if err != nil {
		// If the max-age parameter is not provided or invalid, use the default value
		maxAge = defaultMaxAge
	}
	c.Writer.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(maxAge))
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	filename := dir + "/assets/uploadimage/upload-865640977.jpg"

	c.File(filename)
}

func CacheTestVideo(c *gin.Context) {
	maxAgeStr := c.Query("max_age")
	defaultMaxAge := 60 // 1 hour

	maxAge, err := strconv.Atoi(maxAgeStr)
	if err != nil {
		// If the max-age parameter is not provided or invalid, use the default value
		maxAge = defaultMaxAge
	}
	c.Writer.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(maxAge))
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	filename := dir + "/assets/uploadvideo/upload-2124255261.mp4"
	c.File(filename)
}

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nftime/config"
	"github.com/nftime/db"
	"github.com/nftime/logger"
	"github.com/nftime/model"
	"net/http"
	"strconv"
	"strings"
)

// @Summary Get specific NFT
// @Description Get nft info
// @Tags NFT
// @Accept json
// @Produce json
// @Deprecated True
// @Param nft_id query string true "nft_id"
// @Router /getNFTInfoWithId [get]
//func GetNFTInfoWithId(c *gin.Context) error {
//
//	// nft_owner := c.QueryParam("owner_address")
//	nft_id_str := c.QueryParam("nft_id")
//	nft_id, _ := strconv.ParseUint(nft_id_str, 10, 64)
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
//		Where("n.nft_id=?", nft_id).Scan(&results)
//
//	logger.Info.Println(results)
//	return c.JSON(http.StatusOK, results)
//}

// @Summary GetWorkIdWithNftId
// @Description GetWorkIdWithNftId
// @Tags Work
// @Accept json
// @Produce json
// @Deprecated True
// @Param id path string true "nft id"
// @Router /works/id/nfts/{id} [get]
func GetWorkIdWithNftId(c *gin.Context) {
	nftIdStr := c.Param("id")
	nftId, _ := strconv.ParseUint(nftIdStr, 10, 64)

	streamPlatformDb := db.StreamPlatformDbManager()
	var nfts model.Nft
	var result int
	streamPlatformDb.Model(nfts).Select(`works_id`).
		Where("id=?", nftId).Scan(&result)

	resultStr := strconv.Itoa(result)
	logger.Info.Println(resultStr)
	c.String(http.StatusOK, resultStr)
}

// @Summary GetNftInfoWithWorkId
// @Description GetNftInfoWithWorkId
// @Tags NFT
// @Accept json
// @Produce json
// @Param contract_address path string true "contract_address"
// @Param work_id path string true "work id"
// @Param nft_id path string true "nft id"
// @Router /nftInfo/{contract_address}/{work_id}/{nft_id} [get]
func ResponseMetadataJson(c *gin.Context) {
	contractAddress := c.Param("contract_address")
	workIdStr := c.Param("work_id")
	nftIdStr := c.Param("nft_id")

	//hostUri := "https://secure.nftime.gallery/"
	workId, err := strconv.Atoi(workIdStr)
	if err != nil {
		logger.Error.Printf("err: %v\n", err)
	}

	nftId, err := strconv.Atoi(nftIdStr)
	if err != nil {
		logger.Error.Printf("err: %v\n", err)
	}

	var metadata model.Metadata
	switch strings.ToLower(contractAddress) {
	case strings.ToLower(config.Cfg.Contract.ContractAddress):
		streamPlatformDb := db.StreamPlatformDbManager()
		streamPlatformDb.Select(`w.name, w.description, w.file_path as image, a.name as group_name, a.profile_path as group_icon`).
			Table("works as w").Joins("left join artists as a on w.artist_name = a.name").
			Where("w.id=?", workId).Scan(&metadata)
		metadata.Sendable = true
	case strings.ToLower(config.Cfg.ContractTime.ContractAddress):
		timestorageDb := db.TimeStorageDbManager()
		var worksId uint64
		timestorageDb.Table("nfts").Select("works_id").Where("id=?", nftId).Scan(&worksId)

		timestorageDb.Select(`w.name, w.description, w.file_path as image, a.name as group_name, a.profile_path as group_icon`).
			Table("works as w").Joins("left join artists as a on w.artist_name = a.name").
			Where("w.id=?", worksId).Scan(&metadata)
		metadata.Sendable = false
	}

	c.JSON(http.StatusOK, metadata)
}

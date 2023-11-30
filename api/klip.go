package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nftime/config"
	"github.com/nftime/db"
	"github.com/nftime/logger"
	"github.com/nftime/model"
	"github.com/nftime/service"
	"github.com/umbracle/ethgo/abi"
	"html/template"
	"io"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var KlipRequestMap map[uint64]string

type KlipResponse struct {
	RequestKey     string     `json:"request_key"`
	Status         string     `json:"status"`
	Result         KlipResult `json:"result"`
	ExpirationTime int        `json:"expiration_time"`
	RequestURL     string     `json:"request_url"`
}

type KlipResult struct {
	KlaytnAddress string `json:"klaytn_address"`
}

// @Summary AddNFTWithWorkId
// @Description Add NFT with WorkID
// @Tags NFT
// @Accept json
// @Produce json
// @Deprecated True
// @Param like body model.NftWork true "work data"
// @Router /nfts/works [post]
func AddNFTWithWorkName(c *gin.Context) {

	workBody := model.NftWork{}
	err := c.Bind(&workBody)
	if err != nil {
		logger.Error.Printf("Failed processing Binding: %v\n", err)
	}
	workId, err := strconv.ParseUint(workBody.WorkID, 10, 64)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}

	streamPlatformDb := db.StreamPlatformDbManager()

	streamPlatformDb.Create(&model.Nft{
		WorksID:     uint(workId),
		UserAddress: workBody.Address,
	})

	var nft model.Nft
	var result int

	streamPlatformDb.Model(nft).Select(`MAX(id)`).Scan(&result)

	resultStr := strconv.Itoa(result)
	c.String(http.StatusOK, resultStr)
}

// @Summary MintArtWithoutPaying
// @Description MintArtWithoutPaying
// @Tags Klip
// @Accept json
// @Produce json
// @Param id path string true "work id"
// @Router /klip/mint/swf/work/{id} [get]
func MintArtWithoutPaying(c *gin.Context) {

	workIdStr := c.Param("id")
	//workId, err := strconv.ParseUint(workIdStr, 10, 64)
	//if err != nil {
	//	logger.Error.Printf("err :%v\n", err)
	//}
	klipKey := rand.Uint64()
	reqBodyStr := fmt.Sprintf(`{
		"type": "auth",
		"bapp": {
			"name" : "NFTime",
			"callback": { "success": "https:\/\/%s\/klip\/success-swf\/?key=%s&work_id=%s", "fail": "" }		}
	}`, config.Cfg.Server.Host, strconv.FormatUint(klipKey, 10), workIdStr)
	reqBody := bytes.NewBufferString(reqBodyStr)
	resp, err := http.Post("https://a2a-api.klipwallet.com/v2/a2a/prepare", "Content-Type: application/json", reqBody)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	var jData KlipResponse
	logger.Info.Printf("body: %v\n", body)
	json.Unmarshal(body, &jData)

	logger.Info.Printf("requestKey :%v\n", jData.RequestKey)
	KlipRequestMap[klipKey] = jData.RequestKey
	jData.RequestURL = "https://klipwallet.com/?target=/a2a?request_key="
	jData.RequestURL += jData.RequestKey

	logger.Info.Printf("requestUrl :%v\n", jData.RequestURL)

	// http.Redirect(w, r, jData.RequestQR, http.StatusFound)
	c.Redirect(http.StatusFound, jData.RequestURL)
}

func MintMembershipNft(c *gin.Context) {

	klipKey := rand.Uint64()
	reqBodyStr := fmt.Sprintf(`{
		"type": "auth",
		"bapp": {
			"name" : "NFTime",
			"callback": { "success": "https:\/\/%s\/time\/success-membership\/?key=%s", "fail": "" }		}
	}`, config.Cfg.Server.Host, strconv.FormatUint(klipKey, 10))
	reqBody := bytes.NewBufferString(reqBodyStr)
	resp, err := http.Post("https://a2a-api.klipwallet.com/v2/a2a/prepare", "Content-Type: application/json", reqBody)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	var jData KlipResponse
	logger.Info.Printf("body: %v\n", body)
	json.Unmarshal(body, &jData)

	logger.Info.Printf("requestKey :%v\n", jData.RequestKey)
	KlipRequestMap[klipKey] = jData.RequestKey
	jData.RequestURL = "https://klipwallet.com/?target=/a2a?request_key="
	jData.RequestURL += jData.RequestKey

	logger.Info.Printf("requestUrl :%v\n", jData.RequestURL)

	// http.Redirect(w, r, jData.RequestQR, http.StatusFound)
	c.Redirect(http.StatusFound, jData.RequestURL)
}

func OnSuccessMembershipKlip(c *gin.Context) {
	//c.Query("work_id")
	WEBTOON_BASIC_WORK_ID := uint64(2)
	WEBTOON_ADVANCED_WORK_ID := uint64(3)
	WEBTOON_PLATINUM_WORK_ID := uint64(4)

	WEBTOON_WORK_IDS := []uint64{2, 3, 4}
	klipKeyStr := c.Query("key")
	//

	klipKey, err := strconv.ParseUint(klipKeyStr, 10, 64)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}

	requestKey := KlipRequestMap[klipKey]

	logger.Info.Printf("request key: %v\n", requestKey)

	client1 := &http.Client{}
	reqStr1 := fmt.Sprintf("https://a2a-api.klipwallet.com/v2/a2a/result?request_key=%s", requestKey)
	req1, err := http.NewRequest("GET", reqStr1, nil)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	req1.Header.Add("Content-Type", "application/json")
	resp1, err := client1.Do(req1)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	body1, err := io.ReadAll(resp1.Body)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	var jData1 KlipResponse

	json.Unmarshal(body1, &jData1)

	logger.Info.Printf("Klaytn address: %s\n", jData1.Result.KlaytnAddress)
	address := jData1.Result.KlaytnAddress
	resp1.Body.Close()

	timeStorageDb := db.TimeStorageDbManager()
	streamPlatformDb := db.StreamPlatformDbManager()
	var nfts model.Nft
	artCounts := -1
	webtoonMembershipCounts := -1
	streamPlatformDb.Model(nfts).Select("count(works_id)").Where("user_address =?", address).Scan(&artCounts)

	timeStorageDb.Model(nfts).Select("count(works_id)").Where("works_id IN ? and user_address =?", WEBTOON_WORK_IDS, address).Scan(&webtoonMembershipCounts)

	//log.Printf("[%v] Art Counts : %v, Membership Counts : %v", address, webtoonArtCounts, webtoonMembershipCounts)
	var memberShipWorkId uint64

	if webtoonMembershipCounts <= 0 {

		log.Println("here1")
		var lastNftId uint64

		timeStorageDb.Model(nfts).Select(`MAX(id)`).Scan(&lastNftId)

		logger.Info.Println(lastNftId)

		if artCounts <= 1 {
			memberShipWorkId = WEBTOON_BASIC_WORK_ID
			fmt.Println("2")
		} else if artCounts >= 2 && artCounts < 4 {

			memberShipWorkId = WEBTOON_ADVANCED_WORK_ID
			fmt.Println("3")
		} else {
			memberShipWorkId = WEBTOON_PLATINUM_WORK_ID
			fmt.Println("4")
		}
		timeStorageDb.Create(&model.Nft{
			WorksID:     uint(memberShipWorkId),
			UserAddress: address,
		})
		requestBody := model.MintToAddrParam{Id: strconv.FormatUint(memberShipWorkId, 10), UserAddress: address}
		logger.Info.Printf("here: %v\n", requestBody)
		reqeustBodyByte, err := json.Marshal(requestBody)
		requestBodyReader := bytes.NewReader(reqeustBodyByte)

		if err != nil {
			logger.Error.Printf("err: %v\n", err)
		}
		//reqStr2 := fmt.Sprintf("http://34.212.84.161/mintToAddr?address=%s&work_name=%s", jData1.Result.KlaytnAddress, workName)
		reqStr2 := fmt.Sprintf("https://%s/time/mint", config.Cfg.Server.Host)
		logger.Debug.Printf("reqhost: %v\n", reqStr2)
		resp2, err := http.Post(reqStr2, "application/json", requestBodyReader)
		if err != nil {
			logger.Error.Printf("err :%v\n", err)
		}
		body2, err := io.ReadAll(resp2.Body)
		if err != nil {
			logger.Error.Printf("err :%v\n", err)
		}
		logger.Info.Printf("body2: %s \n", body2)

		resString := string(body2[:])

		logger.Info.Printf("result: %s\n", resString)

		resp2.Body.Close()
	} else if artCounts >= 2 && artCounts < 4 {
		log.Println("here2")

		memberShipWorkId = WEBTOON_ADVANCED_WORK_ID
		timeStorageDb.Model(nfts).Where("user_address=? and works_id = ?", address, WEBTOON_BASIC_WORK_ID).Update("works_id", memberShipWorkId)

	} else if artCounts >= 4 {
		log.Println("here4")

		memberShipWorkId = WEBTOON_PLATINUM_WORK_ID
		timeStorageDb.Model(nfts).Where("user_address=? and works_id = ?", address, WEBTOON_ADVANCED_WORK_ID).Update("works_id", memberShipWorkId)

	} else {
		log.Println("here err")

		c.String(http.StatusBadRequest, "err occured")
		return
	}

	//if jData1이 not used면 success, 아니면 failed

	var results model.WorkResult
	timeStorageDb.Select(`w.id as id,w.name as work_name, w.price as price, w.description as description,
		w.category as work_category, w.file_path as file_path, w.thumbnail_path as thumbnail_path, a.name as artist_name, a.profile_path as profile_path,
		a.address as artist_address`).
		Table("works as w").
		Joins("left join artists as a on w.artist_name = a.name").
		Where("w.id=?", memberShipWorkId).Scan(&results)

	logger.Info.Println(results)

	c.HTML(http.StatusOK, "mbti_result.html", gin.H{
		"CSS":    template.CSS("<link rel='stylesheet' href='/static/nft_page.css'>"),
		"Title":  results.WorkName,
		"Image":  results.FilePath,
		"Artist": results.ArtistName,
		"Ticket": 0,
	})

}
func WebtoonFairMintArtWithoutPaying(c *gin.Context) {

	currentTime := time.Now().UTC()
	targetTimeStr := "2023-10-08T09:00:00Z"
	targetTime, err := time.Parse(time.RFC3339, targetTimeStr)
	if err != nil {
		fmt.Println("Error parsing target time:", err)
		return
	}

	if currentTime.After(targetTime) {
		logger.Error.Printf("WebtoonFair finished")
		c.String(http.StatusBadRequest, "WebtoonFair finished")
		return
	}

	logger.Info.Printf("curTime: %v\n", currentTime)
	workIdStr := c.Param("id")

	klipKey := rand.Uint64()
	reqBodyStr := fmt.Sprintf(`{
		"type": "auth",
		"bapp": {
			"name" : "NFTime",
			"callback": { "success": "https:\/\/%s\/time\/success-webtoon-fair\/?key=%s&work_id=%s", "fail": "" }		}
	}`, config.Cfg.Server.Host, strconv.FormatUint(klipKey, 10), workIdStr)
	reqBody := bytes.NewBufferString(reqBodyStr)
	resp, err := http.Post("https://a2a-api.klipwallet.com/v2/a2a/prepare", "Content-Type: application/json", reqBody)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	var jData KlipResponse
	logger.Info.Printf("body: %v\n", body)
	json.Unmarshal(body, &jData)

	logger.Info.Printf("requestKey :%v\n", jData.RequestKey)
	KlipRequestMap[klipKey] = jData.RequestKey
	jData.RequestURL = "https://klipwallet.com/?target=/a2a?request_key="
	jData.RequestURL += jData.RequestKey

	logger.Info.Printf("requestUrl :%v\n", jData.RequestURL)

	// http.Redirect(w, r, jData.RequestQR, http.StatusFound)
	c.Redirect(http.StatusFound, jData.RequestURL)
}

func MintArtWorkWithoutPaying(c *gin.Context) {
	workIdStr := c.Param("id")
	workId, err := strconv.ParseUint(workIdStr, 10, 64)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	klipKey := rand.Uint64()
	reqBodyStr := fmt.Sprintf(`{
		"type": "auth",
		"bapp": {
			"name" : "NFTime",
			"callback": { "success": "https:\/\/%s\/klip\/success-work\/?key=%s&work_id=%s", "fail": "" }
		}
	}`, config.Cfg.Server.Host, strconv.FormatUint(klipKey, 10), strconv.FormatUint(workId, 10))
	reqBody := bytes.NewBufferString(reqBodyStr)
	resp, err := http.Post("https://a2a-api.klipwallet.com/v2/a2a/prepare", "Content-Type: application/json", reqBody)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	var jData KlipResponse
	logger.Info.Printf("body: %v\n", body)
	json.Unmarshal(body, &jData)

	logger.Info.Printf("requestKey :%v\n", jData.RequestKey)
	KlipRequestMap[klipKey] = jData.RequestKey
	jData.RequestURL = "https://klipwallet.com/?target=/a2a?request_key="
	jData.RequestURL += jData.RequestKey

	logger.Info.Printf("requestUrl :%v\n", jData.RequestURL)

	// http.Redirect(w, r, jData.RequestQR, http.StatusFound)
	c.Redirect(http.StatusFound, jData.RequestURL)
}

func OnSuccessWorkKlip(c *gin.Context) {
	klipKeyStr := c.Query("key")
	workIdStr := c.Query("work_id")
	fmt.Println(klipKeyStr, workIdStr)

	klipKey, err := strconv.ParseUint(klipKeyStr, 10, 64)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	//workId, err := strconv.ParseUint(workIdStr, 10, 64)
	//if err != nil {
	//	logger.Error.Printf("err :%v\n", err)
	//}

	requestKey := KlipRequestMap[klipKey]

	logger.Info.Printf("request key: %v\n", requestKey)

	client1 := &http.Client{}
	reqStr1 := fmt.Sprintf("https://a2a-api.klipwallet.com/v2/a2a/result?request_key=%s", requestKey)
	req1, err := http.NewRequest("GET", reqStr1, nil)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	req1.Header.Add("Content-Type", "application/json")
	resp1, err := client1.Do(req1)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	body1, err := io.ReadAll(resp1.Body)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	var jData1 KlipResponse
	logger.Info.Printf("body1 :%v\n", body1)
	json.Unmarshal(body1, &jData1)

	logger.Info.Printf("Klaytn address: %s\n", jData1.Result.KlaytnAddress)
	address := jData1.Result.KlaytnAddress
	var userAddress string
	var ticketCount int64
	ticketCount = -1
	resp1.Body.Close()
	streamPlatformDb := db.StreamPlatformDbManager()

	streamPlatformDb.Select("user_address").Table("tickets").Where("user_address=?", address).Scan(&userAddress)

	logger.Info.Printf("address: %v\n", userAddress)
	if userAddress == "" {
		ticketInsert := model.Ticket{UserAddress: address, TicketCount: 3, IsMbtiMinted: false}
		streamPlatformDb.Create(&ticketInsert)
	}

	streamPlatformDb.Select(`ticket_count`).
		Table("tickets").
		Where("user_address=?", address).Scan(&ticketCount)

	if ticketCount > 0 {
		//if jData1이 not used면 success, 아니면 failed
		//client2 := &http.Client{}
		requestBody := model.MintToAddrParam{Id: workIdStr, UserAddress: address}

		reqeustBodyByte, err := json.Marshal(requestBody)
		requestBodyReader := bytes.NewReader(reqeustBodyByte)

		if err != nil {
			logger.Error.Printf("err: %v\n", err)
		}
		//reqStr2 := fmt.Sprintf("http://34.212.84.161/mintToAddr?address=%s&work_name=%s", jData1.Result.KlaytnAddress, workName)
		reqStr2 := fmt.Sprintf("https://%s/klip/mint", config.Cfg.Server.Host)
		logger.Debug.Printf("reqhost: %v\n", reqStr2)
		resp2, err := http.Post(reqStr2, "application/json", requestBodyReader)
		if err != nil {
			logger.Error.Printf("err :%v\n", err)
		}
		//resp2, err := client2.Do(req2)
		//if err != nil {
		//	logger.Error.Printf("err :%v\n", err)
		//}
		body2, err := io.ReadAll(resp2.Body)
		if err != nil {
			logger.Error.Printf("err :%v\n", err)
		}
		logger.Info.Printf("body2: %s \n", body2)

		resString := string(body2[:])

		logger.Info.Printf("result: %s\n", resString)

		resp2.Body.Close()
		ticketCount -= 1
		streamPlatformDb.Select("ticket_count").Table("tickets").Where("user_address=?", address).Update("ticket_count", ticketCount)

		c.HTML(http.StatusOK, "work_result.html", gin.H{
			"CSS": template.CSS("<link rel='stylesheet' href='/static/nft_page.css'>"),
		})
	} else {
		c.HTML(http.StatusOK, "work_result.html", gin.H{
			"CSS": template.CSS("<link rel='stylesheet' href='/static/nft_page.css'>"),
		})
	}
}

func OnSuccessSwfKlip(c *gin.Context) {
	//c.Query("work_id")
	klipKeyStr := c.Query("key")
	workIdStr := c.Query("work_id")
	//
	fmt.Println(klipKeyStr, workIdStr)
	decryptedWorkId := service.VerifyQueryParam(workIdStr)
	logger.Info.Printf("decrypedWorkId :%v\n", decryptedWorkId)
	if decryptedWorkId == "" {
		logger.Error.Printf("unable to get workId from nftime")
		c.String(http.StatusBadRequest, "wrong request")
	} else {
		klipKey, err := strconv.ParseUint(klipKeyStr, 10, 64)
		if err != nil {
			logger.Error.Printf("err :%v\n", err)
		}

		requestKey := KlipRequestMap[klipKey]

		logger.Info.Printf("request key: %v\n", requestKey)

		client1 := &http.Client{}
		reqStr1 := fmt.Sprintf("https://a2a-api.klipwallet.com/v2/a2a/result?request_key=%s", requestKey)
		req1, err := http.NewRequest("GET", reqStr1, nil)
		if err != nil {
			logger.Error.Printf("err :%v\n", err)
		}
		req1.Header.Add("Content-Type", "application/json")
		resp1, err := client1.Do(req1)
		if err != nil {
			logger.Error.Printf("err :%v\n", err)
		}
		body1, err := io.ReadAll(resp1.Body)
		if err != nil {
			logger.Error.Printf("err :%v\n", err)
		}
		var jData1 KlipResponse

		json.Unmarshal(body1, &jData1)

		logger.Info.Printf("Klaytn address: %s\n", jData1.Result.KlaytnAddress)
		address := jData1.Result.KlaytnAddress
		var isMinted bool
		resp1.Body.Close()

		streamPlatformDb := db.StreamPlatformDbManager()

		var userAddress string
		streamPlatformDb.Select("user_address").Table("tickets").Where("user_address=?", address).Scan(&userAddress)

		type WorkResult struct {
			Id            int    `json:"id"`
			WorkName      string `json:"work_name"`
			Price         int    `json:"work_price"`
			Description   string `json:"description"`
			WorkCategory  string `json:"category"`
			FilePath      string `json:"file_path"`
			ThumbnailPath string `json:"thumbnail_path"`
			ArtistName    string `json:"artist_name"`
			ProfilePath   string `json:"profile_path"`
			ArtistAddress string `json:"artist_address"`
		}

		if userAddress == "" {
			ticketInsert := model.Ticket{UserAddress: address, TicketCount: 3, IsMbtiMinted: false, IsSwfMinted: false}
			streamPlatformDb.Create(&ticketInsert)
		}
		streamPlatformDb.Select(`is_swf_minted`).
			Table("tickets").
			Where("user_address=?", address).Scan(&isMinted)
		if !isMinted {
			//if jData1이 not used면 success, 아니면 failed
			requestBody := model.MintToAddrParam{Id: decryptedWorkId, UserAddress: address}
			logger.Info.Printf("here: %v\n", requestBody)
			reqeustBodyByte, err := json.Marshal(requestBody)
			requestBodyReader := bytes.NewReader(reqeustBodyByte)

			if err != nil {
				logger.Error.Printf("err: %v\n", err)
			}
			//reqStr2 := fmt.Sprintf("http://34.212.84.161/mintToAddr?address=%s&work_name=%s", jData1.Result.KlaytnAddress, workName)
			reqStr2 := fmt.Sprintf("https://%s/klip/mint", config.Cfg.Server.Host)
			logger.Debug.Printf("reqhost: %v\n", reqStr2)
			resp2, err := http.Post(reqStr2, "application/json", requestBodyReader)
			if err != nil {
				logger.Error.Printf("err :%v\n", err)
			}
			body2, err := io.ReadAll(resp2.Body)
			if err != nil {
				logger.Error.Printf("err :%v\n", err)
			}
			logger.Info.Printf("body2: %s \n", body2)

			resString := string(body2[:])

			logger.Info.Printf("result: %s\n", resString)

			resp2.Body.Close()

			var results model.WorkResult
			streamPlatformDb.Select(`w.id as id,w.name as work_name, w.price as price, w.description as description,
		w.category as work_category, w.file_path as file_path, w.thumbnail_path as thumbnail_path, a.name as artist_name, a.profile_path as profile_path,
		a.address as artist_address`).
				Table("works as w").
				Joins("left join artists as a on w.artist_name = a.name").
				Where("w.id=?", decryptedWorkId).Scan(&results)

			logger.Info.Println(results)
			streamPlatformDb.Select("is_swf_minted").Table("tickets").Where("user_address=?", address).Update("is_swf_minted", true)

			c.HTML(http.StatusOK, "mbti_result.html", gin.H{
				"CSS":    template.CSS("<link rel='stylesheet' href='/static/nft_page.css'>"),
				"Title":  results.WorkName,
				"Image":  results.FilePath,
				"Artist": results.ArtistName,
				"Ticket": 0,
			})

		} else {
			c.HTML(http.StatusOK, "mbti_fail_result.html", gin.H{
				"CSS":     template.CSS("<link rel='stylesheet' href='/static/nft_page.css'>"),
				"Address": address,
			})
		}
	}
}

func OnSuccessWebtoonFairKlip(c *gin.Context) {
	//c.Query("work_id")
	klipKeyStr := c.Query("key")
	workIdStr := c.Query("work_id")
	//
	fmt.Println(klipKeyStr, workIdStr)
	decryptedWorkId := service.VerifyQueryParam(workIdStr)
	logger.Info.Printf("decrypedWorkId :%v\n", decryptedWorkId)
	if decryptedWorkId == "" {
		logger.Error.Printf("unable to get workId from nftime")
		c.String(http.StatusBadRequest, "wrong request")
	} else {
		klipKey, err := strconv.ParseUint(klipKeyStr, 10, 64)
		if err != nil {
			logger.Error.Printf("err :%v\n", err)
		}

		requestKey := KlipRequestMap[klipKey]

		logger.Info.Printf("request key: %v\n", requestKey)

		client1 := &http.Client{}
		reqStr1 := fmt.Sprintf("https://a2a-api.klipwallet.com/v2/a2a/result?request_key=%s", requestKey)
		req1, err := http.NewRequest("GET", reqStr1, nil)
		if err != nil {
			logger.Error.Printf("err :%v\n", err)
		}
		req1.Header.Add("Content-Type", "application/json")
		resp1, err := client1.Do(req1)
		if err != nil {
			logger.Error.Printf("err :%v\n", err)
		}
		body1, err := io.ReadAll(resp1.Body)
		if err != nil {
			logger.Error.Printf("err :%v\n", err)
		}
		var jData1 KlipResponse

		json.Unmarshal(body1, &jData1)

		logger.Info.Printf("Klaytn address: %s\n", jData1.Result.KlaytnAddress)
		address := jData1.Result.KlaytnAddress
		var isMinted bool
		resp1.Body.Close()

		TimeStorageDb := db.TimeStorageDbManager()

		var userAddress string
		TimeStorageDb.Select("user_address").Table("tickets").Where("user_address=?", address).Scan(&userAddress)

		type WorkResult struct {
			Id            int    `json:"id"`
			WorkName      string `json:"work_name"`
			Price         int    `json:"work_price"`
			Description   string `json:"description"`
			WorkCategory  string `json:"category"`
			FilePath      string `json:"file_path"`
			ThumbnailPath string `json:"thumbnail_path"`
			ArtistName    string `json:"artist_name"`
			ProfilePath   string `json:"profile_path"`
			ArtistAddress string `json:"artist_address"`
		}

		if userAddress == "" {
			ticketInsert := model.Ticket{UserAddress: address, TicketCount: 3}
			TimeStorageDb.Create(&ticketInsert)
		}
		TimeStorageDb.Select(`is_webtoon_minted`).
			Table("tickets").
			Where("user_address=?", address).Scan(&isMinted)
		if !isMinted {
			//if jData1이 not used면 success, 아니면 failed
			requestBody := model.MintToAddrParam{Id: decryptedWorkId, UserAddress: address}
			logger.Info.Printf("here: %v\n", requestBody)
			reqeustBodyByte, err := json.Marshal(requestBody)
			requestBodyReader := bytes.NewReader(reqeustBodyByte)

			if err != nil {
				logger.Error.Printf("err: %v\n", err)
			}
			//reqStr2 := fmt.Sprintf("http://34.212.84.161/mintToAddr?address=%s&work_name=%s", jData1.Result.KlaytnAddress, workName)
			reqStr2 := fmt.Sprintf("https://%s/time/mint", config.Cfg.Server.Host)
			logger.Debug.Printf("reqhost: %v\n", reqStr2)
			resp2, err := http.Post(reqStr2, "application/json", requestBodyReader)
			if err != nil {
				logger.Error.Printf("err :%v\n", err)
			}
			body2, err := io.ReadAll(resp2.Body)
			if err != nil {
				logger.Error.Printf("err :%v\n", err)
			}
			logger.Info.Printf("body2: %s \n", body2)

			resString := string(body2[:])

			logger.Info.Printf("result: %s\n", resString)

			resp2.Body.Close()

			var results model.WorkResult
			TimeStorageDb.Select(`w.id as id,w.name as work_name, w.price as price, w.description as description,
		w.category as work_category, w.file_path as file_path, w.thumbnail_path as thumbnail_path, a.name as artist_name, a.profile_path as profile_path,
		a.address as artist_address`).
				Table("works as w").
				Joins("left join artists as a on w.artist_name = a.name").
				Where("w.id=?", decryptedWorkId).Scan(&results)

			logger.Info.Println(results)
			TimeStorageDb.Select("is_webtoon_minted").Table("tickets").Where("user_address=?", address).Update("is_webtoon_minted", true)

			c.HTML(http.StatusOK, "mbti_result.html", gin.H{
				"CSS":    template.CSS("<link rel='stylesheet' href='/static/nft_page.css'>"),
				"Title":  results.WorkName,
				"Image":  results.FilePath,
				"Artist": results.ArtistName,
				"Ticket": 0,
			})

		} else {
			c.HTML(http.StatusOK, "mbti_fail_result.html", gin.H{
				"CSS":     template.CSS("<link rel='stylesheet' href='/static/nft_page.css'>"),
				"Address": address,
			})
		}
	}
}

func OnSuccessKlip(c *gin.Context) {
	klipKeyStr := c.Query("key")
	workIdStr := c.Query("work_id")

	fmt.Println(klipKeyStr, workIdStr)
	klipKey, err := strconv.ParseUint(klipKeyStr, 10, 64)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}

	requestKey := KlipRequestMap[klipKey]

	logger.Info.Printf("request key: %v\n", requestKey)

	client1 := &http.Client{}
	reqStr1 := fmt.Sprintf("https://a2a-api.klipwallet.com/v2/a2a/result?request_key=%s", requestKey)
	req1, err := http.NewRequest("GET", reqStr1, nil)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	req1.Header.Add("Content-Type", "application/json")
	resp1, err := client1.Do(req1)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	body1, err := io.ReadAll(resp1.Body)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	var jData1 KlipResponse
	logger.Info.Printf("body1 :%v\n", body1)
	json.Unmarshal(body1, &jData1)

	logger.Info.Printf("Klaytn address: %s\n", jData1.Result.KlaytnAddress)
	address := jData1.Result.KlaytnAddress
	var isMinted bool
	resp1.Body.Close()
	streamPlatformDb := db.StreamPlatformDbManager()

	var userAddress string
	streamPlatformDb.Select("user_address").Table("tickets").Where("user_address=?", address).Scan(&userAddress)

	type WorkResult struct {
		Id            int    `json:"id"`
		WorkName      string `json:"work_name"`
		Price         int    `json:"work_price"`
		Description   string `json:"description"`
		WorkCategory  string `json:"category"`
		FilePath      string `json:"file_path"`
		ThumbnailPath string `json:"thumbnail_path"`
		ArtistName    string `json:"artist_name"`
		ProfilePath   string `json:"profile_path"`
		ArtistAddress string `json:"artist_address"`
	}

	if userAddress == "" {
		ticketInsert := model.Ticket{UserAddress: address, TicketCount: 3, IsMbtiMinted: false}
		streamPlatformDb.Create(&ticketInsert)
	}
	streamPlatformDb.Select(`is_mbti_minted`).
		Table("tickets").
		Where("user_address=?", address).Scan(&isMinted)
	if !isMinted {
		//if jData1이 not used면 success, 아니면 failed
		requestBody := model.MintToAddrParam{Id: workIdStr, UserAddress: address}
		logger.Info.Printf("here: %v\n", requestBody)
		reqeustBodyByte, err := json.Marshal(requestBody)
		requestBodyReader := bytes.NewReader(reqeustBodyByte)

		if err != nil {
			logger.Error.Printf("err: %v\n", err)
		}
		//reqStr2 := fmt.Sprintf("http://34.212.84.161/mintToAddr?address=%s&work_name=%s", jData1.Result.KlaytnAddress, workName)
		reqStr2 := fmt.Sprintf("https://%s/klip/mint", config.Cfg.Server.Host)
		logger.Debug.Printf("reqhost: %v\n", reqStr2)
		resp2, err := http.Post(reqStr2, "application/json", requestBodyReader)
		if err != nil {
			logger.Error.Printf("err :%v\n", err)
		}
		body2, err := io.ReadAll(resp2.Body)
		if err != nil {
			logger.Error.Printf("err :%v\n", err)
		}
		logger.Info.Printf("body2: %s \n", body2)

		resString := string(body2[:])

		logger.Info.Printf("result: %s\n", resString)

		resp2.Body.Close()

		var results model.WorkResult
		streamPlatformDb.Select(`w.id as id,w.name as work_name, w.price as price, w.description as description, 
	 w.category as work_category, w.file_path as file_path, w.thumbnail_path as thumbnail_path, a.name as artist_name, a.profile_path as profile_path, 
	 a.address as artist_address`).
			Table("works as w").
			Joins("left join artists as a on w.artist_name = a.name").
			Where("w.id=?", workIdStr).Scan(&results)

		logger.Info.Println(results)
		streamPlatformDb.Select("is_mbti_minted").Table("tickets").Where("user_address=?", address).Update("is_mbti_minted", true)

		c.HTML(http.StatusOK, "mbti_result.html", gin.H{
			"CSS":    template.CSS("<link rel='stylesheet' href='/static/nft_page.css'>"),
			"Title":  results.WorkName,
			"Image":  results.FilePath,
			"Artist": results.ArtistName,
			"Ticket": 0,
		})

	} else {
		c.HTML(http.StatusOK, "mbti_fail_result.html", gin.H{
			"CSS":     template.CSS("<link rel='stylesheet' href='/static/nft_page.css'>"),
			"Address": address,
		})
	}

}

// @Summary MintToAddrApple
// @Description MintToAddrApple
// @Tags Klip
// @Accept json
// @Produce json
// @Param info body model.UserInfoParam true "work_id and user_address"
// @Router /klip/mint/apple [post]
func MintToAddrApple(c *gin.Context) {
	streamPlatformDb := db.StreamPlatformDbManager()
	var ticketCount int

	userParamBody := model.UserInfoParam{}
	err := c.Bind(&userParamBody)
	if err != nil {
		logger.Error.Printf("Failed processing Binding: %v\n", err)
	}

	address := userParamBody.UserAddress
	streamPlatformDb.Select(`ticket_count`).
		Table("tickets").
		Where("user_address=?", address).Scan(&ticketCount)
	workIdStr := userParamBody.Id
	workId, err := strconv.ParseUint(workIdStr, 10, 64)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	logger.Info.Printf("Work id : %d\n", workId)

	if ticketCount > 0 {
		var kasBody []byte
		if len(address) == 42 {
			var nfts model.Nft
			var result uint64

			streamPlatformDb.Model(nfts).Select(`MAX(id)`).Scan(&result)

			logger.Info.Println(result)
			newItemId := result + 1

			streamPlatformDb.Create(&model.Nft{
				WorksID:     uint(workId),
				UserAddress: address,
			})

			logger.Info.Printf("Klaytn address: %s\n", address)

			typ := abi.MustNewType("uint256")

			nftId_big := big.NewInt(int64(newItemId))
			nftId_encoded, err := typ.Encode(nftId_big)
			if err != nil {
				panic(err)
			}
			nftId_hex := fmt.Sprintf("%x", nftId_encoded)

			workId_big := big.NewInt(int64(workId))
			workId_encoded, err := typ.Encode(workId_big)
			if err != nil {
				panic(err)
			}
			workId_hex := fmt.Sprintf("%x", workId_encoded)

			addressBase := "0000000000000000000000000000000000000000000000000000000000000000"
			ablen := len(addressBase)
			kalen := len(address)
			addr_hex := fmt.Sprintf("%s%s", addressBase[:(ablen-kalen+2)], address[2:])

			reqCallData := "0x20b7668b"
			reqCallData += addr_hex
			reqCallData += nftId_hex
			reqCallData += workId_hex

			logger.Info.Printf("req call data :%v\n", reqCallData)
			kasClient := &http.Client{}
			kasReqStr := fmt.Sprintf("https://wallet-api.klaytnapi.com/v2/tx/contract/execute")
			jsonStr := fmt.Sprintf(`{
		"from": "0x7c07C1579aD1980863c83876EC4bec43BC8d6dFa",
		"value": "0x0",
		"to": "%s",
		"input": "%s",
		"nonce": 0,
		"gasLimit": 1000000,
		"submit": true
	}`, config.Cfg.Contract.ContractAddress, reqCallData)
			kasReq, err := http.NewRequest("POST", kasReqStr, bytes.NewBufferString(jsonStr))
			if err != nil {
				logger.Error.Printf("err :%v\n", err)
			}
			kasReq.Header.Add("x-chain-id", "8217")
			kasReq.Header.Add("Content-Type", "application/json")
			kasReq.Header.Add("Authorization", "Basic S0FTS0NDRjIxR1VZUUdCOE83Q0JQR09GOm1waHN0cTllSDFTV1d6cXNFX3JrTEM0LTRCMDVFYWhyWmg5SVNFbWI=")
			kasResp, err := kasClient.Do(kasReq)
			if err != nil {
				logger.Error.Printf("err :%v\n", err)
			}
			defer kasResp.Body.Close()
			kasBody, err = io.ReadAll(kasResp.Body)
			if err != nil {
				logger.Error.Printf("err :%v\n", err)
			}
			logger.Info.Printf("kas body: %s \n", kasBody)
		}
		streamPlatformDb.Select("ticket_count").Table("tickets").Where("user_address=?", address).Update("ticket_count", ticketCount-1)

		c.String(http.StatusOK, string(kasBody))
	} else {
		responseStr := fmt.Sprintf("address %s used all tickets", address)
		c.String(http.StatusForbidden, responseStr)
	}
}

// @Summary MintToAddr
// @Description MintToAddr
// @Tags Klip
// @Accept json
// @Produce json
// @Param info body model.UserInfoParam true "work_id and user_address"
// @Router /klip/mint [post]
func MintToAddr(c *gin.Context) {

	WEBTOON_BASIC_WORK_ID := uint64(2)
	WEBTOON_ADVANCED_WORK_ID := uint64(3)
	WEBTOON_PLATINUM_WORK_ID := uint64(4)

	WEBTOON_WORK_IDS := []uint64{2, 3, 4}

	streamPlatformDb := db.StreamPlatformDbManager()

	timeStorageDb := db.TimeStorageDbManager()

	userParamBody := model.MintToAddrParam{}
	err := c.Bind(&userParamBody)
	if err != nil {
		logger.Error.Printf("Failed processing Binding: %v\n", err)
	}

	address := userParamBody.UserAddress

	workIdStr := userParamBody.Id
	workId, err := strconv.ParseUint(workIdStr, 10, 64)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	logger.Info.Printf("Work id : %d\n", workId)

	var nfts model.Nft
	var lastNftId uint64

	streamPlatformDb.Model(nfts).Select(`MAX(id)`).Scan(&lastNftId)

	logger.Info.Println(lastNftId)
	newNftId := lastNftId + 1

	streamPlatformDb.Create(&model.Nft{
		WorksID:     uint(workId),
		UserAddress: address,
	})

	logger.Info.Printf("Klaytn address: %s\n", address)

	kasBody := requestKas(address, newNftId, workId, config.Cfg.Contract.ContractAddress)

	//var webtoonArtWorkIds []uint
	//webtoonArtWorkIds = append(webtoonArtWorkIds, 178, 179, 180, 181, 182, 183, 184, 185, 186)
	logger.Info.Printf("kas body: %s \n", kasBody)

	artCounts := -1
	webtoonMembershipCounts := -1
	streamPlatformDb.Model(nfts).Select("count(works_id)").Where("user_address =?", address).Scan(&artCounts)

	timeStorageDb.Model(nfts).Select("count(works_id)").Where("works_id IN ? and user_address =?", WEBTOON_WORK_IDS, address).Scan(&webtoonMembershipCounts)

	log.Printf("[%v] Art Counts : %v, Membership Counts : %v", address, artCounts, webtoonMembershipCounts)

	if webtoonMembershipCounts <= 0 {

		log.Println("here1")
		timeStorageDb.Model(nfts).Select(`MAX(id)`).Scan(&lastNftId)

		logger.Info.Println(lastNftId)
		newNftId = lastNftId + 1

		var memberShipWorkId uint64

		if artCounts <= 1 {
			memberShipWorkId = WEBTOON_BASIC_WORK_ID
			logger.Info.Println("2")
		} else if artCounts >= 2 && artCounts < 4 {

			memberShipWorkId = WEBTOON_ADVANCED_WORK_ID
			logger.Info.Println("3")
		} else {
			memberShipWorkId = WEBTOON_PLATINUM_WORK_ID
			logger.Info.Println("4")
		}
		timeStorageDb.Create(&model.Nft{
			WorksID:     uint(memberShipWorkId),
			UserAddress: address,
		})

		kasBody := requestKas(address, newNftId, WEBTOON_BASIC_WORK_ID, config.Cfg.ContractTime.ContractAddress)

		log.Printf("kasBody: %v\n", kasBody)
	} else if artCounts >= 2 && artCounts < 4 {
		log.Println("here2")

		timeStorageDb.Model(nfts).Where("user_address=? and works_id = ?", address, WEBTOON_BASIC_WORK_ID).Update("works_id", WEBTOON_ADVANCED_WORK_ID)

	} else if artCounts >= 4 {
		log.Println("here4")

		timeStorageDb.Model(nfts).Where("user_address=? and (works_id = ? or works_id = ?)", address, WEBTOON_BASIC_WORK_ID, WEBTOON_ADVANCED_WORK_ID).Update("works_id", WEBTOON_PLATINUM_WORK_ID)

	} else {
		log.Println("here err")

		c.String(http.StatusBadRequest, "err occured")
		return
	}

	c.String(http.StatusOK, string(kasBody))

}

func requestKas(address string, newNftId uint64, workId uint64, contractAddress string) []byte {
	typ := abi.MustNewType("uint256")

	nftId_big := big.NewInt(int64(newNftId))
	nftId_encoded, err := typ.Encode(nftId_big)
	if err != nil {
		panic(err)
	}
	nftId_hex := fmt.Sprintf("%x", nftId_encoded)

	workId_big := big.NewInt(int64(workId))
	workId_encoded, err := typ.Encode(workId_big)
	if err != nil {
		panic(err)
	}
	workId_hex := fmt.Sprintf("%x", workId_encoded)

	addressBase := "0000000000000000000000000000000000000000000000000000000000000000"
	ablen := len(addressBase)
	kalen := len(address)
	addr_hex := fmt.Sprintf("%s%s", addressBase[:(ablen-kalen+2)], address[2:])

	reqCallData := "0x20b7668b"
	reqCallData += addr_hex
	reqCallData += nftId_hex
	reqCallData += workId_hex

	logger.Info.Printf("req call data :%v\n", reqCallData)
	kasClient := &http.Client{}
	kasReqStr := fmt.Sprintf("https://wallet-api.klaytnapi.com/v2/tx/contract/execute")
	jsonStr := fmt.Sprintf(`{
		"from": "0x7c07C1579aD1980863c83876EC4bec43BC8d6dFa",
		"value": "0x0",
		"to": "%s",
		"input": "%s",
		"nonce": 0,
		"gasLimit": 1000000,
		"submit": true
	}`, contractAddress, reqCallData)
	kasReq, err := http.NewRequest("POST", kasReqStr, bytes.NewBufferString(jsonStr))
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	kasReq.Header.Add("x-chain-id", "8217")
	kasReq.Header.Add("Content-Type", "application/json")
	kasReq.Header.Add("Authorization", "Basic S0FTS0NDRjIxR1VZUUdCOE83Q0JQR09GOm1waHN0cTllSDFTV1d6cXNFX3JrTEM0LTRCMDVFYWhyWmg5SVNFbWI=")
	kasResp, err := kasClient.Do(kasReq)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
		return nil
	}
	defer kasResp.Body.Close()
	kasBody, err := io.ReadAll(kasResp.Body)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
		return nil
	}
	return kasBody
}

// @Summary MintToAddrTime
// @Description MintToAddrTime
// @Tags Klip
// @Accept json
// @Produce json
// @Param info body model.UserInfoParam true "work_id and user_address"
// @Router /time/mint [post]
func MintToAddrTime(c *gin.Context) {
	timeStorageDb := db.TimeStorageDbManager()
	userParamBody := model.MintToAddrParam{}
	err := c.Bind(&userParamBody)
	if err != nil {
		logger.Error.Printf("Failed processing Binding: %v\n", err)
	}

	address := userParamBody.UserAddress

	workIdStr := userParamBody.Id
	workId, err := strconv.ParseUint(workIdStr, 10, 64)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	logger.Info.Printf("Work id : %d\n", workId)

	var nfts model.Nft
	var result uint64

	timeStorageDb.Model(nfts).Select(`MAX(id)`).Scan(&result)

	logger.Info.Println(result)
	newItemId := result + 1

	timeStorageDb.Create(&model.Nft{
		WorksID:     uint(workId),
		UserAddress: address,
	})

	logger.Info.Printf("Klaytn address: %s\n", address)

	typ := abi.MustNewType("uint256")

	nftId_big := big.NewInt(int64(newItemId))
	nftId_encoded, err := typ.Encode(nftId_big)
	if err != nil {
		panic(err)
	}
	nftId_hex := fmt.Sprintf("%x", nftId_encoded)

	workId_big := big.NewInt(int64(workId))
	workId_encoded, err := typ.Encode(workId_big)
	if err != nil {
		panic(err)
	}
	workId_hex := fmt.Sprintf("%x", workId_encoded)

	addressBase := "0000000000000000000000000000000000000000000000000000000000000000"
	ablen := len(addressBase)
	kalen := len(address)
	addr_hex := fmt.Sprintf("%s%s", addressBase[:(ablen-kalen+2)], address[2:])

	reqCallData := "0x20b7668b"
	reqCallData += addr_hex
	reqCallData += nftId_hex
	reqCallData += workId_hex

	logger.Info.Printf("req call data :%v\n", reqCallData)
	kasClient := &http.Client{}
	kasReqStr := fmt.Sprintf("https://wallet-api.klaytnapi.com/v2/tx/contract/execute")
	jsonStr := fmt.Sprintf(`{
		"from": "0x7c07C1579aD1980863c83876EC4bec43BC8d6dFa",
		"value": "0x0",
		"to": "%s",
		"input": "%s",
		"nonce": 0,
		"gasLimit": 1000000,
		"submit": true
	}`, config.Cfg.ContractTime.ContractAddress, reqCallData)
	kasReq, err := http.NewRequest("POST", kasReqStr, bytes.NewBufferString(jsonStr))
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	kasReq.Header.Add("x-chain-id", "8217")
	kasReq.Header.Add("Content-Type", "application/json")
	kasReq.Header.Add("Authorization", "Basic S0FTS0NDRjIxR1VZUUdCOE83Q0JQR09GOm1waHN0cTllSDFTV1d6cXNFX3JrTEM0LTRCMDVFYWhyWmg5SVNFbWI=")
	kasResp, err := kasClient.Do(kasReq)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	defer kasResp.Body.Close()
	kasBody, err := io.ReadAll(kasResp.Body)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	logger.Info.Printf("kas body: %s \n", kasBody)

	c.String(http.StatusOK, string(kasBody))

}

//func MintToAddrWithId(c *gin.Context) error {
//	address := c.QueryParam("address")
//
//	workId_str := c.QueryParam("work_id")
//	workId, err := strconv.ParseUint(workId_str, 10, 64)
//	if err != nil {
//		logger.Error.Printf("err :%v\n", err)
//	}
//
//	nftId_str := c.QueryParam("nft_id")
//	nftId, err := strconv.ParseUint(nftId_str, 10, 64)
//	if err != nil {
//		logger.Error.Printf("err :%v\n", err)
//	}
//
//	streamPlatformDb := db.StreamPlatformDbManager()
//
//	newItemId := uint(nftId)
//
//	var user_id uint
//	db.Select("id").Table("users").Where("address=?", address).Scan(&user_id)
//
//	typ := abi.MustNewType("uint256")
//
//	nftId_big := big.NewInt(int64(newItemId))
//	nftId_encoded, err := typ.Encode(nftId_big)
//	if err != nil {
//		panic(err)
//	}
//	nftId_hex := fmt.Sprintf("%x", nftId_encoded)
//
//	workId_big := big.NewInt(int64(workId))
//	workId_encoded, err := typ.Encode(workId_big)
//	if err != nil {
//		panic(err)
//	}
//	workId_hex := fmt.Sprintf("%x", workId_encoded)
//
//	addressBase := "0000000000000000000000000000000000000000000000000000000000000000"
//	ablen := len(addressBase)
//	kalen := len(address)
//	addr_hex := fmt.Sprintf("%s%s", addressBase[:(ablen-kalen+2)], address[2:])
//
//	reqCallData := "0x20b7668b"
//	reqCallData += addr_hex
//	reqCallData += nftId_hex
//	reqCallData += workId_hex
//
//	logger.Info.Printf("req call data :%v\n", reqCallData)
//
//	kasClient := &http.Client{}
//	kasReqStr := fmt.Sprintf("https://wallet-api.klaytnapi.com/v2/tx/contract/execute")
//	jsonStr := fmt.Sprintf(`{
//		"from": "0x7c07C1579aD1980863c83876EC4bec43BC8d6dFa",
//		"value": "0x0",
//		"to": "%s",
//		"input": "%s",
//		"nonce": 0,
//		"gasLimit": 1000000,
//		"submit": true
//	}`, ContractAddress, reqCallData)
//	kasReq, err := http.NewRequest("POST", kasReqStr, bytes.NewBufferString(jsonStr))
//	if err != nil {
//		logger.Error.Printf("err :%v\n", err)
//	}
//	kasReq.Header.Add("x-chain-id", "8217")
//	kasReq.Header.Add("Content-Type", "application/json")
//	kasReq.Header.Add("Authorization", "Basic S0FTS0NDRjIxR1VZUUdCOE83Q0JQR09GOm1waHN0cTllSDFTV1d6cXNFX3JrTEM0LTRCMDVFYWhyWmg5SVNFbWI=")
//	kasResp, err := kasClient.Do(kasReq)
//	if err != nil {
//		logger.Error.Printf("err :%v\n", err)
//	}
//	defer kasResp.Body.Close()
//	kasBody, err := io.ReadAll(kasResp.Body)
//	if err != nil {
//		logger.Error.Printf("err :%v\n", err)
//	}
//	logger.Info.Printf("kas body: %s \n", kasBody)
//
//	return c.String(http.StatusOK, string(kasBody))
//}

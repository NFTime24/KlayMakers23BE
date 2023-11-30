package api

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nftime/config"
	"github.com/nftime/db"
	"github.com/nftime/logger"
	"github.com/nftime/model"
	"github.com/umbracle/ethgo/abi"
	"io"
	"math/big"
	"net/http"
	"strconv"
)

func GetCertificateList(c *gin.Context) {

	timeStorageDb := db.TimeStorageDbManager()
	var certificateList []model.Certificate
	rows, err := timeStorageDb.Select(`*`).Table(`certificates`).Order(`id`).Rows()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	for rows.Next() {
		timeStorageDb.ScanRows(rows, &certificateList)
	}
	defer rows.Close()
	c.JSON(http.StatusOK, certificateList)

}

func GetUserCertificateList(c *gin.Context) {
	userAddress := c.Query("wallet_address")

	timeStorageDb := db.TimeStorageDbManager()
	var certificateUserList []model.CertificateUserList

	rows, err := timeStorageDb.Table("certificate_users c").
		Select("c.user_wallet_address, ct.certificate_name, ct.company_name, ct.certificate_description, ct.certificate_category, ct.certificate_image, ct.certificate_thumbnail, ct.certificate_website, ct.certificate_start_date, ct.certificate_end_date").
		Joins("LEFT JOIN certificates ct ON c.certificate_id = ct.id").
		Where("user_wallet_address = ?", userAddress).
		Order("c.id asc").
		Rows()

	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	for rows.Next() {
		timeStorageDb.ScanRows(rows, &certificateUserList)
	}
	defer rows.Close()
	c.JSON(http.StatusOK, certificateUserList)

}

// @Summary MintCertificateToAddr
// @Description MintCertificateToAddr
// @Tags Klip
// @Accept json
// @Produce json
// @Param info body model.CertificateIssueParam true "certificate_id, certificate_name, wallet_address"
// @Router /back-office/certificate/issue [post]
func IssueCertificate(c *gin.Context) {
	timeStorageDb := db.TimeStorageDbManager()
	userParamBody := model.CertificateIssueParam{}
	err := c.Bind(&userParamBody)
	if err != nil {
		logger.Error.Printf("Failed processing Binding: %v\n", err)
	}

	address := userParamBody.WalletAddress

	workIdStr := userParamBody.Id
	workId, err := strconv.ParseUint(workIdStr, 10, 64)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	logger.Info.Printf("Work id : %d\n", workId)

	var certificateUser model.CertificateUser
	var result uint64

	timeStorageDb.Model(certificateUser).Select(`MAX(id)`).Scan(&result)

	logger.Info.Println(result)
	newItemId := result + 1

	timeStorageDb.Create(&model.CertificateUser{
		CertificateId:     uint(workId),
		UserWalletAddress: address,
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

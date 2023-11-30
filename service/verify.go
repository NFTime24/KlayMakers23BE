package service

import (
	"encoding/base64"
	"encoding/hex"
	"github.com/nftime/logger"
	"github.com/nftime/util"
	"strings"
)

func VerifyQueryParam(workIdStr string) string {

	prefix := "nftime"

	//redisCli := redis.RedisManager()
	//keyString := fmt.Sprintf("%s:%s", "qr", "aes_key")
	//aesHexKey, err := redisCli.Get(context.Background(), keyString).Result()
	//if err != nil {
	//	logger.Error.Printf("err:%v\n", err)
	//	aesHexKey = "d4e9a43db1f266980f3cca85f00ca47606b9d92b82e8fac632c5b45090eeb5b4"
	//}
	aesHexKey := "d4e9a43db1f266980f3cca85f00ca47606b9d92b82e8fac632c5b45090eeb5b4"

	aesKey, err := hex.DecodeString(aesHexKey)
	if err != nil {
		logger.Error.Printf("Error decoding key from hexadecimal: %v\n", err)
		return ""
	}

	decodedData, err := base64.RawURLEncoding.DecodeString(workIdStr)
	if err != nil {
		logger.Error.Printf("Error decoding data: %v\n", err)
		return ""
	}
	logger.Info.Printf("Decoded data: %v\n", string(decodedData))

	decryptedData, err := util.Decrypt(decodedData, aesKey)
	if err != nil {
		logger.Error.Printf("err: %v\n", err)
		return ""
	}

	decryptedStr := string(decryptedData)
	if strings.HasPrefix(decryptedStr, prefix) {
		return util.RemoveLetters(decryptedStr)
	} else {
		return ""
	}
}

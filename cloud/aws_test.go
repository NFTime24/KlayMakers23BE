package cloud

import (
	"fmt"
	"github.com/nftime/util"
	"os"
	"testing"
)

func TestGetBucketList(t *testing.T) {
	var s3Info S3Info
	s3Info.AwsS3Region = "ap-northeast-2"
	s3Info.BucketName = "nft-slider-bucket"
	s3Info.AwsAccessKey = "AKIA2LTV37TSHHW6RFCK"
	s3Info.AwsSecretKey = "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"
	s3Info.SetS3ConfigByKey()
	s3Info.GetBucketList()
}

func TestGetItemInBucket(t *testing.T) {
	var s3Info S3Info
	s3Info.AwsS3Region = "ap-northeast-2"
	s3Info.BucketName = "nft-slider-bucket"
	s3Info.AwsAccessKey = "AKIA2LTV37TSHHW6RFCK"
	s3Info.AwsSecretKey = "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"
	s3Info.SetS3ConfigByKey()
	prefix := "images/"
	s3Info.GetItems(prefix)
}

func TestUploadItemIntoBucket(t *testing.T) {
	var s3Info S3Info
	s3Info.AwsS3Region = "ap-northeast-2"
	s3Info.BucketName = "nft-slider-bucket"
	s3Info.AwsAccessKey = "AKIA2LTV37TSHHW6RFCK"
	s3Info.AwsSecretKey = "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"

	f, err := os.Open("../uploads/test.png")
	if err != nil {
		t.Errorf("err while reading File: %v\n", t)
	}
	s3Info.SetS3ConfigByKey()

	//fmt.Println(s3Info.S3Client)
	fileName := util.CreateUnixFilename("test")
	result := s3Info.UploadFile(f, fileName, "images/")
	fmt.Println(result.Location)
	//originBucketPath := result.Location
	//fmt.Println(originBucketPath)
	//resizedBucketPath := strings.Replace(originBucketPath, "images/", "resized-images/", 1)

	//fmt.Println(resizedBucketPath)
}

func TestReadFileFromPath(t *testing.T) {
	f, err := os.Open("../uploads/test.png")
	if err != nil {
		t.Errorf("err while reading File: %v\n", t)
	}
	fmt.Println(f)
}

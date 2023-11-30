package api

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"github.com/nftime/config"
	"github.com/nftime/db"
	"github.com/nftime/logger"
	"github.com/nftime/model"
	"github.com/nftime/util"
	"html/template"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func StreamLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "stream_login.html", nil)
}

func TimeLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "time_login.html", nil)
}

func StreamIndexPage(c *gin.Context) {
	type ArtistNames struct {
		ArtistName string `json:"artist_name"`
	}

	streamPlatformDb := db.StreamPlatformDbManager()
	var artistNames []ArtistNames
	rows, err := streamPlatformDb.Select(`name as artist_name`).
		Table(`artists`).Rows()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		streamPlatformDb.ScanRows(rows, &artistNames)
	}
	defer rows.Close()
	logger.Info.Println(artistNames)
	c.HTML(http.StatusOK, "stream_index.html", gin.H{
		"CSS":         template.CSS("<link rel='stylesheet' href='/static/style.css'>"),
		"ArtistNames": artistNames,
	})
}

func StreamUploadWorks(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}

	files := form.File["file"]

	category := form.Value["category"][0]
	workName := form.Value["workname"][0]
	artistName := form.Value["artist"][0]
	priceStr := form.Value["price"][0]
	price, err := strconv.ParseUint(priceStr, 10, 64)
	if err != nil {
		logger.Error.Printf("err :%v\n", err)
	}
	description := form.Value["description"][0]
	logger.Info.Println(files, category, workName, artistName, price, description)
	var originalLocation string
	var cdnLocation string
	var resizedCdnLocation string
	for _, file := range files {
		filename := filepath.Base(file.Filename)
		uploadedFile, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer uploadedFile.Close()

		fileNameParts := strings.Split(filename, ".")
		filename = fileNameParts[0]
		fmt.Println(description)
		fileName := util.CreateUnixFilename(filename)

		config.StreamS3Info.BucketName = config.Cfg.AWS.Stream.WorkBucket
		result := config.StreamS3Info.UploadFile(uploadedFile, fileName, "images/")

		img, _, err := image.Decode(uploadedFile)
		if err != nil {
			log.Fatal(err)
		}
		uploadedFile.Seek(0, 0)
		var buf bytes.Buffer

		// Resize the image to the desired dimensions (e.g., 300 pixels wide)
		resizedImg := resize.Resize(300, 0, img, resize.Lanczos3)

		fmt.Printf("fileName: %v", fileName)
		// Create the output file

		fileExt := filepath.Ext(filename)
		switch strings.ToLower(fileExt) {
		case "png":
			err = png.Encode(&buf, resizedImg)
			if err != nil {
				logger.Error.Printf("error to encode jpeg: %v\n", err)
			}
		case "gif":
			err = gif.Encode(&buf, resizedImg, nil)
			if err != nil {
				logger.Error.Printf("error to encode jpeg: %v\n", err)
			}
		default:
			err = jpeg.Encode(&buf, resizedImg, nil)
			if err != nil {
				logger.Error.Printf("error to encode jpeg: %v\n", err)
			}
		}

		//
		tempThumbnailFile, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer tempThumbnailFile.Close()
		_, err = buf.WriteTo(tempThumbnailFile)
		if err != nil {
			log.Fatal(err)
		}
		buf.Reset()
		tempThumbnailFile.Seek(0, 0)

		fileInfo, _ := tempThumbnailFile.Stat()
		logger.Debug.Printf("file: %v\n", fileInfo.Size())
		resultThumbnail := config.StreamS3Info.UploadFile(tempThumbnailFile, fileName, "resized-images/")

		originalLocation = result.Location
		s3Url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com", config.StreamS3Info.BucketName, config.StreamS3Info.AwsS3Region)
		cdnLocation = strings.Replace(originalLocation, s3Url, config.Cfg.AWS.Stream.WorkBucketCloudFront, 1)
		resizedCdnLocation = resultThumbnail.Location
		resizedCdnLocation = strings.Replace(resizedCdnLocation, s3Url, config.Cfg.AWS.Stream.WorkBucketCloudFront, 1)
		fmt.Println(originalLocation)
		fmt.Println(resizedCdnLocation)

		os.Remove(fileName)

	}
	streamPlatformDb := db.StreamPlatformDbManager()

	streamPlatformDb.Create(&model.Work{
		Name:          workName,
		ArtistName:    artistName,
		Price:         uint(price),
		Description:   description,
		Category:      category,
		FilePath:      cdnLocation,
		ThumbnailPath: resizedCdnLocation,
	})

	c.Redirect(http.StatusFound, "/stream/back-office")

}

func StreamUploadArtists(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}

	files := form.File["file"]

	name := form.Value["name"][0]
	address := form.Value["address"][0]
	instagram := form.Value["instagram"][0]

	introduction := form.Value["introduction"][0]

	var originalLocation string

	for _, file := range files {
		filename := filepath.Base(file.Filename)
		uploadedFile, err := file.Open()

		if err != nil {
			log.Fatal(err)
		}
		defer uploadedFile.Close()

		// Create a temporary file
		tempFile, err := os.CreateTemp("", "upload-*.tmp")
		if err != nil {
			log.Fatal(err)
		}
		defer tempFile.Close()
		fileNameParts := strings.Split(filename, ".")
		filename = fileNameParts[0]
		fileName := util.CreateUnixFilename(filename)

		config.StreamS3Info.BucketName = config.Cfg.AWS.Stream.ProfileBucket

		result := config.StreamS3Info.UploadFile(uploadedFile, fileName, "images/")

		originalLocation = result.Location
		fmt.Println(originalLocation)
	}
	streamPlatformDb := db.StreamPlatformDbManager()

	streamPlatformDb.Create(&model.Artist{
		Name:         name,
		Address:      address,
		ProfilePath:  originalLocation,
		Introduction: introduction,
		Instagram:    instagram,
	})

	c.Redirect(http.StatusFound, "/stream/back-office/artists/upload")

}

func StreamShowArtistsList(c *gin.Context) {

	streamPlatformDb := db.StreamPlatformDbManager()
	var artists []model.Artist
	rows, err := streamPlatformDb.Select(`*`).Table(`artists`).Order(`id`).Rows()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	for rows.Next() {
		streamPlatformDb.ScanRows(rows, &artists)
	}
	defer rows.Close()
	c.HTML(http.StatusOK, `stream_artist_list.html`, gin.H{
		"Artists": artists,
	})
}

func StreamDeleteArtistsList(c *gin.Context) {
	streamPlatformDb := db.StreamPlatformDbManager()
	artistID := c.PostForm("id")
	logger.Info.Printf("artist ID : %v\n", artistID)
	var artist model.Artist
	result := streamPlatformDb.Delete(&artist, artistID)
	if result.Error != nil {
		logger.Error.Printf("err: %v\n", result.Error)
	}
	c.Redirect(http.StatusFound, "/stream/back-office/artists/list")
}

func StreamShowWorksList(c *gin.Context) {

	streamPlatformDb := db.StreamPlatformDbManager()
	var works []model.Work
	rows, err := streamPlatformDb.Select(`*`).Table(`works`).Order(`id`).Rows()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	for rows.Next() {
		streamPlatformDb.ScanRows(rows, &works)
	}
	defer rows.Close()
	c.HTML(http.StatusOK, `stream_work_list.html`, gin.H{
		"Works": works,
	})
}

func StreamDeleteWorksList(c *gin.Context) {

	streamPlatformDb := db.StreamPlatformDbManager()
	workID := c.PostForm("id")
	logger.Info.Printf("work ID : %v\n", workID)
	var work model.Work
	result := streamPlatformDb.Delete(&work, workID)
	if result.Error != nil {
		logger.Error.Printf("err: %v\n", result.Error)
	}
	c.Redirect(http.StatusFound, "/stream/back-office/works/list")
}

func StreamArtistsPage(c *gin.Context) {

	c.HTML(http.StatusOK, "stream_artist.html", gin.H{
		"CSS": template.CSS("<link rel='stylesheet' href='/static/style.css'>"),
	})
}

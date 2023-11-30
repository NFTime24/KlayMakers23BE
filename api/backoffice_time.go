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
	"time"
)

func TimeIndexPage(c *gin.Context) {
	type ArtistNames struct {
		ArtistName string `json:"artist_name"`
	}

	timeStorageDb := db.TimeStorageDbManager()
	var artistNames []ArtistNames
	rows, err := timeStorageDb.Select(`name as artist_name`).
		Table(`artists`).Rows()
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		timeStorageDb.ScanRows(rows, &artistNames)
	}
	defer rows.Close()
	logger.Info.Println(artistNames)
	c.HTML(http.StatusOK, "time_index.html", gin.H{
		"CSS":         template.CSS("<link rel='stylesheet' href='/static/style.css'>"),
		"ArtistNames": artistNames,
	})
}

func TimeUploadWorks(c *gin.Context) {
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

		config.TimeS3Info.BucketName = config.Cfg.AWS.Time.WorkBucket
		result := config.TimeS3Info.UploadFile(uploadedFile, fileName, "images/")

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
		resultThumbnail := config.TimeS3Info.UploadFile(tempThumbnailFile, fileName, "resized-images/")

		originalLocation = result.Location
		s3Url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com", config.TimeS3Info.BucketName, config.TimeS3Info.AwsS3Region)
		cdnLocation = strings.Replace(originalLocation, s3Url, config.Cfg.AWS.Time.WorkBucketCloudFront, 1)
		resizedCdnLocation = resultThumbnail.Location
		resizedCdnLocation = strings.Replace(resizedCdnLocation, s3Url, config.Cfg.AWS.Time.WorkBucketCloudFront, 1)
		fmt.Println(originalLocation)
		fmt.Println(resizedCdnLocation)

		os.Remove(fileName)

	}
	timeStorageDb := db.TimeStorageDbManager()

	timeStorageDb.Create(&model.Work{
		Name:          workName,
		ArtistName:    artistName,
		Price:         uint(price),
		Description:   description,
		Category:      category,
		FilePath:      cdnLocation,
		ThumbnailPath: resizedCdnLocation,
	})

	c.Redirect(http.StatusFound, "/time/back-office")

}

func TimeUploadArtists(c *gin.Context) {
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

		config.TimeS3Info.BucketName = config.Cfg.AWS.Time.ProfileBucket

		result := config.TimeS3Info.UploadFile(uploadedFile, fileName, "images/")

		originalLocation = result.Location
		fmt.Println(originalLocation)
	}
	timeStorageDb := db.TimeStorageDbManager()

	timeStorageDb.Create(&model.Artist{
		Name:         name,
		Address:      address,
		ProfilePath:  originalLocation,
		Introduction: introduction,
		Instagram:    instagram,
	})

	c.Redirect(http.StatusFound, "/time/back-office/artists/upload")

}

func TimeShowArtistsList(c *gin.Context) {

	timeStorageDb := db.TimeStorageDbManager()
	var artists []model.Artist
	rows, err := timeStorageDb.Select(`*`).Table(`artists`).Order(`id`).Rows()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	for rows.Next() {
		timeStorageDb.ScanRows(rows, &artists)
	}
	defer rows.Close()
	c.HTML(http.StatusOK, `time_artist_list.html`, gin.H{
		"Artists": artists,
	})
}

func TimeDeleteArtistsList(c *gin.Context) {
	timeStorageDb := db.TimeStorageDbManager()
	artistID := c.PostForm("id")
	logger.Info.Printf("artist ID : %v\n", artistID)
	var artist model.Artist
	result := timeStorageDb.Delete(&artist, artistID)
	if result.Error != nil {
		logger.Error.Printf("err: %v\n", result.Error)
	}
	c.Redirect(http.StatusFound, "/time/back-office/artists/list")
}

func TimeShowWorksList(c *gin.Context) {

	timeStorageDb := db.TimeStorageDbManager()
	var works []model.Work
	rows, err := timeStorageDb.Select(`*`).Table(`works`).Order(`id`).Rows()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	for rows.Next() {
		timeStorageDb.ScanRows(rows, &works)
	}
	defer rows.Close()
	c.HTML(http.StatusOK, `time_work_list.html`, gin.H{
		"Works": works,
	})
}

func TimeDeleteWorksList(c *gin.Context) {

	timeStorageDb := db.TimeStorageDbManager()
	workID := c.PostForm("id")
	logger.Info.Printf("work ID : %v\n", workID)
	var work model.Work
	result := timeStorageDb.Delete(&work, workID)
	if result.Error != nil {
		logger.Error.Printf("err: %v\n", result.Error)
	}
	c.Redirect(http.StatusFound, "/time/back-office/works/list")
}

func TimeArtistsPage(c *gin.Context) {

	c.HTML(http.StatusOK, "time_artist.html", gin.H{
		"CSS": template.CSS("<link rel='stylesheet' href='/static/style.css'>"),
	})
}

func TimeUploadCertificate(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}

	certificateImage := form.File["certificateImage"]
	certificateCategory := form.Value["certificateCategory"][0]
	certificateName := form.Value["certificateName"][0]
	companyName := form.Value["companyName"][0]
	certificateWebsite := form.Value["certificateWebsite"][0]
	certificateStartDate := form.Value["certificateStartDate"][0]
	certificateEndDate := form.Value["certificateEndDate"][0]
	certificateDescription := form.Value["certificateDescription"][0]

	var originalLocation string
	var cdnLocation string
	var resizedCdnLocation string
	for _, file := range certificateImage {
		filename := filepath.Base(file.Filename)
		uploadedFile, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer uploadedFile.Close()

		fileNameParts := strings.Split(filename, ".")
		filename = fileNameParts[0]
		fileName := util.CreateUnixFilename(filename)

		config.TimeS3Info.BucketName = config.Cfg.AWS.Time.WorkBucket
		result := config.TimeS3Info.UploadFile(uploadedFile, fileName, "images/")

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
		resultThumbnail := config.TimeS3Info.UploadFile(tempThumbnailFile, fileName, "resized-images/")

		originalLocation = result.Location
		s3Url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com", config.TimeS3Info.BucketName, config.TimeS3Info.AwsS3Region)
		cdnLocation = strings.Replace(originalLocation, s3Url, config.Cfg.AWS.Time.WorkBucketCloudFront, 1)
		resizedCdnLocation = resultThumbnail.Location
		resizedCdnLocation = strings.Replace(resizedCdnLocation, s3Url, config.Cfg.AWS.Time.WorkBucketCloudFront, 1)
		fmt.Println(originalLocation)
		fmt.Println(resizedCdnLocation)

		os.Remove(fileName)

	}
	timeStorageDb := db.TimeStorageDbManager()

	parsedStartDateTime, err := time.Parse(certificateStartDate, "2006-01-02T15:04:05")
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return
	}
	parsedEndDateTime, err := time.Parse(certificateEndDate, "2006-01-02T15:04:05")
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return
	}
	timeStorageDb.Create(&model.Certificate{
		CertificateName:        certificateName,
		CompanyName:            companyName,
		CertificateDescription: certificateDescription,
		CertificateCategory:    certificateCategory,
		CertificateImage:       cdnLocation,
		CertificateThumbnail:   resizedCdnLocation,
		CertificateWebsite:     certificateWebsite,
		CertificateStartDate:   parsedStartDateTime,
		CertificateEndDate:     parsedEndDateTime,
	})

	c.Redirect(http.StatusFound, "/time/back-office")
}

func TimeCountStat(c *gin.Context) {

	timeStorageDb := db.TimeStorageDbManager()
	workID := c.PostForm("id")
	logger.Info.Printf("work ID : %v\n", workID)
	var work model.Work
	result := timeStorageDb.Delete(&work, workID)
	if result.Error != nil {
		logger.Error.Printf("err: %v\n", result.Error)
	}
	c.Redirect(http.StatusFound, "/time/back-office/works/list")
}

func RegisterCertificate(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}

	certificate_image := form.File["certificate_image"]

	company_name := form.Value["company_name"][0]
	certificate_category := form.Value["certificate_category"][0]
	certificate_name := form.Value["certificate_name"][0]
	certificate_description := form.Value["certificate_description"][0]
	certificate_website := form.Value["certificate_website"][0]
	certificate_start_date := form.Value["certificate_start_date"][0]
	certificate_end_date := form.Value["certificate_end_date"][0]

	certificate_start_date_time, err := time.Parse("2006-01-02", certificate_start_date)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}

	certificate_end_date_time, err := time.Parse("2006-01-02", certificate_end_date)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}
	logger.Info.Println(certificate_image, company_name, certificate_category, certificate_name, certificate_description, certificate_website, certificate_start_date, certificate_end_date)
	var originalLocation string
	var cdnLocation string
	var resizedCdnLocation string
	for _, file := range certificate_image {
		filename := filepath.Base(file.Filename)
		uploadedFile, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer uploadedFile.Close()

		fileNameParts := strings.Split(filename, ".")
		filename = fileNameParts[0]
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
	timeStorageDb := db.TimeStorageDbManager()

	timeStorageDb.Create(&model.Certificate{
		CertificateName:        certificate_name,
		CompanyName:            company_name,
		CertificateDescription: certificate_description,
		CertificateCategory:    certificate_category,
		CertificateImage:       cdnLocation,
		CertificateThumbnail:   resizedCdnLocation,
		CertificateWebsite:     certificate_website,
		CertificateStartDate:   certificate_start_date_time,
		CertificateEndDate:     certificate_end_date_time,
	})

	c.String(http.StatusOK, "success")
	return
}

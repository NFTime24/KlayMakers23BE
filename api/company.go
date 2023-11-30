package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nftime/db"
	"github.com/nftime/model"
	"net/http"
)

func GetCompanyList(c *gin.Context) {

	timeStorageDb := db.TimeStorageDbManager()
	var companyList []model.Company
	rows, err := timeStorageDb.Select(`*`).Table(`companies`).Order(`id`).Rows()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	for rows.Next() {
		timeStorageDb.ScanRows(rows, &companyList)
	}
	defer rows.Close()
	c.JSON(http.StatusOK, companyList)

}

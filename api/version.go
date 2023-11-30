package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nftime/config"
	"github.com/nftime/logger"
	"net/http"
	"strconv"
	"strings"
)

// @Tags Version
// @Summary CheckVersion
// @Description true: up-to-date, false: outdated
// @Accept json
// @Produce json
// @Param version path string true "current version"
// @Router /versions/{version} [get]
func CheckVersion(c *gin.Context) {
	var updated bool
	var wrongRequest bool
	versionStr := c.Param("version")
	versionServerList := strings.Split(config.Version, ".")
	versionStrList := strings.Split(versionStr, ".")
	fmt.Println(len(versionStrList))
	cnt := strings.Count(versionStr, ".")

	for i, versionComponent := range versionStrList {
		_, err := strconv.Atoi(versionComponent)
		if err != nil {
			logger.Error.Printf("err: %v\n", err)
			wrongRequest = true
		} else {
			if versionComponent < versionServerList[i] {
				updated = false
				break
			} else {
				updated = true
			}
		}
	}
	if cnt != 2 || wrongRequest {
		c.String(http.StatusForbidden, "wrong request")
	} else {
		if updated {
			c.String(http.StatusOK, "true")
		} else {
			c.String(http.StatusOK, "false")
		}
	}
}

package route

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/nftime/api"
	"github.com/nftime/logger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

var idPwMap = map[string]string{
	"admin": "nftime-hashed",
}

func Init() *gin.Engine {

	r := gin.Default()
	r.POST("/back-office/login", api.LoginBackoffice)                    // 1 로그인
	r.POST("/back-office/company/register", api.UploadCompany)           // 2 company 등록
	r.POST("/back-office/certificate/register", api.RegisterCertificate) // 3 certificate 등록
	r.POST("/back-office/certificate/issue", api.IssueCertificate)       // 4 certificate 발행
	r.GET("/company/list", api.GetCompanyList)                           // 5 company list 조회
	r.GET("/certificate/list", api.GetCertificateList)                   // 6 certificate list 조회
	r.GET("/user/certificate/list", api.GetUserCertificateList)          // 7 platform에서 address에 대한 certificate list 조회

	//r.GET("certificate", api.GetCertificateWithId)

	store := cookie.NewStore([]byte("here"))
	r.Use(sessions.Sessions("session-name", store))
	r.Static("/static", "./static")
	r.LoadHTMLGlob("static/*")

	r.POST("/login-stream", StreamLoginHandler)

	r.POST("/login-time", TimeLoginHandler)
	r.GET("/logout-admin", LogoutHandler)
	r.GET("/api/protected", ProtectedHandler)

	r.GET("/nftInfo/:contract_address/:work_id/:nft_id", api.ResponseMetadataJson)

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome To NFTIME!")
	})

	time := r.Group("/time")
	{
		time.GET("/success-webtoon-fair", api.OnSuccessWebtoonFairKlip)              // WebtoonFair
		time.GET("/mint/webtoon-fair/work/:id", api.WebtoonFairMintArtWithoutPaying) // WebtoonFair
		time.POST("/mint", api.MintToAddrTime)                                       // /klip/mint //temporarily deprecated

		time.GET("/success-membership", api.OnSuccessMembershipKlip) // WebtoonFair
		time.GET("/mint/membership", api.MintMembershipNft)          // WebtoonFair
		// Body로 수정필요
	}

	klip := r.Group("/klip")
	{
		klip.GET("/success-swf", api.OnSuccessSwfKlip)                     // SWF
		klip.GET("/success", api.OnSuccessKlip)                            // MBTI
		klip.GET("/success-work", api.OnSuccessWorkKlip)                   // 작품
		klip.POST("/mint", api.MintToAddr)                                 // /klip/mint //temporarily deprecated
		klip.POST("/mint/apple", api.MintToAddrApple)                      // /klip/mint //temporarily deprecated
		klip.GET("/mint/free/work/:id", api.MintArtWithoutPaying)          // MBTI
		klip.GET("/mint/swf/work/:id", api.MintArtWithoutPaying)           // SWF
		klip.GET("/mint-work/free/work/:id", api.MintArtWorkWithoutPaying) // 작품

		// Body로 수정필요
	}

	nft := r.Group("/nfts")
	{
		nft.POST("/works", api.AddNFTWithWorkName) // /nft/work/:id  //temporarily deprecated
	}

	playlist := r.Group("/playlists")
	{
		playlist.GET("/ids/:id", api.GetUserPlayListWithId)
		playlist.POST("/", api.CreatePlaylist)
		playlist.DELETE("/ids/:ids/user/ids/:id", api.DeleteUserPlaylistWithIds)
		playlist.PATCH("/", api.UpdatePlaylistOrder)
	}
	users := r.Group("/users")
	{
		users.POST("/", api.PostUser)                            // /users // TODO: profile path 추가
		users.POST("/logs/:address/:status", api.AddUserLog)     // POST /user/logger //temporarily deprecated
		users.GET("/addresses/:address", api.GetUserWithAddress) // /user/address/:address
		users.POST("/nickname", api.UpdateUserNickname)          // UPDATE /user/nickname/:nickname
		users.PATCH("/nickname", api.PostUserNickname)           // UPDATE /user/nickname/:nickname
		users.POST("/upload-profile", api.UploadProfile)
		users.DELETE("/ids/:id/addresses/:address", api.DeleteUser)
		users.GET("/test/:address", api.TestUser)

	}

	login := r.Group("/login")
	{
		login.POST("/social", api.LoginWithSocial)
	}
	works := r.Group("/works")
	{
		works.GET("/info/:id", api.GetWorkInfoWithID)
		works.GET("/top/:category", api.GetTopWorksWithCategory)
		works.GET("/today", api.GetTodayWorks)
		works.GET("/free", api.GetFreeWorks)
		works.GET("/new", api.GetNewWorks)
		works.GET("/all", api.GetAllWorks)

		works.GET("/artists/names/:name", api.GetArtistWorksWithName)
		works.GET("/select/names/:name", api.GetSelectedWorksWithName)
		works.GET("/artists/ids/:id", api.GetArtistWorksWithId)
		works.GET("/select/ids/:ids", api.GetSelectedWorksWithId)
		works.GET("/stream/select/ids/:ids", api.GetSelectedWorksByteCodeWithId)

		works.GET("/id/nfts/:id", api.GetWorkIdWithNftId)
	}

	artists := r.Group("/artists")
	{
		artists.GET("/active", api.GetActiveArtists)
		artists.GET("/top", api.GetTopArtists)
		artists.GET("/names/:name", api.GetArtistWithName)
		artists.GET("/ids/:id", api.GetArtistWithId)
	}

	fantalks := r.Group("/fantalks")
	{
		fantalks.GET("/artists/:name", api.GetArtistFantalks)
	}

	likes := r.Group("/likes")
	{
		likes.POST("/", api.UpdateLike) //temporarily deprecated
		likes.GET("/list", api.GetLikeList)
		likes.GET("/check/addresses/:address/works/:id", api.CheckLike)
		likes.GET("/works/count", api.GetLikeCount)
	}

	streamBackoffices := r.Group("/stream/back-office")

	{
		streamBackoffices.GET("/", api.StreamLoginPage)
		streamBackoffices.Use(AuthMiddleware)

		streamBackoffices.GET("/works/upload", api.StreamIndexPage)
		streamBackoffices.POST("/works/upload", api.StreamUploadWorks)
		streamBackoffices.GET("/works/list", api.StreamShowWorksList)
		streamBackoffices.POST("/works/list/delete", api.StreamDeleteWorksList)
		streamBackoffices.GET("/artists/upload", api.StreamArtistsPage)
		streamBackoffices.POST("/artists/upload", api.StreamUploadArtists)
		streamBackoffices.GET("/artists/list", api.StreamShowArtistsList)
		streamBackoffices.POST("/artists/list/delete", api.StreamDeleteArtistsList)
	}

	timeBackoffices := r.Group("/time/back-office")
	{
		timeBackoffices.GET("/", api.TimeLoginPage)
		timeBackoffices.Use(AuthMiddleware)

		timeBackoffices.GET("/works/upload", api.TimeIndexPage)
		timeBackoffices.POST("/works/upload", api.TimeUploadWorks)
		timeBackoffices.GET("/works/list", api.TimeShowWorksList)
		timeBackoffices.POST("/works/list/delete", api.TimeDeleteWorksList)
		timeBackoffices.GET("/artists/upload", api.TimeArtistsPage)
		timeBackoffices.POST("/artists/upload", api.TimeUploadArtists)
		timeBackoffices.GET("/artists/list", api.TimeShowArtistsList)
		timeBackoffices.POST("/artists/list/delete", api.TimeDeleteArtistsList)

		timeBackoffices.POST("/certificate/register", api.TimeUploadCertificate)
		timeBackoffices.POST("/stat", api.TimeCountStat)

	}

	versions := r.Group("/versions")
	{
		versions.GET("/:version", api.CheckVersion)
	}
	//TODO: swagger
	setupSwagger(r)
	//TODO: test
	r.POST("/sendCrashReport", api.SendCrashReport)
	r.GET("/test/cache/image", api.CacheTestImage)
	r.GET("/test/cache/video", api.CacheTestVideo)

	r.GET("/location", func(c *gin.Context) {
		c.HTML(http.StatusOK, "location.html", nil)
	})
	r.POST("/location", func(c *gin.Context) {
		// Retrieve the location data from the request
		var location struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		}
		if err := c.ShouldBindJSON(&location); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Process the location data as needed
		fmt.Printf("Latitude: %f\n", location.Latitude)
		fmt.Printf("Longitude: %f\n", location.Longitude)

		c.JSON(http.StatusOK, gin.H{"message": "Location received"})
	})
	return r
}

func setupSwagger(r *gin.Engine) {
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/swagger/stream_index.html")
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
func StreamLoginHandler(c *gin.Context) {
	id := "admin"
	username := c.PostForm("username")
	password := c.PostForm("password")

	hashedPassword, ok := idPwMap[username]
	logger.Info.Println(hashedPassword, ok)
	if !ok || password != hashedPassword || id != username {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login credentials"})
		return
	}

	session := sessions.Default(c)
	session.Set("username", username)
	session.Save()

	c.Redirect(http.StatusSeeOther, "/stream/back-office/works/upload")
	//c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func TimeLoginHandler(c *gin.Context) {
	id := "admin"
	username := c.PostForm("username")
	password := c.PostForm("password")

	hashedPassword, ok := idPwMap[username]
	logger.Info.Println(hashedPassword, ok)
	if !ok || password != hashedPassword || id != username {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login credentials"})
		return
	}

	session := sessions.Default(c)
	session.Set("username", username)
	session.Save()

	c.Redirect(http.StatusSeeOther, "/time/back-office/works/upload")
	//c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("username")
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
func ProtectedHandler(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")

	if username == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hello, " + username.(string) + "! This is a protected resource."})
}
func AuthMiddleware(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")

	if username == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort() // Stops the execution of the middleware chain
		return
	}

	// Continue to the next middleware or the actual handler
	c.Next()
}

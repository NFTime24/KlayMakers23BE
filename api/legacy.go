package api

//kakaoClient := kakao.NewClient(server.KAKAO_CLIENT_ID, server.KAKAO_REDIRECT_URL)
//logger.Info.Printf("cli: %v\n", kakaoClient)
//r.GET("tttt", func(ctx *gin.Context) {
//	data := url.Values{}
//	data.Set("grant_type", "refresh_token")
//	data.Set("client_id", "98408866b1f25eb0a79f70ad926a50e7")
//	data.Set("client_secret", "lGUOoqQNL4OpQwvxDJhNJlQp1AWgjQTu")
//	data.Set("refresh_token", "K55wkUX_PWuc-3Ec6Xe8---XuGh_VYRacu5QpkFQCj102wAAAYluFsbz")
//	resp, err := http.PostForm("https://kauth.kakao.com/oauth/token", data)
//	if err != nil {
//		fmt.Println("Error requesting access token:", err)
//		return
//	}
//	defer resp.Body.Close()
//
//	// Parse the response JSON to extract the new access token
//	var tokenResp struct {
//		AccessToken string `json:"access_token"`
//		TokenType   string `json:"token_type"`
//		ExpiresIn   int    `json:"expires_in"`
//		Scope       string `json:"scope"`
//	}
//
//	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
//	if err != nil {
//		fmt.Println("Error decoding response:", err)
//		return
//	}
//
//	fmt.Println("New Access Token:", tokenResp.AccessToken)
//	fmt.Println("New Access Token:", tokenResp.TokenType)
//
//	fmt.Println("New Access Token:", tokenResp.ExpiresIn)
//	fmt.Println("New Access Token:", tokenResp.Scope)
//})
//
//
//r.GET("/test334", func(ctx *gin.Context) {
//	logger.Info.Printf("here: %v\n", "teste")
//	//logger.Info.Printf("cli: %v\n", kakaoClient)
//	data := url.Values{}
//	data.Set("grant_type", "refresh_token")
//	data.Set("refresh_token", "VO5VYNBQAH5-Wk4v_uhZ5Z-BTJgPZv-XQQWWElyNCj10EQAAAYlt9Nuc")
//	httpClient := config.Client(oauth2.NoContext)
//	resp, err := httpClient.PostForm(config.TokenURL, data)
//	if err != nil {
//		fmt.Println("Error requesting access token:", err)
//		return
//	}
//	defer resp.Body.Close()
//	var tokenResp struct {
//		AccessToken string `json:"access_token"`
//		TokenType   string `json:"token_type"`
//		ExpiresIn   int    `json:"expires_in"`
//		Scope       string `json:"scope"`
//	}
//
//	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
//	if err != nil {
//		fmt.Println("Error decoding response:", err)
//		return
//	}
//
//	fmt.Println("New Access Token:", tokenResp.AccessToken)
//
//})
//r.GET("/user/kakao/auth", func(ctx *gin.Context) {
//	//logger.Info.Printf("cli: %v\n", kakaoClient)
//	kakaoClient.Auth(ctx)
//})
//r.GET("/user/kakao/callback", func(ctx *gin.Context) {
//	kakaoClient.Callback(ctx)
//})
//
//var oauth2Config = oauth2.Config{
//	ClientID:    server.KAKAO_CLIENT_ID,
//	RedirectURL: server.KAKAO_REDIRECT_URL,
//	Scopes:      []string{"account_email"}, // Add other scopes you need
//	Endpoint: oauth2.Endpoint{
//		AuthURL:  "https://kauth.kakao.com/oauth/authorize",
//		TokenURL: "https://kauth.kakao.com/oauth/token",
//	},
//}
//
//r.GET("test", func(ctx *gin.Context) {
//	authURL := oauth2Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
//	ctx.Redirect(http.StatusFound, authURL)
//})

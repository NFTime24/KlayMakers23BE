package model

type FantalkResult struct {
	PostID      uint   `json:"post_id"`
	ArtistName  uint   `json:"artist_name"`
	UserAddress uint   `json:"user_address"`
	PostText    string `json:"post"`
	LikeCount   uint   `json:"like_count"`
	CreateTime  string `json:"created_at"`
	ModifyTime  string `json:"updated_at"`
	Nickname    string `json:"nick_name"`
	ProfilePath string `json:"profile_path"`
}

type Fantalk struct {
	DefaultSetting
	Post        string
	ArtistName  string `gorm:"type:varchar"`
	UserAddress string `gorm:"type:varchar"`
	LikeCount   int
	ProfilePath string `gorm:"type:varchar"`
}

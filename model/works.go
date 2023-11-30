package model

type WorkResult struct {
	Id            int    `json:"id"`
	WorkName      string `json:"work_name"`
	Price         int    `json:"work_price"`
	Description   string `json:"description"`
	WorkCategory  string `json:"category"`
	FilePath      string `json:"file_path"`
	ThumbnailPath string `json:"thumbnail_path"`
	ArtistName    string `json:"artist_name"`
	ProfilePath   string `json:"profile_path"`
	ArtistAddress string `json:"artist_address"`
}

type WorkByteCodeResult struct {
	Id                int    `json:"id"`
	WorkName          string `json:"work_name"`
	Price             int    `json:"work_price"`
	Description       string `json:"description"`
	WorkCategory      string `json:"category"`
	FilePath          string `json:"file_path"`
	ThumbnailPath     string `json:"thumbnail_path"`
	FileByteCode      []byte `json:"file_byte_code"`
	ThumbnailByteCode []byte `json:"thumbnail_byte_code"`
	ArtistName        string `json:"artist_name"`
	ProfilePath       string `json:"profile_path"`
	ArtistAddress     string `json:"artist_address"`
}
type Work struct {
	DefaultSetting
	Name          string `gorm:"not null;primaryKey;type:varchar"`
	ArtistName    string `gorm:"primaryKey;type:varchar"` // dropdown
	Price         uint
	Description   string
	Category      string `gorm:"type:varchar"`
	FilePath      string `gorm:"type:varchar"`
	ThumbnailPath string `gorm:"type:varchar"`
	Display       bool   `gorm:"default:true"`
}

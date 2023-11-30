package model

type Like struct {
	DefaultSetting
	UserAddress string `gorm:"primaryKey;type:varchar"`
	WorkID      uint   `gorm:"primaryKey"`
}

type LikeListResult struct {
	Id int `json:"id"`
}

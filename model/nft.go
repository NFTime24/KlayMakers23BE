package model

type NftWork struct {
	WorkID  string `json:"work-id" binding:"required"`
	Address string `json:"user-address" binding:"required"`
}

type Nft struct {
	DefaultSetting
	WorksID     uint
	UserAddress string `gorm:"type:varchar"`
}

package model

type Ticket struct {
	DefaultSetting
	UserAddress     string `gorm:"primaryKey;type:varchar"`
	TicketCount     uint
	IsMbtiMinted    bool `gorm:"default:false"`
	IsSwfMinted     bool `gorm:"default:false"`
	IsWebtoonMinted bool `gorm:"default:false"`
}

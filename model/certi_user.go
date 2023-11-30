package model

type CertiUser struct {
	UserAddress string `gorm:"type:varchar"`
	Name        string `gorm:"type:varchar"`
	ProfilePath string `gorm:"type:varchar"`
	DefaultSetting
}

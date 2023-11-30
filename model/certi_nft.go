package model

type CertificateUser struct {
	DefaultSetting
	CertificateId     uint
	UserWalletAddress string `gorm:"type:varchar"`
}

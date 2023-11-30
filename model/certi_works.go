package model

import "time"

type Certificate struct {
	CertificateName        string    `json:"certificate_name" gorm:"type:varchar"`
	CompanyName            string    `json:"company_name" gorm:"type:varchar"`
	CertificateDescription string    `json:"certificate_description" gorm:"type:varchar"`
	CertificateCategory    string    `json:"certificate_category" gorm:"type:varchar"`
	CertificateImage       string    `json:"certificate_image" gorm:"type:varchar"`
	CertificateThumbnail   string    `json:"certificate_thumbnail" gorm:"type:varchar"`
	CertificateWebsite     string    `json:"certificate_website" gorm:"type:varchar"`
	CertificateStartDate   time.Time `json:"certificate_start_date" gorm:"type:timestamp without time zone"`
	CertificateEndDate     time.Time `json:"certificate_end_date" gorm:"type:timestamp without time zone"`
	DefaultSetting
}

type CertificateUserList struct {
	UserWalletAddress      string    `json:"user_wallet_address"`
	CertificateName        string    `json:"certificate_name"`
	CompanyName            string    `json:"company_name"`
	CertificateDescription string    `json:"certificate_description"`
	CertificateCategory    string    `json:"certificate_category"`
	CertificateImage       string    `json:"certificate_image"`
	CertificateThumbnail   string    `json:"certificate_thumbnail"`
	CertificateWebsite     string    `json:"certificate_website"`
	CertificateStartDate   time.Time `json:"certificate_start_date"`
	CertificateEndDate     time.Time `json:"certificate_end_date"`
	DefaultSetting
}

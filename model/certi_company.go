package model

import "time"

type Company struct {
	CompanyName        string `gorm:"type:varchar"`
	Address            string `gorm:"type:varchar"`
	CompanyImage       string `gorm:"type:varchar"`
	CompanyDescription string `gorm:"type:varchar"`
	CompanyWebsite     string `gorm:"type:varchar"`
	DefaultSetting
}

type CompanyResult struct {
	Id                 uint      `json:"id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	CompanyName        string    `json:"company_name"`
	CompanyImage       string    `json:"company_image"`
	CompanyDescription string    `json:"company_description"`
	CompanyWebsite     string    `json:"company_website"`
}

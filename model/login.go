package model

type LoginInfo struct {
	Platform string `json:"platform"`
	Target   string `json:"type"`
}
type KakaoUser struct {
	ID          int    `json:"id"`
	ConnectedAt string `json:"connected_at"`
	SynchedAt   string `json:"synched_at"`
	Properties  struct {
		Nickname       string `json:"nickname"`
		ProfileImage   string `json:"profile_image"`
		ThumbnailImage string `json:"thumbnail_image"`
	} `json:"properties"`
	KakaoAccount struct {
		ProfileNicknameNeedsAgreement bool `json:"profile_nickname_needs_agreement"`
		ProfileImageNeedsAgreement    bool `json:"profile_image_needs_agreement"`
		Profile                       struct {
			Nickname          string `json:"nickname"`
			ThumbnailImageURL string `json:"thumbnail_image_url"`
			ProfileImageURL   string `json:"profile_image_url"`
			IsDefaultImage    bool   `json:"is_default_image"`
		} `json:"profile"`
		HasEmail               bool   `json:"has_email"`
		EmailNeedsAgreement    bool   `json:"email_needs_agreement"`
		IsEmailValid           bool   `json:"is_email_valid"`
		IsEmailVerified        bool   `json:"is_email_verified"`
		Email                  string `json:"email"`
		HasAgeRange            bool   `json:"has_age_range"`
		AgeRangeNeedsAgreement bool   `json:"age_range_needs_agreement"`
		AgeRange               string `json:"age_range"`
		HasBirthday            bool   `json:"has_birthday"`
		BirthdayNeedsAgreement bool   `json:"birthday_needs_agreement"`
		Birthday               string `json:"birthday"`
		BirthdayType           string `json:"birthday_type"`
		HasGender              bool   `json:"has_gender"`
		GenderNeedsAgreement   bool   `json:"gender_needs_agreement"`
		Gender                 string `json:"gender"`
	} `json:"kakao_account"`
}

type SocialUser struct {
	DefaultSetting
	NickName     string `gorm:"type:varchar"`
	ProfilePath  string `gorm:"type:varchar"`
	Email        string `gorm:"primaryKey;type:varchar"`
	Gender       string `gorm:"type:varchar"`
	AgeRange     string `gorm:"type:varchar"`
	Birthday     string `gorm:"type:varchar"`
	BirthdayType string `gorm:"type:varchar"`
}

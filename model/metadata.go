package model

type Metadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	GroupName   string `json:"group_name"`
	GroupIcon   string `json:"group_icon"`
	Sendable    bool   `json:"sendable"`
}

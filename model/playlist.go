package model

import (
	"time"
)

type MyIntArray []int64

type Playlist struct {
	DefaultSetting
	Name   string `gorm:"not null;type:varchar;primaryKey"`
	UserId uint   `gorm:"primaryKey"`
	Index  uint
}

type PlaylistWork struct {
	DefaultSetting
	PlaylistId uint
	WorkId     uint
}
type PlaylistResult struct {
	Id        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"playlist_name"`
	Index     uint      `json:"playlist_index"`
	WorkId    []uint    `json:"work_ids"`
}

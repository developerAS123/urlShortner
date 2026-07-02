package models

import (
	"time"
)

type ClickEvent struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	LinkID     uint      `gorm:"index;not null" json:"link_id"`
	ClickedAt  time.Time `gorm:"index;autoCreateTime" json:"clicked_at"`
	IPAddress  string    `json:"-"`
	Country    string    `json:"country"` // to be used in week 2
	City       string    `json:"city"`    // to be used in week 2
	DeviceType string    `json:"device_type"` // mobile/desktop/tablet - week 2
	Browser    string    `json:"browser"` // week 2
	Referrer   string    `json:"referrer"`
	UserAgent  string    `json:"-"`
}

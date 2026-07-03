package models

import (
	"time"
)

type AISummary struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	LinkID      uint      `gorm:"index;not null" json:"link_id"`
	WeekStart   time.Time `gorm:"index;not null" json:"week_start"`
	SummaryText string    `gorm:"type:text;not null" json:"summary_text"`
	GeneratedAt time.Time `gorm:"autoCreateTime" json:"generated_at"`
}

package models

import "time"

type ClusterInfo struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time
	Name      string `json:"name"`
}

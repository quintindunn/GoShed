package models

import "gorm.io/gorm"

type Log struct {
	gorm.Model
	Msg string `gorm:"not null"`
}

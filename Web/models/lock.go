package models

import "gorm.io/gorm"

type RollingCode struct {
	gorm.Model
	Code      string `gorm:"not null"`
	Expiry    int64  `gorm:"not null"`
	Nullified bool   `gorm:"not null;default:false"`
	UUID      string `gorm:"not null;default:''"`
}

type AllocatedCode struct {
	gorm.Model
	Name      string `gorm:"not null"`
	Code      string `gorm:"not null"`
	Expiry    int64  `gorm:"not null"`
	Nullified bool   `gorm:"not null;default:false"`
	UUID      string `gorm:"not null;default:''"`
}

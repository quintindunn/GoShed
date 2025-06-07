package models

import "gorm.io/gorm"

type Config struct {
	gorm.Model
	AdminPin                      string `gorm:"not null;default:8888"`
	NeedAdminPinForUserManagement bool   `gorm:"not null;default:true"`
	UnlockTime                    int64  `gorm:"not null;default:8000"`
}

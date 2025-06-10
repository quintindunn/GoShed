package models

import "gorm.io/gorm"

type Config struct {
	gorm.Model
	AdminPin                      string `gorm:"not null;default:8888"`
	NeedAdminPinForUserManagement bool   `gorm:"not null;default:true"`
	UnlockTime                    int64  `gorm:"not null;default:8000"`
	LockState                     bool   `gorm:"not null;default:false"`
	CodeExpirationCheckInterval   int64  `gorm:"not null;default:5000"`
	RollingCodeLifespanSeconds    int64  `gorm:"not null;default:86400"`
}

/* To add to config:
* First, add the value to the model struct above.
* Second, add the inputs to /Web/templates/configuration.html
* Third, add the values to the payload in /Web/static/javascript/configuration.js
* Fourth, add the values to the ConfigurationRequest, and (if applicable) ConfigurationLog struct in /Web/controllers/configuration.go
* Fifth, add the values to confLog and use utils.SetConfigValue to set the value in /Web/controllers/configuration.go
* Sixth, add the value to the utils.Render function in Configuration in /Web/controllers/configuration.go

* Could this be optimized probably fairly easily to just have to define it in one to two spots, probably. Do I care enough, not at this point in the project.

 */

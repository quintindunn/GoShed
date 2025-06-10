package controllers

import (
	"com.quintindev/WebShed/audit"
	"com.quintindev/WebShed/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Configuration(c *gin.Context) {
	utils.Render(c, http.StatusOK, "configuration", gin.H{
		"needAdminPin":                utils.QueryConfigValue[bool]("need_admin_pin_for_user_management"),
		"unlockTime":                  utils.QueryConfigValue[int64]("unlock_time"),
		"codeExpirationCheckInterval": utils.QueryConfigValue[int64]("code_expiration_check_interval"),
		"rollingCodeLifespanSeconds":  utils.QueryConfigValue[int64]("rolling_code_lifespan_seconds"),
	})
}

type ConfigurationRequest struct {
	AdminPin                      string `json:"adminPin"`
	ChangeAdminPin                bool   `json:"changeAdminPin"`
	NewAdminPin                   string `json:"newAdminPin"`
	NeedAdminPinForUserManagement bool   `json:"needAdminPinForUserManagement"`
	UnlockTime                    int64  `json:"unlockTime"`
	CodeExpirationCheckInterval   int64  `json:"codeExpirationCheckInterval"`
	RollingCodeLifespan           int64  `json:"rollingCodeLifespan"`
}
type ConfigurationLog struct {
	ChangeAdminPin                bool  `json:"changeAdminPin"`
	NeedAdminPinForUserManagement bool  `json:"needAdminPinForUserManagement"`
	UnlockTime                    int64 `json:"unlockTime"`
	CodeExpirationCheckInterval   int64 `json:"codeExpirationCheckInterval"`
	RollingCodeLifespan           int64 `json:"rollingCodeLifespan"`
}

func ConfigurationAPI(c *gin.Context) {
	var json ConfigurationRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "Invalid JSON",
			"errorCode": 1,
		})
		return
	}

	validAdminPin := utils.QueryConfigValue[string]("admin_pin")
	if validAdminPin != json.AdminPin {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":     "Invalid admin code!",
			"errorCode": 2,
		})
		return
	}

	confLog := ConfigurationLog{
		ChangeAdminPin:                json.ChangeAdminPin,
		NeedAdminPinForUserManagement: json.NeedAdminPinForUserManagement,
		UnlockTime:                    json.UnlockTime,
		CodeExpirationCheckInterval:   json.CodeExpirationCheckInterval,
		RollingCodeLifespan:           json.RollingCodeLifespan,
	}
	audit.LogNewConfiguration(fmt.Sprintf("%+v", confLog))

	if json.ChangeAdminPin {
		utils.SetConfigValue[string]("admin_pin", json.NewAdminPin)
	}
	utils.SetConfigValue[bool]("need_admin_pin_for_user_management", json.NeedAdminPinForUserManagement)
	utils.SetConfigValue[int64]("unlock_time", json.UnlockTime)
	utils.SetConfigValue[int64]("code_expiration_check_interval", json.CodeExpirationCheckInterval)
	utils.SetConfigValue[int64]("rolling_code_lifespan_seconds", json.RollingCodeLifespan)

}

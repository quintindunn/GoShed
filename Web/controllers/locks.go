package controllers

import (
	"com.quintindev/WebShed/audit"
	"com.quintindev/WebShed/database"
	"com.quintindev/WebShed/hardwareInterface"
	"com.quintindev/WebShed/models"
	"com.quintindev/WebShed/utils"
	"github.com/google/uuid"

	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Code struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Expiry int64  `json:"expiry"`
	Uuid   string `json:"uuid"`
	Pin    string `json:"adminPin"`
}

func Lock(c *gin.Context) {
	hardwareInterface.GetExpiredCodes()

	var ModelRollingCodes []models.RollingCode
	database.DB.Find(&ModelRollingCodes, "nullified = ?", false)

	var RollingCodes []Code
	for _, code := range ModelRollingCodes {
		RollingCodes = append(RollingCodes, Code{
			Name:   "",
			Code:   code.Code,
			Expiry: code.Expiry,
		})
	}

	var allocatedCodes []models.AllocatedCode
	database.DB.Find(&allocatedCodes, "nullified = ?", false)

	var formattedCodes []gin.H
	for _, code := range allocatedCodes {
		formattedCodes = append(formattedCodes, gin.H{
			"name":   code.Name,
			"code":   code.Code,
			"expiry": time.Unix(code.Expiry, 0).Format("01-02-06 3:04 PM"),
			"UUID":   code.UUID,
		})
	}

	var formattedRollingCodes []gin.H
	for _, code := range RollingCodes {
		formattedRollingCodes = append(formattedRollingCodes, gin.H{
			"code":   code.Code,
			"expiry": time.Unix(code.Expiry, 0).Format("01-02-06 3:04 PM"),
		})
	}

	data := gin.H{
		"codes":                             formattedCodes,
		"rollingCodes":                      formattedRollingCodes,
		"isLocked":                          hardwareInterface.GetGetLocked(),
		"adminPinRequiredForUserManagement": utils.QueryConfigValue[bool]("need_admin_pin_for_user_management"),
	}

	utils.Render(c, 200, "locks", data)
}

type SetLockRequest struct {
	SetLocked bool `json:"setLocked"`
}

func SetLockAPI(c *gin.Context) {
	var json SetLockRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hardwareInterface.PostSetLock(json.SetLocked)

	c.JSON(200, gin.H{
		"isLocked": json.SetLocked,
	})
}

func ResetRollingCodesAPI(c *gin.Context) {
	var codes []Code

	const codeCount = 5
	now := time.Now()

	for i := 0; i < codeCount; i++ {
		num := rand.Intn(1000000)       // 0 to 999999
		str := fmt.Sprintf("%06d", num) // zero-padded string
		codes = append(codes, Code{
			Name:   "",
			Code:   str,
			Expiry: now.Unix() + 86400,
		})
	}

	var ModelRollingCodes []models.RollingCode
	var newCodesAudit []string
	for _, c := range codes {
		ModelRollingCodes = append(ModelRollingCodes, models.RollingCode{
			Code:      c.Code,
			Expiry:    c.Expiry,
			Nullified: false,
		})
		newCodesAudit = append(newCodesAudit, c.Code)
	}

	if err := database.DB.Model(&models.RollingCode{}).Where("nullified = ?", false).Update("nullified", true).Error; err != nil {
		fmt.Println("Failed to nullify codes!")
	}

	if err := database.DB.Create(&ModelRollingCodes).Error; err != nil {
		log.Println("Failed to insert rolling codes:", err)
	}

	c.JSON(200, gin.H{
		"rollingCodes": codes,
	})

	audit.LogRefreshRollingCodes(newCodesAudit)
}

func AddUserCodeAPI(c *gin.Context) {
	var json Code

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validAdminPin := utils.QueryConfigValue[string]("admin_pin")
	needAdminPin := utils.QueryConfigValue[bool]("need_admin_pin_for_user_management")
	if needAdminPin && validAdminPin != json.Pin {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":     "Invalid admin pin!",
			"errorCode": 2,
		})
		return
	}

	var foundCodes []models.AllocatedCode
	database.DB.Model(&models.AllocatedCode{}).
		Where("nullified = ?", false).
		Where("code = ?", json.Code).
		Find(&foundCodes)

	if len(foundCodes) != 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error":     "Code already used!",
			"errorCode": 1,
		})
		return
	}

	json.Expiry /= 1000

	//if (payload.code.length < 4 || payload.code.length > 32) {
	length := len(json.Code)

	if length < 4 || length > 32 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code length must be between 4 and 32"})
	}

	now := time.Now()

	if json.Expiry < now.Unix() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code expiry must be greater than or equal to current time"})
	}

	var allocatedCodes []models.AllocatedCode
	database.DB.Find(&allocatedCodes, "nullified = ?", false)

	var codes []Code

	for _, code := range allocatedCodes {
		codes = append(codes, Code{
			Name:   code.Name,
			Code:   code.Code,
			Expiry: code.Expiry,
			Uuid:   code.UUID,
		})
	}
	id := uuid.New()

	json.Uuid = id.String()
	codes = append(codes, json)
	c.JSON(200, gin.H{
		"authorizedCodes": codes,
	})

	fmt.Printf("locks.go TEMPORARY - Adding User Code: %+v\n", json)
	CodeModel := models.AllocatedCode{
		Name:      json.Name,
		Code:      json.Code,
		Expiry:    json.Expiry,
		Nullified: false,
		UUID:      json.Uuid,
	}

	audit.LogAddNewCode(json.Name, json.Code, json.Expiry)
	if err := database.DB.Create(&CodeModel).Error; err != nil {
		fmt.Println("Error adding code!")
	}
}

type NullifyCodeRequest struct {
	Uuid string `json:"uuid"`
}

func NullifyUserCode(c *gin.Context) {
	var json NullifyCodeRequest

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	targetUuid := json.Uuid

	if targetUuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
	}

	audit.LogRemoveCode(targetUuid)

	if err := database.DB.Model(&models.AllocatedCode{}).
		Where("nullified = ?", false).
		Where("uuid = ?", targetUuid).
		Update("nullified", true).Error; err != nil {
		fmt.Printf("Error nullifying allocated code with UUID %s", targetUuid)
	}

	var allocatedCodes []models.AllocatedCode
	database.DB.Find(&allocatedCodes, "nullified = ?", false)

	var codes []Code

	for _, code := range allocatedCodes {
		codes = append(codes, Code{
			Name:   code.Name,
			Code:   code.Code,
			Expiry: code.Expiry,
			Uuid:   code.UUID,
		})
	}

	c.JSON(200, gin.H{
		"authorizedCodes": codes,
	})
}

type AdminPinSubmission struct {
	Pin string `json:"pin"`
}

func ValidateAdminPin(c *gin.Context) {
	var json AdminPinSubmission

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"valid": false,
			"error": err,
		})
		return
	}

	validAdminPin := utils.QueryConfigValue[string]("admin_pin")

	c.JSON(http.StatusOK, gin.H{
		"valid": json.Pin == validAdminPin,
	})
}

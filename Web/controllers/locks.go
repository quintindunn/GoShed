package controllers

import (
	"com.quintindev/WebShed/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

var codes = []Code{
	{
		Name:   "Temporary1",
		Code:   "987654",
		Expiry: 1750000000,
	},
	{
		Name:   "Temporary2",
		Code:   "456789",
		Expiry: 1760000000,
	},
	{
		Name:   "Temporary3",
		Code:   "654987",
		Expiry: 1770000000,
	},
}

type Code struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Expiry int64  `json:"expiry"`
}

func Lock(c *gin.Context) {

	rollingCodes := []Code{
		{
			Name:   "",
			Code:   "424242",
			Expiry: 1750000000,
		},
		{
			Name:   "",
			Code:   "242424",
			Expiry: 1750000000,
		},
		{
			Name:   "",
			Code:   "131313",
			Expiry: 1750000000,
		},
		{
			Name:   "",
			Code:   "313131",
			Expiry: 1750000000,
		},
	}

	var formattedCodes []gin.H
	for _, code := range codes {
		formattedCodes = append(formattedCodes, gin.H{
			"name":   code.Name,
			"code":   code.Code,
			"expiry": time.Unix(int64(code.Expiry), 0).Format("01-02-06 3:04 PM"),
		})
	}

	var formattedRollingCodes []gin.H
	for _, code := range rollingCodes {
		formattedRollingCodes = append(formattedRollingCodes, gin.H{
			"code":   code.Code,
			"expiry": time.Unix(int64(code.Expiry), 0).Format("01-02-06 3:04 PM"),
		})
	}

	data := gin.H{
		"codes":        formattedCodes,
		"rollingCodes": formattedRollingCodes,
		"isLocked":     false,
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

	fmt.Printf("locks.go TEMPORARY - IsLocked: %v\n", json.SetLocked)

	c.JSON(200, gin.H{
		"isLocked": json.SetLocked,
	})
}

func ResetRollingCodesAPI(c *gin.Context) {
	codes := []Code{}

	const codeCount = 5
	now := time.Now()

	for i := 0; i < codeCount; i++ {
		num := rand.Intn(1000000)       // 0 to 999999
		str := fmt.Sprintf("%06d", num) // zero-padded string
		codes = append(codes, Code{
			Name:   "",
			Code:   str,
			Expiry: now.Unix(),
		})
	}

	c.JSON(200, gin.H{
		"rollingCodes": codes,
	})
}

func AddUserCodeAPI(c *gin.Context) {
	var json Code

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//if (payload.code.length < 4 || payload.code.length > 32) {
	length := len(json.Code)

	if length < 4 || length > 32 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code length must be between 4 and 32"})
	}

	now := time.Now()

	if json.Expiry/1000 < now.Unix() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code expiry must be greater than or equal to current time"})
	}

	codes = append(codes, json)
	c.JSON(200, gin.H{
		"authorizedCodes": codes,
	})

	fmt.Printf("locks.go TEMPORARY - Adding User Code: %+v\n", json)
}

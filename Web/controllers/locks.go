package controllers

import (
	"com.quintindev/WebShed/utils"
	"github.com/gin-gonic/gin"
	"time"
)

type Code struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Expiry int    `json:"expiry"`
}

func Lock(c *gin.Context) {
	codes := []Code{
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
			Name:   "CodeGamma",
			Code:   "654987",
			Expiry: 1770000000,
		},
	}

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

package controllers

import (
	"com.quintindev/WebShed/database"
	"com.quintindev/WebShed/models"
	"com.quintindev/WebShed/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Log struct {
	Date string `json:"date"`
	Msg  string `json:"msg"`
}

func Logs(c *gin.Context) {
	page, found := c.GetQuery("page")

	if !found {
		page = "1"
	}
	pageInt, err := strconv.Atoi(page)

	if err != nil {
		fmt.Println("Error converting page to int:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
		})
		return
	}

	numRes, found := c.GetQuery("maxres")
	if !found {
		numRes = "50"
	}

	numResInt, err := strconv.Atoi(numRes)

	if err != nil {
		fmt.Println("Error converting maxres to int:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
		})
		return
	}
	numResInt--

	fmt.Printf("Page: %d numRes: %d", pageInt, numResInt)

	mx := numResInt * pageInt
	mn := numResInt * (pageInt - 1)

	if mx < mn {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
		})
		return
	}

	if pageInt < 1 || numResInt < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
		})
		return
	}

	var logsModel []models.Log
	result := database.DB.Offset(mn).Limit(mx - mn + 1).Order("created_at desc").Find(&logsModel)
	if result.Error != nil {
		fmt.Println(result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	fmt.Println(len(logsModel))

	var formattedLogs []Log
	for _, log := range logsModel {
		formattedLogs = append(formattedLogs, Log{
			Date: log.CreatedAt.Format("01-02-06 3:04:05 PM"),
			Msg:  log.Msg,
		})
	}

	utils.Render(c, http.StatusOK, "logs", gin.H{
		"logs":        formattedLogs,
		"currentPage": pageInt,
		"nextPage":    pageInt + 1,
		"prevPage":    pageInt - 1,
		"resPerPage":  numResInt + 1,
	})
}

package middleware

import (
	"com.quintindev/WebShed/hardwareInterface"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func VerifyBackendAPI() gin.HandlerFunc {
	return func(c *gin.Context) {

		data, err := hardwareInterface.GetJSONError("/ping")

		if err != nil {
			log.Println("Backend not online!")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":     "Internal Server Error!",
				"errorCode": http.StatusInternalServerError,
			})
			c.Abort()
			return
		}

		obj, ok := data.(map[string]interface{})
		if !ok {
			log.Println("unexpected JSON structure")
		}

		pingReplyVal, ok := obj["msg"]
		if !ok {
			log.Println("key 'msg' not found")
		}

		pingReply, ok := pingReplyVal.(string)
		if !ok {
			log.Println("msg is not a string")
		}

		if pingReply != "pong" {
			log.Println("Backend not online!")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":     "Internal Server Error!",
				"errorCode": http.StatusInternalServerError,
			})
			c.Abort()
		}
	}
}

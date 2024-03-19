package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	calculator := NewCalculator()

	server := gin.Default()
	server.GET("/formula", func(ctx *gin.Context) {
		value := ctx.Query("value")
		answer, err := calculator.ComputeV2(value)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"value": value,
				"err":   fmt.Sprintf("err: %+v", err),
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"answer": answer,
			})
		}
	})
	server.POST("/digits", func(ctx *gin.Context) {
		var postData PostDigits
		if err := ctx.BindJSON(&postData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		calculator.SetDigits(postData.Digits)
		ctx.JSON(http.StatusOK, gin.H{"msg": "OK"})
	})
	server.Run()
}

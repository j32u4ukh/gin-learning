package main

import (
	"fmt"
)

func main() {
	calculator := NewCalculator()
	calculator.SetDigits(5)
	answer, err := calculator.ComputeV1(" 5 * 3 + 1 + 6 / 2")

	if err != nil{
		fmt.Printf("err: %+v\n", err)
		return
	}

	fmt.Printf("answer: %s\n", answer)

	// server := gin.Default()
	// server.GET("/formula", func(ctx *gin.Context) {
	// 	value := ctx.Query("value")
	// 	answer, err := calculator.ComputeV1(value)
	// 	if err != nil {
	// 		ctx.JSON(http.StatusBadRequest, gin.H{
	// 			"value": value,
	// 			"err":   fmt.Sprintf("err: %+v", err),
	// 		})
	// 	} else {
	// 		ctx.JSON(http.StatusOK, gin.H{
	// 			"answer": answer,
	// 		})
	// 	}
	// })
	// server.Run()
}

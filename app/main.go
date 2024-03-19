package main

import "fmt"

func main() {
	calculator := NewCalculator()
	// answer, err := calculator.ComputeV1("5 + (-3)")
	// if err != nil {
	// 	fmt.Printf("err: %+v\n", err)
	// 	return
	// }

	// fmt.Printf("answer: %s\n", answer)
	results, _ := calculator.Parse([]string{"5 + 3 * 7 - 5 / 6"}, "*", "/")
	for _, result := range results{
		fmt.Printf("result: %+v\n", result)
		plusResults, _ := calculator.Parse([]string{result}, "+")

		for _, pr := range plusResults{
			fmt.Printf("\tpr: %+v\n", pr)
		}
	}

	// server := gin.Default()
	// server.GET("/formula", func(ctx *gin.Context) {
	// 	value := ctx.Query("value")
	// 	answer, err := calculator.ComputeV2(value)
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
	// server.POST("/digits", func(ctx *gin.Context) {
	// 	var postData PostDigits
	// 	if err := ctx.BindJSON(&postData); err != nil {
	// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	calculator.SetDigits(postData.Digits)
	// 	ctx.JSON(http.StatusOK, gin.H{"msg": "OK"})
	// })
	// server.Run()
}

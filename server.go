package main

import (
	"GitOperator/logger"
	"GitOperator/routes"
	"net/http"

	"github.com/joho/godotenv"
)

func init() {
	// TO LOAD THE .ENV FILE
	err := godotenv.Load()
	if err != nil {
		logger.FatalLogger.Fatalf("ERROR OCCURED WHILE LOADING .ENV FILE: %v", err)
	}
}

func main() {
	mx := http.NewServeMux()
	routes.InitialiseRoutes(mx)
	logger.DebugLogger.Println("SERVER IS RUNNING ON PORT 3030")
	logger.FatalLogger.Fatalln(http.ListenAndServe(":3030", mx))
}

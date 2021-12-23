package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"net/http"
	"os"
	"path"
)

func main() {
	server := echo.New()
	server.GET(path.Join("/"), Version)

	godotenv.Load()
	port := os.Getenv("PORT")

	addr := fmt.Sprintf("%s:%s", "0.0.0.0", port)
	fmt.Println(addr)
	server.Start(addr)
}

func Version(context echo.Context) error {
	return context.JSON(http.StatusOK, map[string]interface{}{"version": 1})
}

package main

import (
	"fmt"

	"github.com/jeremyKisner/streaming-daemon/internal/database"
)

func main() {
	connector, err := database.CreateDB()
	if err != nil {
		fmt.Println("had error creating connection", err)
		return
	}
	connector.CreateTable()
	defer connector.DB.Close()
}

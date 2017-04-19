package main

import (
	"fmt"
	"os"
	"os/signal"

	_ "github.com/go-sql-driver/mysql" // loads up the go mysql driver.
)

func main() {
	fmt.Println("Welcome to midash")

	cm := make(chan os.Signal, 1)
	signal.Notify(cm, os.Interrupt)
	<-cm
}

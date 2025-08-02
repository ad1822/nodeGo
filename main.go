package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ad1822/nodeGo/internal/runtime"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <script.js>")
	}

	jsrt := runtime.New()
	err := jsrt.RunScript(os.Args[1])
	if err != nil {
		fmt.Println("Script Error:", err)
	}
}

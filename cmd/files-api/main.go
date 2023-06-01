package main

import (
	"fmt"
	"os"

	"github.com/dgurns/files-api/pkg/server"
)

func main() {
	if err := server.Run(); err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"

	"github.com/BON4/employees/internal/server"
)

func main() {
	s, err := server.NewServer()
	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		return
	}

	if err = s.Run(); err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		return
	}
}

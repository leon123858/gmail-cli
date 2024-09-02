package main

import (
	"fmt"
	"github.com/leon123858/gmail-cli/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

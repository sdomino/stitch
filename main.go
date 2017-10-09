package main

import (
	"fmt"
	"os"

	"github.com/sdomino/stitch/cmd"
)

func main() {
	if err := cmd.StitchCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

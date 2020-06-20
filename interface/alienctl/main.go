package main

import (
	"fmt"
	"os"

	aliengamecmd "github.com/ilgooz/aliengame/interface/alienctl/cmd"
)

func main() {
	if err := aliengamecmd.New().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

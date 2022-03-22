package main

import (
	"fmt"
	"mangathrV2/internal/argparse"
	"mangathrV2/internal/config"
	"mangathrV2/internal/utils"
)

func main() {
	fmt.Println("Hello world")

	// Load config object, returns Config struct
	var c config.Config
	if err := c.Load("./examples/config.yml"); err != nil {
		utils.RaiseError(err)
	}
	fmt.Println(c)

	// Load argparse object, returns ArgParse struct
	var a argparse.Argparse
	if err := a.Parse(); err != nil {
		utils.RaiseError(err)
	}
	fmt.Println(a)

	switch a.Command {
	case "download":
		fmt.Println("Downloading", a.Download)
		break
	case "register":
		fmt.Println("Registering", a.Register)
	}

	// Merge Config & ArgParse (ArgParse priority) into ProgramOptions
	// Call (download|register|update|manage|config).go > run(ProgramOptions po) to start program execution
}

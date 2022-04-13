package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"mangathrV2/internal/argparse"
	"mangathrV2/internal/commands/download"
	"mangathrV2/internal/commands/register"
	"mangathrV2/internal/config"
	"mangathrV2/internal/utils"
)

func main() {
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
		download.Run(&a.Download, &c)
		break
	case "register":
		fmt.Println("Registering", a.Register)
		register.Run(&a.Register, &c)
	}

	// Merge Config & ArgParse (ArgParse priority) into ProgramOptions
	// Call (download|register|update|manage|config).go > run(ProgramOptions po) to start program execution
}

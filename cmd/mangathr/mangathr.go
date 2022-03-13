package main

import (
	"fmt"
	"mangathrV2/internal/config"
	"syscall"
)

func main()  {
	fmt.Println("Hello world")

	// Load config object, returns Config struct
	var c config.Config
	err := c.Load("./examples/config.yml")
	if err != nil {
		fmt.Println(err)
		syscall.Exit(1)
	}
	fmt.Println(c)

	// Load argparse object, returns ArgParse struct

	// Merge Config & ArgParse (ArgParse priority) into ProgramOptions
	// Call (download|register|update|manage|config).go > run(ProgramOptions po) to start program execution
}



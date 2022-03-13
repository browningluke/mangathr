package main

import (
	"fmt"
)

func main()  {
	fmt.Println("Hello world")

	// Load config object, returns Config struct
	// Load argparse object, returns ArgParse struct

	// Merge Config & ArgParse (ArgParse priority) into ProgramOptions
	// Call (download|register|update|manage|config).go > run(ProgramOptions po) to start program execution
}



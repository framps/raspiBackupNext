package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/framps/raspiBackupNext/discover"
	"github.com/framps/raspiBackupNext/model"
	"github.com/framps/raspiBackupNext/tools"
)

func main() {

	var debugFlag = flag.Bool("debug", false, "Enable debug messages")
	var collectFlag = flag.Bool("collect", false, "Collect system information")
	var discoverFlag = flag.Bool("discover", false, "Discover system information")
	flag.Parse()

	logger := tools.NewLogger(*debugFlag)
	defer logger.Sync()

	if !*collectFlag && !*discoverFlag {
		*collectFlag = true
	}

	if *collectFlag {
		collectSystem()
		os.Exit(0)
	}

	if *discoverFlag {
		discoverSystem()
		os.Exit(0)
	}

}

func collectSystem() {
	fmt.Printf("%s\n", discover.NewSystem())
}

func discoverSystem() {
	fmt.Printf("%s\n", model.NewSystem())
}

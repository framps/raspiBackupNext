package main

//######################################################################################################################
//
//    Next raspiBackup version written in go
//
//    Copyright (C) 2018 framp at linux-tips-and-tricks dot de
//
//#######################################################################################################################

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
		*discoverFlag = true
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
	system, err := model.NewSystem()
	tools.HandleError(err)
	fmt.Printf("*** From system:\n%s\n", system)
	if err = system.ToJSON("system.model"); err != nil {
		tools.HandleError(err)
	}
	if system, err = model.NewSystemFromJSON("system.model"); system != nil {
		tools.HandleError(err)
	}
	fmt.Printf("*** From json:\n%s\n", system)
}

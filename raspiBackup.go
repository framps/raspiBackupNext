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
	"time"

	"github.com/framps/raspiBackupNext/discover"
	"github.com/framps/raspiBackupNext/model"
	"github.com/framps/raspiBackupNext/tools"
)

func main() {

	var debugFlag = flag.Bool("debug", false, "Enable debug messages")
	var collectFlag = flag.Bool("collect", false, "Collect system information")
	var discoverFlag = flag.Bool("discover", false, "Discover system information")
	var parallelFlag = flag.Bool("parallel", false, "Enable parallel execution")
	flag.Parse()

	tools.NewLogger(*debugFlag)

	if !*collectFlag && !*discoverFlag {
		*discoverFlag = true
	}

	start := time.Now()
	if *collectFlag {
		collectSystem(*parallelFlag)
	}
	if *discoverFlag {
		discoverSystem(*parallelFlag)
	}
	end := time.Now()
	tools.Logger.Debug("Execution time ", end.Sub(start))
	os.Exit(0)
}

func collectSystem(parallelExecution bool) {
	fmt.Printf("=== Collect system ===\n\n%s\n", discover.NewSystem(parallelExecution))
}

func discoverSystem(parallelExecution bool) {
	fmt.Printf("=== Discover system ===\n\n")
	system, err := model.NewSystem(parallelExecution)
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

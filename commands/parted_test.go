package commands

//######################################################################################################################
//
//    Next raspiBackup version written in go
//
//    Copyright (C) 2018 framp at linux-tips-and-tricks dot de
//
//#######################################################################################################################

import (
	"testing"
)

func TestParted(t *testing.T) {
	VerifyData(t, Parted, "parted")
}

package utils

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
)

func PrintBanner() {
	myFigure := figure.NewFigure("crush", "univers", true)
	myFigure.Print()
	fmt.Printf("\n%60s", "version: 3.0\n")
}

func Help() {
	fmt.Println(`Usage: ./crush <command>
command:
	init:	  init database & history info
	monitor:  monitor vul update`)
	PrintLog("info", "Tips: By default, the POCs will not be downloaded locally. If modification is required, the 'downloadPOC' parameter in the config.yaml file can be changed to true.")
}

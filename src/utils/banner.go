package utils

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

func PrintBanner() {
	myFigure := figure.NewFigure("crush", "univers", true)
	myFigure.Print()
	fmt.Printf("\n%60s", "version: 1.0\n")
}

func Help() {
	PrintColor("info", "Usage: ./crush <command>")
	PrintColor("info", `command:
	init:	crawl history vul info
	monitor:  monitor vul update
	TODO:......`)
}
func PrintColor(level, content string, a ...interface{}) {
	var result string
	switch LogLevelMap[level] {
	case INFO:
		result = color.BlueString(content, a...)
	case SUCCESS:
		result = color.GreenString(content, a...)
	case ERROR:
		result = color.RedString(content, a...)
	default:
		result = color.YellowString(content, a...)
	}
	fmt.Println(result)
}

package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

type LogLevel int

const (
	Level = iota
	INFO
	SUCCESS
	ERROR
)

var LogLevelMap = map[string]LogLevel{
	"info":    INFO,
	"success": SUCCESS,
	"error":   ERROR,
}

func WriteToLog(data, filename string) {
	err := os.MkdirAll(filepath.Dir(filename), 0755)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(data + "\n")
	if err != nil {
		fmt.Println("Error writing data:", err)
		return
	}
}

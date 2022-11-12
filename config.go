package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
)

func homeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("APPDATA")
	}
	return os.Getenv("HOME")
}

func projectDir() string {
	return path.Join(homeDir(), SettingDir)
}

func exitF(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
	fmt.Println("")
	os.Exit(1)
}

func exit(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func getConfigMappingPath() string {
	return path.Join(projectDir(), MappingFile)
}

func getConfig() (map[string]string, error) {
	data, err := os.ReadFile(getConfigMappingPath())
	if err != nil {
		return nil, err
	}
	var result map[string]string
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

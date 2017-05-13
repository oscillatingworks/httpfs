package main

import (
	"errors"
	"os"
	"runtime"
)

// UserHomeDir returns the path of your HOME directory
func UserHomeDir() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	return os.Getenv(env)
}

// IsFile tells you whether a file is a file,
// a directory, or does not exist.
// true, nil -> file
// false, nil -> dir
// _, err -> not exist
func IsFile(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		return false, nil
	case mode.IsRegular():
		return true, nil
	}

	return false, errors.New("will not happen")
}

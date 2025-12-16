package main

import (
	"os"
	"path/filepath"
)

func DirExists(path string) (bool, os.FileInfo) {
	isOk, info := PathExists(path)

	return isOk && info.IsDir(), info
}

func PathExists(path string) (bool, os.FileInfo) {
	info, err := os.Stat(path)
	return IsOk(err), info
}

func CreateDir(path string) {
	err := os.Mkdir(path, 0755)
	if !IsOk(err) {
		panic(err)
	}
}

func realPathDir(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return filepath.EvalSymlinks(abs)
}

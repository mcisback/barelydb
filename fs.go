package main

import (
	"os"
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

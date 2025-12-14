package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type JsonKV map[string]any

// Get the path to a database folder
func GetDatabasePath(dbRootDir string, database string) string {
	return filepath.Join(dbRootDir, database)
}

func GetTablePath(dbRootDir string, database string, table string) string {
	return filepath.Join(dbRootDir, database, table+".json")
}

// Check if a database folder exists
func DatabaseExists(dbRootDir string, database string) bool {
	databasePath := GetDatabasePath(dbRootDir, database)

	isOk, info := DirExists(databasePath)

	return isOk && info.IsDir()
}

func TableExists(dbRootDir string, database string, table string) bool {
	tablePath := GetTablePath(dbRootDir, database, table)

	isOk, info := PathExists(tablePath)

	return isOk && info.Mode().IsRegular()
}

// Get the root folder containing all the databases
func getRootDatabaseDirectory(dirName string) string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.Join(wd, dirName)
}

// DatabaseExistsOrCreate creates a database if it doesn't exist and returns its path.
func DatabaseExistsOrCreate(dbRootDir string, database string) string {
	isOk := DatabaseExists(dbRootDir, database)

	databasePath := GetDatabasePath(dbRootDir, database)

	if !isOk {
		CreateDir(databasePath)
	}

	return databasePath
}

func CreateTable(tablePath string) {
	TrySingle(
		os.WriteFile(tablePath, []byte("{}\n"), 0644),
	).OrPanic("Error creating table %s\n", tablePath)
}

func TableExistsOrCreate(dbRootDir string, database string, table string) string {
	isOk := TableExists(dbRootDir, database, table)

	tablePath := GetTablePath(dbRootDir, database, table)

	if !isOk {
		CreateTable(tablePath)
	}

	return tablePath
}

func LoadTable(dbRootDir string, database string, table string) JsonKV {
	tablePath := TableExistsOrCreate(dbRootDir, database, table)

	fileContents := Try(
		os.ReadFile(tablePath),
	).OrPanic("Error reading table %s\n", table)

	var data JsonKV

	TrySingle(
		json.Unmarshal(fileContents, &data),
	).OrPanic("Error unmarshalling table %s\n", table)

	return data
}

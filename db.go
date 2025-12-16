package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	var dbRootDir string

	if strings.HasPrefix(dirName, "/") {
		dbRootDir = dirName
	} else {

		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		dbRootDir = filepath.Join(wd, dirName)
	}

	if ok, _ := DirExists(dbRootDir); !ok {
		CreateDir(dbRootDir)

		fmt.Println("Created rootdatabase directory: ", dbRootDir)
	}

	return dbRootDir
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
	).OrPanic("Error creating table: ", tablePath, "\n")
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
	).OrPanic("Error reading table: ", table, "\n")

	var data JsonKV

	TrySingle(
		json.Unmarshal(fileContents, &data),
	).OrPanic("Error unmarshalling table: ", table, "\n")

	return data
}

func WriteTable(dbRootDir string, database string, table string, tableData JsonKV) {
	tablePath := TableExistsOrCreate(dbRootDir, database, table)

	b := Try(
		json.MarshalIndent(tableData, "", "  "),
	).OrPrint("Error on MarshalIndent in Write Table: ", database, ".", table)

	TrySingle(
		os.WriteFile(tablePath, b, 0644),
	).OrPanic("Error on MarshalIndent in Write Table: ", database, ".", table)
}

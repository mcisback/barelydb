package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

const LISTEN_PORT = ":3838"
const ROOT_DIR_NAME = "barelydb_data"

// TODO: POST,PUT,PATCH,DELETE
// TODO: Implement CRUD operations for the database
// TODO: Implement pagination and filtering
// TODO: Implement error handling and logging?
// TODO: Implement authentication and authorization

func main() {
	app := fiber.New()

	dbRootDir := getRootDatabaseDirectory(ROOT_DIR_NAME)

	if ok, _ := DirExists(dbRootDir); !ok {
		CreateDir(dbRootDir)

		fmt.Println("Created rootdatabase directory: ", dbRootDir)
	}

	fmt.Println("Loading root database directory: ", dbRootDir)

	// Get all records from the database
	app.Get("/:database/:table", func(c *fiber.Ctx) error {
		database := c.Params("database")
		table := c.Params("table")

		if !DatabaseExists(dbRootDir, database) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Database Not Found",
			})
		}

		if !TableExists(dbRootDir, database, table) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Table Not Found",
			})
		}

		tableData := LoadTable(dbRootDir, database, table)

		queryString := c.Queries()

		fmt.Println("GET ", database, "/", table)
		fmt.Println("GET Raw QueryString:", string(c.Request().URI().QueryString()))
		fmt.Println("GET Parsed QueryString:", queryString)
		fmt.Println("GET Table Data: ", tableData)

		return c.JSON(tableData)
	})

	// Get a single record from the database by id
	app.Get("/:database/:table/:id", func(c *fiber.Ctx) error {
		database := c.Params("database")
		table := c.Params("table")
		id := c.Params("id")

		if !DatabaseExists(dbRootDir, database) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Database Not Found",
			})
		}

		if !TableExists(dbRootDir, database, table) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Table Not Found",
			})
		}

		tableData := LoadTable(dbRootDir, database, table)

		queryString := c.Queries()

		fmt.Println("GET ", database, "/", table, "/", id)
		fmt.Println("GET Raw QueryString:", string(c.Request().URI().QueryString()))
		fmt.Println("GET Parsed QueryString:", queryString)
		fmt.Println("GET Table Data: ", tableData)

		raw, ok := tableData[id]
		if !ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Record Not Found",
			})
		}

		data, ok := raw.(map[string]any)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid data format",
			})
		}

		data["id"] = id

		return c.JSON(data)
	})

	log.Fatal(app.Listen(LISTEN_PORT))
}

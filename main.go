package main

import (
	"flag"
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
//
// QS Params:
// page=1&per-page=25
// search-by=column_name&search-term=value
// sort=id,views
// limit=10
// query="email=example@example.com AND|OR name=Giovanni"
// query="email"
// query="email = example@example.com"

func main() {
	app := fiber.New()

	var rootDB string

	flag.StringVar(&rootDB, "root-db", "", "Root database DSN")
	flag.StringVar(&rootDB, "d", "", "Alias for --root-db")

	flag.Parse()

	fmt.Println("root-db =", rootDB)

	if rootDB == "" {
		PrintAndExit("Error: Root database relative or absolute path is required\n")
	}

	dbRootDir := getRootDatabaseDirectory(rootDB)

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

		// data := make(map[string]any)

		// TODO: implement limit
		// limit, ok := queryString["limit"]
		// if ok {
		// 	tableData = tableData[:limit]
		// }

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

	// Create a new record in the table
	app.Post("/:database/:table", func(c *fiber.Ctx) error {
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

		var payload map[string]any

		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON body",
			})
		}

		id, ok := payload["id"]
		if !ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Payload does not contain an 'id' field",
			})
		}

		tableData := LoadTable(dbRootDir, database, table)

		_, ok = tableData[id.(string)]
		if ok {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Record already exists, use PUT or PATCH to update",
			})
		}

		queryString := c.Queries()

		fmt.Println("GET ", database, "/", table, "/", id)
		fmt.Println("GET Raw QueryString:", string(c.Request().URI().QueryString()))
		fmt.Println("GET Parsed QueryString:", queryString)
		fmt.Println("GET Table Data: ", tableData)

		tableData[id.(string)] = payload

		WriteTable(dbRootDir, database, table, tableData)

		return c.JSON(fiber.Map{
			"message":  "Record created successfully",
			"data":     tableData[id.(string)],
			"id":       id,
			"table":    table,
			"database": database,
		})
	})

	// Update a record in a table
	app.Put("/:database/:table/:id", func(c *fiber.Ctx) error {
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

		var payload map[string]any

		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON body",
			})
		}

		tableData := LoadTable(dbRootDir, database, table)

		_, ok := tableData[id]
		if !ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Record doesn't exist",
			})
		}

		queryString := c.Queries()

		fmt.Println("GET ", database, "/", table, "/", id)
		fmt.Println("GET Raw QueryString:", string(c.Request().URI().QueryString()))
		fmt.Println("GET Parsed QueryString:", queryString)
		fmt.Println("GET Table Data: ", tableData)

		tableData[id] = payload

		WriteTable(dbRootDir, database, table, tableData)

		return c.JSON(fiber.Map{
			"message":  "Record updated successfully",
			"data":     tableData[id],
			"id":       id,
			"table":    table,
			"database": database,
		})
	})

	// Patch a record in a table
	app.Patch("/:database/:table/:id", func(c *fiber.Ctx) error {
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

		var payload map[string]any

		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON body",
			})
		}

		tableData := LoadTable(dbRootDir, database, table)

		_, ok := tableData[id]
		if !ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Record doesn't exist",
			})
		}

		queryString := c.Queries()

		fmt.Println("GET ", database, "/", table, "/", id)
		fmt.Println("GET Raw QueryString:", string(c.Request().URI().QueryString()))
		fmt.Println("GET Parsed QueryString:", queryString)
		fmt.Println("GET Table Data: ", tableData)

		existing := tableData[id].(map[string]any)

		for k, v := range payload {
			existing[k] = v
		}

		tableData[id] = existing

		WriteTable(dbRootDir, database, table, tableData)

		return c.JSON(fiber.Map{
			"message":  "Record patched successfully",
			"data":     tableData[id],
			"id":       id,
			"table":    table,
			"database": database,
		})
	})

	log.Fatal(app.Listen(LISTEN_PORT))
}

package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"personal-asset-tracker/models"
)

var (
	//go:embed web/build
	web embed.FS
)

type Application struct {
	logger *log.Logger
	items  *models.ItemModel
}

func main() {
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8080"
	}
	database_file, exists := os.LookupEnv("DATABASE_FILE")
	if !exists {
		database_file = "database.db"
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := sql.Open("sqlite3", database_file)
	if err != nil {
		logger.Fatalf("Failed to open database with error: %v", err)
	}
	defer db.Close()

	itemModel, err := models.NewItemModel(db)
	if err != nil {
		logger.Fatal(err)
	}

	app := Application{
		logger: logger,
		items:  itemModel,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/items/", makeHandler(app.itemHandler))

	dist, _ := fs.Sub(web, "web/build")
	mux.Handle("/", http.FileServer(http.FS(dist)))

	app.logger.Printf("Server listening on %s", port)
	app.logger.Fatal(http.ListenAndServe(":"+port, mux))
}

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}

func (app *Application) itemHandler(w http.ResponseWriter, r *http.Request) {
	itemPath := strings.TrimPrefix(r.URL.Path, "/api/items/")
	app.logger.Printf("%s", itemPath)

	switch r.Method {
	case "GET":
		items, err := app.items.RetrieveItems()
		if err != nil {
			app.logger.Panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]interface{}{
			"items": items,
		})
		if err != nil {
			app.logger.Panic(err)
		}
		return
	case "POST":
		var item models.Item

		err := json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			app.logger.Panic("Unable to parse request body")
		}
	}
}

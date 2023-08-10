package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"os"

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

	http.HandleFunc("/api/item", makeHandler(app.getHandler))

	dist, _ := fs.Sub(web, "web/build")
	http.Handle("/", http.FileServer(http.FS(dist)))

	app.logger.Printf("Server listening on %s", port)
	app.logger.Fatal(http.ListenAndServe(":"+port, nil))
}

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}

func (app *Application) getHandler(w http.ResponseWriter, r *http.Request) {
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
}

func (app *Application) putHandler(w http.ResponseWriter, r *http.Request) {
	var item models.Item

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		app.logger.Panic("Unable to parse request body")
	}
}

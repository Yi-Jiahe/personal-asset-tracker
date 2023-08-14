package main

import (
	"database/sql"
	"embed"
	"encoding/csv"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"fmt"

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

	mux.HandleFunc("/api/items/", app.itemHandler)
	mux.HandleFunc("/api/upload", app.csvHandler)

	dist, _ := fs.Sub(web, "web/build")
	mux.Handle("/", http.FileServer(http.FS(dist)))

	app.logger.Printf("Server listening on %s", port)
	app.logger.Fatal(http.ListenAndServe(":"+port, mux))
}

func (app *Application) itemHandler(w http.ResponseWriter, r *http.Request) {
	itemPath := strings.TrimPrefix(r.URL.Path, "/api/items/")
	app.logger.Printf("%s", itemPath)

	switch r.Method {
	case "GET":
		items, err := app.items.RetrieveItems()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			app.logger.Print(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]interface{}{
			"items": items,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			app.logger.Print(err)
			return
		}
		return
	case "POST":
		var item models.Item

		err := json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			app.logger.Print("Unable to parse request body")
		}

		app.logger.Printf("%+v", item)

		err = app.items.CreateItem(item)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			app.logger.Print(err)
		}

		return
	}
}

func (app *Application) csvHandler(w http.ResponseWriter, r *http.Request) {
	f, _, err := r.FormFile("items.csv")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		app.logger.Print(err)
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.LazyQuotes = true

	// TODO: Add flexibity to csv layout
	headers, err := reader.Read()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		app.logger.Print(err)
		return
	}
	fmt.Printf("%+v", headers)

	items := []models.Item{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			app.logger.Print(err)
			return
		}

		n := len(headers) - len(record)
		if n > 0 {
			record = append(record, make([]string, n)...)
		}

		item := models.Item{
			Item_name: record[1],
		}

		items = append(items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"items": items,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		app.logger.Print(err)
		return
	}
}

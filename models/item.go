package models

import "database/sql"

type Item struct {
	Item_id   int    `json:"item_id"`
	Item_name string `json:"item_name"`
}

type ItemModel struct {
	Db *sql.DB
}

func NewItemModel(db *sql.DB) (*ItemModel, error) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS items (
		item_id INTEGER PRIMARY KEY AUTOINCREMENT,
		item_name string
	)`)
	if err != nil {
		return nil, err
	}

	return &ItemModel{Db: db}, nil
}

func (m ItemModel) CreateItem(item Item) error {
	_, err := m.Db.Exec(`INSERT INTO items (item_name)
						VALUES ($1)`, item.Item_name)
	if err != nil {
		return err
	}
	
	return nil
}

func (m ItemModel) RetrieveItems() ([]Item, error) {
	rows, err := m.Db.Query("SELECT * FROM items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []Item{}
	for rows.Next() {
		var item Item

		err := rows.Scan(&item.Item_id, &item.Item_name)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}



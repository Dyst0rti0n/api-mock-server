package mockdb

import (
	"database/sql"
	"errors"
)

type Resource struct {
	ID   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

func GetAllResources(db *sql.DB) []Resource {
	rows, err := db.Query("SELECT id, name FROM resources")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var resources []Resource
	for rows.Next() {
		var resource Resource
		if err := rows.Scan(&resource.ID, &resource.Name); err != nil {
			panic(err)
		}
		resources = append(resources, resource)
	}
	return resources
}

func AddResource(db *sql.DB, resource Resource) {
	_, err := db.Exec("INSERT INTO resources (id, name) VALUES (?, ?)", resource.ID, resource.Name)
	if err != nil {
		panic(err)
	}
}

func GetResourceByID(db *sql.DB, id string) (Resource, error) {
	row := db.QueryRow("SELECT id, name FROM resources WHERE id = ?", id)
	var resource Resource
	if err := row.Scan(&resource.ID, &resource.Name); err != nil {
		if err == sql.ErrNoRows {
			return Resource{}, errors.New("resource not found")
		}
		panic(err)
	}
	return resource, nil
}

func UpdateResource(db *sql.DB, id string, resource Resource) error {
	_, err := db.Exec("UPDATE resources SET name = ? WHERE id = ?", resource.Name, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("resource not found")
		}
		panic(err)
	}
	return nil
}

func DeleteResource(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM resources WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("resource not found")
		}
		panic(err)
	}
	return nil
}

func ResetDatabase(db *sql.DB) {
	_, err := db.Exec("DROP TABLE IF EXISTS resources")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS resources (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL
	)`)
	if err != nil {
		panic(err)
	}
}

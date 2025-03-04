package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

// Version represents an application version entry
type Version struct {
	ID       int
	AppName  string
	Version  string
	Released time.Time
}

// Database initialization
func initDB() (*sql.DB, error) {
	fmt.Println("Initializing database...")

	db, err := sql.Open("sqlite", "/root/versions.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return nil, err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS versions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		app_name TEXT NOT NULL,
		version TEXT NOT NULL,
		released DATETIME NOT NULL
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Insert a new version
func insertVersion(db *sql.DB, appName, version string) error {
	query := `INSERT INTO versions (app_name, version, released) VALUES (?, ?, ?)`
	_, err := db.Exec(query, appName, version, time.Now())
	return err
}

// Retrieve all versions
func getVersions(db *sql.DB) ([]Version, error) {
	rows, err := db.Query("SELECT id, app_name, version, released FROM versions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []Version
	for rows.Next() {
		var v Version
		var releasedStr string
		err = rows.Scan(&v.ID, &v.AppName, &v.Version, &releasedStr)
		if err != nil {
			return nil, err
		}
		v.Released, _ = time.Parse(time.RFC3339, releasedStr)
		versions = append(versions, v)
	}
	return versions, nil
}

// Delete a version by ID
func deleteVersion(db *sql.DB, id int) error {
	query := `DELETE FROM versions WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

func main() {
	fmt.Println("Starting application...")
	db, err := initDB()
	if err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer db.Close()

	fmt.Println("Adding sample versions...")
	_ = insertVersion(db, "MyApp", "1.0.0")
	_ = insertVersion(db, "MyApp", "1.1.0")
	_ = insertVersion(db, "MyApp", "2.0.0")

	fmt.Println("Fetching versions...")
	versions, err := getVersions(db)
	if err != nil {
		log.Fatal("Error fetching versions:", err)
	}

	for _, v := range versions {
		fmt.Printf("ID: %d, App: %s, Version: %s, Released: %s\n", v.ID, v.AppName, v.Version, v.Released)
	}
}

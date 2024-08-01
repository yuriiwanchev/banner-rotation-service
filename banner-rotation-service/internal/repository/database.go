package repository

import (
	"database/sql"
	"fmt"
	"log"
)

var (
	db     *sql.DB
	signal chan struct{}
)

func InitDB(connStr string) {
	signal = make(chan struct{})
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the database successfully")
	close(signal)
}

func GetDB() *sql.DB {
	<-signal
	return db
}

func CloseDB() {
	db.Close()
}

func CheckTablesExist() (bool, error) {
	query := `
    SELECT EXISTS (
        SELECT 1 
        FROM information_schema.tables 
        WHERE table_name = 'slots'
        OR table_name = 'banners'
        OR table_name = 'user_groups'
		OR table_name = 'slots_banners'
        OR table_name = 'statistics'
    );
    `
	var exists bool
	err := db.QueryRow(query).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func InitSchema() {
	exists, err := CheckTablesExist()
	if err != nil {
		log.Fatalf("Failed to check if tables exist: %v", err)
	}

	if exists {
		log.Println("Database schema already initialized, skipping")
		return
	}

	schema := `
    CREATE TABLE IF NOT EXISTS slots (
        id SERIAL PRIMARY KEY,
        description TEXT NOT NULL
    );

    CREATE TABLE IF NOT EXISTS banners (
        id SERIAL PRIMARY KEY,
        description TEXT NOT NULL
    );

    CREATE TABLE IF NOT EXISTS user_groups (
        id SERIAL PRIMARY KEY,
        description TEXT NOT NULL
    );

	CREATE TABLE IF NOT EXISTS slot_banners (
    slot_id INT NOT NULL REFERENCES slots(id),
    banner_id INT NOT NULL REFERENCES banners(id),
    PRIMARY KEY (slot_id, banner_id)
	);

    CREATE TABLE IF NOT EXISTS statistics (
        id SERIAL PRIMARY KEY,
        slot_id INT NOT NULL REFERENCES slots(id),
        banner_id INT NOT NULL REFERENCES banners(id),
        user_group_id INT NOT NULL REFERENCES user_groups(id),
        clicks INT DEFAULT 0,
        views INT DEFAULT 0
    );
    `

	_, err = db.Exec(schema)
	if err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	log.Println("Database schema initialized successfully")
}

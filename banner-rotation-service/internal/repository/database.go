package repository

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var (
	db *sql.DB
)

func InitDB(connStr string) {
	log.Println("Connecting to the database...")
	var err error
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Println(err)
			log.Println("Failed to connect to database, retrying in 5 seconds...")
			time.Sleep(2 * time.Second)
			continue
		}
		if err = db.Ping(); err == nil {
			break
		}

		log.Println("Database not ready, retrying in 2 seconds...")
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to the database successfully")
}

func GetDB() *sql.DB {
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

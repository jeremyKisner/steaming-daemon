package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var (
	once sync.Once
)

type PostgresConnector struct {
	DB     *sql.DB
	Tables []string
}

func CreateConnection() (PostgresConnector, error) {
	var db PostgresConnector
	once.Do(func() {
		conn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			fmt.Println("error initializing database", err)
			return
		}
		err = conn.Ping()
		if err != nil {
			fmt.Println("ping failed", err)
			return
		}
		db = PostgresConnector{
			DB: conn,
		}
		fmt.Println("database connection started")
	})
	return db, nil
}

func (db *PostgresConnector) Close() {
	if db.DB != nil {
		err := db.DB.Close()
		if err != nil {
			fmt.Println("had issue closing database")
		}
		fmt.Println("database closed")
	}
}

func (db *PostgresConnector) GetTables() []string {
	rows, err := db.DB.QueryContext(context.Background(), "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	fmt.Println("getting tables")
	tables := make([]string, 0)
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(tableName)
		tables = append(tables, tableName)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	db.Tables = tables
	return db.Tables
}

func (db *PostgresConnector) CreateAudioTable() {
	rows, err := db.DB.QueryContext(context.Background(), `CREATE TABLE IF NOT EXISTS audio (
		internal_id SERIAL PRIMARY KEY,
		name VARCHAR(2500) NOT NULL,
		artist VARCHAR(2500) NOT NULL,
		album VARCHAR(2500) NOT NULL,
		pickup_url VARCHAR(2500), 
		plays INT)`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	fmt.Println("created table")
}

type InsertRequest struct {
	Name      string `json:"name"`
	Artist    string `json:"artist"`
	Album     string `json:"album"`
	PickupURL string `json:"pickup_url"`
}

func (db *PostgresConnector) InsertAudio(req InsertRequest) bool {
	result, err := db.DB.Exec(`
	INSERT INTO audio (name, artist, album, pickup_url, plays)
	VALUES ($1, $2, $3, $4, $5)`,
		req.Name, req.Artist, req.Album, req.PickupURL, 0)
	if err != nil {
		// Handle error
		fmt.Println("Error executing INSERT statement:", err)
		return false
	}

	// Get the number of rows affected by the INSERT operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Handle error
		fmt.Println("Error getting rows affected:", err)
		return false
	}

	fmt.Printf("%d rows inserted\n", rowsAffected)
	return true
}

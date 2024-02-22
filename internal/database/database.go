package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jeremyKisner/streaming-daemon/internal/record"
	_ "github.com/lib/pq"
)

var (
	once sync.Once
)

type PostgresConnector struct {
	DB     *sql.DB
	Tables []string
}

// CreateConnection establishes a singleton connection pool.
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

// Close calls to close the connection pool.
func (db *PostgresConnector) Close() {
	if db.DB != nil {
		err := db.DB.Close()
		if err != nil {
			fmt.Println("had issue closing database")
		}
		fmt.Println("database closed")
	}
}

// GetTables loads the current tables available.
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

// CreateAudioTable is used on application startup to create the database for audio records.
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

// InsertNewAudioRecord inserts a new audio record into the database.
func (db *PostgresConnector) InsertNewAudioRecord(req record.AudioRecord) bool {
	// new records should have zero plays
	req.Plays = 0
	return db.InsertAudioRecord(req)
}

// InsertAudioRecord inserts an audio record into the database.
func (db *PostgresConnector) InsertAudioRecord(req record.AudioRecord) bool {
	result, err := db.DB.Exec(`
	INSERT INTO audio (name, artist, album, pickup_url, plays)
	VALUES ($1, $2, $3, $4, $5)`,
		req.Name, req.Artist, req.Album, req.PickupURL, req.Plays)
	if err != nil {
		fmt.Println("error executing INSERT statement:", err)
		return false
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("error getting rows affected", err)
		return false
	}
	fmt.Printf("%d rows inserted\n", rowsAffected)
	return true
}

func (db *PostgresConnector) ExtractAudioByInternalID(internalID int) *record.AudioRecord {
	var name, artist, album, pickupURL string
	var plays int
	err := db.DB.QueryRow("SELECT name, artist, album, pickup_url, plays FROM audio WHERE internal_id = $1", internalID).Scan(&name, &artist, &album, &pickupURL, &plays)
	if err != nil {
		fmt.Println("error getting rows affected:", err)
		return &record.AudioRecord{}
	}
	return &record.AudioRecord{
		Name:      name,
		Artist:    artist,
		Album:     album,
		PickupURL: pickupURL,
		Plays:     plays,
	}
}

package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/cenkalti/backoff/v4"
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
		ping := func() error {
			fmt.Println("pinging database...")
			err := conn.Ping()
			if err != nil {
				fmt.Println("pinging database failed. retrying...")
			}
			return err
		}
		err = backoff.Retry(ping, backoff.NewExponentialBackOff())
		if err != nil {
			fmt.Println("failed to ping, check database connection", err)
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
		description VARCHAR(5000),
		plays INT)`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	fmt.Println("created table")
}

// InsertNewAudioRecord inserts a new audio record into the database.
func (db *PostgresConnector) InsertNewAudioRecord(req record.Audio) bool {
	// new records should have zero plays
	req.Plays = 0
	return db.InsertAudioRecord(req)
}

// InsertAudioRecord inserts an audio record into the database.
func (db *PostgresConnector) InsertAudioRecord(req record.Audio) bool {
	result, err := db.DB.Exec(`
	INSERT INTO audio (name, artist, album, pickup_url, description, plays)
	VALUES ($1, $2, $3, $4, $5, $6)`,
		req.Name, req.Artist, req.Album, req.PickupURL, req.Description, req.Plays)
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

// ExtractAudioByInternalID returns database row using the internal_id.
func (db *PostgresConnector) ExtractAudioByInternalID(internalID int) ([]byte, error) {
	var name, artist, album, pickupURL, description string
	var plays int
	err := db.DB.QueryRow("SELECT name, artist, album, pickup_url, description, plays FROM audio WHERE internal_id = $1", internalID).Scan(&name, &artist, &album, &pickupURL, &description, &plays)
	if err != nil {
		fmt.Println("error getting rows affected:", err)
		return []byte{}, err
	}
	a := &record.Audio{
		ID:          internalID,
		Name:        name,
		Artist:      artist,
		Album:       album,
		PickupURL:   pickupURL,
		Description: description,
		Plays:       plays,
	}
	bts, err := json.Marshal(a)
	if err != nil {
		fmt.Println("had issue marshaling json", err)
		return []byte{}, err
	}
	return bts, nil
}

// IncrementPlays locks to increment a play.
func (db *PostgresConnector) IncrementPlays(a record.Audio) {
	var lock sync.Mutex
	lock.Lock()
	defer lock.Unlock()
	a.Plays += 1
	result, err := db.DB.Exec("UPDATE audio SET plays=$1 WHERE internal_id = $2", a.Plays, a.ID)
	if err != nil {
		fmt.Println("had issue marshaling json", err)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("error getting rows affected", err)
		return
	}
	fmt.Printf("%d rows updated\n", rowsAffected)
}

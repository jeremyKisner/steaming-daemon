package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type PostgresConnector struct {
	DB *sql.DB
}

// CreateDB creates a connection to the database.
// Note: this does not call Close() on the connection.
// Calling functions must independently manage the Close().
func CreateDB() (PostgresConnector, error) {
	conn, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return PostgresConnector{}, err
	}
	err = conn.Ping()
	if err != nil {
		return PostgresConnector{}, err
	}
	db := PostgresConnector{
		DB: conn,
	}
	return db, nil
}

func (db *PostgresConnector) ListTables() {
	rows, err := db.DB.QueryContext(context.Background(), "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	fmt.Println("Tables:")
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(tableName)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func (db *PostgresConnector) CreateTable() {
	rows, err := db.DB.QueryContext(context.Background(), `CREATE TABLE IF NOT EXISTS audio (
		internal_id SERIAL PRIMARY KEY,
		name VARCHAR(2500) NOT NULL,
		artist VARCHAR(2500) NOT NULL,
		album VARCHAR(2500) NOT NULL,
		pickup_url VARCHAR(2500), 
		plays INT`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
}

package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgresDB(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	
	err = db.Ping()
	if err != nil {
		log.Printf("Warning: failed to ping postgres: %v", err)
	} else {
		_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS telemetry (
				id SERIAL PRIMARY KEY,
				agent_id VARCHAR(50),
				timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				cpu_utilization FLOAT,
				memory_utilization FLOAT,
				disk_io FLOAT,
				network_ingress FLOAT,
				network_egress FLOAT,
				temperature FLOAT,
				uptime INT,
				status VARCHAR(50)
			)
		`)
		if err != nil { log.Printf("Error creating telemetry table: %v", err) }
		
		_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS processes (
				id SERIAL PRIMARY KEY,
				agent_id VARCHAR(50),
				timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				pid INT,
				executable_name VARCHAR(100),
				resource_utilization FLOAT
			)
		`)
		if err != nil { log.Printf("Error creating processes table: %v", err) }
	}

	return db, nil
}

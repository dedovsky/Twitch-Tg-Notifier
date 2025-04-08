package Postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
)

var pgxStruct = New()

type Postgres struct {
	ctx  context.Context
	conn *pgxpool.Pool
}

func New() *Postgres {
	url := os.Getenv("POSTGRES_URL")
	if url == "" {
		log.Fatal("POSTGRES_URL environment variable is not set")
	}

	conn, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	p := &Postgres{
		ctx:  context.Background(),
		conn: conn,
	}

	err = p.test()
	if err != nil {
		log.Fatalf("Database test failed: %v", err)
	}

	return p
}

func GetDB() *Postgres {
	return pgxStruct
}

func (p *Postgres) test() error {
	_, err := p.conn.Exec(p.ctx, "SELECT 1")
	if err != nil {
		return err
	}
	// Optionally check for the inventory table
	var exists bool
	err = p.conn.QueryRow(p.ctx, "SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'channels')").Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		log.Println("Warning: 'channels' table does not exist")
	}
	return nil
}

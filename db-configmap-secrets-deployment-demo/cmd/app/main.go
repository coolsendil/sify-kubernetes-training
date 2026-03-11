package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	AppPort    string
	SSLMode    string
}

func mustGetenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func loadConfig() config {
	return config{
		DBHost:     mustGetenv("DB_HOST", "postgres"),
		DBPort:     mustGetenv("DB_PORT", "5432"),
		DBUser:     mustGetenv("DB_USER", "appuser"),
		DBPassword: mustGetenv("DB_PASSWORD", "apppassword"),
		DBName:     mustGetenv("DB_NAME", "appdb"),
		AppPort:    mustGetenv("APP_PORT", "8080"),
		SSLMode:    mustGetenv("DB_SSLMODE", "disable"),
	}
}

func openDB(cfg config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func initDB(db *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS visitors (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
`
	if _, err := db.Exec(schema); err != nil {
		return err
	}

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM visitors`).Scan(&count); err != nil {
		return err
	}

	if count == 0 {
		if _, err := db.Exec(`INSERT INTO visitors (name) VALUES ($1), ($2), ($3)`, "JP", "Alice", "Bob"); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	cfg := loadConfig()

	var db *sql.DB
	var err error

	for i := 1; i <= 20; i++ {
		db, err = openDB(cfg)
		if err == nil {
			break
		}
		log.Printf("database not ready yet (attempt %d/20): %v", i, err)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := initDB(db); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			http.Error(w, fmt.Sprintf("db ping failed: %v", err), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("pong"))
	})
	mux.HandleFunc("/greet", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{"message": "Hello from Sify...",
			"status": "success"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`SELECT id, name, created_at FROM visitors ORDER BY id`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = fmt.Fprintf(w, "Go + PostgreSQL on Kubernetes\n\n")
		_, _ = fmt.Fprintf(w, "Connected to DB: %s on host %s:%s\n\n", cfg.DBName, cfg.DBHost, cfg.DBPort)
		_, _ = fmt.Fprintf(w, "Visitors:\n")
		for rows.Next() {
			var (
				id        int
				name      string
				createdAt time.Time
			)
			if err := rows.Scan(&id, &name, &createdAt); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, _ = fmt.Fprintf(w, "- id=%d name=%s created_at=%s\n", id, name, createdAt.Format(time.RFC3339))
		}
	})

	addr := ":" + cfg.AppPort
	log.Printf("server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

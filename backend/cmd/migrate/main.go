package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func main() {
	fmt.Println("==================================================")
	fmt.Println("🚀 CAFE ERP SYSTEM: DATABASE MIGRATION RUNNER")
	fmt.Println("==================================================")

	// Load env
	_ = godotenv.Load()
	_ = godotenv.Load("../.env_database")
	_ = godotenv.Load("../../.env_database")
	_ = godotenv.Load(".env_database")

	dbUser := getEnv("DB_USERNAME", "postgres")
	dbPass := getEnv("DB_PASSWORD", "samp3321")
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_DATABASE", "cafe_erp")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 1. Connect to maintenance db 'postgres' to check/create database
	adminConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=disable", dbUser, dbPass, dbHost, dbPort)
	fmt.Printf("🔌 Connecting to PostgreSQL server at %s:%s...\n", dbHost, dbPort)

	adminConn, err := pgx.Connect(ctx, adminConnStr)
	if err != nil {
		fmt.Printf("❌ Failed to connect to postgres server: %v\n", err)
		os.Exit(1)
	}

	var exists bool
	checkDbQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", dbName)
	err = adminConn.QueryRow(ctx, checkDbQuery).Scan(&exists)
	if err != nil {
		fmt.Printf("❌ Failed to check database existence: %v\n", err)
		adminConn.Close(ctx)
		os.Exit(1)
	}

	if !exists {
		fmt.Printf("📦 Database '%s' does not exist. Creating database '%s'...\n", dbName, dbName)
		_, err = adminConn.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			fmt.Printf("❌ Failed to create database: %v\n", err)
			adminConn.Close(ctx)
			os.Exit(1)
		}
		fmt.Println("✅ Database created successfully!")
	} else {
		fmt.Printf("ℹ️ Database '%s' already exists.\n", dbName)
	}
	adminConn.Close(ctx)

	// 2. Connect to target database 'cafe_erp'
	targetConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)
	conn, err := pgx.Connect(ctx, targetConnStr)
	if err != nil {
		fmt.Printf("❌ Failed to connect to database '%s': %v\n", dbName, err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	fmt.Printf("✅ Connected to '%s' successfully!\n", dbName)

	// Find migration scripts
	migrationFiles := []string{
		"migrations/000001_full_cafe_erp_schema.sql",
		"migrations/000002_complete_seed_data.sql",
		"migrations/000003_operational_seed_data.sql",
	}

	for _, relPath := range migrationFiles {
		fullPath := relPath
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			fullPath = filepath.Join("..", relPath)
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				fullPath = filepath.Join("backend", relPath)
			}
		}

		fmt.Printf("📄 Executing migration: %s ...\n", fullPath)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			fmt.Printf("❌ Failed to read migration file %s: %v\n", fullPath, err)
			os.Exit(1)
		}

		_, err = conn.Exec(ctx, string(content))
		if err != nil {
			fmt.Printf("❌ Failed to execute migration %s: %v\n", fullPath, err)
			os.Exit(1)
		}
		fmt.Printf("✅ Migration %s applied successfully!\n", filepath.Base(fullPath))
	}

	// Ensure all initial users have valid bcrypt hash for Admin@123
	tag, err := conn.Exec(ctx, `UPDATE users SET password_hash = '$2a$12$wsfbgmY5LXr9BzAcXdu2PeNwsIB5mONs.CQho/8ofs5wIfw9C0HyK' WHERE username IN ('admin', 'sarah.manager', 'jane.cashier', 'chef.john', 'ahmad.warehouse') OR email = 'admin@cafe-erp.com'`)
	if err != nil {
		fmt.Printf("⚠️ Warning updating user password hashes: %v\n", err)
	} else {
		fmt.Printf("🔑 Rows updated: %d\n", tag.RowsAffected())
	}

	rowsUser, err := conn.Query(ctx, "SELECT id, username, email, password_hash, is_active FROM users")
	if err == nil {
		defer rowsUser.Close()
		for rowsUser.Next() {
			var id, u, e, ph string
			var ia bool
			if err := rowsUser.Scan(&id, &u, &e, &ph, &ia); err == nil {
				fmt.Printf("👤 DB USER: username='%s' email='%s' active=%v hash='%s'\n", u, e, ia, ph)
			}
		}
	}

	// 3. Verify tables and foreign key relations
	rows, err := conn.Query(ctx, `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' AND table_type = 'BASE TABLE'
		ORDER BY table_name;
	`)
	if err != nil {
		fmt.Printf("⚠️ Could not list tables: %v\n", err)
	} else {
		defer rows.Close()
		var tables []string
		for rows.Next() {
			var t string
			if err := rows.Scan(&t); err == nil {
				tables = append(tables, t)
			}
		}
		fmt.Printf("\n📊 Total Tables Created: %d\n", len(tables))
		fmt.Printf("📋 Tables: %s\n", strings.Join(tables, ", "))
	}

	// Foreign Keys count
	var fkCount int
	err = conn.QueryRow(ctx, `
		SELECT count(*) 
		FROM information_schema.table_constraints 
		WHERE constraint_type = 'FOREIGN KEY' AND table_schema = 'public';
	`).Scan(&fkCount)
	if err == nil {
		fmt.Printf("🔗 Total Foreign Key Constraints Active: %d\n", fkCount)
	}

	fmt.Println("\n🎉 Database migration & relational integrity setup COMPLETED SUCCESSFULLY!")
	fmt.Println("==================================================")
}

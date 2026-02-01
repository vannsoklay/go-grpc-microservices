#!/bin/bash

# Load environment variables
export DATABASE_URL="postgres://admin:admin123@localhost:5432/mydb?sslmode=disable"

# Run migration
go run cmd/migrate/main.go -path db/migrations -cmd "$1"
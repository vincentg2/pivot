package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/google/uuid"
	"github.com/vincentg2/pivot/backend/internal/config"
	"github.com/vincentg2/pivot/backend/internal/database"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: create-admin EMAIL NICKNAME")
		os.Exit(2)
	}
	fmt.Print("Password (minimum 12 characters): ")
	password, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil || len(password) < 12 {
		fmt.Fprintln(os.Stderr, "a password of at least 12 characters is required")
		os.Exit(2)
	}
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	seedBytes := make([]byte, 16)
	if _, err := rand.Read(seedBytes); err != nil {
		panic(err)
	}
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	pool, err := database.Open(context.Background(), cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	_, err = pool.Exec(context.Background(), `INSERT INTO users(id,email,password_hash,nickname,avatar_seed,role) VALUES($1,$2,$3,$4,$5,'admin')`, uuid.New(), strings.ToLower(strings.TrimSpace(os.Args[1])), string(hash), strings.TrimSpace(os.Args[2]), base64.RawURLEncoding.EncodeToString(seedBytes))
	if err != nil {
		panic(err)
	}
	fmt.Println("Administrator created.")
}

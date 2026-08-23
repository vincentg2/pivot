package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/vincentg2/pivot/backend/internal/config"
	"github.com/vincentg2/pivot/backend/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	pool, err := database.Open(context.Background(), cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	bytes := make([]byte, 10)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	code := "PIVOT-" + base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bytes)
	hash := sha256.Sum256([]byte(strings.ToUpper(code)))
	_, err = pool.Exec(context.Background(), `INSERT INTO invitations(id,code_hash,label,max_uses) VALUES($1,$2,'Local demo',10)`, uuid.New(), hash[:])
	if err != nil {
		panic(err)
	}
	fmt.Printf("Demo invitation (shown once): %s\n", code)
}

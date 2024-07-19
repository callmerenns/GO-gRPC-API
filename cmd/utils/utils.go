package utils

import (
	"fmt"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type Blacklist struct {
	mu     sync.Mutex
	tokens map[string]struct{}
}

func NewBlacklist() *Blacklist {
	return &Blacklist{
		tokens: make(map[string]struct{}), // Inisialisasi peta
	}
}

func (b *Blacklist) Add(token string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	// Simulasi pemeriksaan kesalahan (misalnya, jika ada kesalahan dalam penyimpanan)
	if len(token) == 0 {
		return fmt.Errorf("token cannot be empty")
	}

	// Menambahkan token ke dalam map
	b.tokens[token] = struct{}{}
	return nil
}

func (b *Blacklist) Exists(token string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	_, exists := b.tokens[token]
	return exists
}

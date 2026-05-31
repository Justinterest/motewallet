package verification

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"
)

const (
	CodeLength       = 6
	CodeTTL          = 10 * time.Minute
	ResendCooldown   = 60 * time.Second
)

type entry struct {
	code      string
	expiresAt time.Time
	lastSent  time.Time
}

type Store struct {
	mu    sync.Mutex
	items map[string]*entry
}

func NewStore() *Store {
	return &Store{
		items: make(map[string]*entry),
	}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func generateCode() (string, error) {
	var code string
	for i := 0; i < CodeLength; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code += fmt.Sprintf("%d", n.Int64())
	}
	return code, nil
}

func (s *Store) Issue(email string) (code string, retryAfter time.Duration, err error) {
	key := normalizeEmail(email)
	now := time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	if existing, ok := s.items[key]; ok {
		if since := now.Sub(existing.lastSent); since < ResendCooldown {
			return "", ResendCooldown - since, ErrSendTooFrequent
		}
	}

	code, err = generateCode()
	if err != nil {
		return "", 0, err
	}

	s.items[key] = &entry{
		code:      code,
		expiresAt: now.Add(CodeTTL),
		lastSent:  now,
	}

	return code, 0, nil
}

func (s *Store) Verify(email, code string) error {
	key := normalizeEmail(email)
	now := time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[key]
	if !ok {
		return ErrInvalidCode
	}

	if now.After(item.expiresAt) {
		delete(s.items, key)
		return ErrExpiredCode
	}

	if item.code != strings.TrimSpace(code) {
		return ErrInvalidCode
	}

	delete(s.items, key)
	return nil
}

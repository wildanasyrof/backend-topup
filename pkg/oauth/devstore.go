package oauth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"sync"
	"time"
)

type DevStore interface {
	NewStateAndPKCE() (state, verifier, challenge string, err error)
	Consume(state string) (verifier string, ok bool)
}

type devStore struct {
	mu  sync.Mutex
	exp time.Duration
	// keyed by state; stores code_verifier and expiry
	data map[string]struct {
		verifier string
		expAt    time.Time
	}
}

func NewDevStore(ttl time.Duration) DevStore {
	s := &devStore{exp: ttl, data: make(map[string]struct {
		verifier string
		expAt    time.Time
	})}
	// simple janitor
	go func() {
		t := time.NewTicker(time.Minute)
		defer t.Stop()
		for range t.C {
			now := time.Now()
			s.mu.Lock()
			for k, v := range s.data {
				if now.After(v.expAt) {
					delete(s.data, k)
				}
			}
			s.mu.Unlock()
		}
	}()
	return s
}

func randURLSafe(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (s *devStore) NewStateAndPKCE() (state, codeVerifier, codeChallenge string, err error) {
	state, err = randURLSafe(32)
	if err != nil {
		return
	}
	codeVerifier, err = randURLSafe(64) // RFC 7636: 43–128 chars
	if err != nil {
		return
	}
	h := sha256.Sum256([]byte(codeVerifier))
	codeChallenge = base64.RawURLEncoding.EncodeToString(h[:])
	s.mu.Lock()
	s.data[state] = struct {
		verifier string
		expAt    time.Time
	}{verifier: codeVerifier, expAt: time.Now().Add(s.exp)}
	s.mu.Unlock()
	return
}

func (s *devStore) Consume(state string) (codeVerifier string, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, exists := s.data[state]
	if !exists || time.Now().After(item.expAt) {
		delete(s.data, state)
		return "", false
	}
	delete(s.data, state) // one-time use
	return item.verifier, true
}

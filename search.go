package main

//go:generate msgp

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/inconshreveable/log15"
	"github.com/tinylib/msgp/msgp"
)

// Search structure
type Search struct {
	mu sync.RWMutex

	Items map[string]map[string]string
}

// Result structure
//msgp:ignore Result
type Result struct {
	Found []string `json:"found"`
	Count int      `json:"count"`
}

// NewSearch function
func NewSearch() *Search {
	return &Search{
		Items: make(map[string]map[string]string),
	}
}

// Sync function
func (s *Search) Sync(db *os.File, child bool) error {
	log.Debug("Search DB syncing...")

	if !child {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}

	t1 := time.Now()

	err := msgp.WriteFile(s, db)
	if err != nil {
		return err
	}

	log.Debug("Search DB write to disk", "duration", time.Since(t1).Round(time.Millisecond),
		"size", fmt.Sprintf("%d mb", s.Msgsize()/1024/1024))

	return db.Sync()
}

// Set function
func (s *Search) Set(key, id, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Items[key] == nil {
		s.Items[key] = make(map[string]string)
	}

	s.Items[key][id] = value
}

// Del function
func (s *Search) Del(key, id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Items[key] == nil {
		return
	}

	delete(s.Items[key], id)

	if len(s.Items[key]) == 0 {
		delete(s.Items, key)
	}
}

// Search function
func (s *Search) Search(key, query string, start, stop int) *Result {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := &Result{
		Found: []string{},
		Count: 0,
	}

	if s.Items[key] == nil {
		return result
	}

	for id, value := range s.Items[key] {
		if strings.Contains(value, query) {
			result.Found = append(result.Found, id)
		}
	}

	if start > stop {
		start = stop
	}

	if start < 0 {
		start = 0
	}

	if stop > len(result.Found) {
		stop = len(result.Found)
	}

	result.Count = len(result.Found)

	if start != 0 && stop != 0 {
		result.Found = result.Found[start:stop]
	}

	return result
}

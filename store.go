package main

import (
	"sync"
)

type Store struct {
	mu   sync.RWMutex
	kv   map[string]string
	hash map[string]map[string]string
}

func NewStore() *Store {
	return &Store{
		kv:   make(map[string]string),
		hash: make(map[string]map[string]string),
	}
}

func (s *Store) ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}
	return Value{typ: "string", str: args[0].bulk}
}

func (s *Store) set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'SET' command"}
	}

	key, val := args[0].bulk, args[1].bulk

	s.mu.Lock()
	s.kv[key] = val
	s.mu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func (s *Store) get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'GET' command"}
	}

	key := args[0].bulk

	s.mu.RLock()
	val, ok := s.kv[key]
	s.mu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}
	return Value{typ: "bulk", bulk: val}
}

func (s *Store) hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'HSET' command"}
	}

	hashKey, field, val := args[0].bulk, args[1].bulk, args[2].bulk

	s.mu.Lock()
	if _, exists := s.hash[hashKey]; !exists {
		s.hash[hashKey] = make(map[string]string)
	}
	s.hash[hashKey][field] = val
	s.mu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func (s *Store) hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'HGET' command"}
	}

	hashKey, field := args[0].bulk, args[1].bulk

	s.mu.RLock()
	val, ok := s.hash[hashKey][field]
	s.mu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}
	return Value{typ: "bulk", bulk: val}
}

func (s *Store) hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'HGETALL' command"}
	}

	hashKey := args[0].bulk

	s.mu.RLock()
	entries, ok := s.hash[hashKey]
	s.mu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	// Build [ field1, value1, field2, value2, ... ]
	var out []Value
	for field, val := range entries {
		out = append(out,
			Value{typ: "bulk", bulk: field},
			Value{typ: "bulk", bulk: val},
		)
	}
	return Value{typ: "array", array: out}
}

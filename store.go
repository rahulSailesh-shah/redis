package main

import "sync"

type Store struct {
	mu   sync.RWMutex
	kv   map[string]string
	hash map[string]map[string]string
}

func NewStore() *Store {
	return &Store{
		kv:   map[string]string{},
		hash: map[string]map[string]string{},
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
		return Value{typ: "error", str: "'set' command requires exactly 2 arguments"}
	}

	k := args[0].bulk
	v := args[1].bulk

	s.mu.Lock()
	s.kv[k] = v
	s.mu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func (s *Store) get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "'get' command requires exactly 1 arguments"}
	}

	k := args[0].bulk
	s.mu.RLock()
	value, ok := s.kv[k]
	s.mu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func (s *Store) hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "'hset' command requires exactly 3 arguments"}
	}

	h := args[0].bulk
	k := args[1].bulk
	v := args[2].bulk

	s.mu.Lock()
	if _, ok := s.hash[h]; !ok {
		s.hash[h] = map[string]string{}
	}
	s.hash[h][k] = v
	s.mu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func (s *Store) hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "'hget' command requires exactly 2 arguments"}
	}

	h := args[0].bulk
	k := args[1].bulk

	s.mu.RLock()
	value, ok := s.hash[h][k]
	s.mu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

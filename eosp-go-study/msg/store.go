package main

import (
	"fmt"
	"sync"
)

type Store struct {
	mu   sync.Mutex
	data map[string][]ByteView
}

func NewStore() *Store {
	return &Store{
		mu:   sync.Mutex{},
		data: make(map[string][]ByteView),
	}
}

func (s *Store) addTopic(topic string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.hasTopic(topic) {
		return
	}
	s.data[topic] = make([]ByteView, 0)
}

func (s *Store) hasTopic(topic string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.data[topic]
	return ok
}

func (s *Store) Push(topic string, view ByteView) error {
	if !s.hasTopic(topic) {
		s.addTopic(topic)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[topic] = append(s.data[topic], view)
	return nil
}

func (s *Store) Get(topic string) (ByteView, error) {
	if !s.hasTopic(topic) {
		return nil, fmt.Errorf("[Store] topic %s not exit", topic)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.data[topic]) == 0 {
		return nil, nil
	}
	view := s.data[topic][0]
	s.data[topic] = s.data[topic][1:]
	return view, nil
}

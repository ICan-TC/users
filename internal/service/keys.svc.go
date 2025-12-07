package service

import (
	"os"

	"github.com/ICan-TC/lib/logging"
)

type KeysService struct {
	currentKey      string
	currentKeyIndex int
	keys            *[]string
}

func NewKeysService(keys *[]string) *KeysService {
	s := &KeysService{
		keys:       keys,
		currentKey: (*keys)[0],
	}
	s.Reset()
	return s
}

func (s *KeysService) SetKey(key string) string {
	s.currentKey = key
	l := logging.L()
	l.Info().Msg(s.currentKey)
	os.Setenv("GEMINI_API_KEY", s.currentKey)
	return s.currentKey
}

func (s *KeysService) CurrentKey() string {
	return s.currentKey
}

func (s *KeysService) Reset() {
	s.currentKeyIndex = 0
	s.SetKey((*s.keys)[s.currentKeyIndex])
}

func (s *KeysService) hasNext() bool {
	return s.currentKeyIndex < len(*s.keys)-1
}
func (s *KeysService) hasPrevious() bool {
	return s.currentKeyIndex > 0
}

func (s *KeysService) NextKey() string {
	if s.hasNext() {
		s.currentKeyIndex++
	} else {
		s.currentKeyIndex = 0
	}
	s.SetKey((*s.keys)[s.currentKeyIndex])
	return s.currentKey
}

func (s *KeysService) PreviousKey() string {
	if s.hasPrevious() {
		s.currentKeyIndex--
	} else {
		s.currentKeyIndex = len(*s.keys) - 1
	}
	s.SetKey((*s.keys)[s.currentKeyIndex])
	return s.currentKey
}

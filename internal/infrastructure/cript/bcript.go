package cript

import (
	"golang.org/x/crypto/bcrypt"
)

type Bcripter interface {
	HashPassword(password string) ([]byte, error)
	CompareHashAndPassword(hashPassword []byte, password string) error
}

type Bcript struct{}

func NewBcript() *Bcript {
	return &Bcript{}
}

func (b *Bcript) HashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (b *Bcript) CompareHashAndPassword(hashPassword []byte, password string) error {
	err := bcrypt.CompareHashAndPassword(hashPassword, []byte(password))
	if err != nil {
		return err
	}
	return nil
}

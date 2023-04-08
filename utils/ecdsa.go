package utils

import (
	"fmt"
	"math/big"
)

// Create struct for TX signature
type Signature struct {
	R *big.Int
	S *big.Int
}

// function to print the signature
func (s *Signature) string() string {
	return fmt.Sprintf("%x%x", s.R, s.S)
}

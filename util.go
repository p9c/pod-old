package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

var defaultUser, defaultPass = "user", "pa55word"

func GenerateKey() string {
	k, _ := rand.Int(rand.Reader, big.NewInt(int64(^uint32(0))))
	key := k.Uint64()
	return fmt.Sprintf("%0x", key)
}

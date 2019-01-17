package pod

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
)

var defaultUser, defaultPass = "user", "pa55word"

// GenerateKey gets a crypto-random number and encodes it in hex for generated shared credentials
func GenerateKey() string {
	k, _ := rand.Int(rand.Reader, big.NewInt(int64(^uint32(0))))
	key := k.Uint64()
	return fmt.Sprintf("%0x", key)
}

func ensureDir(fileName string) {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			panic(merr)
		}
	}
}

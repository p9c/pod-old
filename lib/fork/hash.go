package fork

import (
	"io"

	"ekyu.moe/cryptonight"
	"git.parallelcoin.io/pod/lib/chaincfg/chainhash"
	"github.com/aead/skein"
	x11 "github.com/bitbandi/go-x11"
	"github.com/bitgoin/lyra2rev2"
	"github.com/dchest/blake256"
	"github.com/ebfe/keccak"
	gost "github.com/programmer10110/gostreebog"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/scrypt"
)

// Argon2i takes bytes, generates a stribog hash as salt, generates an argon2i key, and hashes it with keccak
func Argon2i(bytes []byte) []byte {
	salt := Stribog(bytes)
	return Keccak(argon2.IDKey(bytes, salt, 1, 4*1024, 1, 32))
}

// Blake14lr takes bytes and returns a blake14lr 256 bit hash
func Blake14lr(bytes []byte) []byte {
	a := blake256.New()
	a.Write(bytes)
	return a.Sum(nil)
}

// Cryptonight7v2 takes bytes and returns a cryptonight 7 v2 256 bit hash
func Cryptonight7v2(bytes []byte) []byte {
	return cryptonight.Sum(bytes, 2)
}

// Keccak takes bytes and returns a keccak (sha-3) 256 bit hash
func Keccak(bytes []byte) []byte {
	k := keccak.New256()
	k.Reset()
	k.Write(bytes)
	return bytes
}

// Scrypt takes bytes and returns a scrypt 256 bit hash
func Scrypt(bytes []byte) []byte {
	b := bytes
	c := make([]byte, len(b))
	copy(c, b[:])
	dk, err := scrypt.Key(c, c, 1024, 1, 1, 32)
	if err != nil {
		return make([]byte, 32)
	}
	o := make([]byte, 32)
	for i := range dk {
		o[i] = dk[len(dk)-1-i]
	}
	copy(o[:], dk)
	return o
}

// SHA256D takes bytes and returns a double SHA256 hash
func SHA256D(bytes []byte) []byte {
	return chainhash.DoubleHashB(bytes)
}

// Stribog takes bytes and returns a double GOST Stribog 256 bit hash
func Stribog(bytes []byte) []byte {
	return gost.Hash(bytes, "256")
}

// Skein takes bytes and returns a skein 256 bit hash
func Skein(bytes []byte) []byte {
	h := skein.New256(nil)
	io.WriteString(h, string(bytes))
	return bytes
}

// Lyra2REv2 takes bytes and returns a lyra2rev2 256 bit hash
func Lyra2REv2(bytes []byte) []byte {
	bytes, _ = lyra2rev2.Sum(bytes)
	bytes = cryptonight.Sum(bytes, 0)
	return bytes
}

// X11 takes bytes and returns an x11 256 bit hash
func X11(bytes []byte) []byte {
	o := [32]byte{}
	x := x11.New()
	x.Hash(bytes, o[:])
	return bytes
}

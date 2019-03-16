package fork

import (
	"crypto/sha256"
	"io"
	"math/big"

	"ekyu.moe/cryptonight"
	"git.parallelcoin.io/dev/pod/pkg/chain/hash"
	"github.com/aead/skein"
	x11 "github.com/bitbandi/go-x11"
	"github.com/bitgoin/lyra2rev2"
	"github.com/dchest/blake256"
	"github.com/ebfe/keccak"
	gost "github.com/programmer10110/gostreebog"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/scrypt"
)

var bTwo = big.NewInt(2)

// Argon2i takes bytes, generates a Lyra2REv2 hash as salt, generates an argon2i key
func Argon2i(
	bytes []byte) []byte {
	return argon2.IDKey(Lyra2REv2(bytes), bytes, 1, 32*1024, 1, 32)
}

// Blake14lr takes bytes and returns a blake14lr 256 bit hash
func Blake14lr(
	bytes []byte) []byte {
	a := blake256.New()
	a.Write(bytes)
	return a.Sum(nil)
}

// Blake2b takes bytes and returns a blake2b 256 bit hash
func Blake2b(
	bytes []byte) []byte {
	b := blake2b.Sum256(bytes)
	return b[:]
}

// Blake2s takes bytes and returns a blake2s 256 bit hash
func Blake2s(
	bytes []byte) []byte {
	b := blake2s.Sum256(bytes)
	return b[:]
}

// Cryptonight7v2 takes bytes and returns a cryptonight 7 v2 256 bit hash
func Cryptonight7v2(
	bytes []byte) []byte {
	return cryptonight.Sum(bytes, 2)
}

// Hash computes the hash of bytes using the named hash
func Hash(
	bytes []byte, name string, height int32) (out chainhash.Hash) {

	// time.Sleep(time.Millisecond * 2000)
	switch name {
	case "blake2b":
		b := Argon2i(Cryptonight7v2(Blake2b(bytes)))
		out.SetBytes(rightShift(Blake2b(b)))
	case "blake14lr":
		b := Argon2i(Cryptonight7v2(Blake14lr(bytes)))
		out.SetBytes(rightShift(Blake14lr(b)))
	case "blake2s":
		b := Argon2i(Cryptonight7v2(Blake2s(bytes)))
		out.SetBytes(rightShift(Blake2s(b)))
	case "lyra2rev2":
		b := Argon2i(Cryptonight7v2(Lyra2REv2(bytes)))
		out.SetBytes(rightShift(Lyra2REv2(b)))
	case "scrypt":
		if GetCurrent(height) > 0 {
			b := Argon2i(Cryptonight7v2(Scrypt(bytes)))
			out.SetBytes(rightShift(Scrypt(b)))
		} else {
			out.SetBytes(Scrypt(bytes))
		}
	case "sha256d": // sha256d
		if GetCurrent(height) > 0 {
			b := Argon2i(Cryptonight7v2(chainhash.DoubleHashB(bytes)))
			out.SetBytes(rightShift(chainhash.DoubleHashB(b)))
		} else {
			out.SetBytes(chainhash.DoubleHashB(bytes))
		}
	case "stribog":
		b := Argon2i(Cryptonight7v2(Stribog(bytes)))
		out.SetBytes(rightShift(Stribog(b)))
	case "skein":
		b := Argon2i(Cryptonight7v2(Skein(bytes)))
		out.SetBytes(rightShift(Skein(b)))
	case "x11":
		b := Argon2i(Cryptonight7v2(X11(bytes)))
		out.SetBytes(rightShift(X11(b)))
	case "keccak":
		b := Argon2i(Cryptonight7v2(Keccak(bytes)))
		out.SetBytes(rightShift(Keccak(b)))
	}
	return
}

// Keccak takes bytes and returns a keccak (sha-3) 256 bit hash
func Keccak(
	bytes []byte) []byte {
	k := keccak.New256()
	k.Reset()
	k.Write(bytes)
	return bytes
}

// Lyra2REv2 takes bytes and returns a lyra2rev2 256 bit hash
func Lyra2REv2(
	bytes []byte) []byte {
	bytes, _ = lyra2rev2.Sum(bytes)
	return bytes
}

// SHA256D takes bytes and returns a double SHA256 hash
func SHA256D(
	bytes []byte) []byte {
	h := sha256.Sum256(bytes)
	h = sha256.Sum256(h[:])
	return h[:]
}

// Scrypt takes bytes and returns a scrypt 256 bit hash
func Scrypt(
	bytes []byte) []byte {
	b := bytes
	c := make([]byte, len(b))
	copy(c, b)
	dk, err := scrypt.Key(c, c, 1024, 1, 1, 32)
	if err != nil {
		return make([]byte, 32)
	}
	o := make([]byte, 32)
	for i := range dk {
		o[i] = dk[len(dk)-1-i]
	}
	copy(o, dk)
	return o
}

// Skein takes bytes and returns a skein 256 bit hash
func Skein(
	bytes []byte) []byte {
	h := skein.New256(nil)
	io.WriteString(h, string(bytes))
	return bytes
}

// Stribog takes bytes and returns a double GOST Stribog 256 bit hash
func Stribog(
	bytes []byte) []byte {
	return gost.Hash(bytes, "256")
}

// X11 takes bytes and returns an x11 256 bit hash
func X11(
	bytes []byte) []byte {
	o := [32]byte{}
	x := x11.New()
	x.Hash(bytes, o[:])
	return bytes
}

func rightShift(
	b []byte) (out []byte) {

	out = make([]byte, 32)
	copy(out, b[1:])
	return
}

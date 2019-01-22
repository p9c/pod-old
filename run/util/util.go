package pu

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"git.parallelcoin.io/pod/module/node"
	"github.com/tucnak/climax"
)

var defaultUser, defaultPass = "user", "pa55word"

// GenKey gets a crypto-random number and encodes it in hex for generated shared credentials
func GenKey() string {
	k, _ := rand.Int(rand.Reader, big.NewInt(int64(^uint32(0))))
	key := k.Uint64()
	return fmt.Sprintf("%0x", key)
}

// EnsureDir checks a file could be written to a path, creates the directories as needed
func EnsureDir(fileName string) {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			panic(merr)
		}
	}
}

// NormalizeAddress reads and corrects an address if it is missing pieces
func NormalizeAddress(addr, defaultPort string, out *string) {
	o := node.NormalizeAddress(addr, defaultPort)
	_, _, err := net.ParseCIDR(o)
	if err != nil {
		ip := net.ParseIP(addr)
		if ip != nil {
			out = &o
		}
	} else {
		out = &o
	}
}

// NormalizeAddresses reads and collects a space separated list of addresses contained in a string
func NormalizeAddresses(addrs string, defaultPort string, out *[]string) {
	O := new([]string)
	addrS := strings.Split(addrs, " ")
	for i := range addrS {
		a := addrS[i]
		o := ""
		NormalizeAddress(a, defaultPort, &o)
		if o != "" {
			*O = append(*O, o)
		}
	}
	// atomically switch out if there was valid addresses
	if len(*O) > 0 {
		out = O
	}
}

// ParseInteger reads a string that should contain a integer and returns the number and any parsing error
func ParseInteger(integer, name string, original *int) (err error) {
	var out int
	out, err = strconv.Atoi(integer)
	if err != nil {
		err = fmt.Errorf("malformed %s `%s` leaving set at `%d` err: %s", name, integer, *original, err.Error())
	} else {
		*original = out
	}
	return
}

// ParseUint32 reads a string that should contain a integer and returns the number and any parsing error
func ParseUint32(integer, name string, original *uint32) (err error) {
	var out int
	out, err = strconv.Atoi(integer)
	if err != nil {
		err = fmt.Errorf("malformed %s `%s` leaving set at `%d` err: %s", name, integer, *original, err.Error())
	} else {
		*original = uint32(out)
	}
	return
}

// ParseFloat reads a string that should contain a floating point number and returns it and any parsing error
func ParseFloat(f, name string, original *float64) (err error) {
	var out float64
	_, err = fmt.Sscanf(f, "%0.f", out)
	if err != nil {
		err = fmt.Errorf("malformed %s `%s` leaving set at `%0.f` err: %s", name, f, *original, err.Error())
	} else {
		*original = out
	}
	return
}

// ParseDuration takes a string of the format `Xd/h/m/s` and returns a time.Duration corresponding with that specification
func ParseDuration(d, name string, out *time.Duration) (err error) {
	var t int
	var ti time.Duration
	switch d[len(d)-1] {
	case 's':
		t, err = strconv.Atoi(d[:len(d)-1])
		ti = time.Duration(t) * time.Second
	case 'm':
		t, err = strconv.Atoi(d[:len(d)-1])
		ti = time.Duration(t) * time.Minute
	case 'h':
		t, err = strconv.Atoi(d[:len(d)-1])
		ti = time.Duration(t) * time.Hour
	case 'd':
		t, err = strconv.Atoi(d[:len(d)-1])
		ti = time.Duration(t) * 24 * time.Hour
	}
	if err != nil {
		err = fmt.Errorf("malformed %s `%s` leaving set at `%s` err: %s", name, d, *out, err.Error())
	} else {
		*out = ti
	}
	return
}

// GenFlag allows a flag to be more concisely declared
func GenFlag(name, short, usage, help string, variable bool) climax.Flag {
	return climax.Flag{
		Name:     name,
		Short:    short,
		Usage:    usage,
		Help:     help,
		Variable: variable,
	}
}

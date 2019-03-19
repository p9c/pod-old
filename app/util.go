package app

import (
	"net"
	"os"
	"path/filepath"
	"strings"

	"git.parallelcoin.io/dev/pod/cmd/node"
)

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

// FileExists reports whether the named file or directory exists.
func FileExists(filePath string) bool {

	_, err := os.Stat(filePath)
	return err == nil
}

// NormalizeAddress reads and corrects an address if it is missing pieces
func NormalizeAddress(addr, defaultPort string, out *string) {

	o := node.NormalizeAddress(addr, defaultPort)
	_, _, err := net.ParseCIDR(o)

	if err != nil {

		ip := net.ParseIP(addr)

		if ip != nil {

			*out = o
		}

	} else {

		*out = o
	}

}

// NormalizeAddresses reads and collects a space separated list of addresses contained in a string
func NormalizeAddresses(addrs string, defaultPort string, out *[]string) {

	O := new([]string)
	addrS := strings.Split(addrs, " ")

	for i := range addrS {

		a := addrS[i]

		// o := ""
		NormalizeAddress(a, defaultPort, &a)

		if a != "" {

			*O = append(*O, a)
		}

	}

	// atomically switch out if there was valid addresses

	if len(*O) > 0 {

		*out = *O
	}

}

// minUint32 is a helper function to return the minimum of two uint32s. This avoids a math import and the need to cast to floats.
func minUint32(a, b uint32) uint32 {

	if a < b {

		return a
	}

	return b
}

// CleanAndExpandPath expands environment variables and leading ~ in the passed path, cleans the result, and returns it.
func CleanAndExpandPath(path string) string {

	// Expand initial ~ to OS specific home directory.
	if strings.HasPrefix(path, "~") {

		homeDir := filepath.Dir(DefaultDataDir)
		path = strings.Replace(path, "~", homeDir, 1)
	}

	// NOTE: The os.ExpandEnv doesn't work with Windows-style %VARIABLE%, but they variables can still be expanded via POSIX-style $VARIABLE.
	return filepath.Clean(os.ExpandEnv(path))
}

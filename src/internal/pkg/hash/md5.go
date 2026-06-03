// Package hash provides password hashing utilities.
package hash

import (
	"crypto/md5"
	"fmt"
)

// MD5 returns the MD5 checksum of the given string as a hex string.
func MD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

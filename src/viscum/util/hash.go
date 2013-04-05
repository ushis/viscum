package util

import (
  "crypto/sha1"
  "encoding/hex"
  "fmt"
)

// Calculates the SHA1 sum of something.
func Sha1Sum(args ...interface{}) (string, error) {
  h := sha1.New()

  if _, err := fmt.Fprint(h, args...); err != nil {
    return "", err
  }
  return hex.EncodeToString(h.Sum(nil)), nil
}

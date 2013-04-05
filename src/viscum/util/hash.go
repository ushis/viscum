package util

import (
  "crypto/sha1"
  "encoding/hex"
)

// Calculates the SHA1 sum of bytes.
func Sha1Sum(b []byte) (string, error) {
  h := sha1.New()

  if _, err := h.Write(b); err != nil {
    return "", err
  }
  return hex.EncodeToString(h.Sum(nil)), nil
}

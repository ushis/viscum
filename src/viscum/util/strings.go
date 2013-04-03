package util

import (
  "bytes"
  "unicode"
)

// Returns a string with squeezed whitespace and removed line breaks.
func Titleize(src string) string {
  sp, end := 0, len(src)

  for ; end > 0; end-- {
    if ! unicode.IsSpace(rune(src[end-1])) {
      break
    }
  }

  var buf bytes.Buffer

  for i, c := range src[:end] {
    if ! unicode.IsSpace(c) {
      buf.WriteRune(c)
      continue
    }
    if sp < i {
      buf.WriteByte(' ')
    }
    sp = i + 1
  }

  return buf.String()
}

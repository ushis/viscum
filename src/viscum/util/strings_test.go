package util

import (
  "testing"
)

const TITLE = "Hello, my name is mud."

const TITLE1 = "  Hello,   my        name is mud.    "

const TITLE2 = `

    Hello,   my

  name
    is mud.
`

func TestTitleize(t *testing.T) {
  for i, s := range []string{TITLE, TITLE1, TITLE2} {
    if res := Titleize(s); res != TITLE {
      t.Errorf("[%d] \"%s\"", i, res)
    }
  }
}

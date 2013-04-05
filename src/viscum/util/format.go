package util

import (
  "bytes"
  "fmt"
  "github.com/moovweb/gokogiri"
  "github.com/moovweb/gokogiri/xml"
  "strings"
)

type Formatter struct {
  buf   bytes.Buffer
  links []string
}

var whitespace = map[byte]bool{
  ' ':  true,
  '\n': true,
  '\r': true,
  '\t': true,
  '\v': true,
}

var ignore = map[string]bool{
  "audio":  true,
  "head":   true,
  "script": true,
  "track":  true,
  "video":  true,
}

var block = map[string]bool{
  "address":    true,
  "article":    true,
  "aside":      true,
  "blockquote": true,
  "body":       true,
  "canvas":     true,
  "del":        true,
  "div":        true,
  "dl":         true,
  "fieldset":   true,
  "figcaption": true,
  "figure":     true,
  "footer":     true,
  "form":       true,
  "header":     true,
  "hgroup":     true,
  "hr":         true,
  "ins":        true,
  "menu":       true,
  "noscript":   true,
  "ol":         true,
  "output":     true,
  "p":          true,
  "section":    true,
  "table":      true,
  "td":         true,
  "tfoot":      true,
  "th":         true,
  "thead":      true,
  "tr":         true,
  "ul":         true,
}

var heading = map[string]bool{
  "h1": true,
  "h2": true,
  "h3": true,
  "h4": true,
  "h5": true,
  "h6": true,
}

var italic = map[string]bool{
  "em": true,
  "i":  true,
}

var bold = map[string]bool{
  "b":      true,
  "strong": true,
}

func Format(src *string) error {
  doc, err := gokogiri.ParseHtml([]byte(*src))

  if err != nil {
    return err
  }
  defer doc.Free()

  f := new(Formatter)
  f.walk(doc.Node)

  for i, link := range f.links {
    f.buf.WriteString(fmt.Sprintf("[%d] %s\n", i, link))
  }

  *src = f.buf.String()
  return nil
}

func (self *Formatter) walk(node xml.Node) {
  for c := node.FirstChild(); c != nil; c = c.NextSibling() {
    self.walk(c)
  }

  if node.NodeType() == xml.XML_ELEMENT_NODE {
    self.handleNode(node)
  }
}

func (self *Formatter) handleNode(node xml.Node) {
  name := node.Name()

  switch {
  case ignore[name]:
    node.SetContent("")
  case name == "pre":
    self.writeCodeBlock(node)
  case heading[name]:
    self.writeBlock(node, "# ")
  case name == "li":
    self.writeBlock(node, "- ")
  case name == "br":
    node.SetContent("\n")
  case italic[name]:
    node.SetContent("/" + node.Content() + "/")
  case bold[name]:
    node.SetContent("*" + node.Content() + "*")
  case name == "img":
    alt, src := node.Attr("alt"), node.Attr("src")

    if len(alt) > 0 && len(src) > 0 {
      node.SetContent(fmt.Sprintf("(%s)[%d]", alt, len(self.links)))
      self.links = append(self.links, src)
    }
  case name == "a":
    href, content := node.Attr("href"), node.Content()

    if len(href) > 0 && len(content) > 0 {
      node.SetContent(fmt.Sprintf("%s[%d]", content, len(self.links)))
      self.links = append(self.links, href)
    }
  case block[name]:
    self.writeBlock(node, "")
  }
}

func (self *Formatter) writeBlock(node xml.Node, prefix string) {
  sp, br, max := 0, 0, 79-len(prefix)
  block := []byte(strings.TrimSpace(node.Content()))
  node.SetContent("")

  if len(block) == 0 {
    return
  }
  self.buf.WriteString(prefix)

  for i, c := range block {
    if i-br > max && sp > br {
      self.buf.WriteByte('\n')
      br = sp

      for j := 0; j < len(prefix); j++ {
        self.buf.WriteByte(' ')
      }
    }
    if whitespace[c] {
      if sp == i {
        sp++
        br++
        continue
      }
      self.buf.Write(block[sp:i])
      self.buf.WriteByte(' ')
      sp = i + 1
    }
  }

  if sp < len(block) {
    self.buf.Write(block[sp:])
  }
  self.buf.Write([]byte{'\n', '\n'})
}

func (self *Formatter) writeCodeBlock(node xml.Node) {
  block := []byte(node.Content())
  node.SetContent("")

  for i := len(block) - 1; i >= 0; i-- {
    if !whitespace[block[i]] {
      self.buf.Write(block[:i+1])
      self.buf.Write([]byte{'\n', '\n'})
      return
    }
  }
}

package util

import (
  "bytes"
  "fmt"
  "github.com/moovweb/gokogiri"
)

func Text(src string) (string, error) {
  doc, err := gokogiri.ParseHtml([]byte(src))

  if err != nil {
    return "", err
  }

  links := []string{}

  nodes, _ := doc.Search(".//br")

  for _, node := range nodes {
    node.Replace("\n")
  }

  // Replace image tags.
  nodes, _ = doc.Search(".//img")

  for _, node := range nodes {
    attr := node.Attributes()
    alt, aok := attr["alt"]
    src, sok := attr["src"]

    if !aok || !sok {
      continue
    }

    node.Replace(fmt.Sprintf("(%s)[%d]", alt.Content(), len(links)))
    links = append(links, src.Content())
  }

  // replace anchors.
  nodes, _ = doc.Search(".//a")

  for _, node := range nodes {
    if href, ok := node.Attributes()["href"]; ok {
      links = append(links, href.Content())
      node.Replace(fmt.Sprintf("%s[%d]", node.Content(), len(links)))
    }
  }

  // Italic
  nodes, _ = doc.Search(".//i|.//em")

  for _, node := range nodes {
    node.Replace("/" + node.Content() + "/")
  }

  // Bold
  nodes, _ = doc.Search(".//b|.//strong|")

  for _, node := range nodes {
    node.Replace("*" + node.Content() + "*")
  }

  // Headlines
  nodes, _ = doc.Search(".//h1|.//h2|.//h3|.//h4|.//h5|.//h6")

  for _, node := range nodes {
    node.SetContent("# " + node.Content())
  }

  // Quotes
  nodes, _ = doc.Search(".//blockquote/*")

  for _, node := range nodes {
    node.SetContent("> " + node.Content())
  }

  var buf bytes.Buffer
  buf.WriteString(doc.Content())
  buf.Write([]byte{'\n', '\n'})

  for i, link := range links {
    buf.WriteString(fmt.Sprintf("[%d] %s\n", i, link))
  }

  return buf.String(), nil
}

package util

import (
	"bytes"
	"encoding/json"

	"github.com/dstotijn/go-notion"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
)

func MDToNotionBlock(source []byte) []notion.Block {
	var buf bytes.Buffer

	buf.WriteString(`{"results": [`)

	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithBlockParsers(),
			parser.WithParagraphTransformers(),
		),
		goldmark.WithExtensions(NotionExtension),
	)
	md.Convert(source, &buf)

	buf.WriteString(`]}`)

	b := bytes.ReplaceAll(buf.Bytes(), []byte("}{"), []byte("},{"))
	b2 := bytes.ReplaceAll(b, []byte("\n"), []byte("\\n"))
	b3 := bytes.ReplaceAll(b2, []byte("}\"}"), []byte("}}"))

	var param notion.BlockChildrenResponse
	json.Unmarshal(b3, &param)
	return param.Results
}

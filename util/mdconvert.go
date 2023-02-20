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

	var param notion.BlockChildrenResponse
	json.Unmarshal(bytes.ReplaceAll(buf.Bytes(), []byte("}{"), []byte("},{")), &param)
	return param.Results
}

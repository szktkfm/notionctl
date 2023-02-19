package util

import (
	"bytes"
	"encoding/json"

	"github.com/dstotijn/go-notion"
	"github.com/yuin/goldmark"
)

func MDToNotionBlock(source []byte) []notion.Block {
	var buf bytes.Buffer

	buf.WriteString(`{"results": [`)

	md := goldmark.New(
		goldmark.WithExtensions(NotionExtension),
	)
	md.Convert(source, &buf)

	buf.Truncate(buf.Len() - 1)
	buf.WriteString(`]}`)

	var got notion.BlockChildrenResponse
	json.Unmarshal(buf.Bytes(), &got)
	return got.Results
}

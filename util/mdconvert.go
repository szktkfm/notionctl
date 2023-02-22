package util

import (
	"bytes"
	"encoding/binary"
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

	// TODO: clear
	b := bytes.ReplaceAll(buf.Bytes(), []byte("}{"), []byte("},{"))
	b = bytes.ReplaceAll(b, []byte("\n"), []byte("\\n"))
	b = bytes.ReplaceAll(b, []byte("}\"}"), []byte("}}"))

	var out []byte
	out = b
	for i, b := range htmlEscapeTable {
		if b != nil {
			bi := make([]byte, 8)
			binary.BigEndian.PutUint64(bi, uint64(i))
			if bytes.HasPrefix(b, []byte("&quot;")) {
				out = bytes.ReplaceAll(out, b, append([]byte("\\"), bi[7:8]...))
			} else {
				out = bytes.ReplaceAll(out, b, bi[7:8])
			}
		}
	}

	var param notion.BlockChildrenResponse
	json.Unmarshal(out, &param)
	return param.Results
}

var htmlEscapeTable = [256][]byte{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, []byte("&quot;"), nil, nil, nil, []byte("&amp;"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, []byte("&lt;"), nil, []byte("&gt;"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil}

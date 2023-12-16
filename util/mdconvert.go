package util

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"

	"github.com/dstotijn/go-notion"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
)

func MDToNotionBlock(source []byte) []notion.Block {
	var buf bytes.Buffer
	convertMarkdownToBuffer(source, &buf)

	formattedJSON := formatForNotion(buf.Bytes())

	var param notion.BlockChildrenResponse
	if err := json.Unmarshal(formattedJSON, &param); err != nil {
		log.Fatal("Cannot unmarshal JSON", err)
	}
	return param.Results
}

func convertMarkdownToBuffer(source []byte, buf *bytes.Buffer) {
	buf.WriteString(`{"results": [`)

	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithBlockParsers(),
			parser.WithParagraphTransformers(),
		),
		goldmark.WithExtensions(NotionExtension),
	)
	md.Convert(source, buf)

	buf.WriteString(`]}`)
}

func formatForNotion(data []byte) []byte {
	data = bytes.ReplaceAll(data, []byte("}{"), []byte("},{"))
	data = bytes.ReplaceAll(data, []byte("\n"), []byte("\\n"))
	data = bytes.ReplaceAll(data, []byte("}\"}"), []byte("}}"))

	return escapeHTMLCharacters(data)
}

func escapeHTMLCharacters(data []byte) []byte {
	for i, b := range htmlEscapeTable {
		if b != nil {
			bi := make([]byte, 8)
			binary.BigEndian.PutUint64(bi, uint64(i))
			if bytes.HasPrefix(b, []byte("&quot;")) {
				data = bytes.ReplaceAll(data, b, append([]byte("\\"), bi[7:8]...))
			} else {
				data = bytes.ReplaceAll(data, b, bi[7:8])
			}
		}
	}
	return data
}

var htmlEscapeTable = [256][]byte{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, []byte("&quot;"), nil, nil, nil, []byte("&amp;"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, []byte("&lt;"), nil, []byte("&gt;"), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil}

package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"testing"
)

func TestGetCmdRunErr(t *testing.T) {

	tests := []struct {
		name       string
		option     *GetOptions
		queryMock  string
		wantErrMsg string
	}{
		{
			name: "test fot not setting env vars",
			option: &GetOptions{
				DB:   "hogedb",
				Page: "",
				Wide: false,
			},
			wantErrMsg: "required env variables NOTION_API_KEY not set",
		},
	}
	for _, tt := range tests {
		os.Unsetenv("NOTION_API_KEY")
		t.Setenv("NOTION_DATABASE", "testdbid")

		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			tt.option.Out = buf
			cmd := newCmdGet(tt.option, buf)
			gotErr := tt.option.Complete(cmd, []string{"test"})
			if tt.wantErrMsg != fmt.Sprint(gotErr) {
				t.Errorf("value is mismatch. want: %s, got: %s", tt.wantErrMsg, gotErr)
			}
		})
	}
}

func TestGetCmdRun(t *testing.T) {

	tests := []struct {
		name      string
		option    *GetOptions
		queryMock string
		want      string
	}{
		{
			name: "test full column",
			option: &GetOptions{
				DB:   "hogedb",
				Page: "",
				Wide: false,
			},
			// TODO: define []notion.Page and unmashal it
			queryMock: `{"object":"list","results":[{"object":"page","id":"711e6ef1-28c6-482e-91f2-75dd26dfd041","created_time":"2022-12-05T11:58:00.000Z","last_edited_time":"2022-12-24T16:44:00.000Z","created_by":{"object":"user","id":"5a639bd5-f786-4565-bc65-5d9281ef3944"},"last_edited_by":{"object":"user","id":"25b9e72d-a868-4007-bf74-841efc304d3e"},"cover":{"type":"external","external":{"url":"https://upload.wikimedia.org/wikipedia/commons/6/62/Tuscankale.jpg"}},"icon":{"type":"emoji","emoji":"🥬"},"parent":{"type":"database_id","database_id":"98079428-d5d0-436f-a316-b2d36da049c2"},"archived":false,"properties":{"Food group":{"id":"B%60Ts","type":"select","select":{"id":"26ab5fc8-7e6b-4d11-b6c5-6864e614c3ed","name":"Vegetable","color":"purple"}},"Description":{"id":"oBRk","type":"rich_text","rich_text":[{"type":"text","text":{"content":"A dark green leafy vegetable","link":null},"annotations":{"bold":false,"italic":false,"strikethrough":false,"underline":false,"code":false,"color":"default"},"plain_text":"A dark green leafy vegetable","href":null}]},"Created time":{"id":"rlOQ","type":"created_time","created_time":"2022-12-05T11:58:00.000Z"},"Tags":{"id":"urJ%5B","type":"multi_select","multi_select":[]},"Price":{"id":"xyUL","type":"number","number":2.5},"Name":{"id":"title","type":"title","title":[{"type":"text","text":{"content":"Tuscan Kale","link":null},"annotations":{"bold":false,"italic":false,"strikethrough":false,"underline":false,"code":false,"color":"default"},"plain_text":"Tuscan Kale","href":null}]}},"url":"https://www.notion.so/Tuscan-Kale-711e6ef128c6482e91f275dd26dfd041"}],"next_cursor":null,"has_more":false,"type":"page","page":{}}`,
			want:      `NAME\s+AGE\s+FOOD GROUP\s+TAGS\s+DESCRIPTION\s+\nTuscan Kale\s+\d*.*\s+Vegetable\s+-\s+A dark green leafy veget\s+\n`,
		},
	}
	for _, tt := range tests {
		t.Setenv("NOTION_API_KEY", "testkey")
		t.Setenv("NOTION_DATABASE", "testdbid")

		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			tt.option.Out = buf

			tt.option.Client = NewMockClient(
				&http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(bytes.NewBufferString(tt.queryMock)),
				},
				"testkey",
			)
			cmd := newCmdGet(tt.option, buf)

			tt.option.Run(cmd)
			got := buf.String()

			r := regexp.MustCompile(tt.want)
			if !r.MatchString(got) {
				t.Errorf("print value is mismatch. want: %s, got: %s", tt.want, got)
			}
		})
	}
}

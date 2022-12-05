package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(newCmdPush())

	// viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>"
	// viper.SetDefault("license", "apache")
}

type PushOptions struct {
	Db string
}

// newでcmdを返す。new関数の中でadd cmdする
func newCmdPush() *cobra.Command {
	o := &PushOptions{}
	cmd := &cobra.Command{
		Use:   "push",
		Short: "push text",
		RunE: func(cmd *cobra.Command, args []string) error {
			//debug
			// fmt.Println(getSecret())
			// flagの穴埋め
			o.Complete(cmd, args)
			// run
			o.Run(cmd, args)
			return nil
		},
	}

	cmd.Flags().StringVar(&o.Db, "db", "", "db id")
	return cmd
}

func (o *PushOptions) Complete(cmd *cobra.Command, args []string) error {
	o.Db = viper.GetString("db") // ToDo: db引数の値を優先する。でなければviperで読み取る

	return nil
}

func (o *PushOptions) Run(cmd *cobra.Command, args []string) error {
	fmt.Println(o.Db)
	buf := &bytes.Buffer{}
	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: &httpTransport{w: buf},
	}
	client := notion.NewClient(getSecret(), notion.WithHTTPClient(httpClient))

	// 一回dbを情報をgetしてきて、そこからparameterをbuildする。
	db, err := client.FindDatabaseByID(context.Background(), o.Db)
	if err != nil {
		return err
	}

	fmt.Println(db.ID)

	params := newParams(db) //ToDo:  param構造体を作る
	// param.build みたいな感じでparameterをbuildして

	fmt.Printf("params: %#v\n", params)
	page, err := client.CreatePage(context.Background(), params)
	fmt.Printf("bufer: %s", buf)

	if err != nil {
		fmt.Println("error")
		fmt.Println(err)
		return err
	}
	// print
	// TODO: tabwriterを別のpackageにする
	const padding = 4
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.TabIndent)
	// printPage(page, w)
	fmt.Println(page)
	w.Flush()
	// fmt.Println(queryResult)
	return nil

}

func newParams(db notion.Database) notion.CreatePageParams {
	dbPageProp := make(notion.DatabasePageProperties)

	for k, dp := range db.Properties {
		// fmt.Printf("properties: %#v\n", k)
		// fmt.Printf("Properties value: %#v\n", dp)
		switch dp.Type {
		case notion.DBPropTypeTitle:
			dbPageProp[k] = notion.DatabasePageProperty{
				Title: getRitchText("test title"),
			}
		case notion.DBPropTypeRichText:
			dbPageProp[k] = notion.DatabasePageProperty{
				RichText: getRitchText("this_is_description"),
			}
			// case notion.DBPropTypeNumber:
			// 	a := 10.0
			// 	dbPageProp[k] = notion.DatabasePageProperty{
			// 		Number: &a,
			// 	}
		}

		// var dpp notion.DatabasePageProperty
		// propJson, _ := json.Marshal(dp)
		// fmt.Printf("json: %s\n", propJson)
		// if err != nil {
		// }

		// if err := json.Unmarshal(propJson, &dpp); err != nil {
		// }
		// dbPageProp[k] = dpp
		fmt.Println(k)
	}

	fmt.Println()
	fmt.Printf("%#v \n", dbPageProp)
	fmt.Println()

	return notion.CreatePageParams{
		ParentType:             notion.ParentTypeDatabase,
		ParentID:               db.ID,
		Title:                  getRitchText("title 2"),
		DatabasePageProperties: &dbPageProp,
		Children: []notion.Block{
			notion.Heading1Block{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "testtest",
						},
					},
				},
			},
		},
	}
}

func getRitchText(content string) []notion.RichText {
	return []notion.RichText{
		{
			Text: &notion.Text{
				Content: content,
			},
		},
	}
}

// func getDatabesePageProperties() *notion.DatabasePageProperties {
// 	dpp := make(notion.DatabasePageProperties)

// 	//TODo: propertiesはdbの設定依存であることに注意
// 	dpp["名前"] = notion.DatabasePageProperty{
// 		Name:  "testpage",
// 		Title: getRitchText(),
// 	}

// 	return &dpp
// }

func printPage(page notion.Page, w io.Writer) error {
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", "page ID", "TITLE", "CREATED TIME", "TAGS")
	// Propertiesから情報を持ってくる処理
	props := page.Properties
	// fmt.Printf("%#v", props.(notion.DatabasePageProperties))

	var title string
	var multiSelect string
	for _, v := range props.(notion.DatabasePageProperties) {
		// fmt.Println(k)
		// fmt.Println(v.ID)
		// pageId := v.ID
		// fmt.Println(v.Type)
		switch v.Type {
		case notion.DBPropTypeTitle:
			title = fmt.Sprintf("%s", v.Title[0].Text.Content) //pythonにおけるmap的な書き方できないんだろうか?
		// case notion.DBPropTypeCreatedTime:
		// 	createdTime := v.CreatedTime
		case notion.DBPropTypeMultiSelect:
			multiSelect = fmt.Sprintf("%s", v.MultiSelect)
		}
		// fmt.Println(v.Title)
		// fmt.Println(v.CreatedTime)
		// fmt.Println("")
	}
	// for key := range props {
	// 	fmt.Println(key)
	// }
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", page.ID, title, page.CreatedTime, multiSelect)
	return nil
}

// RoundTrip implements http.RoundTripper. It multiplexes the read HTTP response
// data to an io.Writer for debugging.
func (t *httpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	res.Body = io.NopCloser(io.TeeReader(res.Body, t.w))

	return res, nil
}

type httpTransport struct {
	w io.Writer
}

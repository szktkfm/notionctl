package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/dstotijn/go-notion"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(newCmdGet())

	// viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>"
	// viper.SetDefault("license", "apache")
}

type GetOptions struct {
	Db   string
	Page string
}

// newでcmdを返す。new関数の中でadd cmdする
func newCmdGet() *cobra.Command {
	o := &GetOptions{}
	cmd := &cobra.Command{
		Use:   "get",
		Short: "config on .notion-go",
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

	cmd.Flags().StringVar(&o.Page, "page", "", "page id")
	cmd.Flags().StringVar(&o.Db, "db", "", "db id")
	return cmd
}

func (o *GetOptions) Complete(cmd *cobra.Command, args []string) error {
	o.Db = viper.GetString("db") // ToDo: db引数の値を優先する。でなければviperで読み取る。
	// viperをつかってconfigから読み取る(configから読み取らずviper packageのmethodをそのまま使ってしまうと、singletonみたいになって、実行順序よって結果が変わるかも)
	o.Page, _ = cmd.Flags().GetString("page")

	return nil
}

func (o *GetOptions) Run(cmd *cobra.Command, args []string) error {
	if len(o.Page) > 0 {
		// ToDo: get markdown and build name
		name := "sample"
		fmt.Printf("get markdown to %s.md\n", name)
		return nil
	}
	client := notion.NewClient(getSecret())
	queryResult, err := client.QueryDatabase(context.Background(), o.Db, nil)

	if err != nil {
		return err
	}
	// print
	// TODO: tabwriterを別のpackageにする
	const padding = 4
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.TabIndent)
	printDatabaseQueryResponce(queryResult, w)
	w.Flush()
	// fmt.Println(queryResult)
	return nil

}

func printDatabaseQueryResponce(res notion.DatabaseQueryResponse, w io.Writer) error {
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", "page ID", "TITLE", "CREATED TIME", "TAGS")
	for _, page := range res.Results {
		// Propertiesから情報を持ってくる処理
		props := page.Properties
		// fmt.Printf("%#v", props.(notion.DatabasePageProperties))

		var title string
		var multiSelect string
		for k, v := range props.(notion.DatabasePageProperties) {
			// fmt.Printf("property key: %s", k)
			if k == "note" {
				// fmt.Printf("%#v \n", v)
			}

			// fmt.Println(k)
			// fmt.Println(v.ID)
			// pageId := v.ID
			// fmt.Println(v.Type)
			switch v.Type {
			case notion.DBPropTypeTitle:
				title = fmt.Sprintf("%s", v.Title[0].Text.Content) //pythonにおけるmap的な書き方できないんだろうか?
				fmt.Printf("title: %#v \n", k)
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
	}
	return nil
}

func readConfig() {
	fmt.Println(viper.GetString("secret"))
}

func getSecret() string {
	return viper.GetString("secret")
}

package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/notionctl/util"
	"mkuznets.com/go/tabwriter"

	"github.com/dstotijn/go-notion"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newCmdGet(&GetOptions{}, os.Stdout))

	// viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>"
	// viper.SetDefault("license", "apache")
}

type GetOptions struct {
	DB   string
	Page string
	Wide bool
	Out  io.Writer
}

type utilFactory struct {
	client http.Client
}

// newでcmdを返す。new関数の中でadd cmdする
func newCmdGet(o *GetOptions, writer io.Writer) *cobra.Command {
	// option := &GetOptions{Out: writer}
	o.Out = writer
	cmd := &cobra.Command{
		Use:   "get",
		Short: "",
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

	// cmd.Flags().StringVar(&o.Page, "page", o.Page, "page id") //このdefault値を引数で渡したoptionの値にする
	// cmd.Flags().StringVar(&o.DB, "db", o.DB, "db id")
	cmd.Flags().BoolVar(&o.Wide, "wide", o.Wide, "wide print")
	return cmd
}

func (o *GetOptions) Complete(cmd *cobra.Command, args []string) error {
	// o.DB = viper.GetString("db") // ToDo: 面倒だからenv variableから読みだそう
	o.DB = os.Getenv("NOTION_DATABASE") // ToDo: 面倒だからenv variableから読みだそう
	// viperをつかってconfigから読み取る(configから読み取らずviper packageのmethodをそのまま使ってしまうと、singletonみたいになって、実行順序よって結果が変わるかも)
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
	queryResult, err := client.QueryDatabase(context.Background(), o.DB, nil)

	if err != nil {
		return err
	}
	// w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)
	w := tabwriter.NewWriter(o.Out, 4, 0, 4, ' ', tabwriter.TabIndent)
	// util.PrintDatabaseQueryResponce(queryResult, w)
	util.NewDatabaseQueryResponcePrinter(queryResult, w).Print()
	w.Flush()
	// fmt.Println(queryResult)
	return nil

}

func getSecret() string {
	// TODO: 環境変数から読み取ろう。
	// return viper.GetString("secret")
	return os.Getenv("NOTION_API_KEY")
}

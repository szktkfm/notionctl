package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/szktkfm/notionctl/util"
	"mkuznets.com/go/tabwriter"

	"github.com/dstotijn/go-notion"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newCmdGet(&GetOptions{}, os.Stdout))
}

type GetOptions struct {
	DB   string
	Page string
	Wide bool
	Out  io.Writer
}

func newCmdGet(o *GetOptions, writer io.Writer) *cobra.Command {
	o.Out = writer
	cmd := &cobra.Command{
		Use:   "get",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Complete(cmd, args)
			o.Run(cmd, args)
			return nil
		},
	}

	cmd.Flags().BoolVar(&o.Wide, "wide", o.Wide, "wide print")
	return cmd
}

func (o *GetOptions) Complete(cmd *cobra.Command, args []string) error {
	o.DB = os.Getenv("NOTION_DATABASE")
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
	w := tabwriter.NewWriter(o.Out, 4, 0, 4, ' ', tabwriter.TabIndent)
	util.NewDatabaseQueryResponcePrinter(queryResult, w).Print()
	w.Flush()
	return nil

}

func getSecret() string {
	return os.Getenv("NOTION_API_KEY")
}

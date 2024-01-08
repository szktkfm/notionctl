package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"mkuznets.com/go/tabwriter"

	"github.com/dstotijn/go-notion"
	"github.com/spf13/cobra"
	"github.com/szktkfm/notionctl/internal/printer"
)

func init() {
	rootCmd.AddCommand(newCmdGet(&GetOptions{}, os.Stdout))
}

type GetOptions struct {
	DB     string
	Page   string
	Wide   bool
	Out    io.Writer
	Client *notion.Client
}

func newCmdGet(o *GetOptions, writer io.Writer) *cobra.Command {
	o.Out = writer
	cmd := &cobra.Command{
		Use:   "get",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Complete(cmd, args)
			if err != nil {
				return err
			}
			err = o.Run(cmd)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&o.Wide, "wide", o.Wide, "wide print")
	return cmd
}

func (o *GetOptions) Complete(cmd *cobra.Command, args []string) error {
	db, err := getDBID()
	if err != nil {
		return err
	}
	o.DB = db

	secret, err := getSecret()
	if err != nil {
		return err
	}

	o.Client = notion.NewClient(secret)
	return nil
}

func (o *GetOptions) Run(cmd *cobra.Command) error {
	if len(o.Page) > 0 {
		// TODO: Building markdown from API responses.
		name := "sample"
		fmt.Printf("get markdown to %s.md\n", name)
		return nil
	}
	queryResult, err := o.Client.QueryDatabase(context.Background(), o.DB, nil)

	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(o.Out, 4, 0, 4, ' ', tabwriter.TabIndent)
	printer.NewDBQueryRespTablePrinter(queryResult, w).Print()
	w.Flush()
	return nil
}

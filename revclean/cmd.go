package revclean

import (
	"github.com/spf13/cobra"
	"fmt"
	"strings"
	"github.com/yoavl/coucher/database"
)

var docid, revision string
var dry bool

func NewCmd(db *database.Database) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "revclean",
		Short: "Cleans up conflicting doc revisions",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("RevClean: " + strings.Join(args, " "))
			CleanRevs(db, docid, revision, dry)
		},
	}

	cmd.PersistentFlags().StringVarP(&docid, "docid", "i", "", "Document id")
	cmd.PersistentFlags().StringVarP(&revision, "revision", "r", "", "Optional revision to pick - NOT SUPPORTED YET")
	cmd.PersistentFlags().BoolVarP(&dry, "dry", "", false, "Dry run")
	return cmd
}

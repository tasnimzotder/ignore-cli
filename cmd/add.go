package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tasnimzotder/ignore-cli/internal/gitignore"
)

var overrideFlag bool

var addCmd = &cobra.Command{
	Use:   "add <template>",
	Short: "Add a .gitignore template to the project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("Adding .gitignore template:", args[0])
		return gitignore.Add(args[0], overrideFlag)
	},
}

func init() {
	addCmd.Flags().BoolVarP(&overrideFlag, "override", "o", false, "Override existing .gitignore file")
}

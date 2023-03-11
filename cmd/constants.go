package cmd

import (
	"fmt"

	"github.com/densmoe/mx-inspector/model"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(constantsCmd)
	constantsCmd.AddCommand(constantsLsCmd)
}

var constantsCmd = &cobra.Command{
	Use:   "constants",
	Short: "constants",
	Long:  "constants",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var constantsLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "Retrieves list of constants",
	Long:  "Retrieves a list of all constants",
	Run: func(cmd *cobra.Command, args []string) {
		model := model.Load(args[0])
		for _, c := range model.Constants {
			fmt.Println(c.Name)
		}
	},
}

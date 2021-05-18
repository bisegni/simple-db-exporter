package cmd

import "github.com/spf13/cobra"

var csvCommand = &cobra.Command{
	Use:   "csv",
	Short: "Export database data into csv file type",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// Initialize the flag
func init() {
	rootCmd.AddCommand(csvCommand)
}

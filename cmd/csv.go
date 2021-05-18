package cmd

import (
	"fmt"
	"os"

	"github.com/bisegni/simple-db-exporter/service"
	"github.com/spf13/cobra"
)

var csvCommand = &cobra.Command{
	Use:   "csv",
	Short: "Export database data into csv file type",
	Long: `Usage:
	simple-db-exporter csv postgres://user:pass@host:port/database_name table-name
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var table string
		var uri string
		var columns []string
		var outputFilePath string
		if len(args) != 2 {
			fmt.Println("bad argument number")
			os.Exit(1)
		}

		switch len(args) {
		default:
		case 2:
			table = args[1]
		case 1:
			uri = args[0]
		}

		columns, _ = cmd.Flags().GetStringArray("column")
		outputFilePath, _ = cmd.Flags().GetString("destination-file")
		if err := service.RunCSVExport(uri, table, columns, outputFilePath); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

// Initialize the flag
func init() {
	csvCommand.Flags().StringP("destination-file", "d", "export.csv", "specify the name of the export file")
	csvCommand.Flags().StringArrayP("column", "c", []string{}, "Specify columns to be exporte")
	csvCommand.Flags().Int32P("max-row-num", "n", -1, "Specify numbero of row to be exported")
	rootCmd.AddCommand(csvCommand)
}

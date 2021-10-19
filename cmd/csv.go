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
		if len(args) != 2 {
			fmt.Println("bad argument number")
			os.Exit(1)
		}
		var table string = args[1]
		var uri string = args[0]
		var columns []string
		var outputFilePath string
		var rowLimit int

		columns, _ = cmd.Flags().GetStringSlice("column")
		outputFilePath, _ = cmd.Flags().GetString("destination-file")
		rowLimit, _ = cmd.Flags().GetInt("max-row-num")
		if err := service.RunCSVExport(uri, table, columns, rowLimit, outputFilePath); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var oracleExport = &cobra.Command{
	Use:   "oracle",
	Short: "Export oracle query as inser postgress sql file",
	Long: `Usage:
	simple-db-exporter oracle {USER}:{PASSWORD}@{HOSTNAME}/{SID} /sql/file/path table_name /dest/file/path 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 4 {
			fmt.Println("bad argument number")
			os.Exit(1)
		}
		var oracle_uri string = args[0]
		var oracleQueryFile string = args[1]
		var destTableName string = args[2]
		var destinationFile string = args[3]

		if err := service.ExportOracleQuery(
			oracle_uri,
			oracleQueryFile,
			destTableName,
			destinationFile); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var fromToExport = &cobra.Command{
	Use:   "export",
	Short: "Export many sql file to the correspective insert file in an output folder",
	Long: `Usage:
	simple-db-exporter export {USER}:{PASSWORD}@{HOSTNAME}/{SID} /sql/folder/path /dest/folder/path 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			fmt.Println("bad argument number")
			os.Exit(1)
		}
		var oracle_uri string = args[0]
		var inpuFolder string = args[1]
		var outputFolder string = args[2]
		err := service.ExportFromFolderToFolder(
			oracle_uri,
			inpuFolder,
			outputFolder,
		)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("done!")
	},
}

// Initialize the flag
func init() {
	csvCommand.Flags().StringP("destination-file", "d", "export.csv", "specify the name of the export file")
	csvCommand.Flags().StringSliceP("column", "c", []string{}, "Specify columns to be exporte")
	csvCommand.Flags().IntP("max-row-num", "n", -1, "Specify numbero of row to be exported")
	rootCmd.AddCommand(csvCommand)
	rootCmd.AddCommand(oracleExport)
	rootCmd.AddCommand(fromToExport)
}

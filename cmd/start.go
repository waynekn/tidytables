package cmd

import (
	"log"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/waynekn/tidytables/db"
	"github.com/waynekn/tidytables/logging"
	"github.com/waynekn/tidytables/tui"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Tidy Tables",
	Long: `The start command initiates Tidy Tables and establishes a connection to the database.
This command requires specific flags to provide the necessary database connection details`,
	Run: func(cmd *cobra.Command, args []string) {

		port := getFlagValue(cmd, "port")
		user := getFlagValue(cmd, "user")
		password := getFlagValue(cmd, "password")
		dbName := getFlagValue(cmd, "name")
		host := getFlagValue(cmd, "host")

		connection, err := db.ConnectToDb(host, port, user, password, dbName)

		if err != nil {
			log.SetFlags(0)
			log.Fatal(color.RedString(err.Error()))
		} else {
			log.SetFlags(0)
			log.Printf(color.GreenString("successfully connected to %v database"), dbName)

		}
		defer connection.Close()

		logFile, err := logging.OpenLogFile()

		if err != nil {
			log.Fatal(color.RedString(err.Error()))
		}

		defer logFile.Close()

		tui.StartTea()

	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().String("port", "5432", "port which the database is running on")

	startCmd.Flags().StringP("host", "H", "localhost", "the host of the database")

	startCmd.Flags().StringP("user", "U", "", "user to connect to the database as")
	startCmd.MarkFlagRequired("user")

	startCmd.Flags().StringP("password", "P", "", "database password")
	startCmd.MarkFlagRequired("password")

	startCmd.Flags().StringP("name", "N", "", "name of the database to connect to")
	startCmd.MarkFlagRequired("name")

	startCmd.Flags().String("db", "postgres", "name of the relational database management system")
	startCmd.MarkFlagRequired("db")
}

func getFlagValue(cmd *cobra.Command, flag string) string {
	value, err := cmd.Flags().GetString(flag)
	if err != nil {
		log.Fatal(color.RedString("An unexpected error occured while getting flags."))
	}
	return value
}

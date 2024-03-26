/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/klippa-app/db-time-machine/db"
	"github.com/klippa-app/db-time-machine/db/dialect"
	"github.com/klippa-app/db-time-machine/internal"
	"github.com/klippa-app/db-time-machine/internal/config"
	"github.com/klippa-app/db-time-machine/internal/hashes"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dbtm",
	Short: "Calculate the current database name, and instantiate and migrate it if necessary.",
	Long: `This command will calculate and return the name of the current
development database an application should connect to.

If necessary it will also instantiate the database by cloning
the previous database and applying all missing migrations.

The database names are calculated from the chained hashes of
each migration file, thus at any time it possible for a
database for each migration to exist. However dbtm will only
ever instantiate a new database for the current migration.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		path := cmd.Flags().Lookup("config").Value.String()

		ctx, err := config.Load(cmd.Context(), path)
		if err != nil {
			panic(err)
		}

		ctx, err = config.MergeFlags(ctx, cmd.Flags())
		if err != nil {
			panic(err)
		}

		ctx, err = hashes.Calculate(ctx)
		if err != nil {
			panic(err)
		}

		ctx = db.AttachContext(ctx, dialect.Postgres())

		cmd.SetContext(ctx)
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		name, err := internal.TimeTravel(ctx, nil)
		if err != nil {
			panic(err)
		}

		fmt.Println(name)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().BoolP("help", "", false, "help for dbtm")

	rootCmd.PersistentFlags().
		StringP("config", "c", "", "config file (default is $PWD/.dbtm.yaml)")

	rootCmd.PersistentFlags().
		StringP("uri", "u", "", "the connection uri for the db")
	rootCmd.PersistentFlags().
		StringP("database", "d", "", "the database name")

	rootCmd.PersistentFlags().
		String("migration-directory", "", "the directory containing the migrations")
	rootCmd.PersistentFlags().
		String("migration-format", "", "a regex for matching migration file names")
	rootCmd.PersistentFlags().
		String("migration-command", "", "the command to run to migrate the database")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

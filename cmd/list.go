/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/klippa-app/db-time-machine/db"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the databases currently in use by the provided config.",
	Long: `	This command lists all the databases based on the config that is provided.
			The list is constructed based off the prefix in the config.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		driver := db.FromContext(ctx)

		dbList, err := driver.List(ctx)
		if err != nil {
			panic(err)
		}

		for i := 0; i < len(dbList); i++ {
			fmt.Println(dbList[i])
		}

		return
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

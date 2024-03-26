/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/klippa-app/db-time-machine/db"
	"github.com/spf13/cobra"
)

var askBool bool

// pruneCmd represents the prune command
var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Prune unused databases",
	Long: `	Prune will look at the current state of the migration folder and compare that to the available databases.
			Any databases that it sees that are not also in the generated hash list will be deleted.
			
			To prevent any big accidental deletes of databases of other projects be sure to use a different prefix in all projects.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		force := cmd.Flags().Lookup("force").Value.String()
		forceBool, err := strconv.ParseBool(force)
		if err != nil {
			return err
		}

		if forceBool {
			askBool = true
			return nil
		}

		ctx := cmd.Context()
		driver := db.FromContext(ctx)

		dbList, err := driver.PruneList(ctx)
		if err != nil {
			return err
		}

		fmt.Println("The following databases will be deleted if you continue: ")
		for _, db := range dbList {
			fmt.Printf("- %s\n", db)
		}

		askBool = ask()

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		driver := db.FromContext(ctx)

		if askBool {
			err := driver.Prune(ctx)
			if err != nil {
				panic(err)
			}
		}

		return
	},
}

func ask() bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, "Do you want to continue? [Y/N]")
		s, _ := reader.ReadString('\n')
		s = strings.TrimSuffix(s, "\n")
		s = strings.ToLower(s)
		if len(s) > 1 {
			continue
		}
		if strings.Compare(s, "n") == 0 {
			return false
		} else if strings.Compare(s, "y") == 0 {
			break
		} else {
			continue
		}
	}
	return true
}

func init() {
	rootCmd.AddCommand(pruneCmd)

	pruneCmd.Flags().Bool("force", false, "use this to force delete unused databases.")
}

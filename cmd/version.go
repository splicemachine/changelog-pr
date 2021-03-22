package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Express the 'version' of splicectl.",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("\nSemVer: %s, BuildDate: %s\nCommitID: %s, GitRef: %s\n", semVer, buildDate, gitCommit, gitRef)

	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

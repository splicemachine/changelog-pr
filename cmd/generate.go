package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/blang/semver/v4"
	log "github.com/maahsome/changelog-pr/common"
	"github.com/maahsome/changelog-pr/provider"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a Changelog from PR descriptions since a specified TAG",
	Long: `EXAMPLE:
	#> changelog-pr generate --path <path/to/git/src> --since-tag v0.1.0 --release-tag v0.2.0
	`,
	Run: func(cmd *cobra.Command, args []string) {
		srcPath, _ := cmd.Flags().GetString("path")
		sinceTag, _ := cmd.Flags().GetString("since-tag")
		releaseTag, _ := cmd.Flags().GetString("release-tag")
		changelogFile, _ := cmd.Flags().GetString("file")

		if !strings.HasPrefix(sinceTag, "v") {
			log.Logger.Fatal("Please provide TAGs in format 'vMAJOR.MINOR.PATCH'")
		}
		if !strings.HasPrefix(releaseTag, "v") {
			log.Logger.Fatal("Please provide TAGs in format 'vMAJOR.MINOR.PATCH'")
		}
		_, sterr := semver.Parse(strings.Replace(sinceTag, "v", "", 1))
		if sterr != nil {
			log.Logger.Fatal(fmt.Sprintf("Error parsing SemVer for %s", sinceTag))
		}
		_, rterr := semver.Parse(strings.Replace(releaseTag, "v", "", 1))
		if rterr != nil {
			log.Logger.Fatal(fmt.Sprintf("Error parsing SemVer for %s", releaseTag))
		}

		glog, err := generateLog(srcPath, sinceTag, releaseTag, changelogFile)
		if err != nil {
			log.Logger.WithError(err).Error("Error generating the changelog")
		}

		fmt.Println(glog)

	},
}

func generateLog(src string, sTag string, rTag string, logFile string) (string, error) {

	var (
		err   error
		chlog string
		gp    provider.Provider
	)

	switch strings.ToLower(gitProvider) {
	case "github":
		gp, err = provider.GetProvider(provider.GITHUB)
		if err != nil {
			return "", errors.New("failed to provision git provider")
		}
	default:
		return "", errors.New("unsupported provider")
	}

	auth := provider.AuthToken{
		GithubToken: ghToken,
	}
	chlog, err = gp.GetChangeLogFromPR(src, sTag, rTag, auth, logFile)
	if err != nil {
		return "", errors.New("failed generation of changelog")
	}

	return chlog, nil
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringP("path", "p", "", "Specify the path to the git source directory")
	generateCmd.Flags().StringP("since-tag", "t", "", "Specify the git TAG to go back to and process PR descriptions")
	generateCmd.Flags().StringP("release-tag", "r", "", "Specify the new release TAG")
	generateCmd.Flags().StringP("file", "f", "", "Specify an output file to save the changelog to")
	generateCmd.MarkFlagRequired("path")
	generateCmd.MarkFlagRequired("since-tag")
	generateCmd.MarkFlagRequired("release-tag")
}

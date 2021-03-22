package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Output the markdown that needs to be added to the PR template",
	Long:  `This command outputs the markdown that is added to the '.github/PULL_REQUEST_TEMPLATE.md' file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`## Changelog Inclusions

<!-- Text Entered in these sections will appear as it is written, MD formatted -->
<!-- Avoid using the # character as the first character on any line, a trim is -->
<!-- performed on each line when checking for markdown section tags -->
- base feature note
  - **BREAKING** note on base feature
  - Basically whatever formatting we have here, just plain-text
	%> command example
- next feature
  - note on next feature
<!-- If there is NO text in a section, no entries will be collected for that section -->

### Additions

### Changes

### Fixes

### Deprecated

### Removed

### Breaking Changes`)
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)
}

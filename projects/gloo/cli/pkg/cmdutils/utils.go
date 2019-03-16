package cmdutils

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// This finds 'current' command in 'parent' command and adds the 'sub' command to it.
func MustAddChildCommand(parent *cobra.Command, current *cobra.Command, sub *cobra.Command) {
	parentCmds := parent.Commands()
	for _, old := range parentCmds {
		if old.Use == current.Use {
			old.AddCommand(sub)
			return
		}
	}
	err := fmt.Errorf("did not find child command to replace")
	log.Fatal(err)
}

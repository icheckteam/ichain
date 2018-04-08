package commands

import (
	"fmt"
	"os"

	"github.com/icheckteam/ichain"
	"github.com/spf13/cobra"
)

var (
	// VersionCmd prints the program's version to stderr and exits.
	VersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print ichain's version",
		Run:   doVersionCmd,
	}
)

func doVersionCmd(cmd *cobra.Command, args []string) {
	v := ichain.Version
	if len(v) == 0 {
		fmt.Fprintln(os.Stderr, "unset")
		return
	}
	fmt.Fprintln(os.Stderr, v)
}

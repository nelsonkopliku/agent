package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/trento-project/agent/version"
)

func NewVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Trento",
		Long:  `All software has versions. This is Trento's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Trento %s version %s\nbuilt with %s %s/%s\n", version.Flavor, version.Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		},
	}

	return versionCmd
}
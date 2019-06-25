package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dsbrng25b/cis/internal/virt"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

const envPrefix = "CIS_"

var (
	version        = "n/a"
	gitCommit      = "n/a"
	buildTime      = "n/a"
	virStoragePool string
	virConnectURL  string
	mgr            *virt.LibvirtManager
)

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "cis",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error
			// flags from environment
			cmd.Flags().VisitAll(func(f *flag.Flag) {
				varName := envPrefix + strings.ToUpper(f.Name)
				if val, ok := os.LookupEnv(varName); !f.Changed && ok {
					f.Value.Set(val)
				}
			})

			// initialize libvirt manager
			mgr, err = virt.NewLibvirtManager(virStoragePool, virConnectURL)
			if err != nil {
				errExit(err)
			}

		},
		BashCompletionFunction: bash_completion_func,
	}
	rootCmd.AddCommand(
		newConfigCmd(),
		newImageCmd(),
		newCreateCmd(),
		newListCmd(),
		newStartCmd(),
		newShutdownCmd(),
		newRemoveCmd(),
		newVersionCmd(),
		newCompletionCmd(),
	)
	return rootCmd
}

func addPoolOption(fs *flag.FlagSet, pool *string) {
	fs.StringVar(pool, "pool", "default", "The storage pool to use")
}

func newCompletionCmd() *cobra.Command {
	var shell string
	cmd := &cobra.Command{
		Use:       "completion <shell>",
		ValidArgs: []string{"bash", "zsh"},
		Args:      cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			shell = args[0]
			switch shell {
			case "bash":
				newRootCmd().GenBashCompletion(os.Stdout)
			case "zsh":
				newRootCmd().GenZshCompletion(os.Stdout)
			default:
				errExit("unknown shell")
			}
		},
	}
	return cmd
}

func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "show version and build information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("version:", version)
			fmt.Println("commit:", gitCommit)
			fmt.Println("build time:", buildTime)
		},
	}
	return cmd
}

func main() {
	_ = newRootCmd().Execute()
}

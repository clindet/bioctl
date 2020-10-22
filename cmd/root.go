package cmd

import (
	"fmt"
	"os"
	"strings"

	arc "github.com/openbiox/ligo/archive"
	"github.com/spf13/cobra"
)

// RootClisT is the bioctl global flags
type RootClisT struct {
	// version of bioctl
	Version string
	Verbose int
	SaveLog bool
	TaskID  string
	LogDir  string
	Clean   bool
	Out     string
	Thread  int

	Uncompress string
	HelpFlags  bool
}

var rootClis = RootClisT{
	Version:   version,
	Verbose:   1,
	HelpFlags: true,
}

var rootCmd = &cobra.Command{
	Use:   "bioctl",
	Short: "A simple command line tool to facilitate the data analysis",
	Long:  `A simple command line tool to facilitate the data analysis. More see here https://github.com/clindet/bioctl.`,
	Run: func(cmd *cobra.Command, args []string) {
		if rootClis.Clean || rootClis.Uncompress != "" {
			initCmd(cmd, args)
		}
		uncompress(cmd, args)
		if rootClis.HelpFlags {
			cmd.Help()
		}
	},
}

// Execute main interface of bget
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if !rootCmd.HasFlags() && !rootCmd.HasSubCommands() {
			rootCmd.Help()
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func init() {
	wd, _ = os.Getwd()
	rootCmd.Version = version
	rootCmd.Flags().StringVarP(&(rootClis.Uncompress), "uncompress", "u", "", "uncompress files.")
	rootCmd.AddCommand(FmtCmd)
	rootCmd.AddCommand(ParCmd)
	rootCmd.AddCommand(PlotCmd)
	rootCmd.AddCommand(RangeCmd)
	rootCmd.AddCommand(RandCmd)
	rootCmd.AddCommand(StatCmd)
	rootCmd.AddCommand(StatDfCmd)
	rootCmd.AddCommand(StatFnCmd)
	rootCmd.AddCommand(ConvertCmd)
	setGlobalFlag(rootCmd)
}

func uncompress(cmd *cobra.Command, args []string) {
	if rootClis.Uncompress != "" {
		files := strings.Split(rootClis.Uncompress, " ")
		for i := range files {
			out := rootClis.Out
			if out == "" {
				out = wd
			}
			err := arc.UnarchiveLog(files[i], out)
			if err != nil {
				log.Warnf("%v", err)
			}
		}
		rootClis.HelpFlags = false
	}
}

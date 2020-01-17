package cmd

import (
	"bufio"
	"io/ioutil"
	"os"

	cvrt "github.com/openbiox/ligo/convert"
	"github.com/openbiox/ligo/flag"
	"github.com/spf13/cobra"
)

// RangeClisT is the type to run bioctl Range
type RangeClisT struct {
}

// RangeClis is the parameters to run par.Tasks
var RangeClis RangeClisT

// RangeCmd is the command line of bioctl Range
var RangeCmd = &cobra.Command{
	Use:   "range",
	Short: "Functions to manipulate intervals.",
	Long:  `Functions to manipulate intervals.`,
	Run: func(cmd *cobra.Command, args []string) {
		RangeCmdRunOptions(cmd, args)
	},
}

func RangeCmdRunOptions(cmd *cobra.Command, args []string) {
	cleanArgs := []string{}
	hasStdin := false
	if cleanArgs, hasStdin = flag.CheckStdInFlag(cmd); hasStdin {
		reader := bufio.NewReader(os.Stdin)
		result, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Fatal(err)
		} else if len(result) > 0 {
		}
	}
	if len(cleanArgs) >= 1 {
		initCmd(cmd, args)
		if rootClis.Verbose == 2 {
			logEnv.Infof("RangeClis: %v", cvrt.Struct2Map(RangeClis))
		}

		rootClis.HelpFlags = false
	}
	if rootClis.HelpFlags {
		cmd.Help()
	}
}

func init() {
	RangeCmd.Example = ``
}

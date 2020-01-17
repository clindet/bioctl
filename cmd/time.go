package cmd

import (
	"bufio"
	"io/ioutil"
	"os"

	cvrt "github.com/openbiox/ligo/convert"
	"github.com/openbiox/ligo/flag"
	"github.com/spf13/cobra"
)

// TimeClisT is the type to run bioctl Time
type TimeClisT struct {
}

// TimeClis is the parameters to run par.Tasks
var TimeClis TimeClisT

// TimeCmd is the command line of bioctl Time
var TimeCmd = &cobra.Command{
	Use:   "time",
	Short: "Functions to manipulate time data.",
	Long:  `Functions to manipulate time data.`,
	Run: func(cmd *cobra.Command, args []string) {
		TimeCmdRunOptions(cmd, args)
	},
}

func TimeCmdRunOptions(cmd *cobra.Command, args []string) {
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
			logEnv.Infof("TimeClis: %v", cvrt.Struct2Map(TimeClis))
		}

		rootClis.HelpFlags = false
	}
	if rootClis.HelpFlags {
		cmd.Help()
	}
}

func init() {
	TimeCmd.Example = ``
}

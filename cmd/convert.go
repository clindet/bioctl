package cmd

import (
	"bufio"
	"io/ioutil"
	"os"

	cvrt "github.com/openbiox/ligo/convert"
	"github.com/openbiox/ligo/flag"
	"github.com/spf13/cobra"
)

// ConvertClisT is the type to run bioctl Convert
type ConvertClisT struct {
}

// ConvertClis is the parameters to run par.Tasks
var ConvertClis ConvertClisT

// ConvertCmd is the command line of bioctl Convert
var ConvertCmd = &cobra.Command{
	Use:   "cvrt",
	Short: "Convert related functions.",
	Long:  `Convert related functions.`,
	Run: func(cmd *cobra.Command, args []string) {
		ConvertCmdRunOptions(cmd, args)
	},
}

func ConvertCmdRunOptions(cmd *cobra.Command, args []string) {
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
			logEnv.Infof("ConvertClis: %v", cvrt.Struct2Map(ConvertClis))
		}

		rootClis.HelpFlags = false
	}
	if rootClis.HelpFlags {
		cmd.Help()
	}
}

func init() {
	ConvertCmd.Example = ``
}

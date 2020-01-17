package cmd

import (
	"bufio"
	"io/ioutil"
	"os"

	cvrt "github.com/openbiox/ligo/convert"
	"github.com/openbiox/ligo/flag"
	"github.com/spf13/cobra"
)

// ExtractClisT is the type to run bioctl Extract
type ExtractClisT struct {
}

// ExtractClis is the parameters to run par.Tasks
var ExtractClis ExtractClisT

// ExtractCmd is the command line of bioctl Extract
var ExtractCmd = &cobra.Command{
	Use:   "extr",
	Short: "Extract information from input stream.",
	Long:  `Extract information from input stream.`,
	Run: func(cmd *cobra.Command, args []string) {
		ExtractCmdRunOptions(cmd, args)
	},
}

func ExtractCmdRunOptions(cmd *cobra.Command, args []string) {
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
			logEnv.Infof("ExtractClis: %v", cvrt.Struct2Map(ExtractClis))
		}

		rootClis.HelpFlags = false
	}
	if rootClis.HelpFlags {
		cmd.Help()
	}
}

func init() {
	ExtractCmd.Example = ``
}

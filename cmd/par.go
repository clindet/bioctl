package cmd

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"

	cvrt "github.com/openbiox/ligo/convert"
	"github.com/openbiox/ligo/flag"
	"github.com/openbiox/ligo/par"
	"github.com/spf13/cobra"
)

// ParClis is the parameters to run par.Tasks
var ParClis par.ClisT

// ParCmd is the command line of bioctl par
var ParCmd = &cobra.Command{
	Use:   "par",
	Short: "Run parallel tasks.",
	Long:  `Run parallel tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		ParClis.Verbose = rootClis.Verbose
		ParClis.SaveLog = rootClis.SaveLog
		ParClis.TaskID = rootClis.TaskID
		ParClis.LogDir = rootClis.LogDir
		parCmdRunOptions(cmd, args)
	},
}

func parCmdRunOptions(cmd *cobra.Command, args []string) {
	cleanArgs := []string{}
	hasStdin := false
	if cleanArgs, hasStdin = flag.CheckStdInFlag(cmd); hasStdin {
		reader := bufio.NewReader(os.Stdin)
		result, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Fatal(err)
		} else if len(result) > 0 {
			ParClis.Script = string(result)
		}
	}
	if len(cleanArgs) >= 1 || hasStdin || ParClis.Script != "" {
		initCmd(cmd, args)
		if rootClis.Verbose == 2 {
			logEnv.Infof("ParClis: %v", cvrt.Struct2Map(ParClis))
		}
		par.Tasks(&ParClis)
		rootClis.HelpFlags = false
	}
	if rootClis.HelpFlags {
		cmd.Help()
	}
}

func init() {
	ParCmd.Flags().StringVarP(&ParClis.Index, "index", "", "", "task index (e.g. 1,2,5-10).")
	ParCmd.Flags().StringVarP(&ParClis.ForceAddIdx, "force-idx", "", "true", "force to add {{index}} at the end of --cmd.")
	ParCmd.Flags().StringVarP(&ParClis.Env, "env", "", "", "environment (key1:value1,key2:value2).")
	ParCmd.Flags().StringVarP(&ParClis.Script, "cmd", "", "", "command string.")
	ParCmd.Flags().IntVarP(&ParClis.Thread, "thread", "t", 1, "thread to process.")
	ParClis.ForceAddIdx = strings.ToLower(ParClis.ForceAddIdx)
	ParCmd.Example = `  # concurent 2 tasks with total 8 tasks
  echo 'touch /tmp/$1 /tmp/$2; sleep ${1}' > job.sh && bioctl par --cmd "sh job.sh" -t 2 --index 1,2,5-10
  # concurent 4 tasks with total 8 tasks and env parse
  bioctl par --cmd 'sh job.sh {{index}} {{key2}}' -t 4 --index 1,2,5-10 --env "key2:123"
	
  # concurent 4 tasks with total 8 tasks (direct) and env parse (more log)
  bioctl par --cmd 'echo {{index}} {{key2}}; sleep {{index}}' -t 4 --index 1,2,5-10 --env "key2:123" --verbose 2 --save-log
	
  # concurent 4 tasks with total 8 tasks, env
  # and not to force add {{index}} at the end of cmd
  bioctl par --cmd 'echo {{key2}}; sleep {{index}}' -t 4 --index 1,2,5-10 --env "key2:123" --force-idx false --save-log
	
  # pipe usage
  echo 'sh job.sh' | bioctl par -t 4 --index 1,2,5-10 -`
}

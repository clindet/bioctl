package cmd

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/openbiox/bioctl/flag"
	clog "github.com/openbiox/bioctl/log"
	"github.com/openbiox/bioctl/par"
	stringo "github.com/openbiox/bioctl/stringo"
	"github.com/spf13/cobra"
)

// ParClis is the parameters to run par.Tasks
var ParClis par.ClisT

// ParCmd is the command line of bioctl par
var ParCmd = &cobra.Command{
	Use:   "par",
	Short: "A simple parallel task manager.",
	Long:  `A simple parallel task manager.`,
	Run: func(cmd *cobra.Command, args []string) {
		ParClis.Quiet = rootClis.Quiet
		clog.SetQuietLog(log, rootClis.Quiet)
		parCmdRunOptions(cmd)
	},
}

func parCmdRunOptions(cmd *cobra.Command) {
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
		par.Tasks(&ParClis)
		rootClis.HelpFlags = false
	}
	if rootClis.HelpFlags {
		cmd.Help()
	}
}

func init() {
	wd, _ = os.Getwd()
	ParCmd.Flags().StringVarP(&ParClis.Index, "index", "", "", "Task index (e.g. 1,2,5-10).")
	ParCmd.Flags().StringVarP(&ParClis.ForceAddIdx, "force-idx", "", "true", "Force to add {{index}} at the end of --cmd.")
	ParCmd.Flags().StringVarP(&ParClis.Env, "env", "", "", "Env (key1:value1,key2:value2).")
	ParCmd.Flags().StringVarP(&ParClis.Script, "cmd", "", "", "Command string.")
	ParCmd.Flags().IntVarP(&ParClis.Thread, "thread", "t", 1, "Thread to process.")
	ParCmd.Flags().StringVarP(&(ParClis.TaskID), "task-id", "", stringo.RandString(15), "Task ID (random).")
	ParCmd.Flags().StringVarP(&(ParClis.LogDir), "log-dir", "", path.Join(wd, "_log"), "Log dir.")
	ParClis.ForceAddIdx = strings.ToLower(ParClis.ForceAddIdx)
	ParCmd.Example = `  # concurent 2 tasks with total 8 tasks
  echo 'echo $1 $2; sleep ${1}' > job.sh && bioctl par --cmd "sh job.sh" -t 2 --index 1,2,5-10
  # concurent 4 tasks with total 8 tasks and env parse
  bioctl par --cmd 'sh job.sh {{index}} {{key2}}' -t 4 --index 1,2,5-10 --env "key2:123"
	
  # concurent 4 tasks with total 8 tasks (direct) and env
  bioctl par --cmd 'echo {{index}} {{key2}}; sleep {{index}}' -t 4 --index 1,2,5-10 --env "key2:123"
	
  # concurent 4 tasks with total 8 tasks, env
  # and not to force add {{index}} at the end of cmd
  bioctl par --cmd 'echo {{key2}}; sleep {{index}}' -t 4 --index 1,2,5-10 --env "key2:123" --force-idx false
	
  # pipe usage
  echo 'sh job.sh' | bioctl par -t 4 --index 1,2,5-10 -`
}

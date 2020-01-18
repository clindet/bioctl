package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"

	cvrt "github.com/openbiox/ligo/convert"
	"github.com/openbiox/ligo/flag"
	"github.com/openbiox/ligo/stringo"
	"github.com/spf13/cobra"
)

// MathClisT is the type to run bioctl Math
type MathClisT struct {
	Max    bool
	Min    bool
	Median bool
	Mean   bool
}

// MathClis is the parameters to run par.Tasks
var MathClis MathClisT

// MathCmd is the command line of bioctl Math
var MathCmd = &cobra.Command{
	Use:   "math",
	Short: "Functions related to math.",
	Long:  `Functions related to math.`,
	Run: func(cmd *cobra.Command, args []string) {
		MathCmdRunOptions(cmd, args)
	},
}

func MathCmdRunOptions(cmd *cobra.Command, args []string) {
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
	var seq []float64
	for _, v := range args {
		i64, _ := strconv.ParseFloat(v, 64)
		seq = append(seq, i64)
	}
	if len(cleanArgs) >= 1 || MathClis.Max || MathClis.Min || MathClis.Mean || MathClis.Median {
		initCmd(cmd, args)
		if rootClis.Verbose == 2 {
			logEnv.Infof("MathClis: %v", cvrt.Struct2Map(MathClis))
		}
		if MathClis.Max {
			fmt.Println(max(seq))
		} else if MathClis.Min {
			fmt.Println(min(seq))
		} else if MathClis.Mean {
			fmt.Println(mean(seq))
		} else if MathClis.Median {
			fmt.Println(stringo.StrReplaceAll(fmt.Sprintf("%v", median(seq)), "[[]|[]]", ""))
		}
		rootClis.HelpFlags = false
	}
	if rootClis.HelpFlags {
		cmd.Help()
	}
}

func max(seq []float64) float64 {
	val := seq[0]
	for i := range seq {
		if seq[i] >= val {
			val = seq[i]
		}
	}
	return val
}

func min(seq []float64) float64 {
	val := seq[0]
	for i := range seq {
		if seq[i] <= val {
			val = seq[i]
		}
	}
	return val
}

func mean(seq []float64) float64 {
	var total float64
	for i := range seq {
		total += seq[i]
	}
	return total / float64(len(seq))
}

func median(seq []float64) []float64 {
	sort.Sort(sort.Float64Slice(seq))
	idx := 0
	if len(seq)%2 == 0 {
		idx = len(seq) / 2
		return []float64{seq[idx], seq[idx+1]}
	} else {
		return []float64{seq[idx+1]}
	}
}

func init() {
	MathCmd.Flags().BoolVarP(&(MathClis.Max), "max", "", false, "returns max value.")
	MathCmd.Flags().BoolVarP(&(MathClis.Min), "min", "", false, "returns min value.")
	MathCmd.Flags().BoolVarP(&(MathClis.Median), "median", "", false, "returns median value.")
	MathCmd.Flags().BoolVarP(&(MathClis.Mean), "mean", "", false, "returns mean value.")

	MathCmd.Example = `  bioctl math --min 1 2 3 4 5 100
  bioctl math --max 1 2 3 4 5
  bioctl math --mean 1 3 5 7 9 26 100
  bioctl math --median 1 10 3 4 100 143 123 12 22.2`
}

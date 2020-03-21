package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/openbiox/ligo/stringo"

	"github.com/montanaflynn/stats"
	cvrt "github.com/openbiox/ligo/convert"
	"github.com/openbiox/ligo/flag"
	"github.com/spf13/cobra"
)

// StatClisT is the type to run bioctl Stat
type StatClisT struct {
	Max        bool
	Min        bool
	Median     bool
	Mean       bool
	Mfreq      bool
	Variance   bool
	Sum        bool
	Pearson    bool
	Percentile float64
	Freq       bool
}

// StatClis is the parameters to run par.Tasks
var StatClis StatClisT

// StatCmd is the command line of bioctl Stat
var StatCmd = &cobra.Command{
	Use:   "stat",
	Short: "Functions related to stat.",
	Long:  `Functions related to stat.`,
	Run: func(cmd *cobra.Command, args []string) {
		StatCmdRunOptions(cmd, args)
	},
}

func StatCmdRunOptions(cmd *cobra.Command, args []string) {
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
	if StatClis.Freq {
		statFreq(args)
		return
	}
	var data, data2 stats.Float64Data
	if !StatClis.Pearson {
		data = stats.LoadRawData(args)
	} else {
		data = stats.LoadRawData(args[0:(len(args) / 2)])
		data2 = stats.LoadRawData(args[len(args)/2 : len(args)])
		fmt.Printf("Vector1: %v\n", data)
		fmt.Printf("Vector2: %v\n", data2)
	}
	if len(cleanArgs) >= 1 || StatClis.Max || StatClis.Min || StatClis.Mean || StatClis.Median {
		initCmd(cmd, args)
		if rootClis.Verbose == 2 {
			logEnv.Infof("StatClis: %v", cvrt.Struct2Map(StatClis))
		}
		var val float64
		var valSlice []float64
		if StatClis.Max {
			val, _ = data.Max()
		} else if StatClis.Min {
			val, _ = data.Min()
		} else if StatClis.Mean {
			val, _ = data.Mean()
		} else if StatClis.Median {
			val, _ = data.Median()
		} else if StatClis.Mfreq {
			valSlice, _ = data.Mode()
		} else if StatClis.Variance {
			val, _ = data.Variance()
		} else if StatClis.Sum {
			val, _ = data.Sum()
		} else if StatClis.Pearson {
			val, _ = stats.Pearson(data, data2)
		} else if StatClis.Percentile != -1 {
			val, _ = stats.Percentile(data, StatClis.Percentile)
		}
		if !StatClis.Mfreq {
			fmt.Println(val)
		} else {
			fmt.Println(stringo.StrReplaceAll(fmt.Sprintf("%v", valSlice), "[[]|[]]", ""))
		}
		rootClis.HelpFlags = false
	}
	if rootClis.HelpFlags {
		cmd.Help()
	}
}

func statFreq(args []string) {
	count := make(map[string]int)
	for _, v := range args {
		count[v]++
	}
	for k, v := range count {
		fmt.Printf("%s: %d\n", k, v)
	}
	return
}

func init() {
	StatCmd.Flags().BoolVarP(&(StatClis.Max), "max", "", false, "returns max value.")
	StatCmd.Flags().BoolVarP(&(StatClis.Min), "min", "", false, "returns min value.")
	StatCmd.Flags().BoolVarP(&(StatClis.Median), "median", "", false, "returns median value.")
	StatCmd.Flags().BoolVarP(&(StatClis.Mean), "mean", "", false, "returns mean value.")
	StatCmd.Flags().BoolVarP(&(StatClis.Mfreq), "mfreq", "", false, "returns the most frequent value.")
	StatCmd.Flags().BoolVarP(&(StatClis.Variance), "var", "", false, "returns the amount of variation in the dataset.")
	StatCmd.Flags().BoolVarP(&(StatClis.Sum), "sum", "", false, "returns the sum of value.")
	StatCmd.Flags().BoolVarP(&(StatClis.Pearson), "pearson", "", false, "returns the pearson product-moment correlation coefficient between two group variables.")
	StatCmd.Flags().Float64VarP(&(StatClis.Percentile), "percentile", "", -1, "returns the relative standing in a slice of floats.")
	StatCmd.Flags().BoolVarP(&(StatClis.Freq), "freq", "", false, "returns the frequency stat.")

	StatCmd.Example = `  bioctl stat --min 1 2 3 4 5 100
  bioctl stat --max 1 2 3 4 5
  bioctl stat --mean 1 3 5 7 9 26 100
  bioctl stat --median 1 10 3 4 100 143 123 12 22.2
  bioctl stat --mfreq 2 10 2 2 10 10 11 12 14
  bioctl stat --var 2 10 2 2 10
  bioctl stat --sum 1 1 2 3 2 14234 12 12 1331 23 12 12
  // vector1: 1-6 vector2: 10-51
  bioctl stat --pearson 1 2 3 4 5 6 10 20 30 40 50 51
  bioctl stat --percentile 30 1 2 3 4 5 6 7 8 9 10
  bioctl stat --freq 1 2 3 2 1 1 1 1 1 1 | sort`
}

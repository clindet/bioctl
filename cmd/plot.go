package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	cvrt "github.com/openbiox/ligo/convert"
	"github.com/openbiox/ligo/flag"
	"github.com/openbiox/ligo/plot"
	"github.com/spf13/cobra"
)

// PlotClisT is the type to run bioctl plot
type PlotClisT struct {
	ThemeName      string
	ShowThemes     bool
	ShowThemesName bool
}

// PlotClis is the parameters to run par.Tasks
var PlotClis PlotClisT

// PlotCmd is the command line of bioctl plot
var PlotCmd = &cobra.Command{
	Use:   "plot",
	Short: "Plots related functions.",
	Long:  `Plots related functions.`,
	Run: func(cmd *cobra.Command, args []string) {
		plotCmdRunOptions(cmd, args)
	},
}

func plotCmdRunOptions(cmd *cobra.Command, args []string) {
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
	if len(cleanArgs) >= 1 || PlotClis.ThemeName != "" || PlotClis.ShowThemes || PlotClis.ShowThemesName {
		initCmd(cmd, args)
		if rootClis.Verbose == 2 {
			logEnv.Infof("PlotClis: %v", cvrt.Struct2Map(PlotClis))
		}
		if PlotClis.ThemeName != "" {
			fmt.Printf("%s\n", strings.Join(plot.GetThemeColors(PlotClis.ThemeName).Colors, "\n"))
		} else if PlotClis.ShowThemes {
			dta, _ := json.MarshalIndent(plot.ThemeColors, "", "  ")
			fmt.Println(string(dta))
		} else if PlotClis.ShowThemesName {
			names := []string{}
			for _, v := range plot.ThemeColors {
				names = append(names, v.Name)
			}
			fmt.Printf("%s\n", strings.Join(names, "\n"))
		}
		rootClis.HelpFlags = false
	}
	if rootClis.HelpFlags {
		cmd.Help()
	}
}

func init() {
	PlotCmd.Flags().StringVarP(&PlotClis.ThemeName, "theme", "", "", "returns a color theme (default returns all).")
	PlotCmd.Flags().BoolVarP(&PlotClis.ShowThemes, "show-themes", "", false, "returns all theme.")
	PlotCmd.Flags().BoolVarP(&PlotClis.ShowThemesName, "show-themes-name", "", false, "returns all theme names.")

	PlotCmd.Example = ``
}

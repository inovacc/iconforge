package cmd

import (

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "iconforge",
	Short: "Cross-platform application icon generator",
	Long: `IconForge - Icons forged for every platform.

Converts SVG to production-ready icons (ICO, ICNS, PNG) and embeds them
into Go binaries, macOS .app bundles, Linux .desktop entries, and
Tauri/Electron projects.`,

}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {


	rootCmd.Version = GetVersionJSON()
	rootCmd.CompletionOptions.DisableDefaultCmd = true


	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

}


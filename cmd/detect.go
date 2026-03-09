package cmd

import (
	"fmt"
	"os"

	"github.com/inovacc/iconforge/internal/detect"
	"github.com/spf13/cobra"
)

var detectDir string

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect framework and show required icon assets",
	Long: `Scan a project directory for known frameworks (Tauri, Electron, Wails, Fyne)
and display the required icon assets for that framework.

Examples:
  iconforge detect
  iconforge detect --dir /path/to/project`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		dir := detectDir
		if dir == "" {
			var err error
			dir, err = os.Getwd()
			if err != nil {
				return fmt.Errorf("get working directory: %w", err)
			}
		}

		fw := detect.DetectFramework(dir)

		if fw == detect.FrameworkNone {
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "No framework detected.")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "\nSupported frameworks:")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "  - Tauri   (tauri.conf.json or src-tauri/)")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "  - Electron (electron in package.json)")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "  - Wails   (wails.json or build/appicon.png)")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "  - Fyne    (fyne.io/fyne in go.mod)")
			return nil
		}

		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Detected framework: %s\n\n", fw)

		switch fw {
		case detect.FrameworkTauri:
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Required assets (src-tauri/icons/):")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "  - icon.ico     (Windows)")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "  - icon.icns    (macOS)")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "  - 32x32.png")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "  - 128x128.png")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "  - 128x128@2x.png (256x256)")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "  - icon.png     (512x512)")
		case detect.FrameworkElectron:
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Required assets (build/):")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "  - icon.ico     (Windows)")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "  - icon.icns    (macOS)")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "  - icon.png     (Linux, 512x512+)")
		case detect.FrameworkWails:
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Required assets (build/):")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "  - appicon.png           (source)")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "  - windows/icon.ico      (Windows)")
		case detect.FrameworkFyne:
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Required assets (project root):")
			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "  - Icon.png     (512x512 or 1024x1024)")
		}

		_, _ = fmt.Fprintln(cmd.OutOrStdout(), "\nRun 'iconforge forge --auto-detect' to generate these automatically.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(detectCmd)
	detectCmd.Flags().StringVar(&detectDir, "dir", "", "Project directory to scan (default: current directory)")
}

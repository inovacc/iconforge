package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/inovacc/iconforge/internal/platform"
	"github.com/spf13/cobra"
)

var (
	embedICOPath    string
	embedOutputPath string
	embedArch       string
	embedMethod     string
	embedViPath     string
	embedAppName    string
	embedVersion    string
	embedCompany    string
	embedCopyright  string
)

var embedCmd = &cobra.Command{
	Use:   "embed",
	Short: "Generate .syso for Go build (Windows)",
	Long: `Generate a .syso resource file from an ICO file for Windows.

The .syso file is automatically linked by the Go linker during build,
embedding the icon and version info into the resulting .exe.

Methods:
  auto    - Use winres (pure Go) with full version info (default)
  simple  - Embed icon only, no version info

Examples:
  iconforge embed --ico build/icons/windows/icon.ico
  iconforge embed --ico icon.ico --name myapp --version 1.0.0
  iconforge embed --ico icon.ico --method simple --arch arm64`,
	RunE: runEmbed,
}

func init() {
	rootCmd.AddCommand(embedCmd)

	embedCmd.Flags().StringVar(&embedICOPath, "ico", "", "Path to ICO file")
	embedCmd.Flags().StringVarP(&embedOutputPath, "output", "o", "resource.syso", "Output .syso file path")
	embedCmd.Flags().StringVar(&embedArch, "arch", "amd64", "Target architecture (amd64, 386, arm, arm64)")
	embedCmd.Flags().StringVar(&embedMethod, "method", "auto", "Method to use: auto, simple")
	embedCmd.Flags().StringVar(&embedViPath, "vi", "", "Path to versioninfo.json (optional)")
	embedCmd.Flags().StringVar(&embedAppName, "name", "app", "Application name for version info")
	embedCmd.Flags().StringVar(&embedVersion, "version", "1.0.0", "Application version")
	embedCmd.Flags().StringVar(&embedCompany, "company", "", "Company name")
	embedCmd.Flags().StringVar(&embedCopyright, "copyright", "", "Copyright notice")

	// Keep deprecated flags as hidden aliases for backward compat
	embedCmd.Flags().String("tool", "", "Deprecated: use --method instead")
	_ = embedCmd.Flags().MarkHidden("tool")

	_ = embedCmd.MarkFlagRequired("ico")
}

func runEmbed(cmd *cobra.Command, _ []string) error {
	method := embedMethod

	if tool, _ := cmd.Flags().GetString("tool"); tool != "" {
		method = tool
	}

	// Map old method names to new ones for backward compat
	switch method {
	case "auto", "goversioninfo":
		return embedWinres(cmd)
	case "simple", "rsrc":
		return embedSimple(cmd)
	default:
		return fmt.Errorf("unsupported method: %s (use 'auto' or 'simple')", method)
	}
}

func embedWinres(cmd *cobra.Command) error {
	outputDir := filepath.Dir(embedOutputPath)
	cfg := platform.WindowsConfig{
		AppName:     embedAppName,
		Description: embedAppName,
		Version:     embedVersion,
		Company:     embedCompany,
		Copyright:   embedCopyright,
		ICOPath:     embedICOPath,
		OutputDir:   outputDir,
	}

	sysoPath, err := platform.GenerateSysoWinres(cfg, embedArch)
	if err != nil {
		return fmt.Errorf("winres: %w", err)
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Generated .syso (winres): %s\n", sysoPath)
	return nil
}

func embedSimple(cmd *cobra.Command) error {
	if err := platform.GenerateSysoFromICO(embedICOPath, embedOutputPath, embedArch); err != nil {
		return fmt.Errorf("winres: %w", err)
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Generated .syso (winres/simple): %s\n", embedOutputPath)
	return nil
}

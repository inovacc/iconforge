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
  auto           - Try goversioninfo (pure Go) first, fall back to rsrc (default)
  goversioninfo  - Pure Go, full version info + icon (no external tool needed)
  rsrc           - External tool, simple icon embedding (github.com/akavel/rsrc)

Examples:
  iconforge embed --ico build/icons/windows/icon.ico
  iconforge embed --ico icon.ico --method goversioninfo --name myapp --version 1.0.0
  iconforge embed --ico icon.ico --method rsrc --arch arm64
  iconforge embed --ico icon.ico --vi versioninfo.json`,
	RunE: runEmbed,
}

func init() {
	rootCmd.AddCommand(embedCmd)

	embedCmd.Flags().StringVar(&embedICOPath, "ico", "", "Path to ICO file")
	embedCmd.Flags().StringVarP(&embedOutputPath, "output", "o", "resource.syso", "Output .syso file path")
	embedCmd.Flags().StringVar(&embedArch, "arch", "amd64", "Target architecture (amd64, 386, arm64)")
	embedCmd.Flags().StringVar(&embedMethod, "method", "auto", "Method to use: auto, goversioninfo, rsrc")
	embedCmd.Flags().StringVar(&embedViPath, "vi", "", "Path to versioninfo.json (optional, for goversioninfo)")
	embedCmd.Flags().StringVar(&embedAppName, "name", "app", "Application name for version info")
	embedCmd.Flags().StringVar(&embedVersion, "version", "1.0.0", "Application version")
	embedCmd.Flags().StringVar(&embedCompany, "company", "", "Company name")
	embedCmd.Flags().StringVar(&embedCopyright, "copyright", "", "Copyright notice")

	// Keep --tool as hidden alias for backward compat
	embedCmd.Flags().String("tool", "", "Deprecated: use --method instead")
	_ = embedCmd.Flags().MarkHidden("tool")

	_ = embedCmd.MarkFlagRequired("ico")
}

func runEmbed(cmd *cobra.Command, _ []string) error {
	method := embedMethod

	// Support deprecated --tool flag for backward compat
	if tool, _ := cmd.Flags().GetString("tool"); tool != "" {
		method = tool
	}

	switch method {
	case "auto":
		return embedAuto(cmd)
	case "goversioninfo":
		return embedGoversioninfo(cmd)
	case "rsrc":
		return embedRsrc(cmd)
	default:
		return fmt.Errorf("unsupported method: %s (use 'auto', 'goversioninfo', or 'rsrc')", method)
	}
}

func embedAuto(cmd *cobra.Command) error {
	// Try goversioninfo (pure Go) first
	if err := embedGoversioninfo(cmd); err == nil {
		return nil
	}

	// Fall back to rsrc
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "goversioninfo failed, falling back to rsrc...\n")
	return embedRsrc(cmd)
}

func embedGoversioninfo(cmd *cobra.Command) error {
	// If a versioninfo.json was provided, use it directly
	if embedViPath != "" {
		if err := platform.GenerateSysoFromJSON(embedViPath, embedOutputPath, embedArch); err != nil {
			return fmt.Errorf("goversioninfo (JSON): %w", err)
		}
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Generated .syso (goversioninfo/JSON): %s\n", embedOutputPath)
		return nil
	}

	// Build from flags
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

	sysoPath, err := platform.GenerateSysoGoversioninfo(cfg, embedArch)
	if err != nil {
		return fmt.Errorf("goversioninfo: %w", err)
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Generated .syso (goversioninfo): %s\n", sysoPath)
	return nil
}

func embedRsrc(cmd *cobra.Command) error {
	if err := platform.GenerateSyso(embedICOPath, embedOutputPath, embedArch); err != nil {
		return fmt.Errorf("rsrc: %w", err)
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Generated .syso (rsrc): %s\n", embedOutputPath)
	return nil
}

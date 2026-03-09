package cmd

import (
	"fmt"

	"github.com/inovacc/iconforge/internal/platform"
	"github.com/spf13/cobra"
)

var (
	embedICOPath    string
	embedOutputPath string
	embedArch       string
	embedTool       string
	embedViPath     string
)

var embedCmd = &cobra.Command{
	Use:   "embed",
	Short: "Generate .syso for Go build (Windows)",
	Long: `Generate a .syso resource file from an ICO file for Windows.

The .syso file is automatically linked by the Go linker during build,
embedding the icon and version info into the resulting .exe.

Supports two tools:
  rsrc           - Simple icon embedding (github.com/akavel/rsrc)
  goversioninfo  - Full version info + icon (github.com/josephspurrier/goversioninfo)

Examples:
  iconforge embed --ico build/icons/windows/icon.ico
  iconforge embed --ico icon.ico --tool goversioninfo --vi versioninfo.json
  iconforge embed --ico icon.ico --arch arm64`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		switch embedTool {
		case "rsrc":
			if err := platform.GenerateSyso(embedICOPath, embedOutputPath, embedArch); err != nil {
				return fmt.Errorf("rsrc: %w", err)
			}
		case "goversioninfo":
			if embedViPath == "" {
				return fmt.Errorf("--vi (versioninfo.json path) is required for goversioninfo")
			}
			if err := platform.GenerateSysoWithGovernVersionInfo(embedViPath, embedOutputPath); err != nil {
				return fmt.Errorf("goversioninfo: %w", err)
			}
		default:
			return fmt.Errorf("unsupported tool: %s (use 'rsrc' or 'goversioninfo')", embedTool)
		}

		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Generated .syso: %s\n", embedOutputPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(embedCmd)

	embedCmd.Flags().StringVar(&embedICOPath, "ico", "", "Path to ICO file")
	embedCmd.Flags().StringVarP(&embedOutputPath, "output", "o", "rsrc.syso", "Output .syso file path")
	embedCmd.Flags().StringVar(&embedArch, "arch", "amd64", "Target architecture (amd64, 386, arm64)")
	embedCmd.Flags().StringVar(&embedTool, "tool", "rsrc", "Tool to use (rsrc or goversioninfo)")
	embedCmd.Flags().StringVar(&embedViPath, "vi", "", "Path to versioninfo.json (goversioninfo only)")

	_ = embedCmd.MarkFlagRequired("ico")
}

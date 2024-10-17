package cmd

import (
	"github.com/blang/semver"
	"github.com/k0kubun/pp"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/spf13/cobra"
)

func NewUpdateCommand() *cobra.Command {
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "updates the PackageLock binary",
		Run: func(cmd *cobra.Command, args []string) {
			pp.Println(AppVersion)
			v := semver.MustParse(AppVersion)
			latest, err := selfupdate.UpdateSelf(v, "hilkopterbob/PackageLock")
			if err != nil {
				pp.Println("Binary update failed:", err)
			}
			if latest.Version.Equals(v) {
				// latest version is the same as current version. It means current binary is up to date.
				pp.Println("Current binary is the latest version", AppVersion)
			} else {
				pp.Println("Successfully updated to version", latest.Version)
				pp.Println("Release note:\n", latest.ReleaseNotes)
			}
		},
	}

	return updateCmd
}

//func runUpdate(app *fiber.App, logger *zap.Logger) {
//}

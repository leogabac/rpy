package python

import (
	"fmt"
	"os"

	"github.com/pterm/pterm"
)

func Remove(version string) error {
	root, err := managedInstallDir(version)
	if err != nil {
		return err
	}

	if _, err := os.Stat(root); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("managed runtime %q does not exist", version)
		}
		return err
	}

	if err := os.RemoveAll(root); err != nil {
		return err
	}

	pterm.Println(pterm.FgLightGreen.Sprint("  OK  ") + "removed " + shortenPath(root))
	return nil
}

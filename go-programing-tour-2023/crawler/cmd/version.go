package cmd

import (
	"fmt"
	"strings"
)

// Version information.
var (
	BuildTS   = "None"
	GitHash   string
	GitBranch = "None"
	Version   = "None"
)

func GetVersion() string {
	if strings.TrimSpace(GitHash) != "" {
		h := GitHash
		if len(h) > 7 {
			h = h[:7]
		}

		return fmt.Sprintf("%s-%s", Version, h)
	}

	return Version
}

// Printer print build version.
func Printer() {
	fmt.Println("Version:          ", GetVersion())
	fmt.Println("Git Branch:       ", GitBranch)
	fmt.Println("Git Commit:       ", GitHash)
	fmt.Println("Build Time (UTC): ", BuildTS)
}

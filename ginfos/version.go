package ginfos

import (
	"fmt"
	"os"
)

var (
	GoVersion    string
	GitBranch    string
	GitCommitID  string
	GitTag       string
	GitBuildTime string
)

func init() {
	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "version") {
		fmt.Println("Go Version:   ", GoVersion)
		fmt.Println("Git Branch:   ", GitBranch)
		fmt.Println("Git CommitID: ", GitCommitID)
		fmt.Println("Git Tag:      ", GitTag)
		fmt.Println("Build Time:   ", GitBuildTime)

		os.Exit(0)
	}
}

func Version() {
}

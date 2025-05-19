package version

import "fmt"

var (
	Version = "dev"
	Commit  = "none"
)

func GetFullVersionInfo() map[string]string {
	return map[string]string{
		"version": Version,
		"commit":  Commit,
	}
}

func VersionCmd(args ...string) int {
	version := fmt.Sprintf("version: %s\ncommit: %s", Version, Commit)
	fmt.Println(version)
	return 0
}

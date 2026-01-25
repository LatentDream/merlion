package logout

import (
	"fmt"
	"merlion/internal/store/cloud"
)

func Cmd(args ...string) int {
	credMgr, err := cloud.NewCredentialsManager()
	if err != nil {
		fmt.Printf("Failed to initialize credentials manager: %v", err)
		return 1
	}

	err = credMgr.RemoveCredentials()
	if err != nil {
		fmt.Println("ERROR: Failed to remove credentials: ", err)
		return 1
	}

	fmt.Println("Credentials has been removed from disk")
	return 0
}

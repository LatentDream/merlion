package logout

import (
	"fmt"
	"merlion/internal/config"
	"merlion/internal/vault/cloud"
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

	config := config.Load()
	for i, vault := range config.Vaults {
		if vault.Provider == cloud.Name {
			config.Vaults = append(config.Vaults[:i], config.Vaults[i+1:]...)
			break
		}
	}
	config.Save()

	fmt.Println("Credentials has been removed from disk")
	return 0
}

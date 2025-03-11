package driver

import (
	"fmt"
	"os"

	"github.com/jackdzi/feederizer/ui/internal/api"
	"github.com/jackdzi/feederizer/ui/internal/config"
)

func (m *model) handleLogin(user string, pass string, initialBoot bool) bool {
	userConfig := config.ReturnConfig()
	if initialBoot {
		if !userConfig.Get("authentication.autoLogin").(bool) {
			return false
		}
		canLogin, err := api.CheckUser(userConfig.Get("authentication.user").(string), userConfig.Get("authentication.pass").(string))
		if err != nil {
			fmt.Println("Error")
			return false
		}
		if canLogin {
			m.user = user
			return true
		}
		return false
	}
	m.user = user
	userConfig.Set("authentication.user", user)
	userConfig.Set("authentication.pass", pass)
	data, err := userConfig.Marshal()
	if err != nil {
		fmt.Println("Error")
		return false
	}
	err = os.WriteFile(config.GetFilePath(), data, 0644)
	if err != nil {
		fmt.Println("Error writing config.toml:", err)
		return false
	}
	return true
}

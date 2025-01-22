package update

import (
	"Typhoon/config"
	"Typhoon/utils"
	"fmt"
	"path"
)

func UpdateSubscriptions() error {
	// Update the subscriptions
	cfg := config.GetConfig(config.ConfigFilePath, false)
	subscriptions := cfg.SubscriptionManage.Subscriptions

	for _, subscription := range subscriptions {
		// Update the subscription
		configPath := path.Join(utils.GetExecutableDir(), "mihomo", subscription.Name, "config.yaml")
		if err := utils.Download(subscription.URL, configPath); err != nil {
			return fmt.Errorf("failed to download file: %v", err)
		}
	}
	return nil
}

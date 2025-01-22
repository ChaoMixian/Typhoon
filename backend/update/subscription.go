package update

import (
	"Typhoon/config"
	"Typhoon/utils"
	"log"
	"path"
)

func UpdateSubscriptions() error {
	// Update the subscriptions
	cfg := config.GetConfig(config.ConfigFilePath, false)
	subscriptions := cfg.SubscriptionManage.Subscriptions

	for _, subscription := range subscriptions {
		// Update the subscription
		configPath := path.Join(utils.GetExecutableDir(), "mihomo", "config", subscription.Name, "config.yaml")
		// if err := utils.Download(subscription.URL, configPath); err != nil {
		// 	return fmt.Errorf("failed to download file: %v", err)
		// }
		err := utils.DownloadWithProgress(subscription.URL, configPath, func(downloaded, total int64) {
			progress := float64(downloaded) / float64(total) * 100
			log.Printf("Progress: %.2f%%", progress)
		})
		if err != nil {
			log.Fatalf("Download failed: %v", err)
		}
	}
	return nil
}

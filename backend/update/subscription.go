package update

import (
	"Typhoon/config"
	"Typhoon/utils"
	"log"
	"path"
)

type SubscriptionUpdateResult struct {
	Name    string `json:"name"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func UpdateSubscriptions() ([]SubscriptionUpdateResult, error) {
	// Update the subscriptions
	cfg, _ := config.GetConfig(config.ConfigFilePath, false)
	subscriptions := cfg.SubscriptionManage.Subscriptions
	results := make([]SubscriptionUpdateResult, 0)

	for i, subscription := range subscriptions {
		// Update the subscription

		result := SubscriptionUpdateResult{
			Name:    subscription.Name,
			Success: true,
		}

		execDir, err := utils.GetExecutableDir()
		if err != nil {
			result.Success = false
			result.Error = err.Error()
			results = append(results, result)
			log.Printf("failed to get executable dir for %s: %v", subscription.Name, err)
			continue
		}
		configPath := path.Join(execDir, "mihomo", "config", subscription.Name, "config.yaml")
		err = utils.DownloadWithProgress(subscription.URL, configPath, func(downloaded, total int64) {
			progress := float64(downloaded) / float64(total) * 100
			log.Printf("Progress: %.2f%%", progress)
		})
		if err != nil {
			result.Success = false
			result.Error = err.Error()
			results = append(results, result)
			log.Printf("failed to download %s: %v", subscription.Name, err) // 记录错误但继续执行
			continue                                                        // 跳过当前订阅，继续处理下一个
		}

		// Modify the subscription path configuration
		subscriptions[i].Path = configPath

		results = append(results, result)
		log.Printf("Subscription %s updated", subscription.Name)
	}
	cfg.SubscriptionManage.Subscriptions = subscriptions

	config.SaveConfig(config.ConfigFilePath, cfg)
	config.ReloadConfig(config.ConfigFilePath)

	return results, nil
}

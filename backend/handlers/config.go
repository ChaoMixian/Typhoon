package handlers

import (
	"Typhoon/config"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ReloadConfigHandler handles API requests to reload configuration
func ReloadConfigHandler(c *gin.Context) {
	// Attempt to reload the configuration. GetConfig will log errors.
	// If reload is true and an error occurs, GetConfig returns the previous instance and the error.
	// If initial load fails (which shouldn't happen here as app is running), GetConfig would Fatalf.
	_, err := config.GetConfig(config.ConfigFilePath, true)
	if err != nil {
		// GetConfig already logged the error.
		// err indicates reload failed but previous config (if any) is still active.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload the configuration. Server may be using previous configuration if available. Check server logs for details."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Configuration reloaded successfully"})
}

// UpdateRequestPayload defines the structure for the configuration update request.
// Fields are pointers to allow distinguishing between a field not being provided
// and a field being provided with an empty/zero value.
type UpdateRequestPayload struct {
	Proxy              *ProxyUpdatePayload               `json:"proxy,omitempty"`
	Logging            *config.LoggingConfigPart         `json:"logging,omitempty"`
	SubscriptionManage *config.SubscriptionManageConfigPart `json:"subscriptionManage,omitempty"`
	API                *config.APIConfigPart             `json:"api,omitempty"`
}

// ProxyUpdatePayload defines the structure for updating proxy-related configurations.
type ProxyUpdatePayload struct {
	CurrentCore *string                     `json:"currentCore,omitempty"`
	Mode        *string                     `json:"mode,omitempty"`
	Mihomo      *config.MihomoConfigPart `json:"mihomo,omitempty"`
	DNS         *config.DNSConfigPart    `json:"dns,omitempty"`
}

// UpdateConfigHandler handles API requests to update configuration
func UpdateConfigHandler(c *gin.Context) {
	filePath := config.ConfigFilePath // Use the default config file path

	var payload UpdateRequestPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// It's important to load the full current config first for any parts that are being partially updated
	// or for fields within parts that are not included in the payload (e.g. updating only Mihomo's port).
	// The new UpdateXYZConfig functions in config.go handle loading the current state,
	// applying the change for the *entire part*, and saving.

	if payload.API != nil {
		if err := config.UpdateAPIConfig(filePath, *payload.API); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update API config: " + err.Error()})
			return
		}
	}

	if payload.Logging != nil {
		if err := config.UpdateLoggingConfig(filePath, *payload.Logging); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Logging config: " + err.Error()})
			return
		}
	}

	if payload.SubscriptionManage != nil {
		if err := config.UpdateSubscriptionManageConfig(filePath, *payload.SubscriptionManage); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update SubscriptionManage config: " + err.Error()})
			return
		}
	}

	if payload.Proxy != nil {
		// For proxy, we need to load the current proxy config and apply individual changes
		// because CurrentCore, Mode, Mihomo, and DNS can be updated independently or together.
		// The UpdateProxyXYZConfig functions expect the full part.
		// So, we load the current config, modify only the parts of Proxy provided in the payload, then call the specific update functions.

		cfg := config.GetConfig(filePath, false) // Get current config (no reload)

		// Create copies of current proxy sub-parts to modify
		// This is crucial because the UpdateProxyMihomoConfig and UpdateProxyDNSConfig expect the full part.
		// If the payload only contains a partial update for Mihomo (e.g. just port), we need to preserve other Mihomo fields.

		// Handling Proxy.CurrentCore
		if payload.Proxy.CurrentCore != nil {
			if err := config.UpdateProxyCurrentCore(filePath, *payload.Proxy.CurrentCore); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Proxy.CurrentCore: " + err.Error()})
				return
			}
		}

		// Handling Proxy.Mode
		if payload.Proxy.Mode != nil {
			if err := config.UpdateProxyMode(filePath, *payload.Proxy.Mode); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Proxy.Mode: " + err.Error()})
				return
			}
		}

		// Handling Proxy.Mihomo
		if payload.Proxy.Mihomo != nil {
			// The UpdateProxyMihomoConfig function takes the *entire* MihomoConfigPart.
			// If the client sends only a few fields for Mihomo, we must load the current Mihomo config,
			// apply the partial updates from the payload, and then pass the complete, updated MihomoConfigPart.
			// This is not ideal. The UpdateProxyMihomoConfig should ideally handle partial updates or
			// the client must always send the full MihomoConfigPart if it wants to change anything in it.

			// For now, let's assume UpdateProxyMihomoConfig replaces the Mihomo part entirely.
			// This means the client MUST send the full MihomoConfigPart if it intends to update it.
			// This simplifies the handler logic.
			if err := config.UpdateProxyMihomoConfig(filePath, *payload.Proxy.Mihomo); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Proxy.Mihomo config: " + err.Error()})
				return
			}
		}

		// Handling Proxy.DNS
		if payload.Proxy.DNS != nil {
			// Similar to Mihomo, assumes client sends the full DNSConfigPart for updates.
			if err := config.UpdateProxyDNSConfig(filePath, *payload.Proxy.DNS); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Proxy.DNS config: " + err.Error()})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Configuration updated successfully"})
}

// Todo: 抽象订阅链接操作，增删改查 (This TODO is still relevant, especially for Subscriptions array)

/*
curl -X PATCH http://localhost:8080/api/v1/updateConfig \
-H "Content-Type: application/json" \
-d '{
    "proxy.currentCore": "xray",
    "proxy.mihomo.listenPort": 8888
}'

curl -X PATCH http://localhost:8080/api/v1/updateConfig \
-H "Content-Type: application/json" \
-d '{
    "proxy.dns.enabled": true,
    "proxy.dns.listen": "127.0.0.1:5353",
    "proxy.dns.upstreamDNS": ["1.1.1.1", "8.8.8.8"],
    "proxy.dns.fakeIPFilter": ["*.example.com", "localhost"]
}'


curl -X PATCH http://localhost:8080/api/v1/updateConfig \
-H "Content-Type: application/json" \
-d '{
    "subscriptionManage.enabled": true,
    "subscriptionManage.subscriptions": [
        {"name": "Sub1", "url": "https://example.com/sub1"},
        {"name": "Sub2", "url": "https://example.com/sub2"}
    ]
}'

*/

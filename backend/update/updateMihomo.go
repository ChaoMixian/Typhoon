package update

import (
	"Typhoon/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GitHubRelease represents a GitHub release
type GitHubRelease struct {
	Assets []struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	} `json:"assets"`
}

// UpdateMihomo performs the complete update process: fetching GitHub release, downloading, decompressing, and replacing.
func UpdateMihomo(owner, repo, destPath, finalPath string, progressCallback func(downloaded, total int64)) error {
	// Step 1: Fetch GitHub release
	release, err := FetchGitHubRelease(owner, repo)
	if err != nil {
		return fmt.Errorf("failed to fetch GitHub release: %v", err)
	}

	// Step 2: Select suitable asset
	downloadURL, err := SelectAssetURL(release)
	if err != nil {
		return fmt.Errorf("failed to select asset: %v", err)
	}

	// Step 3: Define paths
	downloadPath := filepath.Join(os.TempDir(), "mihomo.gz")

	// Step 4: Download the file
	if err := utils.DownloadWithProgress(downloadURL, downloadPath, progressCallback); err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}

	// Step 5: Decompress the file
	if err := utils.DecompressGz(downloadPath, destPath); err != nil {
		return fmt.Errorf("failed to decompress file: %v", err)
	}

	// Step 6: Replace the old version
	if err := utils.ReplaceOldVersion(destPath, finalPath); err != nil {
		return fmt.Errorf("failed to replace old version: %v", err)
	}

	//Step 7: Set execute permission
	if err := os.Chmod(finalPath, 0755); err != nil {
		return fmt.Errorf("failed to set execute permission: %v", err)
	}

	return nil
}

// fetchGitHubRelease fetches the latest release from GitHub
func FetchGitHubRelease(owner, repo string) (*GitHubRelease, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch GitHub release: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %v", resp.Status)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse GitHub release response: %v", err)
	}

	return &release, nil
}

// selectAssetURL selects a .gz download URL based on system architecture
func SelectAssetURL(release *GitHubRelease) (string, error) {
	arch := runtime.GOARCH
	log.Default().Println("arch: ", arch)
	for _, asset := range release.Assets {

		if strings.HasSuffix(asset.Name, ".gz") && ContainsArch(asset.Name, arch) {
			log.Default().Println("asset.Name: ", asset.Name)
			return asset.URL, nil
		}
	}
	return "", fmt.Errorf("no suitable .gz asset found for architecture: %s", arch)
}

// containsArch checks if the asset name matches the architecture
func ContainsArch(assetName, arch string) bool {
	return (arch == "amd64" && utils.ContainsAll(assetName, "amd64", "linux") && !utils.ContainsAny(assetName, "compatible", "go")) ||
		(arch == "arm64" && utils.ContainsAll(assetName, "arm64", "linux") ||
			(arch == "armv5" && utils.ContainsAll(assetName, "armv5", "linux")) ||
			(arch == "armv6" && utils.ContainsAll(assetName, "armv6", "linux")) ||
			(arch == "armv7" && utils.ContainsAll(assetName, "armv7", "linux")))
}

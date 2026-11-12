package main

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// sopsSecretsFile is the sops+age encrypted dotenv file holding the API key
// and default ZIP code.
// It is safe to commit; only holders of a matching age private key
// (default: ~/.config/sops/age/keys.txt) can decrypt it.
const sopsSecretsFile = "secrets.enc.env"

// appConfigDirName is the per-user config directory name; `task install`
// copies secrets.enc.env to $XDG_CONFIG_HOME/weather-tui/ (~/.config/weather-tui/).
const appConfigDirName = "weather-tui"

// findSecretsFile returns the path of the encrypted secrets file, checking
// the current directory first (dev workflow) and then the user config dir
// (installed binary). Returns "" if neither exists.
func findSecretsFile() string {
	candidates := []string{sopsSecretsFile}
	if dir, err := os.UserConfigDir(); err == nil {
		candidates = append(candidates, filepath.Join(dir, appConfigDirName, sopsSecretsFile))
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return ""
}

var (
	sopsOnce      sync.Once
	sopsPlaintext string
)

// decryptSops decrypts the secrets file (see findSecretsFile) with the sops CLI (once per process)
// and returns the dotenv plaintext. The plaintext only ever lives in memory.
// Returns "" if the file is missing, sops isn't installed, or decryption fails.
func decryptSops() string {
	sopsOnce.Do(func() {
		path := findSecretsFile()
		if path == "" {
			return
		}
		sopsPath, err := exec.LookPath("sops")
		if err != nil {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		out, err := exec.CommandContext(ctx, sopsPath, "decrypt",
			"--input-type", "dotenv", "--output-type", "dotenv", path).Output()
		if err != nil {
			return
		}
		sopsPlaintext = string(out)
	})
	return sopsPlaintext
}

// readSecretFromSops returns key from the decrypted sops secrets, or "".
func readSecretFromSops(key string) string {
	return parseDotenvValue(decryptSops(), key)
}

// readAPIKeyFromSops returns OPENWEATHER_API_KEY from the sops secrets.
func readAPIKeyFromSops() string {
	return readSecretFromSops("OPENWEATHER_API_KEY")
}

// getDefaultZip returns the default ZIP code. Lookup order:
//  1. DEFAULT_ZIP_CODE environment variable
//  2. sops+age encrypted secrets.enc.env
//  3. plaintext .env file (legacy)
func getDefaultZip() string {
	if zip := os.Getenv("DEFAULT_ZIP_CODE"); zip != "" {
		return zip
	}
	if zip := readSecretFromSops("DEFAULT_ZIP_CODE"); zip != "" {
		return zip
	}
	if data, err := os.ReadFile(".env"); err == nil {
		return parseDotenvValue(string(data), "DEFAULT_ZIP_CODE")
	}
	return ""
}

// parseDotenvValue returns the value of key from dotenv-formatted content,
// stripping surrounding quotes. Returns "" if not found.
func parseDotenvValue(content, key string) string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	prefix := key + "="
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || !strings.HasPrefix(line, prefix) {
			continue
		}
		value := strings.TrimSpace(strings.TrimPrefix(line, prefix))
		if len(value) >= 2 &&
			((value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'')) {
			value = value[1 : len(value)-1]
		}
		return value
	}
	return ""
}

package main

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"strings"
	"time"
)

// sopsSecretsFile is the sops+age encrypted dotenv file holding the API key.
// It is safe to commit; only holders of a matching age private key
// (default: ~/.config/sops/age/keys.txt) can decrypt it.
const sopsSecretsFile = "secrets.enc.env"

// readAPIKeyFromSops decrypts sopsSecretsFile with the sops CLI and returns
// OPENWEATHER_API_KEY. The plaintext only ever lives in memory.
// Returns "" if the file is missing, sops isn't installed, or decryption fails.
func readAPIKeyFromSops() string {
	if _, err := os.Stat(sopsSecretsFile); err != nil {
		return ""
	}
	sopsPath, err := exec.LookPath("sops")
	if err != nil {
		return ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	out, err := exec.CommandContext(ctx, sopsPath, "decrypt",
		"--input-type", "dotenv", "--output-type", "dotenv", sopsSecretsFile).Output()
	if err != nil {
		return ""
	}
	return parseDotenvValue(string(out), "OPENWEATHER_API_KEY")
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

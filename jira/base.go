package jira

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

type Setting struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

func (s Setting) Print() string {
	return fmt.Sprintf("[Url: %s Token: %s]", s.URL, s.Token)
}

const (
	configFile string = "jira.json"
)

func getConfigPath(projectDir string) string {
	return path.Join(projectDir, configFile)
}

func LoadSetting(projectDir string) (*Setting, error) {
	content, err := os.ReadFile(getConfigPath(projectDir))
	if err != nil {
		return nil, err
	}
	var setting Setting
	err = json.Unmarshal(content, &setting)
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func SaveSetting(projectDir string, setting *Setting) error {
	content, err := json.MarshalIndent(setting, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(getConfigPath(projectDir), content, 0644)
	if err != nil {
		return err
	}
	return nil
}

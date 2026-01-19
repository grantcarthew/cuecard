package config

import (
	"fmt"
	"os"
	"path/filepath"

	"cuelang.org/go/cue/cuecontext"
)

// Config represents the application configuration
type Config struct {
	PromptsDir string       `json:"prompts_dir"`
	Editor     string       `json:"editor"`
	Theme      string       `json:"theme"`
	Window     WindowConfig `json:"window"`
}

// WindowConfig represents window settings
type WindowConfig struct {
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Position string `json:"position"`
}

// DefaultConfig returns a Config with default values
func DefaultConfig() Config {
	return Config{
		PromptsDir: "",
		Editor:     "code",
		Theme:      "system",
		Window: WindowConfig{
			Width:    1024,
			Height:   768,
			Position: "center",
		},
	}
}

// ConfigDir returns the path to the config directory
func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".config", "cuecard"), nil
}

// ConfigPath returns the path to the config file
func ConfigPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.cue"), nil
}

// Load reads and parses the configuration file
func Load() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}
	return LoadFromPath(path)
}

// LoadFromPath reads and parses a configuration file from a specific path
func LoadFromPath(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	return Parse(string(data))
}

// Parse parses CUE configuration content
func Parse(content string) (*Config, error) {
	ctx := cuecontext.New()
	value := ctx.CompileString(content)
	if err := value.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse CUE config: %w", err)
	}

	cfg := DefaultConfig()
	if err := value.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.PromptsDir == "" {
		return fmt.Errorf("prompts_dir is required")
	}

	// Expand home directory if needed
	if len(c.PromptsDir) > 0 && c.PromptsDir[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to expand home directory: %w", err)
		}
		c.PromptsDir = filepath.Join(home, c.PromptsDir[1:])
	}

	// Validate theme
	switch c.Theme {
	case "light", "dark", "system", "":
		// valid
	default:
		return fmt.Errorf("invalid theme: %s (must be light, dark, or system)", c.Theme)
	}

	// Validate window position
	switch c.Window.Position {
	case "remember", "center", "":
		// valid
	default:
		return fmt.Errorf("invalid window position: %s (must be remember or center)", c.Window.Position)
	}

	return nil
}

// Exists checks if the config file exists
func Exists() (bool, error) {
	path, err := ConfigPath()
	if err != nil {
		return false, err
	}
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// Save writes the configuration to the config file
func (c *Config) Save() error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	return c.SaveToPath(path)
}

// SaveToPath writes the configuration to a specific path
func (c *Config) SaveToPath(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	content := c.ToCUE()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

// ToCUE converts the config to CUE format
func (c *Config) ToCUE() string {
	var windowSection string
	if c.Window.Width != 0 || c.Window.Height != 0 || c.Window.Position != "" {
		windowSection = fmt.Sprintf(`window: {
	width:    %d
	height:   %d
	position: %q
}
`, c.Window.Width, c.Window.Height, c.Window.Position)
	}

	return fmt.Sprintf(`prompts_dir: %q
editor:      %q
theme:       %q
%s`, c.PromptsDir, c.Editor, c.Theme, windowSection)
}

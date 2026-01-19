package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    *Config
		wantErr bool
	}{
		{
			name: "valid minimal config",
			content: `prompts_dir: "/home/user/prompts"
editor: "code"`,
			want: &Config{
				PromptsDir: "/home/user/prompts",
				Editor:     "code",
				Theme:      "system",
				Window: WindowConfig{
					Width:    1024,
					Height:   768,
					Position: "center",
				},
			},
			wantErr: false,
		},
		{
			name: "valid full config",
			content: `prompts_dir: "/home/user/prompts"
editor: "nvim"
theme: "dark"
window: {
	width: 800
	height: 600
	position: "remember"
}`,
			want: &Config{
				PromptsDir: "/home/user/prompts",
				Editor:     "nvim",
				Theme:      "dark",
				Window: WindowConfig{
					Width:    800,
					Height:   600,
					Position: "remember",
				},
			},
			wantErr: false,
		},
		{
			name:    "missing prompts_dir",
			content: `editor: "code"`,
			wantErr: true,
		},
		{
			name: "invalid theme",
			content: `prompts_dir: "/home/user/prompts"
editor: "code"
theme: "invalid"`,
			wantErr: true,
		},
		{
			name: "invalid window position",
			content: `prompts_dir: "/home/user/prompts"
editor: "code"
window: {
	position: "invalid"
}`,
			wantErr: true,
		},
		{
			name:    "invalid CUE syntax",
			content: `prompts_dir: "unclosed`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.PromptsDir != tt.want.PromptsDir {
				t.Errorf("PromptsDir = %v, want %v", got.PromptsDir, tt.want.PromptsDir)
			}
			if got.Editor != tt.want.Editor {
				t.Errorf("Editor = %v, want %v", got.Editor, tt.want.Editor)
			}
			if got.Theme != tt.want.Theme {
				t.Errorf("Theme = %v, want %v", got.Theme, tt.want.Theme)
			}
			if got.Window.Width != tt.want.Window.Width {
				t.Errorf("Window.Width = %v, want %v", got.Window.Width, tt.want.Window.Width)
			}
			if got.Window.Height != tt.want.Window.Height {
				t.Errorf("Window.Height = %v, want %v", got.Window.Height, tt.want.Window.Height)
			}
			if got.Window.Position != tt.want.Window.Position {
				t.Errorf("Window.Position = %v, want %v", got.Window.Position, tt.want.Window.Position)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Editor != "code" {
		t.Errorf("default Editor = %v, want code", cfg.Editor)
	}
	if cfg.Theme != "system" {
		t.Errorf("default Theme = %v, want system", cfg.Theme)
	}
	if cfg.Window.Width != 1024 {
		t.Errorf("default Window.Width = %v, want 1024", cfg.Window.Width)
	}
	if cfg.Window.Height != 768 {
		t.Errorf("default Window.Height = %v, want 768", cfg.Window.Height)
	}
	if cfg.Window.Position != "center" {
		t.Errorf("default Window.Position = %v, want center", cfg.Window.Position)
	}
}

func TestToCUE(t *testing.T) {
	cfg := &Config{
		PromptsDir: "/home/user/prompts",
		Editor:     "code",
		Theme:      "dark",
		Window: WindowConfig{
			Width:    800,
			Height:   600,
			Position: "remember",
		},
	}

	cue := cfg.ToCUE()

	// Parse it back to verify round-trip
	parsed, err := Parse(cue)
	if err != nil {
		t.Fatalf("failed to parse generated CUE: %v", err)
	}

	if parsed.PromptsDir != cfg.PromptsDir {
		t.Errorf("round-trip PromptsDir = %v, want %v", parsed.PromptsDir, cfg.PromptsDir)
	}
	if parsed.Editor != cfg.Editor {
		t.Errorf("round-trip Editor = %v, want %v", parsed.Editor, cfg.Editor)
	}
	if parsed.Theme != cfg.Theme {
		t.Errorf("round-trip Theme = %v, want %v", parsed.Theme, cfg.Theme)
	}
}

func TestSaveAndLoad(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.cue")

	cfg := &Config{
		PromptsDir: "/home/user/prompts",
		Editor:     "vim",
		Theme:      "light",
		Window: WindowConfig{
			Width:    1280,
			Height:   720,
			Position: "center",
		},
	}

	// Save
	if err := cfg.SaveToPath(configPath); err != nil {
		t.Fatalf("SaveToPath() error = %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("config file was not created")
	}

	// Load
	loaded, err := LoadFromPath(configPath)
	if err != nil {
		t.Fatalf("LoadFromPath() error = %v", err)
	}

	if loaded.PromptsDir != cfg.PromptsDir {
		t.Errorf("loaded PromptsDir = %v, want %v", loaded.PromptsDir, cfg.PromptsDir)
	}
	if loaded.Editor != cfg.Editor {
		t.Errorf("loaded Editor = %v, want %v", loaded.Editor, cfg.Editor)
	}
	if loaded.Theme != cfg.Theme {
		t.Errorf("loaded Theme = %v, want %v", loaded.Theme, cfg.Theme)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: Config{
				PromptsDir: "/home/user/prompts",
				Editor:     "code",
				Theme:      "system",
			},
			wantErr: false,
		},
		{
			name: "empty prompts_dir",
			cfg: Config{
				PromptsDir: "",
				Editor:     "code",
			},
			wantErr: true,
		},
		{
			name: "valid light theme",
			cfg: Config{
				PromptsDir: "/home/user/prompts",
				Theme:      "light",
			},
			wantErr: false,
		},
		{
			name: "valid dark theme",
			cfg: Config{
				PromptsDir: "/home/user/prompts",
				Theme:      "dark",
			},
			wantErr: false,
		},
		{
			name: "invalid theme",
			cfg: Config{
				PromptsDir: "/home/user/prompts",
				Theme:      "invalid",
			},
			wantErr: true,
		},
		{
			name: "valid remember position",
			cfg: Config{
				PromptsDir: "/home/user/prompts",
				Window: WindowConfig{
					Position: "remember",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid position",
			cfg: Config{
				PromptsDir: "/home/user/prompts",
				Window: WindowConfig{
					Position: "invalid",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

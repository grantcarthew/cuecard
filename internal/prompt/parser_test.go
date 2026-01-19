package prompt

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    *Prompt
		wantErr bool
	}{
		{
			name: "valid prompt with all fields",
			content: `---
title: Test Prompt
description: A test prompt
group: Testing
tags: [test, example]
input: required
input_hint: Enter something
favorite: true
---

This is the content.

${INPUT}`,
			want: &Prompt{
				Title:       "Test Prompt",
				Description: "A test prompt",
				Group:       "Testing",
				Tags:        []string{"test", "example"},
				Input:       "required",
				InputHint:   "Enter something",
				Favorite:    true,
				Content:     "This is the content.\n\n${INPUT}",
			},
			wantErr: false,
		},
		{
			name: "minimal prompt",
			content: `---
title: Minimal
---

Just content.`,
			want: &Prompt{
				Title:   "Minimal",
				Content: "Just content.",
			},
			wantErr: false,
		},
		{
			name:    "no frontmatter",
			content: "Just plain content without frontmatter.",
			want: &Prompt{
				Content: "Just plain content without frontmatter.",
			},
			wantErr: false,
		},
		{
			name: "optional input",
			content: `---
title: Optional Input
input: optional
---

Content here.`,
			want: &Prompt{
				Title:   "Optional Input",
				Input:   "optional",
				Content: "Content here.",
			},
			wantErr: false,
		},
		{
			name: "unclosed frontmatter",
			content: `---
title: Broken
no closing delimiter`,
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
			if got.Title != tt.want.Title {
				t.Errorf("Title = %q, want %q", got.Title, tt.want.Title)
			}
			if got.Description != tt.want.Description {
				t.Errorf("Description = %q, want %q", got.Description, tt.want.Description)
			}
			if got.Group != tt.want.Group {
				t.Errorf("Group = %q, want %q", got.Group, tt.want.Group)
			}
			if got.Input != tt.want.Input {
				t.Errorf("Input = %q, want %q", got.Input, tt.want.Input)
			}
			if got.InputHint != tt.want.InputHint {
				t.Errorf("InputHint = %q, want %q", got.InputHint, tt.want.InputHint)
			}
			if got.Favorite != tt.want.Favorite {
				t.Errorf("Favorite = %v, want %v", got.Favorite, tt.want.Favorite)
			}
			if got.Content != tt.want.Content {
				t.Errorf("Content = %q, want %q", got.Content, tt.want.Content)
			}
			if len(got.Tags) != len(tt.want.Tags) {
				t.Errorf("Tags length = %d, want %d", len(got.Tags), len(tt.want.Tags))
			}
		})
	}
}

func TestHasFrontmatter(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{
			name: "has frontmatter",
			content: `---
title: Test
---
Content`,
			want: true,
		},
		{
			name:    "no frontmatter",
			content: "Just content",
			want:    false,
		},
		{
			name:    "only opening delimiter",
			content: "---\ntitle: Test",
			want:    false,
		},
		{
			name:    "empty",
			content: "",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasFrontmatter(tt.content); got != tt.want {
				t.Errorf("HasFrontmatter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrompt_RequiresInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"required", "required", true},
		{"optional", "optional", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Prompt{Input: tt.input}
			if got := p.RequiresInput(); got != tt.want {
				t.Errorf("RequiresInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrompt_HasOptionalInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"required", "required", false},
		{"optional", "optional", true},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Prompt{Input: tt.input}
			if got := p.HasOptionalInput(); got != tt.want {
				t.Errorf("HasOptionalInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrompt_HasInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"required", "required", true},
		{"optional", "optional", true},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Prompt{Input: tt.input}
			if got := p.HasInput(); got != tt.want {
				t.Errorf("HasInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

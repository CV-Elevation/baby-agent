package skill

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"babyagent/shared"
)

// frontMatter represents the YAML front matter structure
type frontMatter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// LoadSkill loads both metadata and content from a SKILLS.md file by id
// Returns (metadata, content, error)
func LoadSkill(id string) (SkillMetadata, string, error) {
	skillPath := filepath.Join(shared.GetWorkspaceDir(), ".babyagent", "skills", id, "SKILL.md")

	contentBytes, err := os.ReadFile(skillPath)
	if err != nil {
		return SkillMetadata{}, "", fmt.Errorf("failed to read skill file: %w", err)
	}

	text := string(contentBytes)

	// Split by front matter delimiter
	parts := strings.SplitN(text, "---", 3)
	if len(parts) < 3 {
		return SkillMetadata{}, "", errors.New("skill file must have YAML front matter enclosed in `---`")
	}

	// parts[0] is empty (before first ---)
	// parts[1] is the YAML front matter
	// parts[2] is the body content

	// Parse front matter using yaml.Unmarshal
	var fm frontMatter
	if err := yaml.Unmarshal([]byte(parts[1]), &fm); err != nil {
		return SkillMetadata{}, "", fmt.Errorf("failed to parse front matter: %w", err)
	}

	if fm.Name == "" {
		return SkillMetadata{}, "", errors.New("skill must have a 'name' field in front matter")
	}

	if fm.Description == "" {
		return SkillMetadata{}, "", errors.New("skill must have a 'description' field in front matter")
	}

	meta := SkillMetadata{
		ID:          id,
		Name:        fm.Name,
		Description: fm.Description,
	}

	// Extract body content (trim leading newline)
	bodyContent := strings.TrimSpace(parts[2])

	return meta, bodyContent, nil
}

package validate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"sigs.k8s.io/yaml"
)

// loadFile reads a test file from the project root testdata directory.
func loadFile(t *testing.T, relPath string) []byte {
	t.Helper()
	// Determine project root relative to this test file
	root := filepath.Join("..", "..")
	fullPath := filepath.Join(root, relPath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", fullPath, err)
	}
	return data
}

// Test legacy JSON composer format (for backward compatibility)
func TestValidComposerJSON(t *testing.T) {
	v := loadFile(t, "testdata/valid.json")
	if err := ValidateComposerJSON(v); err != nil {
		t.Errorf("expected valid.json to pass, but got: %v", err)
	}
}

func TestInvalidComposerJSON(t *testing.T) {
	v := loadFile(t, "testdata/invalid.json")
	if err := ValidateComposerJSON(v); err == nil {
		t.Errorf("expected invalid.json to fail validation")
	}
}

// Test new YAML image template format
func TestValidImageTemplate(t *testing.T) {
	v := loadFile(t, "image-templates/azl3-x86_64-edge-raw.yml")

	// Parse to generic JSON interface
	var raw interface{}
	if err := yaml.Unmarshal(v, &raw); err != nil {
		t.Errorf("yml parsing error: %v", err)
		return
	}

	// Re‐marshal to JSON bytes
	dataJSON, err := json.Marshal(raw)
	if err != nil {
		t.Errorf("json marshaling error: %v", err)
		return
	}
	if err := ValidateImageJSON(dataJSON); err != nil {
		t.Errorf("expected image-templates/azl3-x86_64-edge-raw.yml to pass, but got: %v", err)
	}
}

func TestInvalidImageTemplate(t *testing.T) {
	v := loadFile(t, "testdata/invalid-image.yml")

	// Parse to generic JSON interface
	var raw interface{}
	if err := yaml.Unmarshal(v, &raw); err != nil {
		t.Errorf("yml parsing error: %v", err)
		return
	}

	// Re‐marshal to JSON bytes
	dataJSON, err := json.Marshal(raw)
	if err != nil {
		t.Errorf("json marshaling error: %v", err)
		return
	}

	if err := ValidateImageJSON(dataJSON); err == nil {
		t.Errorf("expected testdata/invalid-image.yml to pass, but got: %v", err)
	}
}

// Test global config validation
func TestValidConfig(t *testing.T) {
	v := loadFile(t, "testdata/valid-config.yml")

	if v == nil {
		t.Fatal("failed to load testdata/valid-config.yml")
	}
	dataJSON, err := yaml.YAMLToJSON(v)

	if err != nil {
		t.Fatalf("YAML→JSON conversion failed: %v", err)
	}
	if err := ValidateConfigJSON(dataJSON); err != nil {
		t.Errorf("validation failed: %v", err)
	}
}

func TestInvalidConfig(t *testing.T) {
	v := loadFile(t, "testdata/invalid-config.yml")

	// Parse to generic JSON interface
	var raw interface{}
	if err := yaml.Unmarshal(v, &raw); err != nil {
		t.Errorf("yml parsing error: %v", err)
		return
	}

	// Re‐marshal to JSON bytes
	dataJSON, err := yaml.YAMLToJSON(v)
	if err != nil {
		t.Errorf("json marshaling error: %v", err)
		return
	}

	if err := ValidateConfigJSON(dataJSON); err == nil {
		t.Errorf("expected invalid-config.json to fail validation: %v", err)
	} else {
		t.Logf("expected validation error: %v", err)
	}
}

// Test validation of template structure
func TestImageTemplateStructure(t *testing.T) {
	// Test a minimal valid template
	minimalTemplate := `image:
  name: test-image
  version: "1.0.0"

target:
  os: azure-linux
  dist: azl3
  arch: x86_64
  imageType: raw

systemConfigs:
  - name: minimal
    description: Minimal test configuration
    packages:
      - openssh-server
    kernel:
      version: "6.12"
      cmdline: "quiet"
`

	var raw interface{}
	if err := yaml.Unmarshal([]byte(minimalTemplate), &raw); err != nil {
		t.Fatalf("failed to parse minimal template: %v", err)
	}

	dataJSON, err := json.Marshal(raw)
	if err != nil {
		t.Fatalf("failed to marshal to JSON: %v", err)
	}

	if err := ValidateImageJSON(dataJSON); err != nil {
		t.Errorf("minimal template should be valid, but got: %v", err)
	}
}

func TestImageTemplateMissingFields(t *testing.T) {
	// Test template missing required fields
	invalidTemplate := `image:
  name: test-image

target:
  os: azure-linux
  dist: azl3
  arch: x86_64

systemConfigs:
  - name: incomplete
    packages:
      - openssh-server
`

	var raw interface{}
	if err := yaml.Unmarshal([]byte(invalidTemplate), &raw); err != nil {
		t.Fatalf("failed to parse invalid template: %v", err)
	}

	dataJSON, err := json.Marshal(raw)
	if err != nil {
		t.Fatalf("failed to marshal to JSON: %v", err)
	}

	if err := ValidateImageJSON(dataJSON); err == nil {
		t.Errorf("incomplete template should fail validation")
	}
}

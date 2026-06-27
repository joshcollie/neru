package config_test

import (
	"testing"

	"github.com/y3owk1n/neru/internal/config"
)

const (
	dirNormal     = "normal"
	dirReverse    = "reverse"
	bundleExample = "com.example"
)

func TestValidateHints_EnabledRequiresClickableRoles(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Hints.Enabled = true
	cfg.Hints.ClickableRoles = nil

	err := cfg.ValidateHints()
	if err == nil {
		t.Fatal("ValidateHints() expected error when enabled and clickable_roles is empty")
	}
}

func TestValidateHints_BoundaryHighlightGeometry(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Hints.BoundaryHighlight.BorderWidth = -1

	err := cfg.ValidateHints()
	if err == nil {
		t.Fatal("ValidateHints() expected error for negative boundary border width")
	}

	cfg = config.DefaultConfig()
	cfg.Hints.BoundaryHighlight.BorderRadius = -1

	err = cfg.ValidateHints()
	if err != nil {
		t.Fatalf("ValidateHints() expected no error for -1 (auto) border radius, got %v", err)
	}
}

func TestValidateHints_UIPlacement(t *testing.T) {
	validPlacements := []string{
		"top",
		"center",
		"bottom",
	}

	for _, placement := range validPlacements {
		cfg := config.DefaultConfig()
		cfg.Hints.UI.Placement = placement

		err := cfg.ValidateHints()
		if err != nil {
			t.Fatalf("ValidateHints() expected placement %q to be valid, got %v", placement, err)
		}
	}

	cfg := config.DefaultConfig()

	cfg.Hints.UI.Placement = "floating"

	err := cfg.ValidateHints()
	if err == nil {
		t.Fatal("ValidateHints() expected error for invalid hints.ui.placement")
	}
}

func TestValidateHints_PositiveUnitFloat(t *testing.T) {
	// merge_iou_threshold cannot be 0
	cfg := config.DefaultConfig()
	cfg.Hints.Vision.MergeIOUThreshold = 0.0

	err := cfg.ValidateHints()
	if err == nil {
		t.Fatal("ValidateHints() expected error for 0.0 merge_iou_threshold")
	}

	// button_min_confidence cannot be 0
	cfg = config.DefaultConfig()
	cfg.Hints.Vision.ButtonMinConfidence = 0.0

	err = cfg.ValidateHints()
	if err == nil {
		t.Fatal("ValidateHints() expected error for 0.0 button_min_confidence")
	}

	// generic_clickable_min_confidence cannot be 0
	cfg = config.DefaultConfig()
	cfg.Hints.Vision.GenericClickableMinConfidence = 0.0

	err = cfg.ValidateHints()
	if err == nil {
		t.Fatal("ValidateHints() expected error for 0.0 generic_clickable_min_confidence")
	}

	// but minimum_confidence CAN be 0
	cfg = config.DefaultConfig()
	cfg.Hints.Vision.MinimumConfidence = 0.0

	err = cfg.ValidateHints()
	if err != nil {
		t.Fatalf("ValidateHints() expected no error for 0.0 minimum_confidence, got %v", err)
	}
}

func TestValidateHints_LabelDirection(t *testing.T) {
	tests := []struct {
		name      string
		direction string
		wantErr   bool
	}{
		{name: "reverse is valid", direction: dirReverse, wantErr: false},
		{name: "normal is valid", direction: dirNormal, wantErr: false},
		{name: "empty defaults to normal (no error)", direction: "", wantErr: false},
		{name: "unknown value is rejected", direction: "sideways", wantErr: true},
		{name: "uppercase is rejected (case sensitive)", direction: "NORMAL", wantErr: true},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			cfg := config.DefaultConfig()
			cfg.Hints.LabelDirection = testCase.direction

			err := cfg.ValidateHints()
			if testCase.wantErr && err == nil {
				t.Fatalf(
					"ValidateHints() expected error for label_direction=%q, got nil",
					testCase.direction,
				)
			}

			if !testCase.wantErr && err != nil {
				t.Fatalf(
					"ValidateHints() unexpected error for label_direction=%q: %v",
					testCase.direction,
					err,
				)
			}
		})
	}
}

func TestValidateAppConfigs_LabelDirection(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Hints.AppConfigs = []config.AppConfig{
		{
			BundleID:       bundleExample,
			LabelDirection: dirNormal,
		},
	}

	err := cfg.ValidateAppConfigs()
	if err != nil {
		t.Fatalf("ValidateAppConfigs() unexpected error for valid app config: %v", err)
	}

	cfg = config.DefaultConfig()
	cfg.Hints.AppConfigs = []config.AppConfig{
		{
			BundleID:       bundleExample,
			LabelDirection: "diagonal",
		},
	}

	err = cfg.ValidateAppConfigs()
	if err == nil {
		t.Fatal("ValidateAppConfigs() expected error for invalid app label_direction")
	}
}

func TestLabelDirectionForApp(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Hints.LabelDirection = dirReverse
	cfg.Hints.AppConfigs = []config.AppConfig{
		{
			BundleID:       bundleExample,
			LabelDirection: dirNormal,
		},
	}

	// App override takes precedence.
	if got := cfg.Hints.LabelDirectionForApp(bundleExample); got != dirNormal {
		t.Errorf("LabelDirectionForApp(app with override) = %q, want %q", got, dirNormal)
	}

	// Fallback to global config.
	if got := cfg.Hints.LabelDirectionForApp("com.other"); got != dirReverse {
		t.Errorf("LabelDirectionForApp(app without override) = %q, want %q", got, dirReverse)
	}

	// Empty global value normalizes to the default (normal).
	cfg = config.DefaultConfig()
	cfg.Hints.LabelDirection = ""

	if got := cfg.Hints.LabelDirectionForApp(bundleExample); got != dirNormal {
		t.Errorf("LabelDirectionForApp(empty global) = %q, want %q", got, dirNormal)
	}
}

func TestValidateHints_FontSizeFloor(t *testing.T) {
	// The floor was lowered from 6 to 1 for both hints.ui.font_size and
	// hints.search_input_ui.font_size, so 1 must now validate.
	cfg := config.DefaultConfig()
	cfg.Hints.UI.FontSize = 1
	cfg.Hints.SearchInputUI.FontSize = 1

	err := cfg.ValidateHints()
	if err != nil {
		t.Fatalf("ValidateHints() expected font_size 1 to be valid, got %v", err)
	}

	// 0 is still below the floor for hints.ui.font_size.
	cfg = config.DefaultConfig()
	cfg.Hints.UI.FontSize = 0

	err = cfg.ValidateHints()
	if err == nil {
		t.Fatal("ValidateHints() expected error for hints.ui.font_size of 0")
	}

	// hints.search_input_ui.font_size is validated independently.
	cfg = config.DefaultConfig()
	cfg.Hints.SearchInputUI.FontSize = 0

	err = cfg.ValidateHints()
	if err == nil {
		t.Fatal("ValidateHints() expected error for hints.search_input_ui.font_size of 0")
	}

	// The upper bound of 72 is unchanged.
	cfg = config.DefaultConfig()
	cfg.Hints.UI.FontSize = 73

	err = cfg.ValidateHints()
	if err == nil {
		t.Fatal("ValidateHints() expected error for hints.ui.font_size above 72")
	}
}

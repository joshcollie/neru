package config_test

import (
	"testing"

	"github.com/y3owk1n/neru/internal/config"
)

func TestValidateGrid_FontSizeFloor(t *testing.T) {
	// The floor was lowered from 6 to 1 for grid.ui.font_size, matching the
	// existing recursive_grid.ui.font_size minimum.
	cfg := config.DefaultConfig()
	cfg.Grid.UI.FontSize = 1

	err := cfg.ValidateGrid()
	if err != nil {
		t.Fatalf("ValidateGrid() expected font_size 1 to be valid, got %v", err)
	}

	// 0 is still below the floor.
	cfg = config.DefaultConfig()
	cfg.Grid.UI.FontSize = 0

	err = cfg.ValidateGrid()
	if err == nil {
		t.Fatal("ValidateGrid() expected error for grid.ui.font_size of 0")
	}

	// The upper bound of 72 is unchanged.
	cfg = config.DefaultConfig()
	cfg.Grid.UI.FontSize = 73

	err = cfg.ValidateGrid()
	if err == nil {
		t.Fatal("ValidateGrid() expected error for grid.ui.font_size above 72")
	}
}

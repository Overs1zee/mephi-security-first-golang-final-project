package daysteps

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseWhitespace(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"пробел в начале строки", " 678,0h50m"},
		{"пробел перед запятой", "678 ,0h50m"},
		{"пробел после запятой", "678, 0h50m"},
		{"пробел в конце строки", "678,0h50m "},
		{"пробелы вокруг запятой", "678 , 0h50m"},
		{"табуляция перед шагами", "\t678,0h50m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &DaySteps{}

			require.Error(t, ds.Parse(tt.input), "строка %q не должна парситься", tt.input)
			assert.Equal(t, DaySteps{}, *ds)
		})
	}
}

func TestParseDoesNotChangeStateOnError(t *testing.T) {
	ds := &DaySteps{}
	require.NoError(t, ds.Parse("678,0h50m"))

	want := *ds

	invalid := []string{
		"",
		"678",
		"678,0h50m,extra",
		"0,0h50m",
		"-678,0h50m",
		"abc,0h50m",
		"678,0h0m",
		"678,-50m",
		"678,invalid",
	}

	for _, input := range invalid {
		t.Run(input, func(t *testing.T) {
			require.Error(t, ds.Parse(input), "строка %q не должна парситься", input)
			assert.Equal(t, want, *ds, "при ошибке парсинга структура не должна меняться")
		})
	}
}

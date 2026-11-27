package spentenergy

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMeanSpeedInvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		steps    int
		height   float64
		duration time.Duration
	}{
		{
			name:     "отрицательные шаги",
			steps:    -1000,
			height:   1.75,
			duration: time.Hour,
		},
		{
			name:     "нулевые шаги и нулевая продолжительность",
			steps:    0,
			height:   1.75,
			duration: 0,
		},
		{
			name:     "нулевой рост",
			steps:    1000,
			height:   0,
			duration: time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Zero(t, MeanSpeed(tt.steps, tt.height, tt.duration))
		})
	}
}

func TestSpentCaloriesInvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		steps    int
		weight   float64
		height   float64
		duration time.Duration
	}{
		{
			name:     "нулевая продолжительность",
			steps:    1000,
			weight:   75.0,
			height:   1.75,
			duration: 0,
		},
		{
			name:     "нулевой вес",
			steps:    1000,
			weight:   0,
			height:   1.75,
			duration: time.Hour,
		},
		{
			name:     "нулевой рост",
			steps:    1000,
			weight:   75.0,
			height:   0,
			duration: time.Hour,
		},
		{
			name:     "отрицательные шаги",
			steps:    -1000,
			weight:   75.0,
			height:   1.75,
			duration: time.Hour,
		},
	}

	functions := map[string]func(int, float64, float64, time.Duration) (float64, error){
		"бег":    RunningSpentCalories,
		"ходьба": WalkingSpentCalories,
	}

	for _, tt := range tests {
		for name, fn := range functions {
			t.Run(tt.name+"/"+name, func(t *testing.T) {
				calories, err := fn(tt.steps, tt.weight, tt.height, tt.duration)

				require.Error(t, err, "ожидалась ошибка для случая %q", tt.name)
				assert.NotEmpty(t, err.Error(), "ошибка должна быть информативной")
				assert.Zero(t, calories, "при ошибке возвращается 0")
			})
		}
	}
}

func TestWalkingIsHalfOfRunning(t *testing.T) {
	const (
		steps    = 6000
		weight   = 75.0
		height   = 1.75
		duration = time.Hour
	)

	running, err := RunningSpentCalories(steps, weight, height, duration)
	require.NoError(t, err)

	walking, err := WalkingSpentCalories(steps, weight, height, duration)
	require.NoError(t, err)

	assert.InDelta(t, running*walkingCaloriesCoefficient, walking, 1e-9)
}

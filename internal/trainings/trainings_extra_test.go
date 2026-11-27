package trainings

import (
	"testing"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActionInfoCaloriesError(t *testing.T) {
	personal := personaldata.Personal{Weight: 75.0, Height: 1.75}

	tests := []struct {
		name     string
		training Training
	}{
		{
			name: "нулевой вес",
			training: Training{
				Steps:        1000,
				TrainingType: TrainingTypeRunning,
				Duration:     time.Hour,
				Personal:     personaldata.Personal{Weight: 0, Height: personal.Height},
			},
		},
		{
			name: "отрицательный вес",
			training: Training{
				Steps:        1000,
				TrainingType: TrainingTypeWalking,
				Duration:     time.Hour,
				Personal:     personaldata.Personal{Weight: -75.0, Height: personal.Height},
			},
		},
		{
			name: "нулевой рост",
			training: Training{
				Steps:        1000,
				TrainingType: TrainingTypeWalking,
				Duration:     time.Hour,
				Personal:     personaldata.Personal{Weight: personal.Weight, Height: 0},
			},
		},
		{
			name: "нулевая продолжительность",
			training: Training{
				Steps:        1000,
				TrainingType: TrainingTypeRunning,
				Duration:     0,
				Personal:     personal,
			},
		},
		{
			name: "отрицательная продолжительность",
			training: Training{
				Steps:        1000,
				TrainingType: TrainingTypeRunning,
				Duration:     -time.Hour,
				Personal:     personal,
			},
		},
		{
			name: "пустой тип тренировки",
			training: Training{
				Steps:        1000,
				TrainingType: "",
				Duration:     time.Hour,
				Personal:     personal,
			},
		},
		{
			name: "тип тренировки в нижнем регистре",
			training: Training{
				Steps:        1000,
				TrainingType: "бег",
				Duration:     time.Hour,
				Personal:     personal,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.training.ActionInfo()

			require.Error(t, err)
			assert.Empty(t, got)
		})
	}
}

func TestParseDoesNotChangeStateOnError(t *testing.T) {
	training := &Training{}
	require.NoError(t, training.Parse("1000,Бег,1h"))

	want := *training

	invalid := []string{
		"",
		"1000,Бег",
		"1000,Бег,1h,extra",
		"0,Бег,1h",
		"-1000,Бег,1h",
		"abc,Бег,1h",
		"1000,Бег,0h",
		"1000,Бег,-1h",
		"1000,Бег,invalid",
	}

	for _, input := range invalid {
		t.Run(input, func(t *testing.T) {
			require.Error(t, training.Parse(input), "строка %q не должна парситься", input)
			assert.Equal(t, want, *training, "при ошибке парсинга структура не должна меняться")
		})
	}
}

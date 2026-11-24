package trainings

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

// Поддерживаемые типы тренировок.
const (
	TrainingTypeRunning = "Бег"
	TrainingTypeWalking = "Ходьба"
)

// Training хранит данные одной тренировки.
type Training struct {
	Steps        int           // количество шагов.
	TrainingType string        // тип тренировки.
	Duration     time.Duration // длительность тренировки.
	personaldata.Personal
}

// Parse разбирает строку данных в формате "<шаги>,<тип тренировки>,<длительность>".
func (t *Training) Parse(datastring string) (err error) {
	parts := strings.Split(datastring, ",")
	if len(parts) != 3 {
		return fmt.Errorf("ожидается 3 поля, получено %d", len(parts))
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("некорректное количество шагов %q: %w", parts[0], err)
	}
	if steps <= 0 {
		return fmt.Errorf("количество шагов должно быть положительным, получено %d", steps)
	}

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return fmt.Errorf("некорректная длительность %q: %w", parts[2], err)
	}
	if duration <= 0 {
		return fmt.Errorf("длительность должна быть положительной, получено %v", duration)
	}

	t.Steps = steps
	t.TrainingType = parts[1]
	t.Duration = duration

	return nil
}

// ActionInfo возвращает информацию о тренировке.
func (t Training) ActionInfo() (string, error) {
	dist := spentenergy.Distance(t.Steps, t.Height)
	speed := spentenergy.MeanSpeed(t.Steps, t.Height, t.Duration)

	var (
		calories float64
		err      error
	)

	switch t.TrainingType {
	case TrainingTypeRunning:
		calories, err = spentenergy.RunningSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	case TrainingTypeWalking:
		calories, err = spentenergy.WalkingSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		t.TrainingType, t.Duration.Hours(), dist, speed, calories,
	), nil
}

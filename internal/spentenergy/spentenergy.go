package spentenergy

import (
	"fmt"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе.
)

// WalkingSpentCalories возвращает количество калорий, сожжённых при ходьбе.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := checkInputParams(steps, weight, height, duration); err != nil {
		return 0, err
	}

	calories, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}

	return calories * walkingCaloriesCoefficient, nil
}

// RunningSpentCalories возвращает количество калорий, сожжённых при беге.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := checkInputParams(steps, weight, height, duration); err != nil {
		return 0, err
	}

	meanSpeed := MeanSpeed(steps, height, duration)

	return (weight * meanSpeed * duration.Minutes()) / minInH, nil
}

// MeanSpeed возвращает среднюю скорость передвижения в км/ч.
func MeanSpeed(steps int, height float64, duration time.Duration) float64 {
	if steps <= 0 || duration <= 0 {
		return 0
	}

	return Distance(steps, height) / duration.Hours()
}

// Distance возвращает пройденную дистанцию в километрах.
func Distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient

	return (float64(steps) * stepLength) / mInKm
}

// checkInputParams проверяет входные параметры для расчета калорий.
func checkInputParams(steps int, weight, height float64, duration time.Duration) error {
	switch {
	case steps <= 0:
		return fmt.Errorf("количество шагов должно быть положительным, получено %d", steps)
	case weight <= 0:
		return fmt.Errorf("вес должен быть положительным, получено %.2f", weight)
	case height <= 0:
		return fmt.Errorf("рост должен быть положительным, получено %.2f", height)
	case duration <= 0:
		return fmt.Errorf("продолжительность должна быть положительной, получено %v", duration)
	default:
		return nil
	}
}

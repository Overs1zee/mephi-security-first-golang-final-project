package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

// DaySteps хранит данные о дневной активности (количестве шагов).
type DaySteps struct {
	Steps    int           // количество шагов.
	Duration time.Duration // продолжительность прогулки.
	personaldata.Personal
}

// Parse разбирает строку данных в формате "<шаги>,<длительность>".
func (ds *DaySteps) Parse(datastring string) (err error) {
	parts := strings.Split(datastring, ",")
	if len(parts) != 2 {
		return fmt.Errorf("ожидается 2 поля, получено %d", len(parts))
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("некорректное количество шагов %q: %w", parts[0], err)
	}
	if steps <= 0 {
		return fmt.Errorf("количество шагов должно быть положительным, получено %d", steps)
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return fmt.Errorf("некорректная длительность %q: %w", parts[1], err)
	}
	if duration <= 0 {
		return fmt.Errorf("длительность должна быть положительной, получено %v", duration)
	}

	ds.Steps = steps
	ds.Duration = duration

	return nil
}

// ActionInfo возвращает информацию о дневной активности.
func (ds DaySteps) ActionInfo() (string, error) {
	dist := spentenergy.Distance(ds.Steps, ds.Height)

	calories, err := spentenergy.WalkingSpentCalories(ds.Steps, ds.Weight, ds.Height, ds.Duration)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		ds.Steps, dist, calories,
	), nil
}

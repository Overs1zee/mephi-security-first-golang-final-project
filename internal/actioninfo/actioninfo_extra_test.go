package actioninfo

import (
	"bytes"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// captureStdout выполняет fn и возвращает всё, что она вывела в os.Stdout.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	require.NoError(t, err, "Не удалось создать pipe")

	old := os.Stdout
	os.Stdout = w

	done := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r)
		done <- buf.String()
	}()

	fn()

	os.Stdout = old
	require.NoError(t, w.Close())

	return <-done
}

func TestInfoSkipsInvalidEntries(t *testing.T) {
	var logBuf bytes.Buffer
	oldLog := log.Writer()
	log.SetOutput(&logBuf)
	defer log.SetOutput(oldLog)

	errParse := errors.New("parse failed")
	errInfo := errors.New("info failed")

	parser := new(MockDataParser)
	parser.On("Parse", "ok-1").Return(nil).Once()
	parser.On("ActionInfo").Return("info 1", nil).Once()
	parser.On("Parse", "bad-parse").Return(errParse).Once()
	parser.On("Parse", "bad-info").Return(nil).Once()
	parser.On("ActionInfo").Return("", errInfo).Once()
	parser.On("Parse", "ok-2").Return(nil).Once()
	parser.On("ActionInfo").Return("info 2", nil).Once()

	got := captureStdout(t, func() {
		Info([]string{"ok-1", "bad-parse", "bad-info", "ok-2"}, parser)
	})

	parser.AssertExpectations(t)

	assert.Equal(t, "info 1\ninfo 2\n", got, "в stdout должны попасть только успешно обработанные записи")
	assert.Contains(t, logBuf.String(), "parse failed", "ошибка парсинга должна попадать в лог")
	assert.Contains(t, logBuf.String(), "info failed", "ошибка ActionInfo должна попадать в лог")
}

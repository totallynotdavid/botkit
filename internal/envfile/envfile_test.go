package envfile_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/totallynotdavid/botkit/internal/config"
	"github.com/totallynotdavid/botkit/internal/envfile"
	"github.com/totallynotdavid/botkit/internal/flow"
)

func TestASettingWithoutAnExplanationIsRefused(t *testing.T) {
	t.Parallel()

	bots := envfile.Bots()
	bots[0].Entries = append(bots[0].Entries, config.Entry{Key: "LATEX_NEW_SETTING", Default: "1"})

	_, err := envfile.Example(bots)
	if !errors.Is(err, envfile.ErrUndocumented) || !strings.Contains(err.Error(), "LATEX_NEW_SETTING") {
		t.Errorf("Example() error = %v, want ErrUndocumented naming LATEX_NEW_SETTING", err)
	}
}

func TestAnExplanationOfASettingNoBotReadsIsRefused(t *testing.T) {
	t.Parallel()

	_, err := envfile.Example(envfile.Bots()[:1])
	if !errors.Is(err, envfile.ErrUnread) || !strings.Contains(err.Error(), flow.KeyFile) {
		t.Errorf("Example() error = %v, want ErrUnread naming %s", err, flow.KeyFile)
	}
}

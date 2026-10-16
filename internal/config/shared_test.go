package config_test

import (
	"errors"
	"log/slog"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/totallynotdavid/botkit/internal/auth"
	"github.com/totallynotdavid/botkit/internal/config"
)

func TestSharedOverridesEachType(t *testing.T) {
	t.Parallel()

	env := fromMap(map[string]string{
		config.KeyStore:             "/var/lib/bot/state.db",
		config.KeyLogLevel:          "debug",
		config.KeyOwners:            "51900000001@s.whatsapp.net, 51900000002@s.whatsapp.net",
		config.KeyAllowOnly:         "51900000003@s.whatsapp.net",
		config.KeyRateLimitRequests: "20",
		config.KeyRateLimitPeriod:   "30s",
		config.KeyRateLimitCooldown: "2m",
		config.KeyMaxInFlight:       "3",
		config.KeyOwnMessages:       "true",
	})

	got := env.Shared(config.DefaultShared())

	want := config.Shared{
		Store:       "/var/lib/bot/state.db",
		LogLevel:    slog.LevelDebug,
		Owners:      []string{"51900000001@s.whatsapp.net", "51900000002@s.whatsapp.net"},
		AllowOnly:   []string{"51900000003@s.whatsapp.net"},
		RateLimit:   config.RateLimit{Requests: 20, Period: 30 * time.Second, Cooldown: 2 * time.Minute},
		MaxInFlight: 3,
		OwnMessages: true,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Shared() = %#v, want %#v", got, want)
	}

	err := env.Err()
	if err != nil {
		t.Errorf("Err() = %v, want nil", err)
	}
}

func TestSharedKeepsCallerDefaults(t *testing.T) {
	t.Parallel()

	def := config.DefaultShared()
	def.RateLimit.Requests = 60

	got := fromMap(map[string]string{config.KeyMaxInFlight: "2"}).Shared(def)

	if got.RateLimit.Requests != 60 || got.MaxInFlight != 2 {
		t.Errorf("Shared() requests = %d, in flight = %d, want 60 and 2", got.RateLimit.Requests, got.MaxInFlight)
	}
}

func TestSharedBadValuesNameTheirKey(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		key, value string
		wrapped    error
	}{
		{config.KeyLogLevel, "loud", nil},
		{config.KeyRateLimitRequests, "0", config.ErrBelowMinimum},
		{config.KeyRateLimitPeriod, "soon", nil},
		{config.KeyMaxInFlight, "many", nil},
		{config.KeyOwnMessages, "maybe", nil},
		{config.KeyOwners, "user@server:3", auth.ErrInvalidJID},
		{config.KeyAllowOnly, "user@server:3", auth.ErrInvalidJID},
	} {
		env := fromMap(map[string]string{test.key: test.value})
		env.Shared(config.DefaultShared())

		err := env.Err()
		if err == nil || !strings.Contains(err.Error(), test.key+"=") {
			t.Errorf("%s=%q: Err() = %v, want an error naming the key", test.key, test.value, err)

			continue
		}

		if test.wrapped != nil && !errors.Is(err, test.wrapped) {
			t.Errorf("%s=%q: Err() = %v, want it to wrap %v", test.key, test.value, err, test.wrapped)
		}
	}
}

func TestSharedEntriesCarryTheDefaultsTheCodeUses(t *testing.T) {
	t.Parallel()

	def := config.DefaultShared()
	def.Owners = []string{"51900000007@s.whatsapp.net", "51900000008@s.whatsapp.net"}
	def.RateLimit = config.RateLimit{Requests: 20, Period: 90 * time.Second, Cooldown: time.Minute}

	env := fromMap(nil)
	env.Shared(def)

	want := map[string]string{
		config.KeyStore:             "botkit.db",
		config.KeyLogLevel:          "info",
		config.KeyOwners:            "51900000007@s.whatsapp.net,51900000008@s.whatsapp.net",
		config.KeyAllowOnly:         "",
		config.KeyRateLimitRequests: "20",
		config.KeyRateLimitPeriod:   "1m30s",
		config.KeyRateLimitCooldown: "1m",
		config.KeyMaxInFlight:       "10",
		config.KeyOwnMessages:       "false",
	}

	got := map[string]string{}
	for _, entry := range env.Entries() {
		got[entry.Key] = entry.Default
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Entries() = %v, want %v", got, want)
	}
}

package config

import (
	"context"
	"errors"
	"maps"
	"os"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"testing"

	envconfig "github.com/sethvargo/go-envconfig"
)

// envVars returns the env var name of every config field, keyed by field path
// (e.g. "Steam.SleepInterval").
func envVars() map[string]string {
	vars := map[string]string{}
	var walk func(typ reflect.Type, path []string)
	walk = func(typ reflect.Type, path []string) {
		for i := range typ.NumField() {
			field := typ.Field(i)
			fieldPath := append(slices.Clone(path), field.Name)
			tag, hasTag := field.Tag.Lookup("env")
			if !hasTag && field.Type.Kind() == reflect.Struct {
				walk(field.Type, fieldPath)
				continue
			}
			name, _, _ := strings.Cut(tag, ",")
			vars[strings.Join(fieldPath, ".")] = name
		}
	}
	walk(reflect.TypeOf(config{}), nil)
	return vars
}

// Every env var must be named after its field's path (Steam.SleepInterval ->
// STEAMSALTY_STEAM_SLEEPINTERVAL), so code, env and docs use the same names.
func TestEnvNamesMirrorFieldPaths(t *testing.T) {
	for path, name := range envVars() {
		want := "STEAMSALTY_" + strings.ToUpper(strings.ReplaceAll(path, ".", "_"))
		if name != want {
			t.Errorf("%s: env name is %q, want %q", path, name, want)
		}
	}
}

// .env.example and the README must mention every env var, and no unknown ones
// (such as old or misspelled names).
func TestDocsListAllEnvVars(t *testing.T) {
	known := map[string]bool{}
	for _, name := range envVars() {
		known[name] = true
	}
	for _, file := range []string{"../.env.example", "../README.md"} {
		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		mentioned := map[string]bool{}
		for _, name := range regexp.MustCompile(`STEAMSALTY_[A-Z0-9_]+`).FindAllString(string(content), -1) {
			mentioned[name] = true
		}
		for _, name := range slices.Sorted(maps.Keys(known)) {
			if !mentioned[name] {
				t.Errorf("%s does not mention %s", file, name)
			}
		}
		for _, name := range slices.Sorted(maps.Keys(mentioned)) {
			if !known[name] {
				t.Errorf("%s mentions unknown variable %s", file, name)
			}
		}
	}
}

// requiredEnv is the minimal environment the bot starts with.
func requiredEnv() map[string]string {
	return map[string]string{
		"STEAMSALTY_TELEGRAM_APITOKEN":     "1234567890:TESTTOKEN",
		"STEAMSALTY_TELEGRAM_CHATID":       "-1001234567890",
		"STEAMSALTY_STEAM_APIKEY":          "TESTAPIKEY",
		"STEAMSALTY_STEAM_WATCHEDPROFILES": "76561198012345678,76561198087654321",
	}
}

func process(env map[string]string) (config, error) {
	var c config
	err := envconfig.ProcessWith(context.Background(), &envconfig.Config{
		Target:   &c,
		Lookuper: envconfig.MapLookuper(env),
	})
	return c, err
}

func TestProcess_RequiredValuesAndDefaults(t *testing.T) {
	got, err := process(requiredEnv())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var want config
	want.LogLevel = 1
	want.Telegram.APIToken = "1234567890:TESTTOKEN"
	want.Telegram.ChatID = -1001234567890
	want.Steam.APIKey = "TESTAPIKEY"
	want.Steam.WatchedProfiles = []uint64{76561198012345678, 76561198087654321}
	want.Steam.SleepInterval = 60
	want.DeepL.TargetLanguage = "EN-US"
	want.DeepL.FreeTier = true
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got  %+v\nwant %+v", got, want)
	}
}

func TestProcess_MissingRequiredValue(t *testing.T) {
	for _, name := range slices.Sorted(maps.Keys(requiredEnv())) {
		t.Run(name, func(t *testing.T) {
			env := requiredEnv()
			delete(env, name)
			_, err := process(env)
			if !errors.Is(err, envconfig.ErrMissingRequired) || !strings.Contains(err.Error(), name) {
				t.Errorf("err = %v, want a missing required value error for %s", err, name)
			}
		})
	}
}

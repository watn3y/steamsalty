package config

import (
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

// defaultConfig is the config that requiredEnv loads into.
func defaultConfig() config {
	var c config
	c.LogLevel = 1
	c.Telegram.APIToken = "1234567890:TESTTOKEN"
	c.Telegram.ChatID = -1001234567890
	c.Steam.APIKey = "TESTAPIKEY"
	c.Steam.WatchedProfiles = []uint64{76561198012345678, 76561198087654321}
	c.Steam.SleepInterval = 60
	c.DeepL.TargetLanguage = "EN-US"
	c.DeepL.FreeTier = true
	return c
}

func TestLoad_RequiredValuesAndDefaults(t *testing.T) {
	got, err := load(envconfig.MapLookuper(requiredEnv()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := defaultConfig(); !reflect.DeepEqual(got, want) {
		t.Errorf("got  %+v\nwant %+v", got, want)
	}
}

// An empty optional variable (KEY=) must behave like an unset one. Before, it
// replaced the default with 0, false or "", e.g. a sleep interval of 0.
func TestLoad_EmptyValuesUseDefaults(t *testing.T) {
	env := requiredEnv()
	for _, name := range envVars() {
		if _, required := env[name]; !required {
			env[name] = ""
		}
	}
	got, err := load(envconfig.MapLookuper(env))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := defaultConfig(); !reflect.DeepEqual(got, want) {
		t.Errorf("got  %+v\nwant %+v", got, want)
	}
}

// Explicit values must still win over the defaults, including 0 and False.
func TestLoad_ExplicitValuesOverrideDefaults(t *testing.T) {
	env := requiredEnv()
	env["STEAMSALTY_LOGLEVEL"] = "0"
	env["STEAMSALTY_STEAM_SLEEPINTERVAL"] = "1"
	env["STEAMSALTY_DEEPL_ENABLED"] = "True"
	env["STEAMSALTY_DEEPL_TARGETLANGUAGE"] = "de"
	env["STEAMSALTY_DEEPL_APIKEY"] = "a1b2c3d4-56e7-89f0-a1b2-c3d4e5f6a7b8:fx"
	env["STEAMSALTY_DEEPL_FREETIER"] = "False"
	got, err := load(envconfig.MapLookuper(env))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := defaultConfig()
	want.LogLevel = 0
	want.Steam.SleepInterval = 1
	want.DeepL.Enabled = true
	want.DeepL.TargetLanguage = "DE"
	want.DeepL.APIKey = "a1b2c3d4-56e7-89f0-a1b2-c3d4e5f6a7b8:fx"
	want.DeepL.FreeTier = false
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got  %+v\nwant %+v", got, want)
	}
}

// time.Sleep returns immediately for 0 or less, so the watchers would poll
// Steam nonstop.
func TestLoad_RejectsSleepIntervalBelowOne(t *testing.T) {
	for _, value := range []string{"0", "-60"} {
		t.Run(value, func(t *testing.T) {
			env := requiredEnv()
			env["STEAMSALTY_STEAM_SLEEPINTERVAL"] = value
			_, err := load(envconfig.MapLookuper(env))
			if err == nil || !strings.Contains(err.Error(), "STEAMSALTY_STEAM_SLEEPINTERVAL") {
				t.Errorf("err = %v, want an error about STEAMSALTY_STEAM_SLEEPINTERVAL", err)
			}
		})
	}
}

// A required variable that is unset or empty must stop the bot with an error
// naming it.
func TestLoad_MissingRequiredValue(t *testing.T) {
	for _, name := range slices.Sorted(maps.Keys(requiredEnv())) {
		t.Run(name, func(t *testing.T) {
			for _, variant := range []string{"unset", "empty"} {
				t.Run(variant, func(t *testing.T) {
					env := requiredEnv()
					if variant == "unset" {
						delete(env, name)
					} else {
						env[name] = ""
					}
					_, err := load(envconfig.MapLookuper(env))
					if !errors.Is(err, envconfig.ErrMissingRequired) || !strings.Contains(err.Error(), name) {
						t.Errorf("err = %v, want a missing required value error for %s", err, name)
					}
				})
			}
		})
	}
}

package configapp

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	configserver "goxus/src/internal/app/goxus/api/server/config"
	configcron "goxus/src/internal/app/goxus/cron-job/config"

	"goxus/src/internal/app/goxus/cli"

	pgxdb "github.com/nobuenhombre/suikat/pkg/db/connectors/postgres-pgx-db"
	"github.com/nobuenhombre/suikat/pkg/fico"
)

type testConfig struct {
	fileName    string
	fileContent string
	config      *Config
	err         error
}

func getTestConfig() *Config {
	return &Config{
		Hosts: HostsConfig{
			API: configserver.HTTPServerConfig{
				Host: "127.0.0.1",
				Port: "8080",
			},
		},
		Cron: configcron.CronConfig{
			ExampleJob: configcron.CronJobConfig{
				Enabled:  true,
				Schedule: "@every 10m",
			},
		},
	}
}

func getTestConfigWithDB() *Config {
	return &Config{
		DB: pgxdb.Config{
			Host:     "127.0.0.1",
			Port:     "5432",
			Name:     "goxus",
			User:     "goxus",
			Password: "12345",
			SSLMode:  "disable",
		},
		Hosts: HostsConfig{
			API: configserver.HTTPServerConfig{
				Host: "127.0.0.1",
				Port: "8080",
			},
		},
		Cron: configcron.CronConfig{
			ExampleJob: configcron.CronJobConfig{
				Enabled:  true,
				Schedule: "@every 10m",
			},
		},
	}
}

func TestConfigLoad(t *testing.T) {
	tests := []*testConfig{
		{
			fileName:    "config-app_test_load.yaml",
			fileContent: "",
			config:      getTestConfig(),
			err:         nil,
		},
		{
			fileName:    "config-app_test_load_db.yaml",
			fileContent: "",
			config:      getTestConfigWithDB(),
			err:         nil,
		},
	}

	for _, test := range tests {
		t.Run(test.fileName, func(t *testing.T) {
			cfg := new(Config)
			err := cfg.Load(test.fileName)

			if !(reflect.DeepEqual(cfg, test.config) && errors.Is(err, test.err)) {
				t.Errorf(
					"cfg.Load(%#v),\n Expected (cfg = %#v, err = %#v),\n Actual (cfg = %#v, err = %#v).\n",
					test.fileName, test.config, test.err, cfg, err,
				)
			}
		})
	}
}

func TestConfigSave(t *testing.T) {
	tests := []*testConfig{
		{
			fileName: "config-app_test_save.yaml",
			fileContent: "" +
				"hosts:\n" +
				"    api:\n" +
				"        host: 127.0.0.1\n" +
				"        post: \"8080\"\n" +
				"cron:\n" +
				"    example_job:\n" +
				"        enabled: true\n" +
				"        schedule: '@every 10m'\n",
			config: getTestConfig(),
			err:    nil,
		},
		{
			fileName: "config-app_test_save_db.yaml",
			fileContent: "" +
				"db:\n" +
				"    host: 127.0.0.1\n" +
				"    port: \"5432\"\n" +
				"    name: goxus\n" +
				"    user: goxus\n" +
				"    password: \"12345\"\n" +
				"    sslmode: disable\n" +
				"    binaryparameters: \"\"\n" +
				"    statementcachemode: \"\"\n" +
				"    maxconnections: \"\"\n" +
				"    useconnectionpooler: false\n" +
				"hosts:\n" +
				"    api:\n" +
				"        host: 127.0.0.1\n" +
				"        post: \"8080\"\n" +
				"cron:\n" +
				"    example_job:\n" +
				"        enabled: true\n" +
				"        schedule: '@every 10m'\n",
			config: getTestConfigWithDB(),
			err:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.fileName, func(t *testing.T) {
			cfg := test.config
			err := cfg.Save(test.fileName)

			txtConfigFile := fico.TxtFile(test.fileName)
			fileContent, errReadFile := txtConfigFile.Read()

			if errReadFile != nil {
				t.Errorf(
					"txtConfigFile.Read error %#v",
					errReadFile,
				)
			}

			if !(reflect.DeepEqual(fileContent, test.fileContent) && errors.Is(err, test.err)) {
				t.Errorf(
					"cfg.Save(%#v),\n Expected (fileContent = %#v, err = %#v),\n Actual__ (fileContent = %#v, err = %#v).\n",
					test.fileName, test.fileContent, test.err, fileContent, err,
				)
			}
		})
	}
}

// --- New tests below ---

func TestConfigNew(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		dir := t.TempDir()
		cfgPath := filepath.Join(dir, "config.yaml")

		cfg := getTestConfig()
		if err := cfg.Save(cfgPath); err != nil {
			t.Fatalf("Save failed: %v", err)
		}

		svc, err := New(cfgPath)
		if err != nil {
			t.Errorf("New() returned error: %v", err)
		}
		if svc == nil {
			t.Error("New() returned nil Service")
		}
	})

	t.Run("file_not_found", func(t *testing.T) {
		_, err := New("/nonexistent/path/config.yaml")
		if err == nil {
			t.Error("New() expected error for non-existent file")
		}
	})
}

func TestConfigGet(t *testing.T) {
	cfg := getTestConfig()
	got := cfg.Get()
	if got != cfg {
		t.Error("Get() should return the same pointer")
	}
}

func TestConfigLoad_Errors(t *testing.T) {
	t.Run("file_not_found", func(t *testing.T) {
		cfg := new(Config)
		err := cfg.Load("/nonexistent/path/config.yaml")
		if err == nil {
			t.Error("Load() expected error for non-existent file")
		}
	})

	t.Run("invalid_yaml", func(t *testing.T) {
		dir := t.TempDir()
		cfgPath := filepath.Join(dir, "bad.yaml")

		if err := os.WriteFile(cfgPath, []byte("invalid: [yaml: content"), 0644); err != nil {
			t.Fatalf("WriteFile failed: %v", err)
		}

		cfg := new(Config)
		err := cfg.Load(cfgPath)
		if err == nil {
			t.Error("Load() expected error for invalid YAML")
		}
	})
}

func TestConfigSave_Errors(t *testing.T) {
	t.Run("invalid_path", func(t *testing.T) {
		cfg := getTestConfig()
		err := cfg.Save("/nonexistent/dir/config.yaml")
		if err == nil {
			t.Error("Save() expected error for invalid path")
		}
	})
}

func TestProvideConfigApp(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		dir := t.TempDir()
		cfgPath := filepath.Join(dir, "config.yaml")

		cfg := getTestConfig()
		if err := cfg.Save(cfgPath); err != nil {
			t.Fatalf("Save failed: %v", err)
		}

		cliCfg := &cli.Config{Config: cfgPath}
		svc, cleanup, err := ProvideConfigApp(cliCfg)
		if err != nil {
			t.Errorf("ProvideConfigApp() returned error: %v", err)
		}
		if svc == nil {
			t.Error("ProvideConfigApp() returned nil Service")
		}
		cleanup()
	})

	t.Run("error", func(t *testing.T) {
		cliCfg := &cli.Config{Config: "/nonexistent/path/config.yaml"}
		_, cleanup, err := ProvideConfigApp(cliCfg)
		if err == nil {
			t.Error("ProvideConfigApp() expected error for bad config path")
		}
		cleanup()
	})
}

func TestProvideDBConfig(t *testing.T) {
	cfg := getTestConfigWithDB()
	pgxCfg, err := ProvideDBConfig(cfg)
	if err != nil {
		t.Errorf("ProvideDBConfig() returned error: %v", err)
	}
	if pgxCfg == nil {
		t.Fatal("ProvideDBConfig() returned nil")
	}
	if pgxCfg.Host != "127.0.0.1" || pgxCfg.Port != "5432" {
		t.Errorf(
			"ProvideDBConfig() unexpected values: Host=%q, Port=%q",
			pgxCfg.Host, pgxCfg.Port,
		)
	}
}

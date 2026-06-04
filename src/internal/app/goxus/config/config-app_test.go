package configapp

import (
	"errors"
	"reflect"
	"testing"

	configserver "goxus/src/internal/app/goxus/api/server/config"
	configcron "goxus/src/internal/app/goxus/cron-job/config"

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

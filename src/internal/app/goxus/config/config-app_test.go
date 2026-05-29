package configapp

import (
	"errors"
	"reflect"
	"testing"

	"github.com/nobuenhombre/suikat/pkg/fico"
	configserver "goxus/src/internal/app/goxus/api/server/config"
	configcron "goxus/src/internal/app/goxus/cron-job/config"
	configexample "goxus/src/internal/app/goxus/cron-job/jobs/example/config"
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
			ExampleJob: configexample.ExampleJobConfig{
				Enabled:  true,
				Schedule: "@every 10m",
			},
		},
	}
}

func TestConfigLoad(t *testing.T) {
	test := &testConfig{
		fileName:    "config-app_test_load.yaml",
		fileContent: "",
		config:      getTestConfig(),
		err:         nil,
	}

	cfg := new(Config)
	err := cfg.Load(test.fileName)

	if !(reflect.DeepEqual(cfg, test.config) && errors.Is(err, test.err)) {
		t.Errorf(
			"cfg.Load(%#v),\n Expected (cfg = %#v, err = %#v),\n Actual (cfg = %#v, err = %#v).\n",
			test.fileName, test.config, test.err, cfg, err,
		)
	}
}

func TestConfigSave(t *testing.T) {
	test := &testConfig{
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
	}

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
}

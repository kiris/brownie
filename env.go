package brownie

import (
	"github.com/kelseyhightower/envconfig"
)

type env struct {
	SlackToken string `envconfig:"SLACK_TOKEN" required:"true"`
	VerificationToken string `envconfig:"VERIFICATION_TOKEN" required:"true"`
	WorkspaceDir string `envconfig:"WORKSPACE_DIR" required:"true"`
}

func CreateAppFromEnvironmentVariables() (*App, error) {
	env, err := loadEnv()
	if err != nil {
		return nil, err
	}

	return CreateApp(env.SlackToken, env.VerificationToken, env.WorkspaceDir), nil
}

func loadEnv() (env, error) {
	var env env
	err := envconfig.Process("", &env)
	return env, err
}

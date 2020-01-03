package brownie

import "github.com/kelseyhightower/envconfig"

type Env struct {
	SlackToken string `envconfig:"SLACK_TOKEN" required:"true"`
	WorkspaceDir string `envconfig:"WORKSPACE_DIR" required:"true"`
}

func LoadEnv() (Env, error) {
	var env Env
	err := envconfig.Process("", &env)
	return env, err
}

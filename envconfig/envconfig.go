package envconfig

import "github.com/kelseyhightower/envconfig"

type (
	Env struct {
		Debug        bool   `envconfig:"DEBUG" default:"false"`
		Port         int    `envconfig:"PORT" default:"8080" required:"true"`
		DATABASE_URL string `envconfig:"DATABASE_URL" required:"true"`
	}
)

func Process(env *Env) error {
	return envconfig.Process("", env)
}

func New() (*Env, error) {
	var env Env
	err := envconfig.Process("", &env)
	if err != nil {
		return nil, err
	}
	return &env, nil
}

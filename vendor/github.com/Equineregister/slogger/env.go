package slogger

import "errors"

const (
	EnvDev   Env = "dev"
	EnvLocal Env = "local"
	EnvUat   Env = "uat"
	EnvProd  Env = "prod"
)

type Env string

func (e Env) isValid() bool {
	switch e {
	case EnvDev, EnvUat, EnvProd, EnvLocal:
		return true
	}

	return false
}

func (e *Env) Set(value string) error {
	Env := Env(value)
	if !Env.isValid() {
		return errors.New(`Env must be one of "dev", "uat", "prod" or "local"`)
	}

	*e = Env

	return nil
}

func (e Env) String() string {
	return string(e)
}

var EnvLevels = map[string]Env{
	"debug": EnvDev,
	"info":  EnvUat,
	"warn":  EnvProd,
	"local": EnvLocal,
}

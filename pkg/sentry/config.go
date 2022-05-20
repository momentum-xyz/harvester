package sentry

// Config : structure to hold Sentry configuration
type Config struct {
	Dsn              string `yaml:"dsn" envconfig:"SENTRY_DSN"`
	Environment      string `yaml:"env" envconfig:"SENTRY_ENV"`
	Release          string `yaml:"release" envconfig:"SENTRY_RELEASE"`
	AttachStacktrace bool   `yaml:"attachStacktrace" envconfig:"SENTRY_ATTACH_STACK_TRACE"`
	DebugEnable      bool   `yaml:"debugEnable" envconfig:"SENTRY_ENABLE_DEBUG"`
}

package cfgbuild

// A Config requires an init() function to initialize any default values.  Also, fields should be
// tagged with `envvar:"ENV_VAR_NAME"` so that the builder can map env vars to the fields.
type Config interface {
	Init() error
}

package cfgbuild

// A Config requires an describes the functions required for a Config to work with the
//
//	cfgbuild.Builder.  The CfgBuildInit() function can be used initialize a newly instantiated
//
// config such as setting default values.  The CfgBuildValidate() function can validate a create
// config such as by making sure certain values are valid.  Adding the BaseConfig to the
// implementing struct will add default implementations of the functions.  Along with the interface
// functions, fields should be tagged with `envvar:"ENV_VAR_NAME"` so that the builder can map env
// vars to the fields.
type Config interface {
	CfgBuildInit() error
	CfgBuildValidate() error
}

// BaseConfig provides default implmentations for the Config interface.
type BaseConfig struct{}

// CfgBuildInit in BaseConfig does nothing.
func (cfg *BaseConfig) CfgBuildInit() error {
	return nil
}

// CfgBuildValidate in BaseConfig does nothing.
func (cfg *BaseConfig) CfgBuildValidate() error {
	return nil
}

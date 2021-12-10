package plugin

type Plugin interface {
	Name() string
}

type InitPluginFunc func (dir string) (plugins []Plugin, err error)

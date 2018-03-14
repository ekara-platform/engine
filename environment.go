package engine

type Environment interface {
	Labeled

	GetVersion() (Version, error)
	GetName() string
	GetDescription() string
}

func (e environmentDef) GetVersion() (Version, error) {
	return CreateVersion(e.Version)
}

func (e environmentDef) GetName() string {
	return e.Name
}

func (e environmentDef) GetDescription() string {
	return e.Description
}

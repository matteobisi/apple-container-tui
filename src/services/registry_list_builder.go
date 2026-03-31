package services

import "container-tui/src/models"

// RegistryListBuilder builds a registry list command.
type RegistryListBuilder struct{}

func (RegistryListBuilder) Validate() error {
	return nil
}

func (RegistryListBuilder) Build() (models.Command, error) {
	return models.Command{Executable: "container", Args: []string{"registry", "list", "--format", "json"}}, nil
}

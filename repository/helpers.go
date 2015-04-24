package repository

func mapKnownErrors(knownErrors map[string]error, err error) error {
	if err == nil {
		return nil
	}

	errorMessage := err.Error()

	if _, exists := knownErrors[errorMessage]; exists {
		return knownErrors[errorMessage]
	}

	return err
}

package shared

// Converter est l'interface commune que tous les plugins vont implémenter.
type Converter interface {
	// Convert prend la configuration en entrée (YAML ou HCL) et la traite
	Convert(input []byte) ([]byte, error)
}

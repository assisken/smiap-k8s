package differ

type DiffObject interface {
	GetAnnotations() map[string]string
	SetAnnotations(annotations map[string]string)
}

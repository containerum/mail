package misc

// Masker allows to mask some in object in order to not disclose sensitive information
type Masker interface {
	Mask()
}

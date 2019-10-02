package helpers

type ConfigFormatType int

const (
	DeprecatedExtensionsFormat ConfigFormatType = iota
	NewExtensionsFormat
	StronglyTyped
)

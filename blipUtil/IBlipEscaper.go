package blipUtil

type IBlipEscaper interface {
	Escape(inStr string) string
	GetFileType() string
}

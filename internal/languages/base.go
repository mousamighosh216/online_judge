package languages

type Language interface {
	CompileCommand() string
	RunCommand() string
}
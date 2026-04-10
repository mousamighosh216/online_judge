package languages

type CPP struct{}

func (CPP) CompileCommand() string {
	return "g++ main.cpp -O2 -o main"
}

func (CPP) RunCommand() string {
	return "./main"
}
package utils

type stringBuilder struct {
	str string
}

func NewStringBuilder() *stringBuilder {
	return &stringBuilder{}
}

func (s *stringBuilder) String() string {
	return s.str
}

func (s *stringBuilder) NextLine(str string) *stringBuilder {
	s.str += str + "\n"
	return s
}

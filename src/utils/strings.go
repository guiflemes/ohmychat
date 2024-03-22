package utils

import "fmt"

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

type bulletListBuilder[T any] struct{}

func NewBulletListBuilder[T any]() *bulletListBuilder[T] {
	return &bulletListBuilder[T]{}
}

func (b *bulletListBuilder[T]) Build(items []T, transform func(item T) string) string {
	s := NewStringBuilder()
	for _, item := range items {
		str := transform(item)
		s.NextLine(fmt.Sprintf("\u2022 %s", str))
	}
	return s.str
}

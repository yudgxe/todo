package utils

import "fmt"

func Panicf(format string, a ...any) {
	panic(fmt.Sprintf(format, a...))
}

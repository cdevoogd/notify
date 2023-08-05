package log

import (
	"fmt"
	"os"
)

func Println(a ...any) {
	fmt.Fprintln(os.Stdout, a...)
}

func Error(a ...any) {
	fmt.Fprintln(os.Stderr, a...)
}

func Fatal(a ...any) {
	Error(a...)
	os.Exit(1)
}

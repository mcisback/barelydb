package main

import (
	"fmt"
	"os"
)

//go:vet:printf PrintError 1 2
func PrintError(errorMessage ...any) {
	if len(errorMessage) > 0 {
		format, ok := errorMessage[0].(string)

		if ok {
			fmt.Fprintf(os.Stderr, format, errorMessage[1:]...)
		} else {
			fmt.Fprintln(os.Stderr, errorMessage...)
		}
	}
}

type Ok[T any] struct {
	Value *T
	Err   error
}

func Try[T any](v T, err error) Ok[T] {
	return Ok[T]{&v, err}
}

func TrySingle(err error) Ok[any] {
	return Ok[any]{nil, err}
}

func (o Ok[T]) Or(callback func(error)) T {
	if o.Err != nil {
		callback(o.Err)
	}

	return o.GetValue()
}

func (o Ok[T]) OrPanic(errorMessage ...any) T {
	if o.Err != nil {
		PrintError(errorMessage...)

		panic(o.Err)
	}

	return o.GetValue()
}

func (o Ok[T]) OrPrint(errorMessage ...any) T {
	if o.Err != nil {
		PrintError(errorMessage...)
	}

	return o.GetValue()
}

func (o Ok[T]) OrPrintAndExit(errorMessage ...any) T {
	if o.Err != nil {
		PrintError(errorMessage...)

		os.Exit(1)
	}

	return o.GetValue()
}

func (o Ok[T]) GetValue() T {
	var zero T

	if o.Value == nil {
		return zero
	}

	return *o.Value
}

func PrintAndExit(format string, errorMessage ...any) {
	fmt.Fprintf(os.Stderr, format, errorMessage...)

	os.Exit(1)
}

func IsError(err error) bool {
	return err != nil
}

func IsOk(err error) bool {
	return err == nil
}

func IsNil(value any) bool {
	return value == nil
}

func IsEmpty(value string) bool {
	return value == ""
}

func IsNotEmpty(value string) bool {
	return value != ""
}

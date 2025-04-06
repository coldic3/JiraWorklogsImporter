package converter

import (
	"errors"
)

type Converter interface {
	Convert(record []string) (string, error)
	Supports(format string) bool
}

type ConverterFactory struct {
	converters []Converter
}

func NewConverterFactory() *ConverterFactory {
	return &ConverterFactory{
		converters: []Converter{
			&TogglToJiraConverter{},
			&ClockifyToJiraConverter{},
		},
	}
}

func (f *ConverterFactory) GetConverter(format string) (Converter, error) {
	for _, converter := range f.converters {
		if converter.Supports(format) {
			return converter, nil
		}
	}
	return nil, errors.New("no converter found for format: " + format)
}

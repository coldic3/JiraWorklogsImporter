package converter

import (
	"errors"
	"fmt"
)

type Converter interface {
	Convert(record []string) (ConvertedRecord, error)
	Supports(strategy string) bool
}

type ConverterFactory struct {
	converters []Converter
}

type ConvertedRecord struct {
	IssueIdOrKey      string
	ContentText       string
	StartedAtDateTime string
	TimeSpentSeconds  int
}

func NewConverterFactory() *ConverterFactory {
	return &ConverterFactory{
		converters: []Converter{
			&TogglToJiraConverter{},
			&ClockifyToJiraConverter{},
		},
	}
}

func (f *ConverterFactory) GetConverter(strategy string) (Converter, error) {
	for _, converter := range f.converters {
		if converter.Supports(strategy) {
			return converter, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("No converter found for strategy: %s.", strategy))
}

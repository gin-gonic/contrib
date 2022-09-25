package rest

import (
	"github.com/gertd/go-pluralize"
	"regexp"
	"strings"
)

type ModelNameSanitizer func(string) string

func PluralToSingular(m string) string {
	var pluralize = pluralize.NewClient()
	if pluralize.IsPlural(m) {
		return pluralize.Singular(m)
	}
	return m
}

func SpacesToDashes(m string) string {
	const dashSeparator = "-"
	whitespaceRegexp := regexp.MustCompile(`\s+`)
	return whitespaceRegexp.ReplaceAllString(m, dashSeparator)
}

func GetDefaultNameSantizers() []ModelNameSanitizer {
	return []ModelNameSanitizer{
		strings.TrimSpace,
		PluralToSingular,
		SpacesToDashes,
	}
}

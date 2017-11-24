package main

import "regexp"

func initRegexp(pattern string) *regexp.Regexp {
	return regexp.MustCompile(pattern)
}

func matchesRegexp(filenames string, regexp *regexp.Regexp) bool {
	return regexp.MatchString(filenames)
}

func filterStringsByRegexp(links []string, regexp *regexp.Regexp) []string {
	output := make([]string, 0)
	for _, link := range links {
		if matchesRegexp(link, regexp) {
			output = append(output, link)
		}
	}
	return output
}

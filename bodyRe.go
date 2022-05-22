package main

import (
	"regexp"
)

func Re1(body string) {
	strs := regexpMatch(body, `"http([\s\S]*?)m3u8`)
	if len(strs) > 0 {
		s := strs[0]
		result <- s[1:]
	}
}

func Re2(body string) {
	strs := regexpMatch(body, `src=http([\s\S]*?)\.(m3u8|mp4)`)
	if len(strs) > 0 {
		s := strs[0]
		result <- s[4:]
	}
}

func regexpMatch(text, expr string) []string {
	reg, err := regexp.Compile(expr)
	if err != nil {
		return nil
	}
	return reg.FindStringSubmatch(text)
}

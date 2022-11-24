package server

import (
	"fmt"
	"regexp"
	"strings"
)

func CompileRoute(input string) (*regexp.Regexp, option) {
	output := ""
	allInserts := regexp.MustCompile("{[^/]+}").FindAllString(input, -1)
	opt := input
	for _, insert := range allInserts {
		split1 := strings.Split(insert, ": ")
		split2 := strings.Split(insert, ":")
		compiled := insert[1 : len(insert)-1]
		if len(split1) == 2 {
			compiled = "{" + split1[1][:len(split1[1])]
			opt = strings.ReplaceAll(opt, insert, split1[0]+"}")
		} else if len(split2) == 2 {
			compiled = "{" + split2[1][:len(split1[2])]
			opt = strings.ReplaceAll(opt, insert, split2[0]+"}")
		}
		input = strings.ReplaceAll(input, insert, compiled)
	}
	insert := false
	for _, char := range []rune(input) {
		if string(char) == "{" && !insert {
			insert = true
		} else if string(char) == "}" && insert {
			insert = false
		} else if !insert {
			output += regexp.QuoteMeta(string(char))
		} else if insert {
			output += string(char)
		}
	}
	return regexp.MustCompile(fmt.Sprintf("^(%s){1}$", output)), option(opt)
}

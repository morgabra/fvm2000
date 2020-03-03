package fvm

import (
	"fmt"
	"strconv"
	"strings"
)

func Assemble(p string) []byte {
	return parseAsm(p)
}

func parseAsm(p string) []byte {
	p = strings.ToUpper(p)
	tokens := []string{}
	instructions := strings.Split(p, "\n")
	for _, i := range instructions {
		for _, s := range strings.Split(strings.TrimSpace(i), " ") {
			if s != "" {
				tokens = append(tokens, s)
			}
		}
	}

	program, labels := parseLabels(tokens)
	compiled := make([]byte, 256)

	for idx, t := range program {
		if op, ok := opNames[t]; ok {
			compiled[idx] = op
			continue
		}

		if loc, ok := labels[t]; ok {
			switch compiled[idx-1] {
			case BNE:
				compiled[idx] = loc
				continue

			default:
				panic("invalid label use")
			}
		}

		v, err := strconv.ParseUint(t, 10, 8)
		if err != nil {
			panic(fmt.Sprintf("invalid number: %s", t))
		}
		compiled[idx] = byte(v)
	}

	return compiled
}

func parseLabels(tokens []string) ([]string, map[string]byte) {
	labels := make(map[string]byte)
	stripped := []string{}
	for idx, t := range tokens {
		if strings.HasSuffix(t, ":") {
			l := strings.TrimSuffix(t, ":")
			if _, ok := opNames[l]; ok {
				panic(fmt.Sprintf("invalid label name %s", l))
			}
			labels[l] = byte(idx)
			continue
		}
		stripped = append(stripped, t)
	}

	return stripped, labels
}

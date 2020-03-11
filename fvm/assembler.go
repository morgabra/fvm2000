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
	compiled := make([]byte, 2048)

	idx := 0
	pc := 0
	for pc < len(program) {
		t := program[pc]
		if op, ok := opNames[t]; ok {
			compiled[idx] = op
		} else if loc, ok := labels[t]; ok {
			switch compiled[idx-1] {
			default:
				compiled[idx] = byte(loc & 0xFF)
				compiled[idx+1] = byte(loc >> 8)
				idx++
			}
		} else if strings.HasPrefix(t, "0X") {
			v, err := strconv.ParseUint(t[2:], 16, 16)
			if err != nil {
				panic(fmt.Sprintf("invalid address: %s", t))
			}

			compiled[idx] = byte(v & 0xFF)
			compiled[idx+1] = byte(v >> 8)
			idx++
		} else if strings.HasPrefix(t, "#") {
			v, err := strconv.ParseUint(t[1:], 0, 16)
			if err != nil {
				panic(fmt.Sprintf("invalid literal int: %s", t))
			}

			switch program[pc-1] {
			case "MOV":
				compiled[idx-1] = MOVI

			case "ADD":
				compiled[idx-1] = ADDI

			case "SUB":
				compiled[idx-1] = SUBI

			case "MUL":
				compiled[idx-1] = MULI
			}

			compiled[idx] = byte(v & 0xFF)
			compiled[idx+1] = byte(v >> 8)
			idx++
		} else {
			v, err := strconv.ParseUint(t, 10, 8)
			if err != nil {
				panic(fmt.Sprintf("invalid number: %s", t))
			}
			compiled[idx] = byte(v)
		}

		pc++
		idx++
	}

	return compiled
}

func parseLabels(tokens []string) ([]string, map[string]uint16) {
	labels := make(map[string]uint16)
	labelOffset := 0
	stripped := []string{}
	tokenCounter := 0
	instructionsFound := false

	for tokenCounter < len(tokens) {
		t := tokens[tokenCounter]
		if b, ok := builtins[t]; ok {
			labels[t] = b
			stripped = append(stripped, t)
		} else if strings.HasSuffix(t, ":") {
			l := strings.TrimSuffix(t, ":")
			if _, ok := opNames[l]; ok {
				panic(fmt.Sprintf("cannot use reserveed label %s", l))
			}
			if _, ok := builtins[l]; ok {
				panic(fmt.Sprintf("cannot use reserved label: %s", l))
			}
			if _, ok := labels[l]; ok {
				panic(fmt.Sprintf("label aleady defined: %s", l))
			}
			labels[l] = uint16(len(stripped) + labelOffset)
			labelOffset += 2
		} else if t == "#DEFINE" {
			if instructionsFound {
				panic("defines must happen before instructions")
			}

			lbl := tokens[tokenCounter+1]
			addr, err := strconv.ParseUint(tokens[tokenCounter+2], 0, 16)
			if err != nil {
				panic(fmt.Sprintf("invalid address value: %s", tokens[tokenCounter+2]))
			}
			labels[lbl] = uint16(addr)
			tokenCounter += 3
			continue
		} else {
			stripped = append(stripped, t)
		}
		instructionsFound = true
		tokenCounter++
	}

	return stripped, labels
}

package hw3

import (
	"errors"
	"strconv"
	"unicode"
)

type StringPacker struct{}

func (p *StringPacker) Pack(s string) (string, error) {
	return "", errors.New("Pack is not implemented yet")
}

func (p *StringPacker) Unpack(s string) (string, error) {
	codes := make([]rune, 0, len(s))
	multiplierOffset := -1

	for i, code := range s {
		if code == '\\' && (i == 0 || s[i-1] != '\\') {
			continue
		}

		length := len(codes)
		isDigit := unicode.IsDigit(code)
		isEscaped := i > 0 && s[i-1] == '\\' && (i < 2 || s[i-2] != '\\')

		if isDigit {
			if length > 0 {
				if isEscaped {
					p.addToEnd(&codes, code, 1)
				} else if multiplierOffset < 0 {
					multiplierOffset = i
				}
			}
			continue
		}

		if multiplierOffset >= 0 {
			if length > 0 {
				multiplier, err := strconv.Atoi(s[multiplierOffset:i])
				if err == nil {
					p.applyMultiplier(&codes, multiplier)
				} else {
					return "", err
				}
			}
			multiplierOffset = -1
			p.addToEnd(&codes, code, 1)
			continue
		}

		p.addToEnd(&codes, code, 1)
	}

	if multiplierOffset >= 0 {
		multiplier, err := strconv.Atoi(s[multiplierOffset:])
		if err == nil {
			p.applyMultiplier(&codes, multiplier)
		} else {
			return "", err
		}
	}

	return string(codes), nil
}

func (p *StringPacker) applyMultiplier(codes *[]rune, multiplier int) {
	if multiplier > 0 {
		p.addToEnd(codes, (*codes)[len(*codes)-1], multiplier-1)
	} else {
		p.removeFromEnd(codes, 1)
	}
}

func (p *StringPacker) addToEnd(codes *[]rune, code rune, quantity int) {
	if quantity <= 0 {
		return
	}

	for i := 0; i < quantity; i++ {
		*codes = append(*codes, code)
	}
}

func (p *StringPacker) removeFromEnd(codes *[]rune, quantity int) {
	if quantity <= 0 {
		return
	}

	length := len(*codes)
	if length == 0 {
		return
	}

	if quantity >= length {
		quantity = length
	}

	*codes = (*codes)[:length-quantity]
}

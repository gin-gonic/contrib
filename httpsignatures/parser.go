package httpsignatures

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

type parser struct {
	input string
	pos   int
	ch    byte
}

func newParser(input string) (*parser, error) {
	p := &parser{input: input, pos: -1}
	p.readChar()
	if err := p.skipPrefix(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *parser) readChar() {
	if p.pos+1 >= len(p.input) {
		p.ch = 0
	} else {
		p.ch = p.input[p.pos+1]
	}
	p.pos++
}

func (p *parser) peekChar() byte {
	if p.pos+1 >= len(p.input) {
		return 0
	}
	return p.input[p.pos+1]
}

func (p *parser) skipPrefix() error {
	const prefix = "Signature "
	for i := 0; i < len(prefix); i++ {
		ch := prefix[i]
		if ch != p.ch {
			return fmt.Errorf("invalid prefix, expected %c at pos %d, got %c", ch, i, p.ch)
		}
		p.readChar()
	}
	return nil
}

func (p *parser) nextParam() (string, string, error) {
	var (
		key       = bytes.NewBuffer(nil)
		keyParsed = false

		val       = bytes.NewBuffer(nil)
		valParsed = false
	)

	if p.ch == 0 {
		return "", "", io.EOF
	}

	for {
		switch p.ch {
		case ',', 0:
			if !valParsed {
				return "", "", errors.New("unterminated parameter")
			}
			p.readChar()
			return key.String(), val.String(), nil
		case '"':
			if keyParsed {
				if p.peekChar() != ',' && p.peekChar() != 0 {
					if err := val.WriteByte(p.ch); err != nil {
						return "", "", err
					}
				} else {
					valParsed = true
				}
				p.readChar()

			}
		case '=':
			p.readChar()
			if p.ch != '"' {
				return "", "", fmt.Errorf(`expected " chraracter at pos: %d`, p.pos)
			}
			p.readChar()
			keyParsed = true
		default:
			if !keyParsed {
				if err := key.WriteByte(p.ch); err != nil {
					return "", "", err
				}
				p.readChar()
			} else {
				if err := val.WriteByte(p.ch); err != nil {
					return "", "", err
				}
				p.readChar()
			}
		}
	}
}

func (p *parser) parse() (map[string]string, error) {
	var params = make(map[string]string)

	for {
		key, val, err := p.nextParam()
		if err == io.EOF {
			return params, nil
		} else if err != nil {
			return nil, err
		}
		params[key] = val
	}
}

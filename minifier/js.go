// imported from https://github.com/tdewolff/minify/tree/master/js
package minifier

import (
	"bytes"
	"io"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/parse"
	"github.com/tdewolff/parse/js"
)

var (
	spaceBytes   = []byte(" ")
	newlineBytes = []byte("\n")
)

////////////////////////////////////////////////////////////////

// Minifier is a JS minifier.
type Minifier struct{}

// Minify minifies JS data, it reads from r and writes to w.
func Minify(m *minify.M, w io.Writer, r io.Reader, params map[string]string) error {
	return (&Minifier{}).Minify(m, w, r, params)
}

// Minify minifies JS data, it reads from r and writes to w.
func (o *Minifier) Minify(m *minify.M, w io.Writer, r io.Reader, _ map[string]string) error {
	prev := js.LineTerminatorToken
	prevLast := byte(' ')
	lineTerminatorQueued := false
	whitespaceQueued := false

	l := js.NewLexer(r)
	for {
		tt, data := l.Next()
		if tt == js.ErrorToken {
			if l.Err() != io.EOF {
				return l.Err()
			}
			return nil
		} else if tt == js.LineTerminatorToken {
			lineTerminatorQueued = true
		} else if tt == js.WhitespaceToken {
			whitespaceQueued = true
		} else if tt == js.CommentToken {
			if len(data) > 5 && data[1] == '*' && data[2] == '!' {
				if _, err := w.Write(data[:3]); err != nil {
					return err
				}
				comment := parse.TrimWhitespace(parse.ReplaceMultipleWhitespace(data[3 : len(data)-2]))
				if _, err := w.Write(comment); err != nil {
					return err
				}
				if _, err := w.Write(data[len(data)-2:]); err != nil {
					return err
				}
			}
		} else {
			first := data[0]
			if first == '`' {
				if l := len(data); l > 2 {
					html, err := minifyHTMLTemplate(m, data[1:l-2])
					if err == nil && len(html) > 0 {
						html = bytes.Replace(html, []byte("\""), []byte("\\\""), -1) // Escape data
						data = []byte{'"'}
						data = append(data, html...)
						data = append(data, '"')
					}
				}
			}
			if (prev == js.IdentifierToken || prev == js.NumericToken || prev == js.PunctuatorToken || prev == js.StringToken || prev == js.RegexpToken) &&
				(tt == js.IdentifierToken || tt == js.NumericToken || tt == js.PunctuatorToken || tt == js.RegexpToken) {
				if lineTerminatorQueued && (prev != js.PunctuatorToken || prevLast == '}' || prevLast == ']' || prevLast == ')' || prevLast == '+' || prevLast == '-' || prevLast == '"' || prevLast == '\'') &&
					(tt != js.PunctuatorToken || first == '{' || first == '[' || first == '(' || first == '+' || first == '-' || first == '!') {
					if _, err := w.Write(newlineBytes); err != nil {
						return err
					}
				} else if whitespaceQueued && (prev != js.StringToken && prev != js.PunctuatorToken && tt != js.PunctuatorToken || (prevLast == '+' || prevLast == '-') && first == prevLast) {
					if _, err := w.Write(spaceBytes); err != nil {
						return err
					}
				}
			} else if lineTerminatorQueued && (prev == js.IdentifierToken || prev == js.PunctuatorToken && prevLast == '}') && tt == js.StringToken {
				if _, err := w.Write(newlineBytes); err != nil {
					return err
				}
			}
			if _, err := w.Write(data); err != nil {
				return err
			}
			prev = tt
			prevLast = data[len(data)-1]
			lineTerminatorQueued = false
			whitespaceQueued = false
		}
		l.Free(len(data))
	}
}

func minifyHTMLTemplate(m *minify.M, data []byte) ([]byte, error) {
	r := bytes.NewReader(data)
	w := new(bytes.Buffer)
	hm := &html.Minifier{}
	err := hm.Minify(m, w, r, nil)
	b := bytes.Replace(w.Bytes(), []byte{'\n'}, []byte{}, -1)
	return bytes.Replace(b, []byte{'\r'}, []byte{}, -1), err
}

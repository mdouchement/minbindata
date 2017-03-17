package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/juju/errors"
	min "github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/json"
	"github.com/tdewolff/minify/svg"
	"github.com/tdewolff/minify/xml"
)

var filetypeMime = map[string]string{
	"css":    "text/css",
	"js":     "text/javascript",
	"json":   "application/json",
	"svg":    "image/svg+xml",
	"xml":    "text/xml",
	"binary": "application/octet-stream",
}

func minify(base, dst string, recursive bool) error {
	m := min.New()
	m.Add("text/css", &css.Minifier{})
	m.Add("text/javascript", &js.Minifier{})
	// TODO multiline string with backticks
	m.Add("image/svg+xml", &svg.Minifier{})
	m.Add("application/json", &json.Minifier{})
	m.Add("text/xml", &xml.Minifier{})
	m.AddFunc("application/octet-stream", func(m *min.M, w io.Writer, r io.Reader, _ map[string]string) error {
		// Simple copy
		_, err := io.Copy(w, r)
		return err
	})

	return filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !recursive {
			relative := ""
			relative, err = filepath.Rel(base, path)
			if err != nil {
				return errors.Annotate(err, "recursive depth check")
			}
			if strings.Count(relative, string(filepath.Separator)) > 0 {
				return nil
			}
		}

		output, err := getOutputFilename(base, path, dst)
		if err != nil {
			return errors.Annotate(err, "find output path")
		}

		if info.IsDir() {
			if err = os.MkdirAll(output, 0755); err != nil {
				return errors.Annotate(err, "create directory")
			}
			return nil
		}

		r, err := os.Open(path)
		if err != nil {
			return errors.Annotate(err, "input minify")
		}
		defer r.Close()

		w, err := os.Create(output)
		if err != nil {
			return errors.Annotate(err, "input minify")
		}

		err = m.Minify(mime(path), w, r)
		return errors.Annotate(err, "minification")
	})
}

func mime(path string) string {
	ext := filepath.Ext(path)
	if len(ext) > 0 {
		ext = ext[1:]
	}
	if m, ok := filetypeMime[ext]; ok {
		return m
	}
	return filetypeMime["binary"]
}

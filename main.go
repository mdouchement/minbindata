package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	bindata "github.com/jteeuwen/go-bindata"
	"github.com/juju/errors"
	"gopkg.in/urfave/cli.v2"
)

func main() {
	app := &cli.App{
		Name:    "Minify bindata",
		Version: "0.1.0",
		Authors: []*cli.Author{
			&cli.Author{Name: "mdouchement", Email: "https://github.com/mdouchement"},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "pkg",
				Usage: "Package name to use in the generated code",
			},
			&cli.StringSliceFlag{
				Name:  "ignore",
				Usage: "Regex patterns to ignore",
			},
			&cli.StringFlag{
				Name:  "o",
				Usage: "Name of the output file to be generated",
			},
		},
		Action: action,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

func action(c *cli.Context) error {
	patterns := make([]*regexp.Regexp, 0)
	for _, pattern := range c.StringSlice("ignore") {
		patterns = append(patterns, regexp.MustCompile(pattern))
	}

	path := c.Args().First()
	if path == "" {
		return errors.New("Missing <input dir>")
	}
	recursive := false
	if strings.HasSuffix(path, "/...") {
		path = filepath.Clean(path[:len(path)-4])
		recursive = true
	}

	tmp, err := ioutil.TempDir("", "minbindata")
	if err != nil {
		return errors.Annotate(err, "create temporary directory")
	}
	defer os.RemoveAll(tmp)

	if err = minify(path, tmp, recursive); err != nil {
		return err
	}

	cfg := bindata.NewConfig()
	cfg.Input = []bindata.InputConfig{
		{
			Path:      filepath.Clean(tmp),
			Recursive: recursive,
		},
	}
	cfg.Prefix = filepath.Clean(tmp)
	if v := c.String("pkg"); v != "" {
		cfg.Package = v
	}
	if v := c.String("o"); v != "" {
		cfg.Output = v
	}
	if len(patterns) > 0 {
		cfg.Ignore = patterns
	}

	err = bindata.Translate(cfg)
	return errors.Annotate(err, "bindata")
}

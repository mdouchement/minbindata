# Minify binned data

Concats [minify](github.com/tdewolff/minify) and [go-bindata](github.com/jteeuwen/go-bindata) in one CLI application.

## Requirements

- Golang 1.7.x,1.8.x

## Installation

```sh
go get github.com/mdouchement/minbindata
```

## Usage

```sh
minbindata -h
```

```sh
minbindata --ignore .DS_Store -pkg web -o web/assets.go public/...
```


## License

**MIT**

## Contributing

1. Fork it
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
5. Push to the branch (git push origin my-new-feature)
6. Create new Pull Request

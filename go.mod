module github.com/ekara-platform/engine

go 1.13

require (
	github.com/ekara-platform/model v1.0.0
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/fatih/color v1.7.0
	github.com/gobwas/glob v0.2.3
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/oklog/ulid v1.3.1
	github.com/stretchr/testify v1.2.2
	golang.org/x/crypto v0.0.0-20181015023909-0c41d7ab0a0e
	golang.org/x/net v0.0.0-20181011144130-49bb7cea24b1 // indirect
	gopkg.in/src-d/go-billy.v4 v4.3.0 // indirect
	gopkg.in/src-d/go-git.v4 v4.7.1
	gopkg.in/yaml.v2 v2.2.1
)

replace github.com/ekara-platform/model => ../model

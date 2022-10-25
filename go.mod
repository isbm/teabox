module gitlab.com/isbm/teabox

go 1.19

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/isbm/go-nanoconf v0.0.0-20210917204429-663038ee6e05
	github.com/karrick/godirwalk v1.17.0
)

require github.com/go-yaml/yaml v2.1.0+incompatible // indirect

replace github.com/isbm/crtview v1.6.2 => ../crtview

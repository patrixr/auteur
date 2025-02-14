.PHONY: all release test templ clean tailwind build example tidy run tag

BIN_NAME := auteur
MODULE_NAME := github.com/patrixr/auteur
BUILD_FOLDER := ./out
VERSION_CMD := grep "version:" auteur.yaml | cut -d: -f2 | tr -d ' '
VERSION := $(shell ${VERSION_CMD})

tag: test
	git tag -a "v`${VERSION_CMD}`" -m "Release version `${VERSION_CMD}`"
	git push origin v`${VERSION_CMD}`

release:
	gh release create  v`${VERSION_CMD}`

test:
	ENV=test go test -json -v ./... | go run  github.com/mfridman/tparse@latest -all

templ:
	go run github.com/a-h/templ/cmd/templ@latest generate

clean:
	$(RM) ./**/*_templ.go

example:
	AUTEUR_DEV_MODE=true go run ./examples/blog/main.go

tidy:
	go mod tidy

build:
	go build -ldflags="-X '${MODULE_NAME}/cmd.Version=${VERSION}'" -o ${BUILD_FOLDER}/${BIN_NAME} ./

run:
	go run main.go

auteur: run
	echo "auteur.tronica.io" > dist/CNAME

air:
	go get github.com/air-verse/air
	AUTEUR_WEBROOT="/" go run github.com/air-verse/air

serve:
	npx serve ./dist

version: build
	${BUILD_FOLDER}/${BIN_NAME} version

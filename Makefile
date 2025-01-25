.PHONY: all release test templ clean tailwind build example tidy run

BIN_NAME := auteur
BUILD_FOLDER := ./out

release:
	git tag -a "v`cat ./VERSION`" -m "Release version `cat ./VERSION`"
	git push origin v`cat ./VERSION`

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
	go build -o ${BUILD_FOLDER}/${BIN_NAME} ./

run:
	go run main.go

air:
	go run github.com/air-verse/air

serve:
	npx serve ./dist

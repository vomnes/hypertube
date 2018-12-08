.PHONY: build doc fmt lint run test vendor_clean vendor_get

GOPATH := ${PWD}/_vendor:${GOPATH}
export GOPATH

default: build

build:
	go build -v -o ./bin/api ./src/*.go

doc:
	godoc -http=:6060 -index -goroot .

fmt:
	go fmt ./

lint:
	golint .

launch:
	./bin/api -port=4003

run: build launch

test:
	@go test ./src/auth/.
	@go test ./src/mail/.
	@go test ./src/profile/.
	@go test ./src/movie/.
	@go test ./src/comment/.

test-d:
	@go test ./src/$(TARGET)/.

test-verbose:
	@go test -v ./src/auth/.
	@go test -v ./src/mail/.
	@go test -v ./src/profile/.
	@go test -v ./src/movie/.
	@go test -v ./src/comment/.

test-target:
	go test -v -run '$(TARGET)' ./src/...

cover:
	go test -cover ./src/...

ROUTES=src/routes/
cover-html:
	go test -coverprofile=cover.out ./src/...
	sed -i -e "s#.*\($(ROUTES).*\.go\)#\./\\1#" cover.out #Format cover.out
	go tool cover -html=cover.out -o coverage.html
	rm cover.out

cover-func:
	@go test -coverprofile=cover.out ./src/...
	@sed -i -e "s#.*\($(ROUTES).*\.go\)#\./\\1#" cover.out #Format cover.out
	@go tool cover -func=cover.out
	@rm cover.out

vendor_clean:
	rm -rf ./_vendor/src

# We have to set GOPATH to just the _vendor
# directory to ensure that `go get` doesn't
# update packages in our primary GOPATH instead.
# This will happen if you already have the package
# installed in GOPATH since `go get` will use
# that existing location as the destination.

vendor_get: vendor_clean
	GOPATH=${PWD}/_vendor go get -d -u -v \
		github.com/rs/cors \
		github.com/gorilla/mux \
		github.com/gorilla/handlers \
		gopkg.in/mgo.v2 \
		github.com/kylelemons/godebug/pretty \
		github.com/dgrijalva/jwt-go \
		github.com/mailjet/mailjet-apiv3-go \
		github.com/markbates/goth \
		github.com/gorilla/sessions \
		golang.org/x/crypto/bcrypt \
		github.com/markbates/goth/gothic \
		github.com/markbates/goth/providers/gplus \
		gopkg.in/mgo.v2 \
		golang.org/x/text/unicode/rangetable

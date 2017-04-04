MAIN_VERSION:=$(shell git describe --abbrev=0 --tags || echo "0.1.0")
VERSION:=${MAIN_VERSION}\#$(shell git log -n 1 --pretty=format:"%h")
PACKAGES:=$(shell go list ./... | sed -n '1!p' | grep -v /vendor/)
LDFLAGS:=-ldflags "-X github.com/ederavilaprado/golang-web-architecture-template/app.Version=${VERSION}"

default: dev

depends:
	../../../../bin/glide up

test:
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES), \
		go test -p=1 -cover -covermode=count -coverprofile=coverage.out ${pkg}; \
		tail -n +2 coverage.out >> coverage-all.out;)

cover: test
	go tool cover -html=coverage-all.out

run:
	go run ${LDFLAGS} main.go

dev:
	CompileDaemon -exclude-dir ".git" -color -build "go build -o _build_hot_reload" -command "./_build_hot_reload"

# build: clean
# 	go build ${LDFLAGS} -a -o api main.go

clean:
	rm -rf _build_hot_reload coverage.out coverage-all.out
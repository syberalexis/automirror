DIST_FOLDER := "dist"
TAG_NAME := $(shell git tag -l --contains HEAD)

default: binary

dist:
	mkdir $(DIST_FOLDER)

assets: dist
	GO111MODULE=$(GO111MODULE) GOOS= go generate -x -v $(GOOPTS)

binary: assets
	go build -o $(DIST_FOLDER)/automirror_$(ARCH)

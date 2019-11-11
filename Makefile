DIST_FOLDER := "dist"
TAG_NAME := $(shell git tag -l --contains HEAD)
ARCH := $(shell go version | awk '{print $4}' | cut -d'/' -f2)

default: binary

dist:
	mkdir $(DIST_FOLDER)

binary: dist
	go build -o $(DIST_FOLDER)/automirror_$(TAG_NAME)_$(ARCH)

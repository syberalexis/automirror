DIST_FOLDER := "dist"
TAG_NAME := $(shell git tag -l --contains HEAD)
PROJECT_FOLDER := $(shell pwd)
PROJECT_NAME := $(shell basename $(PROJECT_FOLDER))
ARCH := $(shell go version | awk '{print $4}' | cut -d'/' -f2)

ifdef $(TAG_NAME)
VERSION = $(TAG_NAME)
else
VERSION = "dev"
endif

default: build

dist:
	mkdir $(DIST_FOLDER)

build: dist
	go build -o $(DIST_FOLDER)/$(PROJECT_NAME)_$(VERSION)_$(ARCH) cmd/$(PROJECT_NAME)/main.go

clean:
	rm -r $(DIST_FOLDER)
DIST_FOLDER := "dist"

default: binary

dist:
	mkdir $(DIST_FOLDER)

binary: dist
	go build -o $(DIST_FOLDER)

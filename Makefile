include $(GOROOT)/src/Make.inc

TARG=gresec
GOFILES=\
	gresec.go\
	node.go\
	nodemap.go\
	http.go\

include $(GOROOT)/src/Make.cmd

fmt:
	for SOURCE_FILE in *.go; do gofmt < $$SOURCE_FILE >$$SOURCE_FILE.fmt && mv $$SOURCE_FILE.fmt $$SOURCE_FILE; done

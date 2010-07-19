include $(GOROOT)/src/Make.$(GOARCH)

TARG=wsbench
GOFMT=gofmt -spaces=true -tabindent=false -tabwidth=2

GOFILES=\
	wsbench.go\

include $(GOROOT)/src/Make.pkg

format:
	${GOFMT} -w wsbench.go
	${GOFMT} -w wsbench_test.go

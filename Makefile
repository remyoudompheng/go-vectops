GOFILES = $(wildcard *.go)

.PHONY: vgo

vgo: $(GOFILES)
	go build -o vgo

clean:
	rm -f vgo

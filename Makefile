GOFILES = $(wildcard *.go)

vgo: $(GOFILES)
	go build -o vgo

clean:
	rm -f vgo

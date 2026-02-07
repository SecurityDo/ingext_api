.PHONY: build clean

build:
	go build -o ingext cmd/ingext/main.go

clean:
	rm -f ingext

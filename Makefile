GOCMD=go
BINARY_NAME=powerman
PACKAGES=.
DESTDIR=/usr/bin

.PHONY: test run install clean

all: test build
test:
	CGO_ENABLED=0 $(GOCMD) test -v ./...
build:
	$(GOCMD) build -o build/$(BINARY_NAME) -v $(PACKAGES)
run: build
	./build/$(BINARY_NAME)
install: build
	install -d $(DESTDIR)
	install build/$(BINARY_NAME) $(DESTDIR)/$(BINARY_NAME)
	install init/powerman.service /usr/lib/systemd/user/powerman.service
clean:
	$(GOCMD) clean
	rm -f build/$(BINARY_NAME)
	test -d build && rmdir build


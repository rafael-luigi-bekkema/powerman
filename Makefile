GOCMD=go
BINARY_NAME=powerman
PACKAGES=.
DESTDIR=/usr/bin

.PHONY: test run install clean

all: test build
test:
	$(GOCMD) test -v ./...
build:
	$(GOCMD) build -o build/$(BINARY_NAME) -v $(PACKAGES)
run: build
	./$(BINARY_NAME)
install: build
	install -d $(DESTDIR)
	install build/$(BINARY_NAME) $(DESTDIR)/$(BINARY_NAME)
clean:
	$(GOCMD) clean
	rm -f $(BINARY_NAME)


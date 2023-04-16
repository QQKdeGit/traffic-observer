NAME  = fake-traffic
ENTRY = .

MAKEFLAGS += --always-make

build:
	CGO_ENABLED=0 \
	go build \
	-gcflags="all=-trimpath=." \
	-asmflags="all=-trimpath=." \
	-ldflags="-s -w" \
	-o bin/$(NAME) $(ENTRY)
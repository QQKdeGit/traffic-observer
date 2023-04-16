NAME  = httpserver
ENTRY = .

MAKEFLAGS += --always-make

all: build container

build:
	CGO_ENABLED=0 \
	go build \
	-gcflags="all=-trimpath=." \
	-asmflags="all=-trimpath=." \
	-ldflags="-s -w" \
	-o bin/$(NAME) $(ENTRY)

container:
	docker build -t $(NAME) .

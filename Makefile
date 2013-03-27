# Viscum Makefile

# Source and build directory.
SRC   	:= 	src
BUILD		:=	build

# Build targets.
CLIENT  := 	$(BUILD)/viscum
SERVER 	:= 	$(BUILD)/viscumd

# External dependencies.
PKGS    := 	github.com/jteeuwen/go-pkg-rss   \
            github.com/jbarham/gopgsqldriver \
            code.google.com/p/goconf/conf


.PHONY: all

all: server client

clean:
	gd clean $(SRC)

server: env
	gd $(SRC) -o $(SERVER) -M ^server

client: env
	gd $(SRC) -o $(CLIENT) -M ^client

env: builddir pkgs

builddir:
	mkdir -p build

pkgs:
	for pkg in $(PKGS); do go get $$pkg; done

fmt:
	gd $(SRC) -fmt -w2

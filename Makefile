# Viscum Makefile

# Source and build directory.
SRC     :=  src
BUILD   :=  build

# Build targets.
CLIENT  :=  $(BUILD)/viscum
SERVER  :=  $(BUILD)/viscumd

# External dependencies.
PKGS    :=  github.com/jteeuwen/go-pkg-rss   \
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

install:
	install -dm755 $(DESTDIR)/usr/sbin
	install -m755  $(SERVER) $(DESTDIR)/usr/sbin/viscumd
	install -dm755 $(DESTDIR)/usr/bin
	install -m755  $(CLIENT) $(DESTDIR)/usr/bin/viscum
	install -dm755 $(DESTDIR)/usr/share/viscum
	install -m644 -t $(DESTDIR)/usr/share/viscum share/*
	install -dm755 $(DESTDIR)/etc/viscum
	install -m600  etc/viscumd.conf $(DESTDIR)/etc/viscum/viscumd.conf
	install -m644  etc/viscum.conf  $(DESTDIR)/etc/viscum/viscum.conf
	install -dm755 $(DESTDIR)/usr/lib/systemd/system
	install -m644 -t $(DESTDIR)/usr/lib/systemd/system etc/systemd/*.service
	install -dm755 $(DESTDIR)/usr/lib/tmpfiles.d
	install -m644 -t $(DESTDIR)/usr/lib/tmpfiles.d etc/systemd/*.conf

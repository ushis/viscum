# Viscum Makefile

# Source and build directory.
SRC     :=  src
BUILD   :=  build

# Build targets.
CLIENT  :=  $(BUILD)/viscum
SERVER  :=  $(BUILD)/viscumd

# External dependencies.
PKGS    :=  github.com/jteeuwen/go-pkg-xmlx  \
            github.com/ushis/gopgsqldriver   \
            code.google.com/p/goconf/conf    \
            github.com/moovweb/gokogiri

.PHONY: all

all: builddir pkgs
	gd $(SRC) -o $(SERVER) -M ^server
	gd $(SRC) -o $(CLIENT) -M ^client

clean:
	gd clean $(SRC)

builddir:
	mkdir -p build

pkgs:
	for pkg in $(PKGS); do go get $$pkg; done

fmt:
	gd $(SRC) -fmt -w2

test:
	gd $(SRC) test

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

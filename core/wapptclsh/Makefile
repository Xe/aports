#!/usr/bin/make

CC = gcc -Os
TCLLIB = /usr/lib/libtcl8.6.so -lm -lz -lpthread -ldl
TCLINC = /usr/include
TCLSH = tclsh

all: wapptclsh

wapptclsh: wapptclsh.c
	$(CC) -I. -I$(TCLINC) -o $@ wapptclsh.c $(TCLLIB)

wapptclsh.c:	wapptclsh.c.in wapp.tcl wapptclsh.tcl tclsqlite3.c mkccode.tcl
	$(TCLSH) mkccode.tcl wapptclsh.c.in >$@

clean:	
	rm wapptclsh wapptclsh.c

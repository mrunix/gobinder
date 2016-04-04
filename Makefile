
TOOLCHAIN_PREFIX := powerpc-linux-gnu-

SWIG_FLAGS  = -go -gccgo -intgosize 32
CFLAGS      = -DBINDER_IPC_32BIT -I$(CURDIR)/gobinder

.PHONY: all build clean 

all: build

build: gobinder/libgobinder.a gobindertest/gobindertest

gobinder/gobinder_wrap.c: gobinder/gobinder.i
	swig $(SWIG_FLAGS) -o gobinder/gobinder_wrap.c $<

gobinder/gobinder_wrap.o: gobinder/gobinder_wrap.c
	$(TOOLCHAIN_PREFIX)gcc $(CFLAGS) -c $< -o $@

gobinder/binder_impl.o: gobinder/binder_impl.c gobinder/binder_impl.h gobinder/linux/binder.h
	$(TOOLCHAIN_PREFIX)gcc $(CFLAGS) -c $< -o $@

gobinder/gobinder.o: gobinder/gobinder.go
	$(TOOLCHAIN_PREFIX)gccgo -c $< -o $@

gobinder/libgobinder.a: gobinder/gobinder_wrap.o gobinder/binder_impl.o gobinder/gobinder.o
	cd gobinder && $(TOOLCHAIN_PREFIX)ar -r libgobinder.a gobinder_wrap.o binder_impl.o gobinder.o

gobindertest/gobindertest: gobinder/libgobinder.a
	$(TOOLCHAIN_PREFIX)gccgo -static gobindertest/main.go -L./gobinder -lgobinder -o gobindertest/gobindertest

clean:
	rm -rf *.out
	rm -rf gobinder/gobinder_wrap.c
	rm -rf gobinder/gobinder.go
	rm -rf gobinder/binder_impl.o
	rm -rf gobinder/gobinder_wrap.o
	rm -rf gobinder/gobinder.o
	rm -rf gobinder/libgobinder.a
	rm -rf gobindertest/gobindertest


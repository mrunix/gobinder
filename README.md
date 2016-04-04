# gobinder
A Binder IPC for go lang

Use swig to generate a golang interface, GCCGO only.
see https://github.com/mrunix/gobinder/blob/master/Makefile

tested on PowerPC platform.

If port to another platform, should modify Makefile

<code>
TOOLCHAIN_PREFIX := powerpc-linux-gnu-

SWIG_FLAGS  = -go -gccgo -intgosize 32
CFLAGS      = -DBINDER_IPC_32BIT -I$(CURDIR)/gobinder
</code>

-intgosize 32  -- 32bit CPU
-intgosize 64  -- 64bit CPU

-DBINDER_IPC_32BIT  -- Binder Protocol Version 7


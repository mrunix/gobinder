package main

import (
	"fmt"
	"gobinder"
	"os"
	"unsafe"
)

var SVC_MGR_NAME string = "android.os.IServiceManager"

var PING_TRANSACTION uint = ((('_') << 24) | (('P') << 16) | (('N') << 8) | ('G'))
var SVC_MGR_GET_SERVICE uint = 1
var SVC_MGR_CHECK_SERVICE uint = 2
var SVC_MGR_ADD_SERVICE uint = 3
var SVC_MGR_LIST_SERVICES uint = 4

func usage() {
	fmt.Printf("usage: %s verb name\n", os.Args[0])
	os.Exit(2)
}

func svcmgr_lookup(bs gobinder.Struct_SS_binder_state, target uint, name string) uint {
	var handle uint
	var iodata [512]byte
	var msg gobinder.Struct_SS_binder_io = gobinder.New_binder_io()
	var reply gobinder.Struct_SS_binder_io = gobinder.New_binder_io()

	gobinder.Bio_init(msg, uintptr(unsafe.Pointer(&iodata)), 512, 4)
	gobinder.Bio_put_uint32(msg, uint(0))
	gobinder.Bio_put_string16_x(msg, SVC_MGR_NAME)
	gobinder.Bio_put_string16_x(msg, name)

	if gobinder.Binder_call(bs, msg, reply, target, SVC_MGR_CHECK_SERVICE) != 0 {
		gobinder.Free_binder_io(msg)
		gobinder.Free_binder_io(reply)
		return 0
	}

	handle = gobinder.Bio_get_ref(reply)
	if handle != 0 {
		gobinder.Binder_acquire(bs, handle)
	}

	gobinder.Binder_done(bs, msg, reply)

	gobinder.Free_binder_io(msg)
	gobinder.Free_binder_io(reply)
	return handle
}

func main() {
	if len(os.Args) != 3 {
		usage()
	}

	var svcmgr uint = 0

	bs := gobinder.Binder_open(128 * 1024)
	if bs == nil {
		fmt.Println("failed to open binder driver.")
		os.Exit(-1)
	}

	if os.Args[1] == "lookup" {
		handle := svcmgr_lookup(bs, svcmgr, os.Args[2])
		fmt.Printf("lookup(%s) = %x\n", os.Args[2], handle)
	}
}

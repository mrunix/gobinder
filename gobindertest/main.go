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

/*
#include <stdio.h>
#include <stdlib.h>
#include <errno.h>

#include "binder.h"

int main(int argc, char** argv) {
  struct binder_state* bs;
  uint32_t svcmgr = BINDER_SERVICE_MANAGER;
  uint32_t handle;

  bs = binder_open(128 * 1024);
  if (!bs) {
    fprintf(stderr, "failed to open binder driver\n");
    return -1;
  }

  argc--;
  argv++;
  while (argc > 0) {
    if (!strcmp(argv[0], "alt")) {
      handle = svcmgr_lookup(bs, svcmgr, "alt_svc_mgr");
      if (!handle) {
        fprintf(stderr, "cannot find alt_svc_mgr\n");
        return -1;
      }
      svcmgr = handle;
      fprintf(stderr, "svcmgr is via %x\n", handle);
    } else if (!strcmp(argv[0], "lookup")) {
      if (argc < 2) {
        fprintf(stderr, "argument required\n");
        return -1;
      }
      handle = svcmgr_lookup(bs, svcmgr, argv[1]);
      fprintf(stderr, "lookup(%s) = %x\n", argv[1], handle);
      argc--;
      argv++;
    } else if (!strcmp(argv[0], "publish")) {
      if (argc < 2) {
        fprintf(stderr, "argument required\n");
        return -1;
      }
      svcmgr_publish(bs, svcmgr, argv[1], &token);
      argc--;
      argv++;
    } else {
      fprintf(stderr, "unknown command %s\n", argv[0]);
      return -1;
    }
    argc--;
    argv++;
  }
  return 0;
}

uint32_t svcmgr_lookup(struct binder_state* bs, uint32_t target, const char* name) {
  uint32_t handle;
  unsigned iodata[512 / 4];
  struct binder_io msg, reply;

  bio_init(&msg, iodata, sizeof(iodata), 4);
  bio_put_uint32(&msg, 0);  // strict mode header
  bio_put_string16_x(&msg, SVC_MGR_NAME);
  bio_put_string16_x(&msg, name);

  if (binder_call(bs, &msg, &reply, target, SVC_MGR_CHECK_SERVICE))
    return 0;

  handle = bio_get_ref(&reply);

  if (handle)
    binder_acquire(bs, handle);

  binder_done(bs, &msg, &reply);

  return handle;
}

int svcmgr_publish(struct binder_state* bs, uint32_t target, const char* name, void* ptr) {
  int status;
  unsigned iodata[512 / 4];
  struct binder_io msg, reply;

  bio_init(&msg, iodata, sizeof(iodata), 4);
  bio_put_uint32(&msg, 0);  // strict mode header
  bio_put_string16_x(&msg, SVC_MGR_NAME);
  bio_put_string16_x(&msg, name);
  bio_put_obj(&msg, ptr);

  if (binder_call(bs, &msg, &reply, target, SVC_MGR_ADD_SERVICE))
    return -1;

  status = bio_get_uint32(&reply);

  binder_done(bs, &msg, &reply);

  return status;
}

unsigned token;
*/

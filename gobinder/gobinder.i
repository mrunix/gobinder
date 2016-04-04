/* File : gobinder.i */

%module gobinder

%{
#include "binder_impl.h"
%}

typedef unsigned int uint32_t;

%include typemaps.i
extern struct binder_io *new_binder_io();
extern void free_binder_io(struct binder_io *p);

extern struct binder_state* binder_open(size_t mapsize);
extern void binder_close(struct binder_state* bs);

extern int binder_call(struct binder_state* bs,
                       struct binder_io* msg, struct binder_io* reply,
                       uint32_t target, uint32_t code);

extern void binder_done(struct binder_state* bs,
                        struct binder_io* msg, struct binder_io* reply);

extern void binder_acquire(struct binder_state* bs, uint32_t target);
extern void binder_release(struct binder_state* bs, uint32_t target);

extern void binder_link_to_death(struct binder_state* bs, uint32_t target, struct binder_death* death);

extern void binder_loop(struct binder_state* bs, binder_handler func);

extern int binder_become_context_manager(struct binder_state* bs);

extern void bio_init(struct binder_io* bio, void* data,
                     size_t maxdata, size_t maxobjects);

extern void bio_put_obj(struct binder_io* bio, void* ptr);
extern void bio_put_ref(struct binder_io* bio, uint32_t handle);
extern void bio_put_uint32(struct binder_io* bio, uint32_t n);
extern void bio_put_string16(struct binder_io* bio, const uint16_t* str);
extern void bio_put_string16_x(struct binder_io* bio, const char* _str);

extern uint32_t bio_get_uint32(struct binder_io* bio);
extern uint16_t* bio_get_string16(struct binder_io* bio, size_t* sz);
extern uint32_t bio_get_ref(struct binder_io* bio);


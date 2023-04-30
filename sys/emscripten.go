package sys

import (
	"fmt"
	"log"

	wasmtime "github.com/bytecodealliance/wasmtime-go/v7"
)

// env.hiwire_get_length_helper => wasm function #14
func env_hiwire_get_length_helper(p0 int32) int32 {
	log.Printf("call to --> env_hiwire_get_length_helper:0x%x", p0)
	return 0
}

// env.emscripten_exit_with_live_runtime => wasm function #123
func env_emscripten_exit_with_live_runtime() {
	log.Printf("call to --> env_emscripten_exit_with_live_runtime:")
}

// wasi_snapshot_preview1.fd_seek => wasm function #194
func wasi_snapshot_preview1_fd_seek(p0 int32, p1 int64, p2 int32, p3 int32) int32 {
	log.Printf("call to --> wasi_snapshot_preview1_fd_seek:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
	return 0
}

// env.__cxa_current_primary_exception => wasm function #247
func env___cxa_current_primary_exception() int32 {
	log.Printf("call to --> env___cxa_current_primary_exception:")
	return 0
}

// env.strftime_l => wasm function #261
func env_strftime_l(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32) int32 {
	log.Printf("call to --> env_strftime_l:0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4)
	return 0
}

// env.JsArray_New => wasm function #4
func env_JsArray_New() int32 {
	log.Printf("call to --> env_JsArray_New:")
	return 0
}

// env.__assert_fail => wasm function #19
func env___assert_fail(p0 int32, p1 int32, p2 int32, p3 int32) {
	log.Printf("call to --> env___assert_fail:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
}

// env.__syscall_openat => wasm function #150
func env___syscall_openat(p0 int32, p1 int32, p2 int32, p3 int32) int32 {
	log.Printf("call to --> env___syscall_openat:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
	return 0
}

// env.invoke_iiiiiii => wasm function #249
func env_invoke_iiiiiii(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32, p6 int32) int32 {
	log.Printf("call to --> env_invoke_iiiiiii:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5, p6)
	return 0
}

// env.hiwire_call_bound => wasm function #58
func env_hiwire_call_bound(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env_hiwire_call_bound:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.destroy_proxies_js => wasm function #122
func env_destroy_proxies_js(p0 int32) int32 {
	log.Printf("call to --> env_destroy_proxies_js:0x%x", p0)
	return 0
}

// env.emscripten_memcpy_big => wasm function #167
func env_emscripten_memcpy_big(p0 int32, p1 int32, p2 int32) {
	log.Printf("call to --> env_emscripten_memcpy_big:0x%x,0x%x,0x%x", p0, p1, p2)
}

// env.__cxa_find_matching_catch_2 => wasm function #223
func env___cxa_find_matching_catch_2() int32 {
	log.Printf("call to --> env___cxa_find_matching_catch_2:")
	return 0
}

// env.hiwire_CallMethod => wasm function #10
func env_hiwire_CallMethod(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env_hiwire_CallMethod:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.hiwire_CallMethod_NoArgs => wasm function #12
func env_hiwire_CallMethod_NoArgs(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_hiwire_CallMethod_NoArgs:0x%x,0x%x", p0, p1)
	return 0
}

// env.hiwire_get_bool => wasm function #73
func env_hiwire_get_bool(p0 int32) int32 {
	log.Printf("call to --> env_hiwire_get_bool:0x%x", p0)
	return 0
}

// env.__syscall_newfstatat => wasm function #183
func env___syscall_newfstatat(p0 int32, p1 int32, p2 int32, p3 int32) int32 {
	log.Printf("call to --> env___syscall_newfstatat:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
	return 0
}

// env.__syscall_recvfrom => wasm function #277
func env___syscall_recvfrom(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32) int32 {
	log.Printf("call to --> env___syscall_recvfrom:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5)
	return 0
}

// env.hiwire_get_length_string => wasm function #15
func env_hiwire_get_length_string(p0 int32) int32 {
	log.Printf("call to --> env_hiwire_get_length_string:0x%x", p0)
	return 0
}

// env.JsException_new_helper => wasm function #53
func env_JsException_new_helper(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env_JsException_new_helper:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.ffi_closure_alloc_js => wasm function #142
func env_ffi_closure_alloc_js(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_ffi_closure_alloc_js:0x%x,0x%x", p0, p1)
	return 0
}

// wasi_snapshot_preview1.environ_get => wasm function #147
func wasi_snapshot_preview1_environ_get(p0 int32, p1 int32) int32 {
	log.Printf("call to --> wasi_snapshot_preview1_environ_get:0x%x,0x%x", p0, p1)
	return 0
}

// env.JsArray_Get => wasm function #48
func env_JsArray_Get(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_JsArray_Get:0x%x,0x%x", p0, p1)
	return 0
}

// env.js2python_convert => wasm function #78
func env_js2python_convert(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env_js2python_convert:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env._python2js_destroy_cache => wasm function #115
func env__python2js_destroy_cache(p0 int32) {
	log.Printf("call to --> env__python2js_destroy_cache:0x%x", p0)
}

// wasi_snapshot_preview1.fd_pwrite => wasm function #203
func wasi_snapshot_preview1_fd_pwrite(p0 int32, p1 int32, p2 int32, p3 int64, p4 int32) int32 {
	log.Printf("call to --> wasi_snapshot_preview1_fd_pwrite:0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4)
	return 0
}

// env.invoke_jiii => wasm function #266
func env_invoke_jiii(p0 int32, p1 int32, p2 int32, p3 int32) int64 {
	log.Printf("call to --> env_invoke_jiii:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
	return 0
}

// env.hiwire_greater_than => wasm function #71
func env_hiwire_greater_than(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_hiwire_greater_than:0x%x,0x%x", p0, p1)
	return 0
}

// env.hiwire_greater_than_equal => wasm function #72
func env_hiwire_greater_than_equal(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_hiwire_greater_than_equal:0x%x,0x%x", p0, p1)
	return 0
}

// env.descr_set_trampoline_call => wasm function #126
func env_descr_set_trampoline_call(p0 int32, p1 int32, p2 int32, p3 int32) int32 {
	log.Printf("call to --> env_descr_set_trampoline_call:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
	return 0
}

// env.__syscall_poll => wasm function #198
func env___syscall_poll(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env___syscall_poll:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.js2python => wasm function #26
func env_js2python(p0 int32) int32 {
	log.Printf("call to --> env_js2python:0x%x", p0)
	return 0
}

// env.JsArray_Set => wasm function #52
func env_JsArray_Set(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env_JsArray_Set:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.__syscall_renameat => wasm function #208
func env___syscall_renameat(p0 int32, p1 int32, p2 int32, p3 int32) int32 {
	log.Printf("call to --> env___syscall_renameat:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
	return 0
}

// env.invoke_vi => wasm function #240
func env_invoke_vi(p0 int32, p1 int32) {
	log.Printf("call to --> env_invoke_vi:0x%x,0x%x", p0, p1)
}

// env.__syscall_recvmsg => wasm function #278
func env___syscall_recvmsg(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32) int32 {
	log.Printf("call to --> env___syscall_recvmsg:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5)
	return 0
}

// env.destroy_proxies => wasm function #25
func env_destroy_proxies(p0 int32, p1 int32) {
	log.Printf("call to --> env_destroy_proxies:0x%x,0x%x", p0, p1)
}

// wasi_snapshot_preview1.environ_sizes_get => wasm function #146
func wasi_snapshot_preview1_environ_sizes_get(p0 int32, p1 int32) int32 {
	log.Printf("call to --> wasi_snapshot_preview1_environ_sizes_get:0x%x,0x%x", p0, p1)
	return 0
}

// env._mktime_js => wasm function #173
func env__mktime_js(p0 int32) int32 {
	log.Printf("call to --> env__mktime_js:0x%x", p0)
	return 0
}

// env.invoke_viiidi => wasm function #227
func env_invoke_viiidi(p0 int32, p1 int32, p2 int32, p3 int32, p4 float64, p5 int32) {
	log.Printf("call to --> env_invoke_viiidi:0x%x,0x%x,0x%x,0x%x,%4.4f,0x%x", p0, p1, p2, p3, p4, p5)
}

// env.__syscall_connect => wasm function #272
func env___syscall_connect(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32) int32 {
	log.Printf("call to --> env___syscall_connect:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5)
	return 0
}

// env.hiwire_assign_to_ptr => wasm function #31
func env_hiwire_assign_to_ptr(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_hiwire_assign_to_ptr:0x%x,0x%x", p0, p1)
	return 0
}

// env.hiwire_write_to_file => wasm function #91
func env_hiwire_write_to_file(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_hiwire_write_to_file:0x%x,0x%x", p0, p1)
	return 0
}

// env.__syscall_getdents64 => wasm function #191
func env___syscall_getdents64(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env___syscall_getdents64:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.__cxa_decrement_exception_refcount => wasm function #245
func env___cxa_decrement_exception_refcount(p0 int32) {
	log.Printf("call to --> env___cxa_decrement_exception_refcount:0x%x", p0)
}

// env.__syscall_getsockname => wasm function #274
func env___syscall_getsockname(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32) int32 {
	log.Printf("call to --> env___syscall_getsockname:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5)
	return 0
}

// env.JsObject_Entries => wasm function #75
func env_JsObject_Entries(p0 int32) int32 {
	log.Printf("call to --> env_JsObject_Entries:0x%x", p0)
	return 0
}

// env.__syscall_linkat => wasm function #193
func env___syscall_linkat(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32) int32 {
	log.Printf("call to --> env___syscall_linkat:0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4)
	return 0
}

// env.invoke_jii => wasm function #233
func env_invoke_jii(p0 int32, p1 int32, p2 int32) int64 {
	log.Printf("call to --> env_invoke_jii:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.invoke_iiiiii => wasm function #243
func env_invoke_iiiiii(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32) int32 {
	log.Printf("call to --> env_invoke_iiiiii:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5)
	return 0
}

// env.hiwire_call => wasm function #6
func env_hiwire_call(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_hiwire_call:0x%x,0x%x", p0, p1)
	return 0
}

// env.JsObject_DeleteString => wasm function #66
func env_JsObject_DeleteString(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_JsObject_DeleteString:0x%x,0x%x", p0, p1)
	return 0
}

// env.JsDoubleProxy_unwrap_helper => wasm function #80
func env_JsDoubleProxy_unwrap_helper(p0 int32) int32 {
	log.Printf("call to --> env_JsDoubleProxy_unwrap_helper:0x%x", p0)
	return 0
}

// env.exit => wasm function #128
func env_exit(p0 int32) {
	log.Printf("call to --> env_exit:0x%x", p0)
}

// env.JsProxy_compute_typeflags => wasm function #22
func env_JsProxy_compute_typeflags(p0 int32) int32 {
	log.Printf("call to --> env_JsProxy_compute_typeflags:0x%x", p0)
	return 0
}

// env.hiwire_int => wasm function #105
func env_hiwire_int(p0 int32) int32 {
	log.Printf("call to --> env_hiwire_int:0x%x", p0)
	return 0
}

// env._dlopen_js => wasm function #161
func env__dlopen_js(p0 int32) int32 {
	log.Printf("call to --> env__dlopen_js:0x%x", p0)
	return 0
}

// env.__syscall_mkdirat => wasm function #195
func env___syscall_mkdirat(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env___syscall_mkdirat:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.JsProxy_GetIter_js => wasm function #42
func env_JsProxy_GetIter_js(p0 int32) int32 {
	log.Printf("call to --> env_JsProxy_GetIter_js:0x%x", p0)
	return 0
}

// env.hiwire_string_utf8 => wasm function #95
func env_hiwire_string_utf8(p0 int32) int32 {
	log.Printf("call to --> env_hiwire_string_utf8:0x%x", p0)
	return 0
}

// env.pyproxy_new_ex => wasm function #100
func env_pyproxy_new_ex(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env_pyproxy_new_ex:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.gethostbyaddr => wasm function #135
func env_gethostbyaddr(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env_gethostbyaddr:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.JsObject_New => wasm function #29
func env_JsObject_New() int32 {
	log.Printf("call to --> env_JsObject_New:")
	return 0
}

// env._emscripten_dlopen_js => wasm function #162
func env__emscripten_dlopen_js(p0 int32, p1 int32, p2 int32, p3 int32) {
	log.Printf("call to --> env__emscripten_dlopen_js:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
}

// env._mmap_js => wasm function #170
func env__mmap_js(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32, p6 int32) int32 {
	log.Printf("call to --> env__mmap_js:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5, p6)
	return 0
}

// env.__syscall_getcwd => wasm function #190
func env___syscall_getcwd(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env___syscall_getcwd:0x%x,0x%x", p0, p1)
	return 0
}

// env.invoke_diii => wasm function #259
func env_invoke_diii(p0 int32, p1 int32, p2 int32, p3 int32) float64 {
	log.Printf("call to --> env_invoke_diii:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
	return 0
}

// env.js2python_init => wasm function #17
func env_js2python_init() int32 {
	log.Printf("call to --> env_js2python_init:")
	return 0
}

// env.hiwire_assign_from_ptr => wasm function #89
func env_hiwire_assign_from_ptr(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_hiwire_assign_from_ptr:0x%x,0x%x", p0, p1)
	return 0
}

// env.__syscall_fcntl64 => wasm function #148
func env___syscall_fcntl64(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env___syscall_fcntl64:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.__syscall_fadvise64 => wasm function #199
func env___syscall_fadvise64(p0 int32, p1 int64, p2 int64, p3 int32) int32 {
	log.Printf("call to --> env___syscall_fadvise64:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
	return 0
}

// env.hiwire_is_promise => wasm function #46
func env_hiwire_is_promise(p0 int32) int32 {
	log.Printf("call to --> env_hiwire_is_promise:0x%x", p0)
	return 0
}

// env.wrap_async_generator => wasm function #63
func env_wrap_async_generator(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_wrap_async_generator:0x%x,0x%x", p0, p1)
	return 0
}

// env.JsSet_Add => wasm function #121
func env_JsSet_Add(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_JsSet_Add:0x%x,0x%x", p0, p1)
	return 0
}

// env.abort => wasm function #129
func env_abort() {
	log.Printf("call to --> env_abort:")
}

// env.__syscall_listen => wasm function #276
func env___syscall_listen(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32) int32 {
	log.Printf("call to --> env___syscall_listen:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5)
	return 0
}

// env.JsArray_Check => wasm function #55
func env_JsArray_Check(p0 int32) int32 {
	log.Printf("call to --> env_JsArray_Check:0x%x", p0)
	return 0
}

// env.__syscall_fchdir => wasm function #180
func env___syscall_fchdir(p0 int32) int32 {
	log.Printf("call to --> env___syscall_fchdir:0x%x", p0)
	return 0
}

// env.invoke_v => wasm function #234
func env_invoke_v(p0 int32) {
	log.Printf("call to --> env_invoke_v:0x%x", p0)
}

// env.invoke_viijii => wasm function #250
func env_invoke_viijii(p0 int32, p1 int32, p2 int32, p3 int64, p4 int32, p5 int32) {
	log.Printf("call to --> env_invoke_viijii:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5)
}

// env.__syscall_getsockopt => wasm function #275
func env___syscall_getsockopt(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32) int32 {
	log.Printf("call to --> env___syscall_getsockopt:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5)
	return 0
}

// env.JsArray_Push_unchecked => wasm function #5
func env_JsArray_Push_unchecked(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_JsArray_Push_unchecked:0x%x,0x%x", p0, p1)
	return 0
}

// env.hiwire_resolve_promise => wasm function #56
func env_hiwire_resolve_promise(p0 int32) int32 {
	log.Printf("call to --> env_hiwire_resolve_promise:0x%x", p0)
	return 0
}

// env.hiwire_read_from_file => wasm function #92
func env_hiwire_read_from_file(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_hiwire_read_from_file:0x%x,0x%x", p0, p1)
	return 0
}

// env.emscripten_promise_destroy => wasm function #165
func env_emscripten_promise_destroy(p0 int32) {
	log.Printf("call to --> env_emscripten_promise_destroy:0x%x", p0)
}

// env.__syscall_accept4 => wasm function #270
func env___syscall_accept4(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32) int32 {
	log.Printf("call to --> env___syscall_accept4:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5)
	return 0
}

// env._python2js_buffer_inner => wasm function #101
func env__python2js_buffer_inner(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32, p6 int32) int32 {
	log.Printf("call to --> env__python2js_buffer_inner:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5, p6)
	return 0
}

// env.pyproxy_new => wasm function #107
func env_pyproxy_new(p0 int32) int32 {
	log.Printf("call to --> env_pyproxy_new:0x%x", p0)
	return 0
}

// env._python2js_ucs4 => wasm function #111
func env__python2js_ucs4(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env__python2js_ucs4:0x%x,0x%x", p0, p1)
	return 0
}

// env._localtime_js => wasm function #174
func env__localtime_js(p0 int32, p1 int32) {
	log.Printf("call to --> env__localtime_js:0x%x,0x%x", p0, p1)
}

// env.getnameinfo => wasm function #137
func env_getnameinfo(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32, p6 int32) int32 {
	log.Printf("call to --> env_getnameinfo:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5, p6)
	return 0
}

// env.__syscall_dup => wasm function #159
func env___syscall_dup(p0 int32) int32 {
	log.Printf("call to --> env___syscall_dup:0x%x", p0)
	return 0
}

// env.__syscall_stat64 => wasm function #186
func env___syscall_stat64(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env___syscall_stat64:0x%x,0x%x", p0, p1)
	return 0
}

// env.invoke_iiiiij => wasm function #251
func env_invoke_iiiiij(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int64) int32 {
	log.Printf("call to --> env_invoke_iiiiij:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5)
	return 0
}

// env._Py_CheckEmscriptenSignals_Helper => wasm function #132
func env__Py_CheckEmscriptenSignals_Helper() int32 {
	log.Printf("call to --> env__Py_CheckEmscriptenSignals_Helper:")
	return 0
}

// env.ffi_prep_closure_loc_js => wasm function #144
func env_ffi_prep_closure_loc_js(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32) int32 {
	log.Printf("call to --> env_ffi_prep_closure_loc_js:0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4)
	return 0
}

// env._gmtime_js => wasm function #175
func env__gmtime_js(p0 int32, p1 int32) {
	log.Printf("call to --> env__gmtime_js:0x%x,0x%x", p0, p1)
}

// env.invoke_viiiii => wasm function #242
func env_invoke_viiiii(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32) {
	log.Printf("call to --> env_invoke_viiiii:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5)
}

// env.hiwire_is_comlink_proxy => wasm function #35
func env_hiwire_is_comlink_proxy(p0 int32) int32 {
	log.Printf("call to --> env_hiwire_is_comlink_proxy:0x%x", p0)
	return 0
}

// env.JsObject_GetString => wasm function #64
func env_JsObject_GetString(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_JsObject_GetString:0x%x,0x%x", p0, p1)
	return 0
}

// env.__syscall_fdatasync => wasm function #185
func env___syscall_fdatasync(p0 int32) int32 {
	log.Printf("call to --> env___syscall_fdatasync:0x%x", p0)
	return 0
}

// env.__cxa_throw => wasm function #241
func env___cxa_throw(p0 int32, p1 int32, p2 int32) {
	log.Printf("call to --> env___cxa_throw:0x%x,0x%x,0x%x", p0, p1, p2)
}

// env.new_error => wasm function #1
func env_new_error(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env_new_error:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.descr_get_trampoline_call => wasm function #125
func env_descr_get_trampoline_call(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env_descr_get_trampoline_call:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.invoke_viii => wasm function #235
func env_invoke_viii(p0 int32, p1 int32, p2 int32, p3 int32) {
	log.Printf("call to --> env_invoke_viii:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
}

// env.invoke_iiiiiiii => wasm function #254
func env_invoke_iiiiiiii(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32, p6 int32, p7 int32) int32 {
	log.Printf("call to --> env_invoke_iiiiiiii:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5, p6, p7)
	return 0
}

// env.hiwire_constructor_name => wasm function #33
func env_hiwire_constructor_name(p0 int32) int32 {
	log.Printf("call to --> env_hiwire_constructor_name:0x%x", p0)
	return 0
}

// env._emscripten_get_progname => wasm function #202
func env__emscripten_get_progname(p0 int32, p1 int32) {
	log.Printf("call to --> env__emscripten_get_progname:0x%x,0x%x", p0, p1)
}

// env.__cxa_end_catch => wasm function #238
func env___cxa_end_catch() {
	log.Printf("call to --> env___cxa_end_catch:")
}

// env.invoke_iiiiijj => wasm function #253
func env_invoke_iiiiijj(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int64, p6 int64) int32 {
	log.Printf("call to --> env_invoke_iiiiijj:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5, p6)
	return 0
}

// env.__syscall_getpeername => wasm function #273
func env___syscall_getpeername(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32) int32 {
	log.Printf("call to --> env___syscall_getpeername:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5)
	return 0
}

// env.JsSet_New => wasm function #120
func env_JsSet_New() int32 {
	log.Printf("call to --> env_JsSet_New:")
	return 0
}

// env.__syscall_fstat64 => wasm function #151
func env___syscall_fstat64(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env___syscall_fstat64:0x%x,0x%x", p0, p1)
	return 0
}

// env.emscripten_get_now_res => wasm function #179
func env_emscripten_get_now_res() float64 {
	log.Printf("call to --> env_emscripten_get_now_res:")
	return 0
}

// env.__syscall_fstatfs64 => wasm function #212
func env___syscall_fstatfs64(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env___syscall_fstatfs64:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.invoke_viif => wasm function #267
func env_invoke_viif(p0 int32, p1 int32, p2 int32, p3 float32) {
	log.Printf("call to --> env_invoke_viif:0x%x,0x%x,0x%x,%4.4f", p0, p1, p2, p3)
}

// env.console_error_obj => wasm function #94
func env_console_error_obj(p0 int32) {
	log.Printf("call to --> env_console_error_obj:0x%x", p0)
}

// env._JsArray_PushEntry_helper => wasm function #117
func env__JsArray_PushEntry_helper(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env__JsArray_PushEntry_helper:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.__syscall__newselect => wasm function #209
func env___syscall__newselect(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32) int32 {
	log.Printf("call to --> env___syscall__newselect:0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4)
	return 0
}

// env.invoke_iiiiiiiiiii => wasm function #255
func env_invoke_iiiiiiiiiii(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32, p6 int32, p7 int32, p8 int32, p9 int32, p10 int32) int32 {
	log.Printf("call to --> env_invoke_iiiiiiiiiii:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10)
	return 0
}

// env.invoke_viid => wasm function #268
func env_invoke_viid(p0 int32, p1 int32, p2 int32, p3 float64) {
	log.Printf("call to --> env_invoke_viid:0x%x,0x%x,0x%x,%4.4f", p0, p1, p2, p3)
}

// env.js2python_immutable => wasm function #21
func env_js2python_immutable(p0 int32) int32 {
	log.Printf("call to --> env_js2python_immutable:0x%x", p0)
	return 0
}

// env.JsArray_Splice => wasm function #84
func env_JsArray_Splice(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_JsArray_Splice:0x%x,0x%x", p0, p1)
	return 0
}

// env.emscripten_promise_resolve => wasm function #164
func env_emscripten_promise_resolve(p0 int32, p1 int32, p2 int32) {
	log.Printf("call to --> env_emscripten_promise_resolve:0x%x,0x%x,0x%x", p0, p1, p2)
}

// env.invoke_iiiii => wasm function #239
func env_invoke_iiiii(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32) int32 {
	log.Printf("call to --> env_invoke_iiiii:0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4)
	return 0
}

// env.get_async_js_call_done_callback => wasm function #62
func env_get_async_js_call_done_callback(p0 int32) int32 {
	log.Printf("call to --> env_get_async_js_call_done_callback:0x%x", p0)
	return 0
}

// env.ffi_call_js => wasm function #141
func env_ffi_call_js(p0 int32, p1 int32, p2 int32, p3 int32) {
	log.Printf("call to --> env_ffi_call_js:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
}

// wasi_snapshot_preview1.proc_exit => wasm function #145
func wasi_snapshot_preview1_proc_exit(p0 int32) {
	panic("exit called " + fmt.Sprint(p0))
	log.Printf("call to --> wasi_snapshot_preview1_proc_exit:0x%x", p0)
}

// env._tzset_js => wasm function #171
func env__tzset_js(p0 int32, p1 int32, p2 int32) {
	log.Printf("call to --> env__tzset_js:0x%x,0x%x,0x%x", p0, p1, p2)
}

// env.JsObjMap_subscript_js => wasm function #38
func env_JsObjMap_subscript_js(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_JsObjMap_subscript_js:0x%x,0x%x", p0, p1)
	return 0
}

// env.hiwire_to_bool => wasm function #45
func env_hiwire_to_bool(p0 int32) int32 {
	log.Printf("call to --> env_hiwire_to_bool:0x%x", p0)
	return 0
}

// env.JsArray_reverse_helper => wasm function #82
func env_JsArray_reverse_helper(p0 int32) int32 {
	log.Printf("call to --> env_JsArray_reverse_helper:0x%x", p0)
	return 0
}

// env.system => wasm function #138
func env_system(p0 int32) int32 {
	log.Printf("call to --> env_system:0x%x", p0)
	return 0
}

// env.JsProxy_subscript_js => wasm function #44
func env_JsProxy_subscript_js(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_JsProxy_subscript_js:0x%x,0x%x", p0, p1)
	return 0
}

// env.pyproxy_AsPyObject => wasm function #85
func env_pyproxy_AsPyObject(p0 int32) int32 {
	log.Printf("call to --> env_pyproxy_AsPyObject:0x%x", p0)
	return 0
}

// env.JsArray_count_helper => wasm function #87
func env_JsArray_count_helper(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_JsArray_count_helper:0x%x,0x%x", p0, p1)
	return 0
}

// env._Py_emscripten_runtime => wasm function #131
func env__Py_emscripten_runtime() int32 {
	log.Printf("call to --> env__Py_emscripten_runtime:")
	return 0
}

// env.invoke_viiiiiii => wasm function #260
func env_invoke_viiiiiii(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32, p6 int32, p7 int32) {
	log.Printf("call to --> env_invoke_viiiiiii:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5, p6, p7)
}

// env.python2js_buffer_init => wasm function #18
func env_python2js_buffer_init() int32 {
	log.Printf("call to --> env_python2js_buffer_init:")
	return 0
}

// env.strftime => wasm function #140
func env_strftime(p0 int32, p1 int32, p2 int32, p3 int32) int32 {
	log.Printf("call to --> env_strftime:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
	return 0
}

// env.ffi_closure_free_js => wasm function #143
func env_ffi_closure_free_js(p0 int32) {
	log.Printf("call to --> env_ffi_closure_free_js:0x%x", p0)
}

// env.__syscall_ftruncate64 => wasm function #189
func env___syscall_ftruncate64(p0 int32, p1 int64) int32 {
	log.Printf("call to --> env___syscall_ftruncate64:0x%x,0x%x", p0, p1)
	return 0
}

// env.hiwire_CallMethodString => wasm function #11
func env_hiwire_CallMethodString(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env_hiwire_CallMethodString:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env._PyImport_InitFunc_TrampolineCall => wasm function #127
func env__PyImport_InitFunc_TrampolineCall(p0 int32) int32 {
	log.Printf("call to --> env__PyImport_InitFunc_TrampolineCall:0x%x", p0)
	return 0
}

// env.getaddrinfo => wasm function #133
func env_getaddrinfo(p0 int32, p1 int32, p2 int32, p3 int32) int32 {
	log.Printf("call to --> env_getaddrinfo:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
	return 0
}

// env.emscripten_get_heap_max => wasm function #215
func env_emscripten_get_heap_max() int32 {
	log.Printf("call to --> env_emscripten_get_heap_max:")
	return 0
}

// env.JsObject_Keys => wasm function #76
func env_JsObject_Keys(p0 int32) int32 {
	log.Printf("call to --> env_JsObject_Keys:0x%x", p0)
	return 0
}

// env.proxy_cache_get => wasm function #96
func env_proxy_cache_get(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_proxy_cache_get:0x%x,0x%x", p0, p1)
	return 0
}

// env.__syscall_unlinkat => wasm function #206
func env___syscall_unlinkat(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env___syscall_unlinkat:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.__syscall_statfs64 => wasm function #211
func env___syscall_statfs64(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env___syscall_statfs64:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.destroy_proxy => wasm function #24
func env_destroy_proxy(p0 int32, p1 int32) {
	log.Printf("call to --> env_destroy_proxy:0x%x,0x%x", p0, p1)
}

// env.create_promise_handles => wasm function #57
func env_create_promise_handles(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env_create_promise_handles:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.JsMap_clear_js => wasm function #88
func env_JsMap_clear_js(p0 int32) int32 {
	log.Printf("call to --> env_JsMap_clear_js:0x%x", p0)
	return 0
}

// env.invoke_iii => wasm function #222
func env_invoke_iii(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env_invoke_iii:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.invoke_viijj => wasm function #269
func env_invoke_viijj(p0 int32, p1 int32, p2 int32, p3 int64, p4 int64) {
	log.Printf("call to --> env_invoke_viijj:0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4)
}

// env.__syscall_faccessat => wasm function #155
func env___syscall_faccessat(p0 int32, p1 int32, p2 int32, p3 int32) int32 {
	log.Printf("call to --> env___syscall_faccessat:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
	return 0
}

// env._dlsym_js => wasm function #166
func env__dlsym_js(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env__dlsym_js:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.__cxa_increment_exception_refcount => wasm function #246
func env___cxa_increment_exception_refcount(p0 int32) {
	log.Printf("call to --> env___cxa_increment_exception_refcount:0x%x", p0)
}

// env.invoke_iiiiid => wasm function #252
func env_invoke_iiiiid(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 float64) int32 {
	log.Printf("call to --> env_invoke_iiiiid:0x%x,0x%x,0x%x,0x%x,0x%x,%4.4f", p0, p1, p2, p3, p4, p5)
	return 0
}

// env.hiwire_CallMethod_OneArg => wasm function #13
func env_hiwire_CallMethod_OneArg(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env_hiwire_CallMethod_OneArg:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.handle_next_result_js => wasm function #20
func env_handle_next_result_js(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env_handle_next_result_js:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.hiwire_is_generator => wasm function #59
func env_hiwire_is_generator(p0 int32) int32 {
	log.Printf("call to --> env_hiwire_is_generator:0x%x", p0)
	return 0
}

// env.__syscall_symlinkat => wasm function #214
func env___syscall_symlinkat(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env___syscall_symlinkat:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.__syscall_sendto => wasm function #280
func env___syscall_sendto(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32) int32 {
	log.Printf("call to --> env___syscall_sendto:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5)
	return 0
}

// env.JsArray_Push => wasm function #28
func env_JsArray_Push(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_JsArray_Push:0x%x,0x%x", p0, p1)
	return 0
}

// env.hiwire_construct => wasm function #79
func env_hiwire_construct(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_hiwire_construct:0x%x,0x%x", p0, p1)
	return 0
}

// env.JsArray_index_helper => wasm function #86
func env_JsArray_index_helper(p0 int32, p1 int32, p2 int32, p3 int32) int32 {
	log.Printf("call to --> env_JsArray_index_helper:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
	return 0
}

// env.__syscall_fchmodat => wasm function #182
func env___syscall_fchmodat(p0 int32, p1 int32, p2 int32, p3 int32) int32 {
	log.Printf("call to --> env___syscall_fchmodat:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
	return 0
}

// env.JsObjMap_length_js => wasm function #37
func env_JsObjMap_length_js(p0 int32) int32 {
	log.Printf("call to --> env_JsObjMap_length_js:0x%x", p0)
	return 0
}

// env.hiwire_not_equal => wasm function #70
func env_hiwire_not_equal(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_hiwire_not_equal:0x%x,0x%x", p0, p1)
	return 0
}

// env.emscripten_date_now => wasm function #176
func env_emscripten_date_now() float64 {
	log.Printf("call to --> env_emscripten_date_now:")
	return 0
}

// env.invoke_viiifi => wasm function #226
func env_invoke_viiifi(p0 int32, p1 int32, p2 int32, p3 int32, p4 float32, p5 int32) {
	log.Printf("call to --> env_invoke_viiifi:0x%x,0x%x,0x%x,0x%x,%4.4f,0x%x", p0, p1, p2, p3, p4, p5)
}

// env.JsObjMap_GetIter_js => wasm function #36
func env_JsObjMap_GetIter_js(p0 int32) int32 {
	log.Printf("call to --> env_JsObjMap_GetIter_js:0x%x", p0)
	return 0
}

// env._JsArray_PostProcess_helper => wasm function #118
func env__JsArray_PostProcess_helper(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env__JsArray_PostProcess_helper:0x%x,0x%x", p0, p1)
	return 0
}

// env.gethostbyname => wasm function #134
func env_gethostbyname(p0 int32) int32 {
	log.Printf("call to --> env_gethostbyname:0x%x", p0)
	return 0
}

// wasi_snapshot_preview1.fd_read => wasm function #153
func wasi_snapshot_preview1_fd_read(p0 int32, p1 int32, p2 int32, p3 int32) int32 {
	log.Printf("call to --> wasi_snapshot_preview1_fd_read:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
	return 0
}

// env._agen_handle_result_js => wasm function #27
func env__agen_handle_result_js(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32) int32 {
	log.Printf("call to --> env__agen_handle_result_js:0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4)
	return 0
}

// env.JsObject_Values => wasm function #77
func env_JsObject_Values(p0 int32) int32 {
	log.Printf("call to --> env_JsObject_Values:0x%x", p0)
	return 0
}

// env._python2js_cache_lookup => wasm function #102
func env__python2js_cache_lookup(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env__python2js_cache_lookup:0x%x,0x%x", p0, p1)
	return 0
}

// env.__syscall_truncate64 => wasm function #216
func env___syscall_truncate64(p0 int32, p1 int64) int32 {
	log.Printf("call to --> env___syscall_truncate64:0x%x,0x%x", p0, p1)
	return 0
}

// env.JsString_InternFromCString => wasm function #8
func env_JsString_InternFromCString(p0 int32) int32 {
	log.Printf("call to --> env_JsString_InternFromCString:0x%x", p0)
	return 0
}

// env.proxy_cache_set => wasm function #97
func env_proxy_cache_set(p0 int32, p1 int32, p2 int32) {
	log.Printf("call to --> env_proxy_cache_set:0x%x,0x%x,0x%x", p0, p1, p2)
}

// env.emscripten_promise_create => wasm function #163
func env_emscripten_promise_create() int32 {
	log.Printf("call to --> env_emscripten_promise_create:")
	return 0
}

// env.invoke_ii => wasm function #229
func env_invoke_ii(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_invoke_ii:0x%x,0x%x", p0, p1)
	return 0
}

// env.__syscall_chmod => wasm function #157
func env___syscall_chmod(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env___syscall_chmod:0x%x,0x%x", p0, p1)
	return 0
}

// env._munmap_js => wasm function #168
func env__munmap_js(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32) int32 {
	log.Printf("call to --> env__munmap_js:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5)
	return 0
}

// env.invoke_i => wasm function #225
func env_invoke_i(p0 int32) int32 {
	log.Printf("call to --> env_invoke_i:0x%x", p0)
	return 0
}

// env.__cxa_begin_catch => wasm function #236
func env___cxa_begin_catch(p0 int32) int32 {
	log.Printf("call to --> env___cxa_begin_catch:0x%x", p0)
	return 0
}

// env.__syscall_sendmsg => wasm function #279
func env___syscall_sendmsg(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32) int32 {
	log.Printf("call to --> env___syscall_sendmsg:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5)
	return 0
}

// env.create_once_callable => wasm function #81
func env_create_once_callable(p0 int32) int32 {
	log.Printf("call to --> env_create_once_callable:0x%x", p0)
	return 0
}

// env.hiwire_into_file => wasm function #93
func env_hiwire_into_file(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_hiwire_into_file:0x%x,0x%x", p0, p1)
	return 0
}

// env._PyCFunctionWithKeywords_TrampolineCall => wasm function #124
func env__PyCFunctionWithKeywords_TrampolineCall(p0 int32, p1 int32, p2 int32, p3 int32) int32 {
	log.Printf("call to --> env__PyCFunctionWithKeywords_TrampolineCall:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
	return 0
}

// env.invoke_j => wasm function #230
func env_invoke_j(p0 int32) int64 {
	log.Printf("call to --> env_invoke_j:0x%x", p0)
	return 0
}

// env.hiwire_throw_error => wasm function #2
func env_hiwire_throw_error(p0 int32) {
	log.Printf("call to --> env_hiwire_throw_error:0x%x", p0)
}

// env._msync_js => wasm function #169
func env__msync_js(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32) int32 {
	log.Printf("call to --> env__msync_js:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5)
	return 0
}

// env.__call_sighandler => wasm function #204
func env___call_sighandler(p0 int32, p1 int32) {
	log.Printf("call to --> env___call_sighandler:0x%x,0x%x", p0, p1)
}

// env.__resumeException => wasm function #224
func env___resumeException(p0 int32) {
	log.Printf("call to --> env___resumeException:0x%x", p0)
}

// env.invoke_viiiiiiiiii => wasm function #263
func env_invoke_viiiiiiiiii(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32, p6 int32, p7 int32, p8 int32, p9 int32, p10 int32) {
	log.Printf("call to --> env_invoke_viiiiiiiiii:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10)
}

// env.JsObjMap_contains_js => wasm function #40
func env_JsObjMap_contains_js(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_JsObjMap_contains_js:0x%x,0x%x", p0, p1)
	return 0
}

// env._python2js_add_to_cache => wasm function #104
func env__python2js_add_to_cache(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env__python2js_add_to_cache:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.getprotobyname => wasm function #136
func env_getprotobyname(p0 int32) int32 {
	log.Printf("call to --> env_getprotobyname:0x%x", p0)
	return 0
}

// env.__syscall_symlink => wasm function #213
func env___syscall_symlink(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env___syscall_symlink:0x%x,0x%x", p0, p1)
	return 0
}

// env.hiwire_init => wasm function #16
func env_hiwire_init() int32 {
	log.Printf("call to --> env_hiwire_init:")
	return 0
}

// env.array_to_js => wasm function #99
func env_array_to_js(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_array_to_js:0x%x,0x%x", p0, p1)
	return 0
}

// env.hiwire_double => wasm function #108
func env_hiwire_double(p0 float64) int32 {
	log.Printf("call to --> env_hiwire_double:%4.4f", p0)
	return 0
}

// env.invoke_vii => wasm function #232
func env_invoke_vii(p0 int32, p1 int32, p2 int32) {
	log.Printf("call to --> env_invoke_vii:0x%x,0x%x,0x%x", p0, p1, p2)
}

// env.invoke_iiiiiiiiiiiii => wasm function #257
func env_invoke_iiiiiiiiiiiii(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32, p6 int32, p7 int32, p8 int32, p9 int32, p10 int32, p11 int32, p12 int32) int32 {
	log.Printf("call to --> env_invoke_iiiiiiiiiiiii:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12)
	return 0
}

// env.__syscall_ioctl => wasm function #149
func env___syscall_ioctl(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env___syscall_ioctl:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.emscripten_get_now => wasm function #178
func env_emscripten_get_now() float64 {
	log.Printf("call to --> env_emscripten_get_now:")
	return 0
}

// wasi_snapshot_preview1.fd_fdstat_get => wasm function #192
func wasi_snapshot_preview1_fd_fdstat_get(p0 int32, p1 int32) int32 {
	log.Printf("call to --> wasi_snapshot_preview1_fd_fdstat_get:0x%x,0x%x", p0, p1)
	return 0
}

// env.__syscall_pipe => wasm function #197
func env___syscall_pipe(p0 int32) int32 {
	log.Printf("call to --> env___syscall_pipe:0x%x", p0)
	return 0
}

// env.hiwire_decref => wasm function #7
func env_hiwire_decref(p0 int32) {
	log.Printf("call to --> env_hiwire_decref:0x%x", p0)
}

// env.hiwire_less_than_equal => wasm function #68
func env_hiwire_less_than_equal(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_hiwire_less_than_equal:0x%x,0x%x", p0, p1)
	return 0
}

// env.hiwire_reversed_iterator => wasm function #83
func env_hiwire_reversed_iterator(p0 int32) int32 {
	log.Printf("call to --> env_hiwire_reversed_iterator:0x%x", p0)
	return 0
}

// env.JsMap_Set => wasm function #116
func env_JsMap_Set(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env_JsMap_Set:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.invoke_fiii => wasm function #258
func env_invoke_fiii(p0 int32, p1 int32, p2 int32, p3 int32) float32 {
	log.Printf("call to --> env_invoke_fiii:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
	return 0
}

// env.JsObject_Dir => wasm function #54
func env_JsObject_Dir(p0 int32) int32 {
	log.Printf("call to --> env_JsObject_Dir:0x%x", p0)
	return 0
}

// env.hiwire_is_function => wasm function #65
func env_hiwire_is_function(p0 int32) int32 {
	log.Printf("call to --> env_hiwire_is_function:0x%x", p0)
	return 0
}

// env.python2js_custom__create_jscontext => wasm function #113
func env_python2js_custom__create_jscontext(p0 int32, p1 int32, p2 int32, p3 int32) int32 {
	log.Printf("call to --> env_python2js_custom__create_jscontext:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
	return 0
}

// env.getentropy => wasm function #130
func env_getentropy(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_getentropy:0x%x,0x%x", p0, p1)
	return 0
}

// env.fail_test => wasm function #0
func env_fail_test() {
	log.Printf("call to --> env_fail_test:")
}

// env.JsObject_SetString => wasm function #30
func env_JsObject_SetString(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env_JsObject_SetString:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.python2js__default_converter => wasm function #103
func env_python2js__default_converter(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_python2js__default_converter:0x%x,0x%x", p0, p1)
	return 0
}

// env.__cxa_find_matching_catch_3 => wasm function #221
func env___cxa_find_matching_catch_3(p0 int32) int32 {
	log.Printf("call to --> env___cxa_find_matching_catch_3:0x%x", p0)
	return 0
}

// env.invoke_viiiiii => wasm function #265
func env_invoke_viiiiii(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32, p6 int32) {
	log.Printf("call to --> env_invoke_viiiiii:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5, p6)
}

// env.hiwire_HasMethod => wasm function #9
func env_hiwire_HasMethod(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_hiwire_HasMethod:0x%x,0x%x", p0, p1)
	return 0
}

// env.JsProxy_GetAsyncIter_js => wasm function #43
func env_JsProxy_GetAsyncIter_js(p0 int32) int32 {
	log.Printf("call to --> env_JsProxy_GetAsyncIter_js:0x%x", p0)
	return 0
}

// env.__syscall_rmdir => wasm function #207
func env___syscall_rmdir(p0 int32) int32 {
	log.Printf("call to --> env___syscall_rmdir:0x%x", p0)
	return 0
}

// env.invoke_ji => wasm function #231
func env_invoke_ji(p0 int32, p1 int32) int64 {
	log.Printf("call to --> env_invoke_ji:0x%x,0x%x", p0, p1)
	return 0
}

// env.invoke_jiiii => wasm function #256
func env_invoke_jiiii(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32) int64 {
	log.Printf("call to --> env_invoke_jiiii:0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4)
	return 0
}

// env.hiwire_incref => wasm function #34
func env_hiwire_incref(p0 int32) int32 {
	log.Printf("call to --> env_hiwire_incref:0x%x", p0)
	return 0
}

// env._emscripten_get_now_is_monotonic => wasm function #177
func env__emscripten_get_now_is_monotonic() int32 {
	log.Printf("call to --> env__emscripten_get_now_is_monotonic:")
	return 0
}

// env.__syscall_utimensat => wasm function #217
func env___syscall_utimensat(p0 int32, p1 int32, p2 int32, p3 int32) int32 {
	log.Printf("call to --> env___syscall_utimensat:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
	return 0
}

// env.__cxa_uncaught_exceptions => wasm function #244
func env___cxa_uncaught_exceptions() int32 {
	log.Printf("call to --> env___cxa_uncaught_exceptions:")
	return 0
}

// env.emscripten_asm_const_int => wasm function #3
func env_emscripten_asm_const_int(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env_emscripten_asm_const_int:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.hiwire_is_async_generator => wasm function #61
func env_hiwire_is_async_generator(p0 int32) int32 {
	log.Printf("call to --> env_hiwire_is_async_generator:0x%x", p0)
	return 0
}

// env.JsBuffer_DecodeString_js => wasm function #90
func env_JsBuffer_DecodeString_js(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_JsBuffer_DecodeString_js:0x%x,0x%x", p0, p1)
	return 0
}

// env.invoke_viiii => wasm function #228
func env_invoke_viiii(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32) {
	log.Printf("call to --> env_invoke_viiii:0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4)
}

// env._python2js_ucs1 => wasm function #109
func env__python2js_ucs1(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env__python2js_ucs1:0x%x,0x%x", p0, p1)
	return 0
}

// env.JsMap_New => wasm function #112
func env_JsMap_New() int32 {
	log.Printf("call to --> env_JsMap_New:")
	return 0
}

// env.__syscall_fchown32 => wasm function #184
func env___syscall_fchown32(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env___syscall_fchown32:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.__cxa_rethrow_primary_exception => wasm function #248
func env___cxa_rethrow_primary_exception(p0 int32) {
	log.Printf("call to --> env___cxa_rethrow_primary_exception:0x%x", p0)
}

// env.JsMap_GetIter_js => wasm function #41
func env_JsMap_GetIter_js(p0 int32) int32 {
	log.Printf("call to --> env_JsMap_GetIter_js:0x%x", p0)
	return 0
}

// env._python2js_addto_postprocess_list => wasm function #106
func env__python2js_addto_postprocess_list(p0 int32, p1 int32, p2 int32, p3 int32) {
	log.Printf("call to --> env__python2js_addto_postprocess_list:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
}

// env.__syscall_chdir => wasm function #156
func env___syscall_chdir(p0 int32) int32 {
	log.Printf("call to --> env___syscall_chdir:0x%x", p0)
	return 0
}

// env.__syscall_fallocate => wasm function #200
func env___syscall_fallocate(p0 int32, p1 int32, p2 int64, p3 int64) int32 {
	log.Printf("call to --> env___syscall_fallocate:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
	return 0
}

// env.invoke_viiiiiiiiiiiiiii => wasm function #264
func env_invoke_viiiiiiiiiiiiiii(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32, p6 int32, p7 int32, p8 int32, p9 int32, p10 int32, p11 int32, p12 int32, p13 int32, p14 int32, p15 int32) {
	log.Printf("call to --> env_invoke_viiiiiiiiiiiiiii:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11, p12, p13, p14, p15)
}

// env.hiwire_less_than => wasm function #67
func env_hiwire_less_than(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_hiwire_less_than:0x%x,0x%x", p0, p1)
	return 0
}

// env.hiwire_equal => wasm function #69
func env_hiwire_equal(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_hiwire_equal:0x%x,0x%x", p0, p1)
	return 0
}

// wasi_snapshot_preview1.fd_close => wasm function #152
func wasi_snapshot_preview1_fd_close(p0 int32) int32 {
	log.Printf("call to --> wasi_snapshot_preview1_fd_close:0x%x", p0)
	return 0
}

// env.__syscall_mknodat => wasm function #196
func env___syscall_mknodat(p0 int32, p1 int32, p2 int32, p3 int32) int32 {
	log.Printf("call to --> env___syscall_mknodat:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
	return 0
}

// env.hiwire_get_buffer_info => wasm function #32
func env_hiwire_get_buffer_info(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32) {
	log.Printf("call to --> env_hiwire_get_buffer_info:0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4)
}

// env.JsArray_slice => wasm function #49
func env_JsArray_slice(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32) int32 {
	log.Printf("call to --> env_JsArray_slice:0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4)
	return 0
}

// env.JsArray_Delete => wasm function #51
func env_JsArray_Delete(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_JsArray_Delete:0x%x,0x%x", p0, p1)
	return 0
}

// env._emscripten_throw_longjmp => wasm function #219
func env__emscripten_throw_longjmp() {
	log.Printf("call to --> env__emscripten_throw_longjmp:")
}

// env.JsObjMap_ass_subscript_js => wasm function #39
func env_JsObjMap_ass_subscript_js(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env_JsObjMap_ass_subscript_js:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env._python2js_ucs2 => wasm function #110
func env__python2js_ucs2(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env__python2js_ucs2:0x%x,0x%x", p0, p1)
	return 0
}

// env._timegm_js => wasm function #172
func env__timegm_js(p0 int32) int32 {
	log.Printf("call to --> env__timegm_js:0x%x", p0)
	return 0
}

// env._setitimer_js => wasm function #210
func env__setitimer_js(p0 int32, p1 float64) int32 {
	log.Printf("call to --> env__setitimer_js:0x%x,%4.4f", p0, p1)
	return 0
}

// env.pyproxy_Check => wasm function #23
func env_pyproxy_Check(p0 int32) int32 {
	log.Printf("call to --> env_pyproxy_Check:0x%x", p0)
	return 0
}

// env.hiwire_to_string => wasm function #47
func env_hiwire_to_string(p0 int32) int32 {
	log.Printf("call to --> env_hiwire_to_string:0x%x", p0)
	return 0
}

// env.JsArray_slice_assign => wasm function #50
func env_JsArray_slice_assign(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32, p6 int32) int32 {
	log.Printf("call to --> env_JsArray_slice_assign:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5, p6)
	return 0
}

// env.wrap_generator => wasm function #60
func env_wrap_generator(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_wrap_generator:0x%x,0x%x", p0, p1)
	return 0
}

// env.hiwire_int_from_digits => wasm function #119
func env_hiwire_int_from_digits(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_hiwire_int_from_digits:0x%x,0x%x", p0, p1)
	return 0
}

// env.getloadavg => wasm function #139
func env_getloadavg(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_getloadavg:0x%x,0x%x", p0, p1)
	return 0
}

// wasi_snapshot_preview1.fd_write => wasm function #154
func wasi_snapshot_preview1_fd_write(p0 int32, p1 int32, p2 int32, p3 int32) int32 {
	log.Printf("call to --> wasi_snapshot_preview1_fd_write:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
	return 0
}

// env.invoke_iiii => wasm function #220
func env_invoke_iiii(p0 int32, p1 int32, p2 int32, p3 int32) int32 {
	log.Printf("call to --> env_invoke_iiii:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
	return 0
}

// env.__syscall_socket => wasm function #281
func env___syscall_socket(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32) int32 {
	log.Printf("call to --> env___syscall_socket:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5)
	return 0
}

// env.hiwire_typeof => wasm function #74
func env_hiwire_typeof(p0 int32) int32 {
	log.Printf("call to --> env_hiwire_typeof:0x%x", p0)
	return 0
}

// env.__syscall_dup3 => wasm function #160
func env___syscall_dup3(p0 int32, p1 int32, p2 int32) int32 {
	log.Printf("call to --> env___syscall_dup3:0x%x,0x%x,0x%x", p0, p1, p2)
	return 0
}

// env.__syscall_lstat64 => wasm function #187
func env___syscall_lstat64(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env___syscall_lstat64:0x%x,0x%x", p0, p1)
	return 0
}

// env.__syscall_readlinkat => wasm function #205
func env___syscall_readlinkat(p0 int32, p1 int32, p2 int32, p3 int32) int32 {
	log.Printf("call to --> env___syscall_readlinkat:0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3)
	return 0
}

// env.__syscall_bind => wasm function #271
func env___syscall_bind(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32) int32 {
	log.Printf("call to --> env___syscall_bind:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5)
	return 0
}

// env.hiwire_call_OneArg => wasm function #98
func env_hiwire_call_OneArg(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env_hiwire_call_OneArg:0x%x,0x%x", p0, p1)
	return 0
}

// env.__syscall_fchownat => wasm function #158
func env___syscall_fchownat(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32) int32 {
	log.Printf("call to --> env___syscall_fchownat:0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4)
	return 0
}

// wasi_snapshot_preview1.fd_sync => wasm function #188
func wasi_snapshot_preview1_fd_sync(p0 int32) int32 {
	log.Printf("call to --> wasi_snapshot_preview1_fd_sync:0x%x", p0)
	return 0
}

// env.emscripten_resize_heap => wasm function #218
func env_emscripten_resize_heap(p0 int32) int32 {
	log.Printf("call to --> env_emscripten_resize_heap:0x%x", p0)
	return 0
}

// env._python2js_handle_postprocess_list => wasm function #114
func env__python2js_handle_postprocess_list(p0 int32, p1 int32) {
	log.Printf("call to --> env__python2js_handle_postprocess_list:0x%x,0x%x", p0, p1)
}

// env.__syscall_fchmod => wasm function #181
func env___syscall_fchmod(p0 int32, p1 int32) int32 {
	log.Printf("call to --> env___syscall_fchmod:0x%x,0x%x", p0, p1)
	return 0
}

// wasi_snapshot_preview1.fd_pread => wasm function #201
func wasi_snapshot_preview1_fd_pread(p0 int32, p1 int32, p2 int32, p3 int64, p4 int32) int32 {
	log.Printf("call to --> wasi_snapshot_preview1_fd_pread:0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4)
	return 0
}

// env.__cxa_rethrow => wasm function #237
func env___cxa_rethrow() {
	log.Printf("call to --> env___cxa_rethrow:")
}

// env.invoke_iiiiiiiiiiii => wasm function #262
func env_invoke_iiiiiiiiiiii(p0 int32, p1 int32, p2 int32, p3 int32, p4 int32, p5 int32, p6 int32, p7 int32, p8 int32, p9 int32, p10 int32, p11 int32) int32 {
	log.Printf("call to --> env_invoke_iiiiiiiiiiii:0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x,0x%x", p0, p1, p2, p3, p4, p5, p6, p7, p8, p9, p10, p11)
	return 0
}

func addEmscriptenFuncs(store wasmtime.Storelike, result map[string]*wasmtime.Func, rt *Runtime) {
	result["env.destroy_proxies"] = wasmtime.WrapFunc(store, env_destroy_proxies)
	result["wasi_snapshot_preview1.environ_sizes_get"] = wasmtime.WrapFunc(store, wasi_snapshot_preview1_environ_sizes_get)
	result["env._mktime_js"] = wasmtime.WrapFunc(store, env__mktime_js)
	result["env.invoke_viiidi"] = wasmtime.WrapFunc(store, env_invoke_viiidi)
	result["env.__syscall_connect"] = wasmtime.WrapFunc(store, env___syscall_connect)
	result["env.__syscall_recvmsg"] = wasmtime.WrapFunc(store, env___syscall_recvmsg)
	result["env.hiwire_assign_to_ptr"] = wasmtime.WrapFunc(store, env_hiwire_assign_to_ptr)
	result["env.hiwire_write_to_file"] = wasmtime.WrapFunc(store, env_hiwire_write_to_file)
	result["env.__syscall_getdents64"] = wasmtime.WrapFunc(store, env___syscall_getdents64)
	result["env.__cxa_decrement_exception_refcount"] = wasmtime.WrapFunc(store, env___cxa_decrement_exception_refcount)
	result["env.__syscall_getsockname"] = wasmtime.WrapFunc(store, env___syscall_getsockname)
	result["env.JsObject_Entries"] = wasmtime.WrapFunc(store, env_JsObject_Entries)
	result["env.__syscall_linkat"] = wasmtime.WrapFunc(store, env___syscall_linkat)
	result["env.invoke_jii"] = wasmtime.WrapFunc(store, env_invoke_jii)
	result["env.invoke_iiiiii"] = wasmtime.WrapFunc(store, env_invoke_iiiiii)
	result["env.hiwire_call"] = wasmtime.WrapFunc(store, env_hiwire_call)
	result["env.JsObject_DeleteString"] = wasmtime.WrapFunc(store, env_JsObject_DeleteString)
	result["env.JsDoubleProxy_unwrap_helper"] = wasmtime.WrapFunc(store, env_JsDoubleProxy_unwrap_helper)
	result["env.exit"] = wasmtime.WrapFunc(store, env_exit)
	result["env.JsProxy_compute_typeflags"] = wasmtime.WrapFunc(store, env_JsProxy_compute_typeflags)
	result["env.hiwire_int"] = wasmtime.WrapFunc(store, env_hiwire_int)
	result["env._dlopen_js"] = wasmtime.WrapFunc(store, env__dlopen_js)
	result["env.__syscall_mkdirat"] = wasmtime.WrapFunc(store, env___syscall_mkdirat)
	result["env.JsProxy_GetIter_js"] = wasmtime.WrapFunc(store, env_JsProxy_GetIter_js)
	result["env.hiwire_string_utf8"] = wasmtime.WrapFunc(store, env_hiwire_string_utf8)
	result["env.pyproxy_new_ex"] = wasmtime.WrapFunc(store, env_pyproxy_new_ex)
	result["env.gethostbyaddr"] = wasmtime.WrapFunc(store, env_gethostbyaddr)
	result["env.JsObject_New"] = wasmtime.WrapFunc(store, env_JsObject_New)
	result["env._emscripten_dlopen_js"] = wasmtime.WrapFunc(store, env__emscripten_dlopen_js)
	result["env._mmap_js"] = wasmtime.WrapFunc(store, env__mmap_js)
	result["env.__syscall_getcwd"] = wasmtime.WrapFunc(store, env___syscall_getcwd)
	result["env.invoke_diii"] = wasmtime.WrapFunc(store, env_invoke_diii)
	result["env.js2python_init"] = wasmtime.WrapFunc(store, env_js2python_init)
	result["env.hiwire_assign_from_ptr"] = wasmtime.WrapFunc(store, env_hiwire_assign_from_ptr)
	result["env.__syscall_fcntl64"] = wasmtime.WrapFunc(store, env___syscall_fcntl64)
	result["env.__syscall_fadvise64"] = wasmtime.WrapFunc(store, env___syscall_fadvise64)
	result["env.hiwire_is_promise"] = wasmtime.WrapFunc(store, env_hiwire_is_promise)
	result["env.wrap_async_generator"] = wasmtime.WrapFunc(store, env_wrap_async_generator)
	result["env.JsSet_Add"] = wasmtime.WrapFunc(store, env_JsSet_Add)
	result["env.abort"] = wasmtime.WrapFunc(store, env_abort)
	result["env.__syscall_listen"] = wasmtime.WrapFunc(store, env___syscall_listen)
	result["env.JsArray_Check"] = wasmtime.WrapFunc(store, env_JsArray_Check)
	result["env.__syscall_fchdir"] = wasmtime.WrapFunc(store, env___syscall_fchdir)
	result["env.invoke_v"] = wasmtime.WrapFunc(store, env_invoke_v)
	result["env.invoke_viijii"] = wasmtime.WrapFunc(store, env_invoke_viijii)
	result["env.__syscall_getsockopt"] = wasmtime.WrapFunc(store, env___syscall_getsockopt)
	result["env.JsArray_Push_unchecked"] = wasmtime.WrapFunc(store, env_JsArray_Push_unchecked)
	result["env.hiwire_resolve_promise"] = wasmtime.WrapFunc(store, env_hiwire_resolve_promise)
	result["env.hiwire_read_from_file"] = wasmtime.WrapFunc(store, env_hiwire_read_from_file)
	result["env.emscripten_promise_destroy"] = wasmtime.WrapFunc(store, env_emscripten_promise_destroy)
	result["env.__syscall_accept4"] = wasmtime.WrapFunc(store, env___syscall_accept4)
	result["env._python2js_buffer_inner"] = wasmtime.WrapFunc(store, env__python2js_buffer_inner)
	result["env.pyproxy_new"] = wasmtime.WrapFunc(store, env_pyproxy_new)
	result["env._python2js_ucs4"] = wasmtime.WrapFunc(store, env__python2js_ucs4)
	result["env._localtime_js"] = wasmtime.WrapFunc(store, env__localtime_js)
	result["env.getnameinfo"] = wasmtime.WrapFunc(store, env_getnameinfo)
	result["env.__syscall_dup"] = wasmtime.WrapFunc(store, env___syscall_dup)
	result["env.__syscall_stat64"] = wasmtime.WrapFunc(store, env___syscall_stat64)
	result["env.invoke_iiiiij"] = wasmtime.WrapFunc(store, env_invoke_iiiiij)
	result["env._Py_CheckEmscriptenSignals_Helper"] = wasmtime.WrapFunc(store, env__Py_CheckEmscriptenSignals_Helper)
	result["env.ffi_prep_closure_loc_js"] = wasmtime.WrapFunc(store, env_ffi_prep_closure_loc_js)
	result["env._gmtime_js"] = wasmtime.WrapFunc(store, env__gmtime_js)
	result["env.invoke_viiiii"] = wasmtime.WrapFunc(store, env_invoke_viiiii)
	result["env.hiwire_is_comlink_proxy"] = wasmtime.WrapFunc(store, env_hiwire_is_comlink_proxy)
	result["env.JsObject_GetString"] = wasmtime.WrapFunc(store, env_JsObject_GetString)
	result["env.__syscall_fdatasync"] = wasmtime.WrapFunc(store, env___syscall_fdatasync)
	result["env.__cxa_throw"] = wasmtime.WrapFunc(store, env___cxa_throw)
	result["env.new_error"] = wasmtime.WrapFunc(store, env_new_error)
	result["env.descr_get_trampoline_call"] = wasmtime.WrapFunc(store, env_descr_get_trampoline_call)
	result["env.invoke_viii"] = wasmtime.WrapFunc(store, env_invoke_viii)
	result["env.invoke_iiiiiiii"] = wasmtime.WrapFunc(store, env_invoke_iiiiiiii)
	result["env.hiwire_constructor_name"] = wasmtime.WrapFunc(store, env_hiwire_constructor_name)
	result["env._emscripten_get_progname"] = wasmtime.WrapFunc(store, env__emscripten_get_progname)
	result["env.__cxa_end_catch"] = wasmtime.WrapFunc(store, env___cxa_end_catch)
	result["env.invoke_iiiiijj"] = wasmtime.WrapFunc(store, env_invoke_iiiiijj)
	result["env.__syscall_getpeername"] = wasmtime.WrapFunc(store, env___syscall_getpeername)
	result["env.JsSet_New"] = wasmtime.WrapFunc(store, env_JsSet_New)
	result["env.__syscall_fstat64"] = wasmtime.WrapFunc(store, env___syscall_fstat64)
	result["env.emscripten_get_now_res"] = wasmtime.WrapFunc(store, env_emscripten_get_now_res)
	result["env.__syscall_fstatfs64"] = wasmtime.WrapFunc(store, env___syscall_fstatfs64)
	result["env.invoke_viif"] = wasmtime.WrapFunc(store, env_invoke_viif)
	result["env.console_error_obj"] = wasmtime.WrapFunc(store, env_console_error_obj)
	result["env._JsArray_PushEntry_helper"] = wasmtime.WrapFunc(store, env__JsArray_PushEntry_helper)
	result["env.__syscall__newselect"] = wasmtime.WrapFunc(store, env___syscall__newselect)
	result["env.invoke_iiiiiiiiiii"] = wasmtime.WrapFunc(store, env_invoke_iiiiiiiiiii)
	result["env.invoke_viid"] = wasmtime.WrapFunc(store, env_invoke_viid)
	result["env.js2python_immutable"] = wasmtime.WrapFunc(store, env_js2python_immutable)
	result["env.JsArray_Splice"] = wasmtime.WrapFunc(store, env_JsArray_Splice)
	result["env.emscripten_promise_resolve"] = wasmtime.WrapFunc(store, env_emscripten_promise_resolve)
	result["env.invoke_iiiii"] = wasmtime.WrapFunc(store, env_invoke_iiiii)
	result["env.get_async_js_call_done_callback"] = wasmtime.WrapFunc(store, env_get_async_js_call_done_callback)
	result["env.ffi_call_js"] = wasmtime.WrapFunc(store, env_ffi_call_js)
	result["wasi_snapshot_preview1.proc_exit"] = wasmtime.WrapFunc(store, wasi_snapshot_preview1_proc_exit)
	result["env._tzset_js"] = wasmtime.WrapFunc(store, env__tzset_js)
	result["env.JsObjMap_subscript_js"] = wasmtime.WrapFunc(store, env_JsObjMap_subscript_js)
	result["env.hiwire_to_bool"] = wasmtime.WrapFunc(store, env_hiwire_to_bool)
	result["env.JsArray_reverse_helper"] = wasmtime.WrapFunc(store, env_JsArray_reverse_helper)
	result["env.system"] = wasmtime.WrapFunc(store, env_system)
	result["env.JsProxy_subscript_js"] = wasmtime.WrapFunc(store, env_JsProxy_subscript_js)
	result["env.pyproxy_AsPyObject"] = wasmtime.WrapFunc(store, env_pyproxy_AsPyObject)
	result["env.JsArray_count_helper"] = wasmtime.WrapFunc(store, env_JsArray_count_helper)
	result["env._Py_emscripten_runtime"] = wasmtime.WrapFunc(store, env__Py_emscripten_runtime)
	result["env.invoke_viiiiiii"] = wasmtime.WrapFunc(store, env_invoke_viiiiiii)
	result["env.python2js_buffer_init"] = wasmtime.WrapFunc(store, env_python2js_buffer_init)
	result["env.strftime"] = wasmtime.WrapFunc(store, env_strftime)
	result["env.ffi_closure_free_js"] = wasmtime.WrapFunc(store, env_ffi_closure_free_js)
	result["env.__syscall_ftruncate64"] = wasmtime.WrapFunc(store, env___syscall_ftruncate64)
	result["env.hiwire_CallMethodString"] = wasmtime.WrapFunc(store, env_hiwire_CallMethodString)
	result["env._PyImport_InitFunc_TrampolineCall"] = wasmtime.WrapFunc(store, env__PyImport_InitFunc_TrampolineCall)
	result["env.getaddrinfo"] = wasmtime.WrapFunc(store, env_getaddrinfo)
	result["env.emscripten_get_heap_max"] = wasmtime.WrapFunc(store, env_emscripten_get_heap_max)
	result["env.JsObject_Keys"] = wasmtime.WrapFunc(store, env_JsObject_Keys)
	result["env.proxy_cache_get"] = wasmtime.WrapFunc(store, env_proxy_cache_get)
	result["env.__syscall_unlinkat"] = wasmtime.WrapFunc(store, env___syscall_unlinkat)
	result["env.__syscall_statfs64"] = wasmtime.WrapFunc(store, env___syscall_statfs64)
	result["env.destroy_proxy"] = wasmtime.WrapFunc(store, env_destroy_proxy)
	result["env.create_promise_handles"] = wasmtime.WrapFunc(store, env_create_promise_handles)
	result["env.JsMap_clear_js"] = wasmtime.WrapFunc(store, env_JsMap_clear_js)
	result["env.invoke_iii"] = wasmtime.WrapFunc(store, env_invoke_iii)
	result["env.invoke_viijj"] = wasmtime.WrapFunc(store, env_invoke_viijj)
	result["env.__syscall_faccessat"] = wasmtime.WrapFunc(store, env___syscall_faccessat)
	result["env._dlsym_js"] = wasmtime.WrapFunc(store, env__dlsym_js)
	result["env.__cxa_increment_exception_refcount"] = wasmtime.WrapFunc(store, env___cxa_increment_exception_refcount)
	result["env.invoke_iiiiid"] = wasmtime.WrapFunc(store, env_invoke_iiiiid)
	result["env.hiwire_CallMethod_OneArg"] = wasmtime.WrapFunc(store, env_hiwire_CallMethod_OneArg)
	result["env.handle_next_result_js"] = wasmtime.WrapFunc(store, env_handle_next_result_js)
	result["env.hiwire_is_generator"] = wasmtime.WrapFunc(store, env_hiwire_is_generator)
	result["env.__syscall_symlinkat"] = wasmtime.WrapFunc(store, env___syscall_symlinkat)
	result["env.__syscall_sendto"] = wasmtime.WrapFunc(store, env___syscall_sendto)
	result["env.JsArray_Push"] = wasmtime.WrapFunc(store, env_JsArray_Push)
	result["env.hiwire_construct"] = wasmtime.WrapFunc(store, env_hiwire_construct)
	result["env.JsArray_index_helper"] = wasmtime.WrapFunc(store, env_JsArray_index_helper)
	result["env.__syscall_fchmodat"] = wasmtime.WrapFunc(store, env___syscall_fchmodat)
	result["env.JsObjMap_length_js"] = wasmtime.WrapFunc(store, env_JsObjMap_length_js)
	result["env.hiwire_not_equal"] = wasmtime.WrapFunc(store, env_hiwire_not_equal)
	result["env.emscripten_date_now"] = wasmtime.WrapFunc(store, env_emscripten_date_now)
	result["env.invoke_viiifi"] = wasmtime.WrapFunc(store, env_invoke_viiifi)
	result["env.JsObjMap_GetIter_js"] = wasmtime.WrapFunc(store, env_JsObjMap_GetIter_js)
	result["env._JsArray_PostProcess_helper"] = wasmtime.WrapFunc(store, env__JsArray_PostProcess_helper)
	result["env.gethostbyname"] = wasmtime.WrapFunc(store, env_gethostbyname)
	result["wasi_snapshot_preview1.fd_read"] = wasmtime.WrapFunc(store, wasi_snapshot_preview1_fd_read)
	result["env._agen_handle_result_js"] = wasmtime.WrapFunc(store, env__agen_handle_result_js)
	result["env.JsObject_Values"] = wasmtime.WrapFunc(store, env_JsObject_Values)
	result["env._python2js_cache_lookup"] = wasmtime.WrapFunc(store, env__python2js_cache_lookup)
	result["env.__syscall_truncate64"] = wasmtime.WrapFunc(store, env___syscall_truncate64)
	result["env.JsString_InternFromCString"] = wasmtime.WrapFunc(store, env_JsString_InternFromCString)
	result["env.proxy_cache_set"] = wasmtime.WrapFunc(store, env_proxy_cache_set)
	result["env.emscripten_promise_create"] = wasmtime.WrapFunc(store, env_emscripten_promise_create)
	result["env.invoke_ii"] = wasmtime.WrapFunc(store, env_invoke_ii)
	result["env.__syscall_chmod"] = wasmtime.WrapFunc(store, env___syscall_chmod)
	result["env._munmap_js"] = wasmtime.WrapFunc(store, env__munmap_js)
	result["env.invoke_i"] = wasmtime.WrapFunc(store, env_invoke_i)
	result["env.__cxa_begin_catch"] = wasmtime.WrapFunc(store, env___cxa_begin_catch)
	result["env.__syscall_sendmsg"] = wasmtime.WrapFunc(store, env___syscall_sendmsg)
	result["env.create_once_callable"] = wasmtime.WrapFunc(store, env_create_once_callable)
	result["env.hiwire_into_file"] = wasmtime.WrapFunc(store, env_hiwire_into_file)
	result["env._PyCFunctionWithKeywords_TrampolineCall"] = wasmtime.WrapFunc(store, env__PyCFunctionWithKeywords_TrampolineCall)
	result["env.invoke_j"] = wasmtime.WrapFunc(store, env_invoke_j)
	result["env.hiwire_throw_error"] = wasmtime.WrapFunc(store, env_hiwire_throw_error)
	result["env._msync_js"] = wasmtime.WrapFunc(store, env__msync_js)
	result["env.__call_sighandler"] = wasmtime.WrapFunc(store, env___call_sighandler)
	result["env.__resumeException"] = wasmtime.WrapFunc(store, env___resumeException)
	result["env.invoke_viiiiiiiiii"] = wasmtime.WrapFunc(store, env_invoke_viiiiiiiiii)
	result["env.JsObjMap_contains_js"] = wasmtime.WrapFunc(store, env_JsObjMap_contains_js)
	result["env._python2js_add_to_cache"] = wasmtime.WrapFunc(store, env__python2js_add_to_cache)
	result["env.getprotobyname"] = wasmtime.WrapFunc(store, env_getprotobyname)
	result["env.__syscall_symlink"] = wasmtime.WrapFunc(store, env___syscall_symlink)
	result["env.hiwire_init"] = wasmtime.WrapFunc(store, env_hiwire_init)
	result["env.array_to_js"] = wasmtime.WrapFunc(store, env_array_to_js)
	result["env.hiwire_double"] = wasmtime.WrapFunc(store, env_hiwire_double)
	result["env.invoke_vii"] = wasmtime.WrapFunc(store, env_invoke_vii)
	result["env.invoke_iiiiiiiiiiiii"] = wasmtime.WrapFunc(store, env_invoke_iiiiiiiiiiiii)
	result["env.__syscall_ioctl"] = wasmtime.WrapFunc(store, env___syscall_ioctl)
	result["env.emscripten_get_now"] = wasmtime.WrapFunc(store, env_emscripten_get_now)
	result["wasi_snapshot_preview1.fd_fdstat_get"] = wasmtime.WrapFunc(store, wasi_snapshot_preview1_fd_fdstat_get)
	result["env.__syscall_pipe"] = wasmtime.WrapFunc(store, env___syscall_pipe)
	result["env.hiwire_decref"] = wasmtime.WrapFunc(store, env_hiwire_decref)
	result["env.hiwire_less_than_equal"] = wasmtime.WrapFunc(store, env_hiwire_less_than_equal)
	result["env.hiwire_reversed_iterator"] = wasmtime.WrapFunc(store, env_hiwire_reversed_iterator)
	result["env.JsMap_Set"] = wasmtime.WrapFunc(store, env_JsMap_Set)
	result["env.invoke_fiii"] = wasmtime.WrapFunc(store, env_invoke_fiii)
	result["env.JsObject_Dir"] = wasmtime.WrapFunc(store, env_JsObject_Dir)
	result["env.hiwire_is_function"] = wasmtime.WrapFunc(store, env_hiwire_is_function)
	result["env.python2js_custom__create_jscontext"] = wasmtime.WrapFunc(store, env_python2js_custom__create_jscontext)
	result["env.getentropy"] = wasmtime.WrapFunc(store, env_getentropy)
	result["env.fail_test"] = wasmtime.WrapFunc(store, env_fail_test)
	result["env.JsObject_SetString"] = wasmtime.WrapFunc(store, env_JsObject_SetString)
	result["env.python2js__default_converter"] = wasmtime.WrapFunc(store, env_python2js__default_converter)
	result["env.__cxa_find_matching_catch_3"] = wasmtime.WrapFunc(store, env___cxa_find_matching_catch_3)
	result["env.invoke_viiiiii"] = wasmtime.WrapFunc(store, env_invoke_viiiiii)
	result["env.hiwire_HasMethod"] = wasmtime.WrapFunc(store, env_hiwire_HasMethod)
	result["env.JsProxy_GetAsyncIter_js"] = wasmtime.WrapFunc(store, env_JsProxy_GetAsyncIter_js)
	result["env.__syscall_rmdir"] = wasmtime.WrapFunc(store, env___syscall_rmdir)
	result["env.invoke_ji"] = wasmtime.WrapFunc(store, env_invoke_ji)
	result["env.invoke_jiiii"] = wasmtime.WrapFunc(store, env_invoke_jiiii)
	result["env.hiwire_incref"] = wasmtime.WrapFunc(store, env_hiwire_incref)
	result["env._emscripten_get_now_is_monotonic"] = wasmtime.WrapFunc(store, env__emscripten_get_now_is_monotonic)
	result["env.__syscall_utimensat"] = wasmtime.WrapFunc(store, env___syscall_utimensat)
	result["env.__cxa_uncaught_exceptions"] = wasmtime.WrapFunc(store, env___cxa_uncaught_exceptions)
	result["env.emscripten_asm_const_int"] = wasmtime.WrapFunc(store, env_emscripten_asm_const_int)
	result["env.hiwire_is_async_generator"] = wasmtime.WrapFunc(store, env_hiwire_is_async_generator)
	result["env.JsBuffer_DecodeString_js"] = wasmtime.WrapFunc(store, env_JsBuffer_DecodeString_js)
	result["env.invoke_viiii"] = wasmtime.WrapFunc(store, env_invoke_viiii)
	result["env._python2js_ucs1"] = wasmtime.WrapFunc(store, env__python2js_ucs1)
	result["env.JsMap_New"] = wasmtime.WrapFunc(store, env_JsMap_New)
	result["env.__syscall_fchown32"] = wasmtime.WrapFunc(store, env___syscall_fchown32)
	result["env.__cxa_rethrow_primary_exception"] = wasmtime.WrapFunc(store, env___cxa_rethrow_primary_exception)
	result["env.JsMap_GetIter_js"] = wasmtime.WrapFunc(store, env_JsMap_GetIter_js)
	result["env._python2js_addto_postprocess_list"] = wasmtime.WrapFunc(store, env__python2js_addto_postprocess_list)
	result["env.__syscall_chdir"] = wasmtime.WrapFunc(store, env___syscall_chdir)
	result["env.__syscall_fallocate"] = wasmtime.WrapFunc(store, env___syscall_fallocate)
	result["env.invoke_viiiiiiiiiiiiiii"] = wasmtime.WrapFunc(store, env_invoke_viiiiiiiiiiiiiii)
	result["env.hiwire_less_than"] = wasmtime.WrapFunc(store, env_hiwire_less_than)
	result["env.hiwire_equal"] = wasmtime.WrapFunc(store, env_hiwire_equal)
	result["wasi_snapshot_preview1.fd_close"] = wasmtime.WrapFunc(store, wasi_snapshot_preview1_fd_close)
	result["env.__syscall_mknodat"] = wasmtime.WrapFunc(store, env___syscall_mknodat)
	result["env.hiwire_get_buffer_info"] = wasmtime.WrapFunc(store, env_hiwire_get_buffer_info)
	result["env.JsArray_slice"] = wasmtime.WrapFunc(store, env_JsArray_slice)
	result["env.JsArray_Delete"] = wasmtime.WrapFunc(store, env_JsArray_Delete)
	result["env._emscripten_throw_longjmp"] = wasmtime.WrapFunc(store, env__emscripten_throw_longjmp)
	result["env.JsObjMap_ass_subscript_js"] = wasmtime.WrapFunc(store, env_JsObjMap_ass_subscript_js)
	result["env._python2js_ucs2"] = wasmtime.WrapFunc(store, env__python2js_ucs2)
	result["env._timegm_js"] = wasmtime.WrapFunc(store, env__timegm_js)
	result["env._setitimer_js"] = wasmtime.WrapFunc(store, env__setitimer_js)
	result["env.pyproxy_Check"] = wasmtime.WrapFunc(store, env_pyproxy_Check)
	result["env.hiwire_to_string"] = wasmtime.WrapFunc(store, env_hiwire_to_string)
	result["env.JsArray_slice_assign"] = wasmtime.WrapFunc(store, env_JsArray_slice_assign)
	result["env.wrap_generator"] = wasmtime.WrapFunc(store, env_wrap_generator)
	result["env.hiwire_int_from_digits"] = wasmtime.WrapFunc(store, env_hiwire_int_from_digits)
	result["env.getloadavg"] = wasmtime.WrapFunc(store, env_getloadavg)
	result["wasi_snapshot_preview1.fd_write"] = wasmtime.WrapFunc(store, wasi_snapshot_preview1_fd_write)
	result["env.invoke_iiii"] = wasmtime.WrapFunc(store, env_invoke_iiii)
	result["env.__syscall_socket"] = wasmtime.WrapFunc(store, env___syscall_socket)
	result["env.hiwire_typeof"] = wasmtime.WrapFunc(store, env_hiwire_typeof)
	result["env.__syscall_dup3"] = wasmtime.WrapFunc(store, env___syscall_dup3)
	result["env.__syscall_lstat64"] = wasmtime.WrapFunc(store, env___syscall_lstat64)
	result["env.__syscall_readlinkat"] = wasmtime.WrapFunc(store, env___syscall_readlinkat)
	result["env.__syscall_bind"] = wasmtime.WrapFunc(store, env___syscall_bind)
	result["env.hiwire_call_OneArg"] = wasmtime.WrapFunc(store, env_hiwire_call_OneArg)
	result["env.__syscall_fchownat"] = wasmtime.WrapFunc(store, env___syscall_fchownat)
	result["wasi_snapshot_preview1.fd_sync"] = wasmtime.WrapFunc(store, wasi_snapshot_preview1_fd_sync)
	result["env.emscripten_resize_heap"] = wasmtime.WrapFunc(store, env_emscripten_resize_heap)
	result["env._python2js_handle_postprocess_list"] = wasmtime.WrapFunc(store, env__python2js_handle_postprocess_list)
	result["env.__syscall_fchmod"] = wasmtime.WrapFunc(store, env___syscall_fchmod)
	result["wasi_snapshot_preview1.fd_pread"] = wasmtime.WrapFunc(store, wasi_snapshot_preview1_fd_pread)
	result["env.__cxa_rethrow"] = wasmtime.WrapFunc(store, env___cxa_rethrow)
	result["env.invoke_iiiiiiiiiiii"] = wasmtime.WrapFunc(store, env_invoke_iiiiiiiiiiii)
	result["env.hiwire_get_length_helper"] = wasmtime.WrapFunc(store, env_hiwire_get_length_helper)
	result["env.emscripten_exit_with_live_runtime"] = wasmtime.WrapFunc(store, env_emscripten_exit_with_live_runtime)
	result["wasi_snapshot_preview1.fd_seek"] = wasmtime.WrapFunc(store, wasi_snapshot_preview1_fd_seek)
	result["env.__cxa_current_primary_exception"] = wasmtime.WrapFunc(store, env___cxa_current_primary_exception)
	result["env.strftime_l"] = wasmtime.WrapFunc(store, env_strftime_l)
	result["env.JsArray_New"] = wasmtime.WrapFunc(store, env_JsArray_New)
	result["env.__assert_fail"] = wasmtime.WrapFunc(store, env___assert_fail)
	result["env.__syscall_openat"] = wasmtime.WrapFunc(store, env___syscall_openat)
	result["env.invoke_iiiiiii"] = wasmtime.WrapFunc(store, env_invoke_iiiiiii)
	result["env.hiwire_call_bound"] = wasmtime.WrapFunc(store, env_hiwire_call_bound)
	result["env.destroy_proxies_js"] = wasmtime.WrapFunc(store, env_destroy_proxies_js)
	result["env.emscripten_memcpy_big"] = wasmtime.WrapFunc(store, env_emscripten_memcpy_big)
	result["env.__cxa_find_matching_catch_2"] = wasmtime.WrapFunc(store, env___cxa_find_matching_catch_2)
	result["env.hiwire_CallMethod"] = wasmtime.WrapFunc(store, env_hiwire_CallMethod)
	result["env.hiwire_CallMethod_NoArgs"] = wasmtime.WrapFunc(store, env_hiwire_CallMethod_NoArgs)
	result["env.hiwire_get_bool"] = wasmtime.WrapFunc(store, env_hiwire_get_bool)
	result["env.__syscall_newfstatat"] = wasmtime.WrapFunc(store, env___syscall_newfstatat)
	result["env.__syscall_recvfrom"] = wasmtime.WrapFunc(store, env___syscall_recvfrom)
	result["env.hiwire_get_length_string"] = wasmtime.WrapFunc(store, env_hiwire_get_length_string)
	result["env.JsException_new_helper"] = wasmtime.WrapFunc(store, env_JsException_new_helper)
	result["env.ffi_closure_alloc_js"] = wasmtime.WrapFunc(store, env_ffi_closure_alloc_js)
	result["wasi_snapshot_preview1.environ_get"] = wasmtime.WrapFunc(store, wasi_snapshot_preview1_environ_get)
	result["env.JsArray_Get"] = wasmtime.WrapFunc(store, env_JsArray_Get)
	result["env.js2python_convert"] = wasmtime.WrapFunc(store, env_js2python_convert)
	result["env._python2js_destroy_cache"] = wasmtime.WrapFunc(store, env__python2js_destroy_cache)
	result["wasi_snapshot_preview1.fd_pwrite"] = wasmtime.WrapFunc(store, wasi_snapshot_preview1_fd_pwrite)
	result["env.invoke_jiii"] = wasmtime.WrapFunc(store, env_invoke_jiii)
	result["env.hiwire_greater_than"] = wasmtime.WrapFunc(store, env_hiwire_greater_than)
	result["env.hiwire_greater_than_equal"] = wasmtime.WrapFunc(store, env_hiwire_greater_than_equal)
	result["env.descr_set_trampoline_call"] = wasmtime.WrapFunc(store, env_descr_set_trampoline_call)
	result["env.__syscall_poll"] = wasmtime.WrapFunc(store, env___syscall_poll)
	result["env.js2python"] = wasmtime.WrapFunc(store, env_js2python)
	result["env.JsArray_Set"] = wasmtime.WrapFunc(store, env_JsArray_Set)
	result["env.__syscall_renameat"] = wasmtime.WrapFunc(store, env___syscall_renameat)
	result["env.invoke_vi"] = wasmtime.WrapFunc(store, env_invoke_vi)

}

func addEmscriptenGlobals(store wasmtime.Storelike, result map[string]*wasmtime.Global) {
	var valType *wasmtime.ValType
	var gType *wasmtime.GlobalType
	var g *wasmtime.Global
	var err error

	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glUniform1fv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glStencilFuncSeparate"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glRenderbufferStorage"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetUniformLocation"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetFramebufferAttachmentParameteriv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glDeleteBuffers"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glCopyTexSubImage2D"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.__cxa_rethrow"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glViewport"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glTexParameterf"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetShaderiv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetFloatv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glDepthRangef"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, false)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["env.__table_base"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glUniformMatrix4fv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glUniform2iv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glStencilMask"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glReleaseShaderCompiler"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetVertexAttribPointerv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGenFramebuffers"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glBlendEquationSeparate"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glDeleteQueriesEXT"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func._emscripten_out"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetShaderInfoLog"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glBlendEquation"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glIsBuffer"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glDeleteProgram"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetQueryObjectui64vEXT"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetQueryivEXT"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glTexImage2D"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetShaderSource"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glQueryCounterEXT"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glStencilMaskSeparate"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glLinkProgram"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetProgramiv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glCompressedTexImage2D"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glClearDepthf"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["env.__stack_pointer"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glIsQueryEXT"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_console_error"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glIsProgram"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetProgramInfoLog"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGenRenderbuffers"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glClearColor"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.__cxa_end_catch"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glUniform3i"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetUniformiv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetRenderbufferParameteriv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glBufferData"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glBeginQueryEXT"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetError"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGenTextures"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glEnable"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glCopyTexImage2D"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glAttachShader"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glVertexAttrib4f"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glUseProgram"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glUniform3iv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetTexParameterfv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glDrawArrays"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glDisable"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glDeleteShader"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glClearStencil"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.abort"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_console_warn"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glVertexAttrib3f"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glStencilOp"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetVertexAttribfv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glCheckFramebufferStatus"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glDrawElementsInstancedANGLE"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGenQueriesEXT"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glVertexAttrib2fv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glShaderBinary"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetBufferParameteriv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func._emscripten_err"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glValidateProgram"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetBooleanv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetAttribLocation"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGenBuffers"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glDeleteFramebuffers"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGenVertexArraysOES"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glVertexAttrib1f"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glUniform3fv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glUniform2i"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetIntegerv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.mem.__stack_high"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glVertexAttribPointer"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glVertexAttrib1fv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glIsRenderbuffer"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGenerateMipmap"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glFrontFace"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.mem.__stack_low"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glUniform4i"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glUniform3f"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glBlendFuncSeparate"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glBindTexture"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glEndQueryEXT"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_console_log"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glUniform2f"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glTexParameteriv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glTexParameteri"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glDrawElements"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glBindRenderbuffer"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glDrawArraysInstancedANGLE"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glBindVertexArrayOES"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glVertexAttrib3fv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glUniformMatrix3fv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glPixelStorei"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glFramebufferRenderbuffer"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glActiveTexture"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, false)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["env.__memory_base"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetActiveAttrib"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glFramebufferTexture2D"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glClear"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glUniform1iv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glLineWidth"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glFinish"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glDetachShader"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glDepthMask"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glCompressedTexSubImage2D"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.mem.__heap_base"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glVertexAttribDivisorANGLE"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glDeleteVertexArraysOES"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glStencilFunc"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glIsFramebuffer"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glIsEnabled"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glUniformMatrix2fv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glTexParameterfv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glShaderSource"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glReadPixels"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glCreateShader"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glUniform4f"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glUniform1f"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetString"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glFlush"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glCreateProgram"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glBlendFunc"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glIsVertexArrayOES"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetQueryObjectuivEXT"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glUniform1i"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glPolygonOffset"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glBlendColor"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetQueryObjecti64vEXT"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glTexSubImage2D"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glStencilOpSeparate"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glDeleteTextures"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glUniform2fv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetVertexAttribiv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetTexParameteriv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glCullFace"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glBindBuffer"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glSampleCoverage"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetActiveUniform"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glEnableVertexAttribArray"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glDisableVertexAttribArray"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glDrawBuffersWEBGL"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glUniform4iv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glIsTexture"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetShaderPrecisionFormat"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glColorMask"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.__cxa_throw"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glScissor"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glIsShader"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glDepthFunc"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glBufferSubData"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetQueryObjectivEXT"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glVertexAttrib2f"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glHint"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetAttachedShaders"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glDeleteRenderbuffers"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glBindFramebuffer"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glVertexAttrib4fv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glUniform4fv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glGetUniformfv"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glCompileShader"] = g
	valType = wasmtime.NewValType(wasmtime.KindI32)
	gType = wasmtime.NewGlobalType(valType, true)
	g, err = wasmtime.NewGlobal(store, gType, wasmtime.ValI32(0))
	if err != nil {
		panic(err.Error())
	}
	result["GOT.func.emscripten_glBindAttribLocation"] = g

}

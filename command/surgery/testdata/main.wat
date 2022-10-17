(module
  (type (;0;) (func (param i32 i32 i32 i32) (result i32)))
  (type (;1;) (func (param i32 i32 i32) (result i32)))
  (type (;2;) (func))
  (type (;3;) (func (param i32 i32)))
  (type (;4;) (func (param i32)))
  (type (;5;) (func (param i32) (result i32)))
  (type (;6;) (func (param i32 i32) (result i32)))
  (import "wasi_snapshot_preview1" "fd_write" (func $runtime.fd_write (type 0)))
  (import "env" "runtime.alloc" (func $runtime.alloc (type 1)))
  (import "env" "runtime.realloc" (func $runtime.realloc (type 1)))
  (func $runtime.lookupPanic (type 2)
    call $runtime.runtimePanic
    unreachable)
  (func $runtime.runtimePanic (type 2)
    i32.const 65536
    i32.const 22
    call $runtime.printstring
    i32.const 65558
    i32.const 18
    call $runtime.printstring
    call $runtime.printnl
    unreachable
    unreachable)
  (func $runtime.printstring (type 3) (param i32 i32)
    local.get 1
    i32.const 0
    local.get 1
    i32.const 0
    i32.gt_s
    select
    local.set 1
    block  ;; label = @1
      loop  ;; label = @2
        local.get 1
        i32.eqz
        br_if 1 (;@1;)
        local.get 0
        i32.load8_u
        call $runtime.putchar
        local.get 0
        i32.const 1
        i32.add
        local.set 0
        local.get 1
        i32.const -1
        i32.add
        local.set 1
        br 0 (;@2;)
      end
    end)
  (func $runtime.printnl (type 2)
    i32.const 10
    call $runtime.putchar)
  (func $runtime.putchar (type 4) (param i32)
    (local i32 i32)
    block  ;; label = @1
      i32.const 0
      i32.load offset=65592
      local.tee 1
      i32.const 119
      i32.gt_u
      br_if 0 (;@1;)
      i32.const 0
      local.get 1
      i32.const 1
      i32.add
      local.tee 2
      i32.store offset=65592
      local.get 1
      i32.const 65596
      i32.add
      local.get 0
      i32.store8
      block  ;; label = @2
        block  ;; label = @3
          local.get 0
          i32.const 255
          i32.and
          i32.const 10
          i32.eq
          br_if 0 (;@3;)
          local.get 1
          i32.const 119
          i32.ne
          br_if 1 (;@2;)
        end
        i32.const 0
        local.get 2
        i32.store offset=65588
        i32.const 1
        i32.const 65584
        i32.const 1
        i32.const 65716
        call $runtime.fd_write
        drop
        i32.const 0
        i32.const 0
        i32.store offset=65592
      end
      return
    end
    call $runtime.lookupPanic
    unreachable)
  (func $malloc (type 5) (param i32) (result i32)
    local.get 0
    i32.const 0
    local.get 0
    call $runtime.alloc)
  (func $free (type 4) (param i32))
  (func $calloc (type 6) (param i32 i32) (result i32)
    local.get 1
    local.get 0
    i32.mul
    i32.const 0
    local.get 1
    call $runtime.alloc)
  (func $realloc (type 6) (param i32 i32) (result i32)
    local.get 0
    local.get 1
    local.get 0
    call $runtime.realloc)
  (func $_start (type 2)
    (local i32 i32 i32)
    i32.const 6
    i32.const 0
    local.get 0
    call $runtime.alloc
    local.tee 0
    i32.const 4
    i32.add
    i32.const 0
    i32.load16_u offset=65580 align=1
    i32.store16 align=1
    local.get 0
    i32.const 0
    i32.load offset=65576 align=1
    i32.store align=1
    i32.const 48
    call $runtime.putchar
    i32.const 120
    call $runtime.putchar
    i32.const 8
    local.set 1
    block  ;; label = @1
      loop  ;; label = @2
        local.get 1
        i32.eqz
        br_if 1 (;@1;)
        local.get 0
        i32.const 28
        i32.shr_u
        local.tee 2
        i32.const 48
        i32.or
        local.get 2
        i32.const 87
        i32.add
        local.get 2
        i32.const 10
        i32.lt_u
        select
        call $runtime.putchar
        local.get 1
        i32.const -1
        i32.add
        local.set 1
        local.get 0
        i32.const 4
        i32.shl
        local.set 0
        br 0 (;@2;)
      end
    end)
  (func $zap (type 1) (param i32 i32 i32) (result i32)
    block  ;; label = @1
      local.get 1
      i32.eqz
      br_if 0 (;@1;)
      local.get 0
      return
    end
    call $runtime.lookupPanic
    unreachable)
  (func $dummy (type 2))
  (func $__wasm_call_dtors (type 2)
    call $dummy
    call $dummy)
  (func $malloc.command_export (type 5) (param i32) (result i32)
    local.get 0
    call $malloc
    call $__wasm_call_dtors)
  (func $free.command_export (type 4) (param i32)
    local.get 0
    call $free
    call $__wasm_call_dtors)
  (func $calloc.command_export (type 6) (param i32 i32) (result i32)
    local.get 0
    local.get 1
    call $calloc
    call $__wasm_call_dtors)
  (func $realloc.command_export (type 6) (param i32 i32) (result i32)
    local.get 0
    local.get 1
    call $realloc
    call $__wasm_call_dtors)
  (func $_start.command_export (type 2)
    call $_start
    call $__wasm_call_dtors)
  (func $zap.command_export (type 1) (param i32 i32 i32) (result i32)
    local.get 0
    local.get 1
    local.get 2
    call $zap
    call $__wasm_call_dtors)
  (table (;0;) 1 1 funcref)
  (memory (;0;) 2)
  (global $__stack_pointer (mut i32) (i32.const 65536))
  (export "memory" (memory 0))
  (export "malloc" (func $malloc.command_export))
  (export "free" (func $free.command_export))
  (export "calloc" (func $calloc.command_export))
  (export "realloc" (func $realloc.command_export))
  (export "_start" (func $_start.command_export))
  (export "zap" (func $zap.command_export))
  (data $.rodata (i32.const 65536) "panic: runtime error: index out of rangefoobie")
  (data $.data (i32.const 65584) "<\00\01\00\00\00\00\00"))

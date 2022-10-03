(module
  (type (;0;) (func))
  (type (;1;) (func (param i32)))
  (type (;2;) (func (param i32 i32)))
  (type (;3;) (func (param i32) (result i32)))
  (type (;4;) (func (param i32 i32 i32 i32) (result i32)))
  (type (;5;) (func (param i32 i64 i32)))
  (type (;6;) (func (param i32 i32) (result i32)))
  (type (;7;) (func (param i64 i32) (result i32)))
  (type (;8;) (func (result i32)))
  (type (;9;) (func (param i64 i32 i32 i32) (result i64)))
  (type (;10;) (func (param i64 i32 i32 i32 i32)))
  (type (;11;) (func (param i64 i32)))
  (type (;12;) (func (param i32 i32 i32) (result i64)))
  (type (;13;) (func (param i64 i32 i32 i64 i32)))
  (type (;14;) (func (param i64 i32 i32) (result i64)))
  (type (;15;) (func (param i32 i64 i32 i32 i32 i32 i32 i32)))
  (type (;16;) (func (param i32 i32 i32 i32 i32)))
  (type (;17;) (func (param i32 i64 i32 i32 i32)))
  (type (;18;) (func (param i32 i64 i32 i32 i32 i32 i32)))
  (type (;19;) (func (param i32 i64)))
  (type (;20;) (func (param i32 i32 i32)))
  (type (;21;) (func (param i32 f64)))
  (type (;22;) (func (param i64 i32 i32 i32 i32 i32)))
  (import "wasi_snapshot_preview1" "fd_write" (func $runtime.fd_write (type 4)))
  (import "parigot_abi" "outputString" (func $github.com/iansmith/parigot/abi.OutputString (type 2)))
  (import "env" "syscall/js.valueGet" (func $syscall/js.valueGet (type 9)))
  (import "env" "syscall/js.valuePrepareString" (func $syscall/js.valuePrepareString (type 5)))
  (import "env" "syscall/js.valueLoadString" (func $syscall/js.valueLoadString (type 10)))
  (import "env" "syscall/js.finalizeRef" (func $syscall/js.finalizeRef (type 11)))
  (import "env" "syscall/js.stringVal" (func $syscall/js.stringVal (type 12)))
  (import "env" "syscall/js.valueSet" (func $syscall/js.valueSet (type 13)))
  (import "env" "syscall/js.valueLength" (func $syscall/js.valueLength (type 7)))
  (import "env" "syscall/js.valueIndex" (func $syscall/js.valueIndex (type 14)))
  (import "env" "syscall/js.valueCall" (func $syscall/js.valueCall (type 15)))
  (func $_*internal/task.gcData_.swap (type 1) (param i32)
    (local i32)
    block (result i32)  ;; label = @1
      global.get 1
      i32.const 2
      i32.eq
      if  ;; label = @2
        global.get 2
        global.get 2
        i32.load
        i32.const 4
        i32.sub
        i32.store
        global.get 2
        i32.load
        i32.load
        local.set 1
      end
      global.get 1
      i32.const 1
      local.get 0
      select
      i32.eqz
      if  ;; label = @2
        local.get 0
        i32.load
        local.set 1
        local.get 0
        i32.const 66380
        i32.load
        i32.store
        i32.const 66380
        local.get 1
        i32.store
        return
      end
      local.get 1
      i32.const 0
      global.get 1
      select
      i32.eqz
      if  ;; label = @2
        call $runtime.nilPanic
        i32.const 0
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        unreachable
      end
      return
    end
    local.set 0
    global.get 2
    i32.load
    local.get 0
    i32.store
    global.get 2
    global.get 2
    i32.load
    i32.const 4
    i32.add
    i32.store)
  (func $runtime.nilPanic (type 0)
    (local i32)
    block (result i32)  ;; label = @1
      global.get 1
      i32.const 2
      i32.eq
      if (result i32)  ;; label = @2
        global.get 2
        global.get 2
        i32.load
        i32.const 4
        i32.sub
        i32.store
        global.get 2
        i32.load
        i32.load
      else
        local.get 0
      end
      i32.const 0
      global.get 1
      select
      i32.eqz
      if  ;; label = @2
        i32.const 65592
        i32.const 23
        call $runtime.runtimePanic
        i32.const 0
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        unreachable
      end
      return
    end
    local.set 0
    global.get 2
    i32.load
    local.get 0
    i32.store
    global.get 2
    global.get 2
    i32.load
    i32.const 4
    i32.add
    i32.store)
  (func $_*internal/task.Queue_.Pop (type 8) (result i32)
    (local i32)
    i32.const 66356
    i32.load
    local.tee 0
    if  ;; label = @1
      i32.const 66356
      local.get 0
      i32.load
      i32.store
      i32.const 66360
      i32.load
      local.get 0
      i32.eq
      if  ;; label = @2
        i32.const 66360
        i32.const 0
        i32.store
      end
      local.get 0
      i32.const 0
      i32.store
    end
    local.get 0)
  (func $_*internal/task.Task_.Resume (type 1) (param i32)
    (local i32 i32 i32 i32)
    global.get 1
    i32.const 2
    i32.eq
    if  ;; label = @1
      global.get 2
      global.get 2
      i32.load
      i32.const 12
      i32.sub
      i32.store
      global.get 2
      i32.load
      local.tee 3
      i32.load
      local.set 0
      local.get 3
      i32.load offset=4
      local.set 2
      local.get 3
      i32.load offset=8
      local.set 3
    end
    block (result i32)  ;; label = @1
      global.get 1
      i32.const 2
      i32.eq
      if  ;; label = @2
        global.get 2
        global.get 2
        i32.load
        i32.const 4
        i32.sub
        i32.store
        global.get 2
        i32.load
        i32.load
        local.set 1
      end
      local.get 2
      local.get 0
      i32.eqz
      global.get 1
      select
      local.set 2
      block  ;; label = @2
        block  ;; label = @3
          global.get 1
          i32.eqz
          if  ;; label = @4
            local.get 2
            br_if 1 (;@3;)
            i32.const 66224
            i32.load
            local.set 3
            local.get 0
            i32.const 16
            i32.add
            local.set 2
          end
          local.get 1
          i32.const 0
          global.get 1
          select
          i32.eqz
          if  ;; label = @4
            local.get 2
            call $_*internal/task.gcData_.swap
            i32.const 0
            global.get 1
            i32.const 1
            i32.eq
            br_if 3 (;@1;)
            drop
          end
          global.get 1
          i32.eqz
          if  ;; label = @4
            i32.const 66224
            local.get 0
            i32.store
            local.get 0
            i32.const 20
            i32.add
            local.set 4
            local.get 0
            i32.const 36
            i32.add
            i32.load8_u
            if  ;; label = @5
              local.get 4
              call $tinygo_rewind
              br 3 (;@2;)
            end
            local.get 4
            call $tinygo_launch
            local.get 0
            i32.const 1
            i32.store8 offset=36
            br 2 (;@2;)
          end
        end
        local.get 1
        i32.const 1
        i32.eq
        i32.const 1
        global.get 1
        select
        if  ;; label = @3
          call $runtime.nilPanic
          i32.const 1
          global.get 1
          i32.const 1
          i32.eq
          br_if 2 (;@1;)
          drop
        end
        global.get 1
        i32.eqz
        if  ;; label = @3
          unreachable
        end
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        i32.const 66224
        local.get 3
        i32.store
      end
      local.get 1
      i32.const 2
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        local.get 2
        call $_*internal/task.gcData_.swap
        i32.const 2
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        local.get 0
        i32.const 32
        i32.add
        i32.load
        local.get 0
        i32.const 28
        i32.add
        i32.load
        i32.ge_u
        if  ;; label = @3
          return
        end
      end
      local.get 1
      i32.const 3
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        i32.const 65536
        i32.const 14
        call $runtime.runtimePanic
        i32.const 3
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        unreachable
      end
      return
    end
    local.set 1
    global.get 2
    i32.load
    local.get 1
    i32.store
    global.get 2
    global.get 2
    i32.load
    i32.const 4
    i32.add
    i32.store
    global.get 2
    i32.load
    local.tee 1
    local.get 0
    i32.store
    local.get 1
    local.get 2
    i32.store offset=4
    local.get 1
    local.get 3
    i32.store offset=8
    global.get 2
    global.get 2
    i32.load
    i32.const 12
    i32.add
    i32.store)
  (func $runtime.runtimePanic (type 2) (param i32 i32)
    (local i32)
    global.get 1
    i32.const 2
    i32.eq
    if  ;; label = @1
      global.get 2
      global.get 2
      i32.load
      i32.const 8
      i32.sub
      i32.store
      global.get 2
      i32.load
      local.tee 1
      i32.load
      local.set 0
      local.get 1
      i32.load offset=4
      local.set 1
    end
    block (result i32)  ;; label = @1
      global.get 1
      i32.const 2
      i32.eq
      if  ;; label = @2
        global.get 2
        global.get 2
        i32.load
        i32.const 4
        i32.sub
        i32.store
        global.get 2
        i32.load
        i32.load
        local.set 2
      end
      local.get 2
      i32.const 0
      global.get 1
      select
      i32.eqz
      if  ;; label = @2
        i32.const 65570
        i32.const 22
        call $runtime.printstring
        i32.const 0
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
      end
      local.get 2
      i32.const 1
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        local.get 0
        local.get 1
        call $runtime.printstring
        i32.const 1
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
      end
      local.get 2
      i32.const 2
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        call $runtime.printnl
        i32.const 2
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        unreachable
      end
      return
    end
    local.set 2
    global.get 2
    i32.load
    local.get 2
    i32.store
    global.get 2
    global.get 2
    i32.load
    i32.const 4
    i32.add
    i32.store
    global.get 2
    i32.load
    local.tee 2
    local.get 0
    i32.store
    local.get 2
    local.get 1
    i32.store offset=4
    global.get 2
    global.get 2
    i32.load
    i32.const 8
    i32.add
    i32.store)
  (func $internal/task.start (type 1) (param i32)
    (local i32 i32 i32 i32 i32 i32 i32 i32)
    global.get 1
    i32.const 2
    i32.eq
    if  ;; label = @1
      global.get 2
      global.get 2
      i32.load
      i32.const 28
      i32.sub
      i32.store
      global.get 2
      i32.load
      local.tee 1
      i32.load
      local.set 0
      local.get 1
      i32.load offset=4
      local.set 3
      local.get 1
      i32.load offset=12
      local.set 4
      local.get 1
      i32.load offset=16
      local.set 5
      local.get 1
      i32.load offset=20
      local.set 7
      local.get 1
      i32.load offset=24
      local.set 8
      local.get 1
      i32.load offset=8
      local.set 2
    end
    block (result i32)  ;; label = @1
      global.get 1
      i32.const 2
      i32.eq
      if  ;; label = @2
        global.get 2
        global.get 2
        i32.load
        i32.const 4
        i32.sub
        i32.store
        global.get 2
        i32.load
        i32.load
        local.set 6
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        global.get $__stack_pointer
        i32.const 32
        i32.sub
        local.tee 3
        global.set $__stack_pointer
        local.get 3
        i32.const 28
        i32.add
        local.tee 7
        i32.const 0
        i32.store
        local.get 3
        i32.const 20
        i32.add
        local.tee 8
        i64.const 0
        i64.store align=4
        local.get 3
        i64.const 4
        i64.store offset=12 align=4
        i32.const 66380
        i32.load
        local.set 4
        i32.const 66380
        local.get 3
        i32.const 8
        i32.add
        local.tee 2
        i32.store
        local.get 3
        local.get 4
        i32.store offset=8
      end
      local.get 6
      i32.const 0
      global.get 1
      select
      i32.eqz
      if  ;; label = @2
        i32.const 48
        call $runtime.alloc
        local.set 1
        i32.const 0
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
        local.get 1
        local.set 2
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        local.get 2
        i32.const 0
        i32.store offset=24
        local.get 2
        local.get 0
        i32.store offset=20
        local.get 3
        i32.const 16
        i32.add
        local.tee 0
        local.get 2
        i32.store
        local.get 3
        i32.const 24
        i32.add
        local.set 5
      end
      local.get 6
      i32.const 1
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        i32.const 16384
        call $runtime.alloc
        local.set 1
        i32.const 1
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
        local.get 1
        local.set 0
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        local.get 5
        local.get 0
        i32.store
        local.get 7
        local.get 0
        i32.store
        local.get 8
        local.get 0
        i32.store
        local.get 2
        local.get 0
        i32.store offset=28
        local.get 0
        i32.const -1204030091
        i32.store
        local.get 2
        local.get 0
        i32.const 16384
        i32.add
        i32.store offset=32
      end
      local.get 6
      i32.const 2
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        local.get 2
        call $runtime.runqueuePushBack
        i32.const 2
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        i32.const 66380
        local.get 4
        i32.store
        local.get 3
        i32.const 32
        i32.add
        global.set $__stack_pointer
      end
      return
    end
    local.set 1
    global.get 2
    i32.load
    local.get 1
    i32.store
    global.get 2
    global.get 2
    i32.load
    i32.const 4
    i32.add
    i32.store
    global.get 2
    i32.load
    local.tee 1
    local.get 0
    i32.store
    local.get 1
    local.get 3
    i32.store offset=4
    local.get 1
    local.get 2
    i32.store offset=8
    local.get 1
    local.get 4
    i32.store offset=12
    local.get 1
    local.get 5
    i32.store offset=16
    local.get 1
    local.get 7
    i32.store offset=20
    local.get 1
    local.get 8
    i32.store offset=24
    global.get 2
    global.get 2
    i32.load
    i32.const 28
    i32.add
    i32.store)
  )

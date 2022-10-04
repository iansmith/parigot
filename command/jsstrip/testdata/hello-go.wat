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
  (func $runtime.alloc (type 3) (param i32) (result i32)
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
      local.tee 3
      i32.load
      local.set 0
      local.get 3
      i32.load offset=4
      local.set 1
      local.get 3
      i32.load offset=12
      local.set 4
      local.get 3
      i32.load offset=16
      local.set 5
      local.get 3
      i32.load offset=20
      local.set 6
      local.get 3
      i32.load offset=24
      local.set 7
      local.get 3
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
        local.set 8
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        local.get 0
        i32.eqz
        if  ;; label = @3
          i32.const 66376
          return
        end
        local.get 0
        i32.const 15
        i32.add
        i32.const 4
        i32.shr_u
        local.set 6
        i32.const 66368
        i32.load
        local.tee 4
        local.set 5
        i32.const 0
        local.set 7
        i32.const 0
        local.set 2
      end
      loop  ;; label = @2
        local.get 1
        local.get 4
        local.get 5
        i32.ne
        global.get 1
        select
        local.set 1
        block  ;; label = @3
          block  ;; label = @4
            global.get 1
            i32.eqz
            if  ;; label = @5
              local.get 1
              br_if 1 (;@4;)
              local.get 2
              i32.const 255
              i32.and
              local.set 4
              i32.const 1
              local.set 1
            end
            block  ;; label = @5
              global.get 1
              i32.eqz
              if  ;; label = @6
                block  ;; label = @7
                  local.get 4
                  br_table 4 (;@3;) 0 (;@7;) 2 (;@5;)
                end
                i32.const 66380
                local.set 1
              end
              loop  ;; label = @6
                block  ;; label = @7
                  global.get 1
                  i32.eqz
                  if  ;; label = @8
                    local.get 1
                    i32.load
                    local.tee 1
                    i32.eqz
                    local.tee 2
                    br_if 1 (;@7;)
                    local.get 1
                    i32.const 8
                    i32.add
                    local.tee 2
                    local.get 1
                    i32.load offset=4
                    i32.const 2
                    i32.shl
                    i32.add
                    local.set 4
                  end
                  local.get 8
                  i32.const 0
                  global.get 1
                  select
                  i32.eqz
                  if  ;; label = @8
                    local.get 2
                    local.get 4
                    call $runtime.markRoots
                    i32.const 0
                    global.get 1
                    i32.const 1
                    i32.eq
                    br_if 7 (;@1;)
                    drop
                  end
                  global.get 1
                  i32.eqz
                  br_if 1 (;@6;)
                end
              end
              local.get 8
              i32.const 1
              i32.eq
              i32.const 1
              global.get 1
              select
              if  ;; label = @6
                i32.const 65536
                i32.const 66656
                call $runtime.markRoots
                i32.const 1
                global.get 1
                i32.const 1
                i32.eq
                br_if 5 (;@1;)
                drop
              end
              loop  ;; label = @6
                global.get 1
                i32.eqz
                if  ;; label = @7
                  i32.const 66377
                  i32.load8_u
                  local.tee 2
                  i32.eqz
                  if  ;; label = @8
                    i32.const 0
                    local.set 2
                    i32.const 0
                    local.set 1
                    loop  ;; label = @9
                      local.get 1
                      i32.const 66372
                      i32.load
                      i32.ge_u
                      if  ;; label = @10
                        i32.const 2
                        local.set 1
                        br 7 (;@3;)
                      else
                        block  ;; label = @11
                          block  ;; label = @12
                            block  ;; label = @13
                              block  ;; label = @14
                                local.get 1
                                call $_runtime.gcBlock_.state
                                i32.const 255
                                i32.and
                                i32.const 1
                                i32.sub
                                br_table 1 (;@13;) 0 (;@14;) 2 (;@12;) 3 (;@11;)
                              end
                              local.get 2
                              i32.const 1
                              i32.and
                              local.set 4
                              i32.const 0
                              local.set 2
                              local.get 4
                              i32.eqz
                              br_if 2 (;@11;)
                            end
                            local.get 1
                            call $_runtime.gcBlock_.markFree
                            i32.const 1
                            local.set 2
                            br 1 (;@11;)
                          end
                          i32.const 0
                          local.set 2
                          i32.const 66364
                          i32.load
                          local.get 1
                          i32.const 2
                          i32.shr_u
                          i32.add
                          local.tee 4
                          i32.load8_u
                          i32.const 2
                          local.get 1
                          i32.const 1
                          i32.shl
                          i32.const 6
                          i32.and
                          i32.shl
                          i32.const -1
                          i32.xor
                          i32.and
                          local.set 3
                          local.get 4
                          local.get 3
                          i32.store8
                        end
                        local.get 1
                        i32.const 1
                        i32.add
                        local.set 1
                        br 1 (;@9;)
                      end
                      unreachable
                    end
                    unreachable
                  end
                  i32.const 66377
                  i32.const 0
                  i32.store8
                  i32.const 0
                  local.set 1
                end
                loop  ;; label = @7
                  global.get 1
                  i32.eqz
                  if  ;; label = @8
                    i32.const 66372
                    i32.load
                    local.get 1
                    i32.le_u
                    local.tee 2
                    br_if 2 (;@6;)
                    local.get 1
                    call $_runtime.gcBlock_.state
                    i32.const 255
                    i32.and
                    i32.const 3
                    i32.ne
                    local.set 2
                  end
                  global.get 1
                  i32.const 1
                  local.get 2
                  select
                  i32.const 0
                  local.get 8
                  i32.const 2
                  i32.eq
                  i32.const 1
                  global.get 1
                  select
                  select
                  if  ;; label = @8
                    local.get 1
                    call $runtime.startMark
                    i32.const 2
                    global.get 1
                    i32.const 1
                    i32.eq
                    br_if 7 (;@1;)
                    drop
                  end
                  global.get 1
                  i32.eqz
                  if  ;; label = @8
                    local.get 1
                    i32.const 1
                    i32.add
                    local.set 1
                    br 1 (;@7;)
                  end
                end
              end
            end
            global.get 1
            i32.eqz
            if  ;; label = @5
              memory.size
              memory.grow
              i32.const -1
              i32.eq
              local.tee 1
              i32.eqz
              if  ;; label = @6
                memory.size
                local.set 1
                i32.const 66228
                i32.load
                local.set 4
                i32.const 66228
                local.get 1
                i32.const 16
                i32.shl
                i32.store
                i32.const 66364
                i32.load
                local.set 1
                call $runtime.calculateHeapAddresses
                i32.const 66364
                i32.load
                local.get 1
                local.get 4
                local.get 1
                i32.sub
                memory.copy
                br 2 (;@4;)
              end
            end
            local.get 8
            i32.const 3
            i32.eq
            i32.const 1
            global.get 1
            select
            if  ;; label = @5
              i32.const 65550
              i32.const 13
              call $runtime.runtimePanic
              i32.const 3
              global.get 1
              i32.const 1
              i32.eq
              br_if 4 (;@1;)
              drop
            end
            global.get 1
            i32.eqz
            if  ;; label = @5
              unreachable
            end
          end
          local.get 1
          local.get 2
          global.get 1
          select
          local.set 1
        end
        global.get 1
        i32.eqz
        if  ;; label = @3
          block  ;; label = @4
            block  ;; label = @5
              i32.const 66372
              i32.load
              local.get 5
              i32.eq
              if  ;; label = @6
                i32.const 0
                local.set 5
                br 1 (;@5;)
              end
              local.get 5
              call $_runtime.gcBlock_.state
              i32.const 255
              i32.and
              if  ;; label = @6
                local.get 5
                i32.const 1
                i32.add
                local.set 5
                br 1 (;@5;)
              end
              local.get 5
              i32.const 1
              i32.add
              local.set 2
              local.get 7
              i32.const 1
              i32.add
              local.tee 7
              local.get 6
              i32.ne
              if  ;; label = @6
                local.get 2
                local.set 5
                br 2 (;@4;)
              end
              i32.const 66368
              local.get 2
              i32.store
              local.get 2
              local.get 6
              i32.sub
              local.tee 2
              i32.const 1
              call $_runtime.gcBlock_.setState
              local.get 5
              local.get 6
              i32.sub
              i32.const 2
              i32.add
              local.set 1
              loop  ;; label = @6
                local.get 1
                i32.const 66368
                i32.load
                i32.ne
                if  ;; label = @7
                  local.get 1
                  i32.const 2
                  call $_runtime.gcBlock_.setState
                  local.get 1
                  i32.const 1
                  i32.add
                  local.set 1
                  br 1 (;@6;)
                end
              end
              local.get 2
              i32.const 4
              i32.shl
              i32.const 66656
              i32.add
              local.tee 1
              i32.const 0
              local.get 0
              memory.fill
              local.get 1
              return
            end
            i32.const 0
            local.set 7
          end
          i32.const 66368
          i32.load
          local.set 4
          local.get 1
          local.set 2
          br 1 (;@2;)
        end
      end
      unreachable
    end
    local.set 3
    global.get 2
    i32.load
    local.get 3
    i32.store
    global.get 2
    global.get 2
    i32.load
    i32.const 4
    i32.add
    i32.store
    global.get 2
    i32.load
    local.tee 3
    local.get 0
    i32.store
    local.get 3
    local.get 1
    i32.store offset=4
    local.get 3
    local.get 2
    i32.store offset=8
    local.get 3
    local.get 4
    i32.store offset=12
    local.get 3
    local.get 5
    i32.store offset=16
    local.get 3
    local.get 6
    i32.store offset=20
    local.get 3
    local.get 7
    i32.store offset=24
    global.get 2
    global.get 2
    i32.load
    i32.const 28
    i32.add
    i32.store
    i32.const 0)
  (func $runtime.runqueuePushBack (type 1) (param i32)
    (local i32 i32)
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
      i32.eqz
      if  ;; label = @2
        block  ;; label = @3
          i32.const 66360
          i32.load
          if  ;; label = @4
            i32.const 66360
            i32.load
            local.tee 2
            i32.eqz
            br_if 1 (;@3;)
            local.get 2
            local.get 0
            i32.store
          end
          i32.const 66360
          local.get 0
          i32.store
          local.get 0
          i32.eqz
          br_if 0 (;@3;)
          local.get 0
          i32.const 0
          i32.store
          i32.const 66356
          i32.load
          i32.eqz
          if  ;; label = @4
            i32.const 66356
            local.get 0
            i32.store
          end
          return
        end
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
  (func $internal/task.Pause (type 0)
    (local i32 i32 i32)
    global.get 1
    i32.const 2
    i32.eq
    if  ;; label = @1
      global.get 2
      global.get 2
      i32.load
      i32.const 4
      i32.sub
      i32.store
      global.get 2
      i32.load
      i32.load
      local.set 0
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
      global.get 1
      i32.eqz
      if  ;; label = @2
        i32.const 66224
        i32.load
        local.tee 0
        i32.eqz
        local.set 2
      end
      block  ;; label = @2
        block  ;; label = @3
          block  ;; label = @4
            global.get 1
            i32.eqz
            if  ;; label = @5
              local.get 2
              br_if 1 (;@4;)
              local.get 0
              i32.const 28
              i32.add
              i32.load
              i32.load
              i32.const -1204030091
              i32.ne
              local.tee 0
              br_if 2 (;@3;)
              i32.const 66224
              i32.load
              local.tee 0
              i32.eqz
              br_if 1 (;@4;)
              local.get 0
              i32.const 28
              i32.add
              local.set 0
            end
            local.get 1
            i32.const 0
            global.get 1
            select
            i32.eqz
            if  ;; label = @5
              local.get 0
              call $tinygo_unwind
              i32.const 0
              global.get 1
              i32.const 1
              i32.eq
              br_if 4 (;@1;)
              drop
            end
            global.get 1
            i32.eqz
            if  ;; label = @5
              i32.const 66224
              i32.load
              local.tee 0
              br_if 3 (;@2;)
            end
          end
          local.get 1
          i32.const 1
          i32.eq
          i32.const 1
          global.get 1
          select
          if  ;; label = @4
            call $runtime.nilPanic
            i32.const 1
            global.get 1
            i32.const 1
            i32.eq
            br_if 3 (;@1;)
            drop
          end
          global.get 1
          i32.eqz
          if  ;; label = @4
            unreachable
          end
        end
        local.get 1
        i32.const 2
        i32.eq
        i32.const 1
        global.get 1
        select
        if  ;; label = @3
          i32.const 65536
          i32.const 14
          call $runtime.runtimePanic
          i32.const 2
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
        local.get 0
        i32.const 28
        i32.add
        i32.load
        i32.const -1204030091
        i32.store
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
    local.get 0
    i32.store
    global.get 2
    global.get 2
    i32.load
    i32.const 4
    i32.add
    i32.store)
  (func $runtime.markRoots (type 2) (param i32 i32)
    (local i32 i32 i32)
    global.get 1
    i32.const 2
    i32.eq
    if  ;; label = @1
      global.get 2
      global.get 2
      i32.load
      i32.const 16
      i32.sub
      i32.store
      global.get 2
      i32.load
      local.tee 2
      i32.load
      local.set 0
      local.get 2
      i32.load offset=4
      local.set 1
      local.get 2
      i32.load offset=8
      local.set 4
      local.get 2
      i32.load offset=12
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
        local.set 3
      end
      loop  ;; label = @2
        block  ;; label = @3
          global.get 1
          i32.eqz
          if  ;; label = @4
            local.get 0
            local.get 1
            i32.ge_u
            br_if 1 (;@3;)
            local.get 0
            i32.load
            local.tee 4
            call $runtime.looksLikePointer
            i32.const 1
            i32.and
            i32.eqz
            local.set 2
          end
          block  ;; label = @4
            global.get 1
            i32.eqz
            if  ;; label = @5
              local.get 2
              br_if 1 (;@4;)
              local.get 4
              i32.const 66656
              i32.sub
              i32.const 4
              i32.shr_u
              local.tee 4
              call $_runtime.gcBlock_.state
              i32.const 255
              i32.and
              i32.eqz
              local.tee 2
              br_if 1 (;@4;)
              local.get 4
              call $_runtime.gcBlock_.findHead
              local.tee 4
              call $_runtime.gcBlock_.state
              i32.const 255
              i32.and
              i32.const 3
              i32.eq
              local.tee 2
              br_if 1 (;@4;)
            end
            local.get 3
            i32.const 0
            global.get 1
            select
            i32.eqz
            if  ;; label = @5
              local.get 4
              call $runtime.startMark
              i32.const 0
              global.get 1
              i32.const 1
              i32.eq
              br_if 4 (;@1;)
              drop
            end
          end
          global.get 1
          i32.eqz
          if  ;; label = @4
            local.get 0
            i32.const 4
            i32.add
            local.set 0
            br 2 (;@2;)
          end
        end
      end
      return
    end
    local.set 3
    global.get 2
    i32.load
    local.get 3
    i32.store
    global.get 2
    global.get 2
    i32.load
    i32.const 4
    i32.add
    i32.store
    global.get 2
    i32.load
    local.tee 3
    local.get 0
    i32.store
    local.get 3
    local.get 1
    i32.store offset=4
    local.get 3
    local.get 4
    i32.store offset=8
    local.get 3
    local.get 2
    i32.store offset=12
    global.get 2
    global.get 2
    i32.load
    i32.const 16
    i32.add
    i32.store)
  (func $_runtime.gcBlock_.state (type 3) (param i32) (result i32)
    i32.const 66364
    i32.load
    local.get 0
    i32.const 2
    i32.shr_u
    i32.add
    i32.load8_u
    local.get 0
    i32.const 1
    i32.shl
    i32.const 6
    i32.and
    i32.shr_u
    i32.const 3
    i32.and)
  (func $_runtime.gcBlock_.markFree (type 1) (param i32)
    (local i32 i32)
    i32.const 66364
    i32.load
    local.get 0
    i32.const 2
    i32.shr_u
    i32.add
    local.tee 1
    i32.load8_u
    local.set 2
    local.get 1
    local.get 2
    i32.const 3
    local.get 0
    i32.const 1
    i32.shl
    i32.const 6
    i32.and
    i32.shl
    i32.const -1
    i32.xor
    i32.and
    i32.store8)
  (func $runtime.startMark (type 1) (param i32)
    (local i32 i32 i32 i32 i32)
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
        local.set 5
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        global.get $__stack_pointer
        i32.const -64
        i32.add
        local.tee 3
        global.set $__stack_pointer
        local.get 3
        i32.const 0
        i32.const 64
        memory.fill
        local.get 3
        local.get 0
        i32.store
        local.get 0
        i32.const 3
        call $_runtime.gcBlock_.setState
        i32.const 1
        local.set 1
        block  ;; label = @3
          loop  ;; label = @4
            local.get 1
            i32.const 0
            i32.gt_s
            if  ;; label = @5
              local.get 1
              i32.const 1
              i32.sub
              local.tee 1
              i32.const 15
              i32.gt_u
              br_if 2 (;@3;)
              local.get 3
              local.get 1
              i32.const 2
              i32.shl
              i32.add
              i32.load
              local.tee 0
              call $_runtime.gcBlock_.findNext
              i32.const 4
              i32.shl
              local.get 0
              i32.const 4
              i32.shl
              local.tee 4
              i32.sub
              local.set 0
              local.get 4
              i32.const 66656
              i32.add
              local.set 4
              loop  ;; label = @6
                local.get 0
                i32.eqz
                br_if 2 (;@4;)
                block  ;; label = @7
                  local.get 4
                  i32.load
                  local.tee 2
                  call $runtime.looksLikePointer
                  i32.const 1
                  i32.and
                  i32.eqz
                  br_if 0 (;@7;)
                  local.get 2
                  i32.const 66656
                  i32.sub
                  i32.const 4
                  i32.shr_u
                  local.tee 2
                  call $_runtime.gcBlock_.state
                  i32.const 255
                  i32.and
                  i32.eqz
                  br_if 0 (;@7;)
                  local.get 2
                  call $_runtime.gcBlock_.findHead
                  local.tee 2
                  call $_runtime.gcBlock_.state
                  i32.const 255
                  i32.and
                  i32.const 3
                  i32.eq
                  br_if 0 (;@7;)
                  local.get 2
                  i32.const 3
                  call $_runtime.gcBlock_.setState
                  local.get 1
                  i32.const 16
                  i32.eq
                  if  ;; label = @8
                    i32.const 66377
                    i32.const 1
                    i32.store8
                    i32.const 16
                    local.set 1
                    br 1 (;@7;)
                  end
                  local.get 1
                  i32.const 15
                  i32.gt_u
                  br_if 4 (;@3;)
                  local.get 3
                  local.get 1
                  i32.const 2
                  i32.shl
                  i32.add
                  local.get 2
                  i32.store
                  local.get 1
                  i32.const 1
                  i32.add
                  local.set 1
                end
                local.get 0
                i32.const 4
                i32.sub
                local.set 0
                local.get 4
                i32.const 4
                i32.add
                local.set 4
                br 0 (;@6;)
              end
              unreachable
            end
          end
          local.get 3
          i32.const -64
          i32.sub
          global.set $__stack_pointer
          return
        end
      end
      local.get 5
      i32.const 0
      global.get 1
      select
      i32.eqz
      if  ;; label = @2
        call $runtime.lookupPanic
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
  (func $runtime.calculateHeapAddresses (type 0)
    (local i32 i32)
    i32.const 66228
    i32.load
    local.tee 0
    i32.const 66592
    i32.sub
    i32.const 65
    i32.div_u
    local.set 1
    i32.const 66364
    local.get 0
    local.get 1
    i32.sub
    local.tee 0
    i32.store
    i32.const 66372
    local.get 0
    i32.const 66656
    i32.sub
    i32.const 4
    i32.shr_u
    i32.store)
  (func $_runtime.gcBlock_.setState (type 2) (param i32 i32)
    (local i32 i32)
    i32.const 66364
    i32.load
    local.get 0
    i32.const 2
    i32.shr_u
    i32.add
    local.tee 2
    i32.load8_u
    local.set 3
    local.get 2
    local.get 3
    local.get 1
    local.get 0
    i32.const 1
    i32.shl
    i32.const 6
    i32.and
    i32.shl
    i32.or
    i32.store8)
  (func $runtime.printstring (type 2) (param i32 i32)
    (local i32 i32)
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
      local.tee 2
      i32.load
      local.set 0
      local.get 2
      i32.load offset=4
      local.set 1
      local.get 2
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
        local.set 3
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        local.get 1
        i32.const 0
        local.get 1
        i32.const 0
        i32.gt_s
        local.tee 2
        select
        local.set 1
      end
      loop  ;; label = @2
        block  ;; label = @3
          global.get 1
          i32.eqz
          if  ;; label = @4
            local.get 1
            i32.eqz
            br_if 1 (;@3;)
            local.get 0
            i32.load8_u
            local.set 2
          end
          local.get 3
          i32.const 0
          global.get 1
          select
          i32.eqz
          if  ;; label = @4
            local.get 2
            call $runtime.putchar
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
            local.get 0
            i32.const 1
            i32.add
            local.set 0
            local.get 1
            i32.const 1
            i32.sub
            local.set 1
            br 2 (;@2;)
          end
        end
      end
      return
    end
    local.set 3
    global.get 2
    i32.load
    local.get 3
    i32.store
    global.get 2
    global.get 2
    i32.load
    i32.const 4
    i32.add
    i32.store
    global.get 2
    i32.load
    local.tee 3
    local.get 0
    i32.store
    local.get 3
    local.get 1
    i32.store offset=4
    local.get 3
    local.get 2
    i32.store offset=8
    global.get 2
    global.get 2
    i32.load
    i32.const 12
    i32.add
    i32.store)
  (func $runtime.printnl (type 0)
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
        i32.const 10
        call $runtime.putchar
        i32.const 0
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
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
  (func $runtime.putchar (type 1) (param i32)
    (local i32 i32 i32)
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
      global.get 1
      i32.eqz
      if  ;; label = @2
        i32.const 66232
        i32.load
        local.tee 3
        i32.const 119
        i32.gt_u
        local.set 1
      end
      block  ;; label = @2
        global.get 1
        i32.eqz
        if  ;; label = @3
          local.get 1
          br_if 1 (;@2;)
          i32.const 66232
          local.get 3
          i32.const 1
          i32.add
          local.tee 1
          i32.store
          local.get 3
          i32.const 66236
          i32.add
          local.get 0
          i32.store8
          local.get 0
          i32.const 255
          i32.and
          i32.const 10
          i32.eq
          local.set 0
        end
        block  ;; label = @3
          global.get 1
          i32.eqz
          if  ;; label = @4
            local.get 0
            i32.eqz
            local.get 3
            i32.const 119
            i32.ne
            i32.and
            br_if 1 (;@3;)
            i32.const 66188
            local.get 1
            i32.store
          end
          local.get 2
          i32.const 0
          global.get 1
          select
          i32.eqz
          if  ;; label = @4
            i32.const 1
            i32.const 66184
            i32.const 1
            i32.const 66384
            call $runtime.fd_write
            drop
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
            i32.const 66232
            i32.const 0
            i32.store
          end
        end
        global.get 1
        i32.eqz
        if  ;; label = @3
          return
        end
      end
      local.get 2
      i32.const 1
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        call $runtime.lookupPanic
        i32.const 1
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
  (func $runtime.lookupPanic (type 0)
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
        i32.const 65615
        i32.const 18
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
  (func $_runtime.gcBlock_.findNext (type 3) (param i32) (result i32)
    (local i32)
    block  ;; label = @1
      local.get 0
      call $_runtime.gcBlock_.state
      i32.const 255
      i32.and
      i32.const 1
      i32.ne
      if  ;; label = @2
        local.get 0
        call $_runtime.gcBlock_.state
        i32.const 255
        i32.and
        i32.const 3
        i32.ne
        br_if 1 (;@1;)
      end
      local.get 0
      i32.const 1
      i32.add
      local.set 0
    end
    local.get 0
    i32.const 4
    i32.shl
    i32.const 66656
    i32.add
    local.set 1
    loop  ;; label = @1
      block  ;; label = @2
        i32.const 66364
        i32.load
        local.get 1
        i32.le_u
        br_if 0 (;@2;)
        local.get 0
        call $_runtime.gcBlock_.state
        i32.const 255
        i32.and
        i32.const 2
        i32.ne
        br_if 0 (;@2;)
        local.get 1
        i32.const 16
        i32.add
        local.set 1
        local.get 0
        i32.const 1
        i32.add
        local.set 0
        br 1 (;@1;)
      end
    end
    local.get 0)
  (func $runtime.looksLikePointer (type 3) (param i32) (result i32)
    (local i32)
    i32.const 0
    local.set 1
    local.get 0
    i32.const 66656
    i32.ge_u
    if (result i32)  ;; label = @1
      i32.const 66364
      i32.load
      local.get 0
      i32.gt_u
    else
      local.get 1
    end)
  (func $_runtime.gcBlock_.findHead (type 3) (param i32) (result i32)
    (local i32)
    loop  ;; label = @1
      local.get 0
      call $_runtime.gcBlock_.state
      local.set 1
      local.get 0
      i32.const 1
      i32.sub
      local.set 0
      local.get 1
      i32.const 255
      i32.and
      i32.const 2
      i32.eq
      br_if 0 (;@1;)
    end
    local.get 0
    i32.const 1
    i32.add)
  (func $runtime.slicePanic (type 0)
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
        i32.const 65633
        i32.const 18
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
  (func $runtime.hash32 (type 4) (param i32 i32 i32 i32) (result i32)
    local.get 2
    i32.const -2128831035
    i32.mul
    local.set 2
    loop  ;; label = @1
      local.get 1
      if  ;; label = @2
        local.get 1
        i32.const 1
        i32.sub
        local.set 1
        local.get 0
        i32.load8_u
        local.get 2
        i32.xor
        i32.const 16777619
        i32.mul
        local.set 2
        local.get 0
        i32.const 1
        i32.add
        local.set 0
        br 1 (;@1;)
      end
    end
    local.get 2)
  (func $malloc (type 3) (param i32) (result i32)
    (local i32 i32 i32)
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
      local.tee 2
      i32.load
      local.set 0
      local.get 2
      i32.load offset=4
      local.set 3
      local.get 2
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
        local.set 1
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        global.get $__stack_pointer
        i32.const 16
        i32.sub
        local.tee 3
        global.set $__stack_pointer
        local.get 3
        i64.const 1
        i64.store offset=4 align=4
        i32.const 66380
        i32.load
        local.set 2
        i32.const 66380
        local.get 3
        i32.store
        local.get 3
        local.get 2
        i32.store
      end
      local.get 1
      i32.const 0
      global.get 1
      select
      i32.eqz
      if  ;; label = @2
        local.get 0
        call $runtime.alloc
        local.set 1
        i32.const 0
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
        i32.const 66380
        local.get 2
        i32.store
        local.get 3
        i32.const 8
        i32.add
        local.get 0
        i32.store
        local.get 3
        i32.const 16
        i32.add
        global.set $__stack_pointer
        local.get 0
        return
      end
      unreachable
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
    global.get 2
    global.get 2
    i32.load
    i32.const 12
    i32.add
    i32.store
    i32.const 0)
  (func $free (type 1) (param i32)
    nop)
  (func $calloc (type 6) (param i32 i32) (result i32)
    (local i32 i32 i32)
    global.get 1
    i32.const 2
    i32.eq
    if  ;; label = @1
      global.get 2
      global.get 2
      i32.load
      i32.const 16
      i32.sub
      i32.store
      global.get 2
      i32.load
      local.tee 3
      i32.load
      local.set 0
      local.get 3
      i32.load offset=4
      local.set 1
      local.get 3
      i32.load offset=8
      local.set 4
      local.get 3
      i32.load offset=12
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
        local.set 2
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        global.get $__stack_pointer
        i32.const 16
        i32.sub
        local.tee 4
        global.set $__stack_pointer
        local.get 4
        i64.const 1
        i64.store offset=4 align=4
        i32.const 66380
        i32.load
        local.set 3
        i32.const 66380
        local.get 4
        i32.store
        local.get 4
        local.get 3
        i32.store
        local.get 0
        local.get 1
        i32.mul
        local.set 0
      end
      local.get 2
      i32.const 0
      global.get 1
      select
      i32.eqz
      if  ;; label = @2
        local.get 0
        call $runtime.alloc
        local.set 2
        i32.const 0
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
        local.get 2
        local.set 1
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        i32.const 66380
        local.get 3
        i32.store
        local.get 4
        i32.const 8
        i32.add
        local.get 1
        i32.store
        local.get 4
        i32.const 16
        i32.add
        global.set $__stack_pointer
        local.get 1
        return
      end
      unreachable
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
    local.get 2
    local.get 4
    i32.store offset=8
    local.get 2
    local.get 3
    i32.store offset=12
    global.get 2
    global.get 2
    i32.load
    i32.const 16
    i32.add
    i32.store
    i32.const 0)
  (func $realloc (type 6) (param i32 i32) (result i32)
    (local i32 i32 i32 i32 i32 i32)
    global.get 1
    i32.const 2
    i32.eq
    if  ;; label = @1
      global.get 2
      global.get 2
      i32.load
      i32.const 24
      i32.sub
      i32.store
      global.get 2
      i32.load
      local.tee 2
      i32.load
      local.set 0
      local.get 2
      i32.load offset=8
      local.set 3
      local.get 2
      i32.load offset=12
      local.set 4
      local.get 2
      i32.load offset=16
      local.set 5
      local.get 2
      i32.load offset=20
      local.set 6
      local.get 2
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
        local.set 7
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
        i64.const 0
        i64.store offset=20 align=4
        local.get 3
        i64.const 3
        i64.store offset=12 align=4
        i32.const 66380
        i32.load
        local.set 4
        i32.const 66380
        local.get 3
        i32.const 8
        i32.add
        local.tee 5
        i32.store
        local.get 3
        local.get 4
        i32.store offset=8
      end
      block  ;; label = @2
        block  ;; label = @3
          global.get 1
          i32.eqz
          if  ;; label = @4
            local.get 0
            br_if 1 (;@3;)
            local.get 3
            i32.const 16
            i32.add
            local.set 0
          end
          local.get 7
          i32.const 0
          global.get 1
          select
          i32.eqz
          if  ;; label = @4
            local.get 1
            call $runtime.alloc
            local.set 2
            i32.const 0
            global.get 1
            i32.const 1
            i32.eq
            br_if 3 (;@1;)
            drop
            local.get 2
            local.set 1
          end
          global.get 1
          i32.eqz
          if  ;; label = @4
            local.get 0
            local.get 1
            i32.store
            br 2 (;@2;)
          end
        end
        global.get 1
        i32.eqz
        if  ;; label = @3
          local.get 1
          i32.const 66656
          local.get 0
          i32.sub
          local.get 0
          i32.const 66656
          i32.sub
          i32.const 4
          i32.shr_u
          call $_runtime.gcBlock_.findNext
          i32.const 4
          i32.shl
          i32.add
          local.tee 5
          i32.le_u
          if  ;; label = @4
            local.get 0
            local.set 1
            br 2 (;@2;)
          end
          local.get 3
          i32.const 20
          i32.add
          local.set 6
        end
        local.get 7
        i32.const 1
        i32.eq
        i32.const 1
        global.get 1
        select
        if  ;; label = @3
          local.get 1
          call $runtime.alloc
          local.set 2
          i32.const 1
          global.get 1
          i32.const 1
          i32.eq
          br_if 2 (;@1;)
          drop
          local.get 2
          local.set 1
        end
        global.get 1
        i32.eqz
        if  ;; label = @3
          local.get 6
          local.get 1
          i32.store
          local.get 1
          local.get 0
          local.get 5
          memory.copy
        end
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        i32.const 66380
        local.get 4
        i32.store
        local.get 3
        i32.const 24
        i32.add
        local.get 1
        i32.store
        local.get 3
        i32.const 32
        i32.add
        global.set $__stack_pointer
        local.get 1
        return
      end
      unreachable
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
    local.get 2
    local.get 3
    i32.store offset=8
    local.get 2
    local.get 4
    i32.store offset=12
    local.get 2
    local.get 5
    i32.store offset=16
    local.get 2
    local.get 6
    i32.store offset=20
    global.get 2
    global.get 2
    i32.load
    i32.const 24
    i32.add
    i32.store
    i32.const 0)
  (func $runtime.deadlock (type 0)
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
        local.set 0
      end
      local.get 0
      i32.const 0
      global.get 1
      select
      i32.eqz
      if  ;; label = @2
        call $internal/task.Pause
        i32.const 0
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
      end
      local.get 0
      i32.const 1
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        i32.const 34
        i32.const 65672
        call $runtime._panic
        i32.const 1
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
  (func $runtime._panic (type 2) (param i32 i32)
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
        i32.const 65563
        i32.const 7
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
        call $runtime.printitf
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
  (func $runtime.printitf (type 2) (param i32 i32)
    (local i32 i32 i32 i32 i32 i64 i32)
    global.get 1
    i32.const 2
    i32.eq
    if  ;; label = @1
      global.get 2
      global.get 2
      i32.load
      i32.const 36
      i32.sub
      i32.store
      global.get 2
      i32.load
      local.tee 6
      i32.load
      local.set 0
      local.get 6
      i32.load offset=4
      local.set 1
      local.get 6
      i32.load offset=8
      local.set 2
      local.get 6
      i32.load offset=12
      local.set 3
      local.get 6
      i32.load offset=16
      local.set 5
      local.get 6
      i64.load offset=20 align=4
      local.set 7
      local.get 6
      i32.load offset=28
      local.set 8
      local.get 6
      i32.load offset=32
      local.set 6
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
        local.set 4
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        global.get $__stack_pointer
        i32.const 208
        i32.sub
        local.tee 2
        global.set $__stack_pointer
        local.get 2
        i64.const 60129542144
        i64.store offset=144
        local.get 2
        i64.const 0
        i64.store offset=200
        local.get 2
        i64.const 0
        i64.store offset=192
        local.get 2
        i64.const 0
        i64.store offset=184
        local.get 2
        i64.const 0
        i64.store offset=176
        local.get 2
        i64.const 0
        i64.store offset=168
        local.get 2
        i64.const 0
        i64.store offset=160
        local.get 2
        i64.const 0
        i64.store offset=152
        i32.const 66380
        i32.load
        local.set 8
        i32.const 66380
        local.get 2
        i32.const 144
        i32.add
        i32.store
        local.get 2
        local.get 8
        i32.store offset=144
        local.get 0
        i32.const 34
        i32.eq
        local.set 3
      end
      block  ;; label = @2
        block  ;; label = @3
          block  ;; label = @4
            block  ;; label = @5
              block (result i32)  ;; label = @6
                block  ;; label = @7
                  block  ;; label = @8
                    block  ;; label = @9
                      block  ;; label = @10
                        block  ;; label = @11
                          block  ;; label = @12
                            block  ;; label = @13
                              block  ;; label = @14
                                block  ;; label = @15
                                  block  ;; label = @16
                                    block  ;; label = @17
                                      global.get 1
                                      i32.eqz
                                      if  ;; label = @18
                                        local.get 3
                                        br_if 1 (;@17;)
                                        local.get 0
                                        i32.const 63
                                        i32.eq
                                        local.tee 3
                                        br_if 10 (;@8;)
                                        local.get 0
                                        i32.const 68
                                        i32.eq
                                        local.tee 3
                                        br_if 11 (;@7;)
                                        local.get 0
                                        i32.const 79
                                        i32.eq
                                        local.tee 3
                                        br_if 6 (;@12;)
                                        local.get 0
                                        i32.const 127
                                        i32.eq
                                        local.tee 3
                                        br_if 4 (;@14;)
                                        local.get 0
                                        i32.const 2021
                                        i32.eq
                                        local.tee 3
                                        br_if 8 (;@10;)
                                        local.get 0
                                        i32.const 2181
                                        i32.eq
                                        local.tee 3
                                        br_if 9 (;@9;)
                                        local.get 0
                                        i32.const 2533
                                        i32.eq
                                        local.tee 3
                                        br_if 7 (;@11;)
                                        local.get 0
                                        i32.const 3045
                                        i32.eq
                                        local.tee 3
                                        br_if 2 (;@16;)
                                        local.get 0
                                        i32.const 4069
                                        i32.eq
                                        local.tee 3
                                        br_if 3 (;@15;)
                                        local.get 2
                                        i64.const 0
                                        i64.store offset=112
                                        local.get 2
                                        local.get 0
                                        i32.store offset=104
                                        local.get 2
                                        local.get 1
                                        i32.store offset=108
                                        local.get 2
                                        i32.const 196
                                        i32.add
                                        local.get 2
                                        i32.const 112
                                        i32.add
                                        i32.store
                                        local.get 2
                                        i32.const 200
                                        i32.add
                                        local.get 2
                                        i32.const 112
                                        i32.add
                                        i32.store
                                        local.get 2
                                        i32.const 188
                                        i32.add
                                        local.get 2
                                        i32.const 112
                                        i32.add
                                        i32.store
                                        local.get 2
                                        local.get 0
                                        i32.store offset=112
                                        local.get 2
                                        local.get 1
                                        i32.store offset=116
                                        local.get 2
                                        i32.const 192
                                        i32.add
                                        local.tee 3
                                        local.get 2
                                        i32.const 104
                                        i32.add
                                        local.tee 5
                                        i32.store
                                      end
                                      local.get 4
                                      i32.const 0
                                      global.get 1
                                      select
                                      i32.eqz
                                      if  ;; label = @18
                                        i32.const 40
                                        call $runtime.putchar
                                        i32.const 0
                                        global.get 1
                                        i32.const 1
                                        i32.eq
                                        br_if 17 (;@1;)
                                        drop
                                      end
                                      global.get 1
                                      i32.eqz
                                      if  ;; label = @18
                                        local.get 2
                                        i32.const 204
                                        i32.add
                                        local.get 2
                                        i32.const 120
                                        i32.add
                                        local.tee 5
                                        i32.store
                                        local.get 2
                                        i32.const 136
                                        i32.add
                                        i32.const 0
                                        i32.store
                                        local.get 2
                                        i32.const 128
                                        i32.add
                                        i64.const 0
                                        i64.store
                                        local.get 2
                                        i64.const 0
                                        i64.store offset=120
                                        local.get 0
                                        i64.extend_i32_u
                                        local.set 7
                                        i32.const 19
                                        local.set 3
                                        i32.const 19
                                        local.set 0
                                        br 13 (;@5;)
                                      end
                                    end
                                    global.get 1
                                    i32.eqz
                                    if  ;; label = @17
                                      local.get 1
                                      i32.load
                                      local.set 0
                                      local.get 1
                                      i32.load offset=4
                                      local.set 1
                                    end
                                    local.get 4
                                    i32.const 1
                                    i32.eq
                                    i32.const 1
                                    global.get 1
                                    select
                                    if  ;; label = @17
                                      local.get 0
                                      local.get 1
                                      call $runtime.printstring
                                      i32.const 1
                                      global.get 1
                                      i32.const 1
                                      i32.eq
                                      br_if 16 (;@1;)
                                      drop
                                    end
                                    global.get 1
                                    i32.eqz
                                    br_if 14 (;@2;)
                                  end
                                  global.get 1
                                  i32.eqz
                                  if  ;; label = @16
                                    local.get 1
                                    i32.eqz
                                    br_if 12 (;@4;)
                                    local.get 2
                                    i32.const 152
                                    i32.add
                                    local.get 1
                                    i32.load
                                    local.tee 0
                                    i32.store
                                    local.get 1
                                    i32.load offset=4
                                    local.set 5
                                    local.get 2
                                    i32.const 48
                                    i32.add
                                    local.set 3
                                  end
                                  local.get 4
                                  i32.const 2
                                  i32.eq
                                  i32.const 1
                                  global.get 1
                                  select
                                  if  ;; label = @16
                                    local.get 3
                                    i32.const 66136
                                    i32.const 20
                                    local.get 0
                                    local.get 5
                                    call $runtime.stringConcat
                                    i32.const 2
                                    global.get 1
                                    i32.const 1
                                    i32.eq
                                    br_if 15 (;@1;)
                                    drop
                                  end
                                  global.get 1
                                  i32.eqz
                                  if  ;; label = @16
                                    local.get 2
                                    i32.const 156
                                    i32.add
                                    local.get 2
                                    i32.load offset=48
                                    local.tee 0
                                    i32.store
                                    local.get 2
                                    i32.load offset=52
                                    local.set 5
                                    local.get 2
                                    i32.const 40
                                    i32.add
                                    local.set 3
                                  end
                                  local.get 4
                                  i32.const 3
                                  i32.eq
                                  i32.const 1
                                  global.get 1
                                  select
                                  if  ;; label = @16
                                    local.get 3
                                    local.get 0
                                    local.get 5
                                    i32.const 66156
                                    i32.const 4
                                    call $runtime.stringConcat
                                    i32.const 3
                                    global.get 1
                                    i32.const 1
                                    i32.eq
                                    br_if 15 (;@1;)
                                    drop
                                  end
                                  global.get 1
                                  i32.eqz
                                  if  ;; label = @16
                                    local.get 2
                                    i32.const 160
                                    i32.add
                                    local.get 2
                                    i32.load offset=40
                                    local.tee 0
                                    i32.store
                                    local.get 2
                                    i32.load offset=44
                                    local.set 3
                                    local.get 2
                                    i32.const 32
                                    i32.add
                                    local.set 5
                                    local.get 1
                                    i32.load offset=8
                                    local.set 1
                                  end
                                  local.get 4
                                  i32.const 4
                                  i32.eq
                                  i32.const 1
                                  global.get 1
                                  select
                                  if  ;; label = @16
                                    local.get 5
                                    local.get 1
                                    call $_syscall/js.Type_.String
                                    i32.const 4
                                    global.get 1
                                    i32.const 1
                                    i32.eq
                                    br_if 15 (;@1;)
                                    drop
                                  end
                                  global.get 1
                                  i32.eqz
                                  if  ;; label = @16
                                    local.get 2
                                    i32.const 164
                                    i32.add
                                    local.get 2
                                    i32.load offset=32
                                    local.tee 1
                                    i32.store
                                    local.get 2
                                    i32.load offset=36
                                    local.set 6
                                    local.get 2
                                    i32.const 24
                                    i32.add
                                    local.set 5
                                  end
                                  local.get 4
                                  i32.const 5
                                  i32.eq
                                  i32.const 1
                                  global.get 1
                                  select
                                  if  ;; label = @16
                                    local.get 5
                                    local.get 0
                                    local.get 3
                                    local.get 1
                                    local.get 6
                                    call $runtime.stringConcat
                                    i32.const 5
                                    global.get 1
                                    i32.const 1
                                    i32.eq
                                    br_if 15 (;@1;)
                                    drop
                                  end
                                  global.get 1
                                  i32.eqz
                                  if  ;; label = @16
                                    local.get 2
                                    i32.const 168
                                    i32.add
                                    local.tee 1
                                    local.get 2
                                    i32.load offset=24
                                    local.tee 0
                                    i32.store
                                    local.get 2
                                    i32.load offset=28
                                    local.set 3
                                    br 3 (;@13;)
                                  end
                                end
                                global.get 1
                                i32.eqz
                                if  ;; label = @15
                                  local.get 1
                                  i32.eqz
                                  br_if 11 (;@4;)
                                  local.get 2
                                  i32.const 172
                                  i32.add
                                  local.get 1
                                  i32.load offset=8
                                  local.tee 0
                                  i32.store
                                  local.get 1
                                  i64.load
                                  local.set 7
                                  local.get 2
                                  i32.const 16
                                  i32.add
                                  local.set 3
                                end
                                local.get 4
                                i32.const 6
                                i32.eq
                                i32.const 1
                                global.get 1
                                select
                                if  ;; label = @15
                                  local.get 3
                                  local.get 7
                                  local.get 0
                                  call $_syscall/js.Error_.Error
                                  i32.const 6
                                  global.get 1
                                  i32.const 1
                                  i32.eq
                                  br_if 14 (;@1;)
                                  drop
                                end
                                global.get 1
                                i32.eqz
                                if  ;; label = @15
                                  local.get 2
                                  i32.const 176
                                  i32.add
                                  local.tee 1
                                  local.get 2
                                  i32.load offset=16
                                  local.tee 0
                                  i32.store
                                  local.get 2
                                  i32.load offset=20
                                  local.set 3
                                  br 2 (;@13;)
                                end
                              end
                              global.get 1
                              i32.eqz
                              if  ;; label = @14
                                local.get 1
                                i64.load
                                local.set 7
                                local.get 1
                                i32.load offset=8
                                local.set 1
                                local.get 2
                                i32.const 8
                                i32.add
                                local.set 0
                              end
                              local.get 4
                              i32.const 7
                              i32.eq
                              i32.const 1
                              global.get 1
                              select
                              if  ;; label = @14
                                local.get 0
                                local.get 7
                                local.get 1
                                call $_syscall/js.Error_.Error
                                i32.const 7
                                global.get 1
                                i32.const 1
                                i32.eq
                                br_if 13 (;@1;)
                                drop
                              end
                              global.get 1
                              i32.eqz
                              if  ;; label = @14
                                local.get 2
                                i32.load offset=12
                                local.set 3
                                local.get 2
                                i32.load offset=8
                                local.set 0
                              end
                            end
                            global.get 1
                            i32.eqz
                            if  ;; label = @13
                              local.get 2
                              i32.const 180
                              i32.add
                              local.tee 1
                              local.get 0
                              i32.store
                            end
                            local.get 4
                            i32.const 8
                            i32.eq
                            i32.const 1
                            global.get 1
                            select
                            if  ;; label = @13
                              local.get 0
                              local.get 3
                              call $runtime.printstring
                              i32.const 8
                              global.get 1
                              i32.const 1
                              i32.eq
                              br_if 12 (;@1;)
                              drop
                            end
                            global.get 1
                            i32.eqz
                            br_if 10 (;@2;)
                          end
                          local.get 0
                          local.get 2
                          i32.const 56
                          i32.add
                          global.get 1
                          select
                          local.set 0
                          local.get 4
                          i32.const 9
                          i32.eq
                          i32.const 1
                          global.get 1
                          select
                          if  ;; label = @12
                            local.get 0
                            local.get 1
                            call $_struct_syscall/js.Value__.String$invoke
                            i32.const 9
                            global.get 1
                            i32.const 1
                            i32.eq
                            br_if 11 (;@1;)
                            drop
                          end
                          global.get 1
                          i32.eqz
                          if  ;; label = @12
                            local.get 2
                            i32.load offset=60
                            local.set 3
                            local.get 2
                            i32.load offset=56
                            br 6 (;@6;)
                          end
                        end
                        local.get 0
                        local.get 2
                        i32.const -64
                        i32.sub
                        global.get 1
                        select
                        local.set 0
                        local.get 4
                        i32.const 10
                        i32.eq
                        i32.const 1
                        global.get 1
                        select
                        if  ;; label = @11
                          local.get 0
                          local.get 1
                          call $_*struct_syscall/js.Value__.String
                          i32.const 10
                          global.get 1
                          i32.const 1
                          i32.eq
                          br_if 10 (;@1;)
                          drop
                        end
                        global.get 1
                        i32.eqz
                        if  ;; label = @11
                          local.get 2
                          i32.load offset=68
                          local.set 3
                          local.get 2
                          i32.load offset=64
                          br 5 (;@6;)
                        end
                      end
                      local.get 0
                      local.get 2
                      i32.const 72
                      i32.add
                      global.get 1
                      select
                      local.set 0
                      local.get 4
                      i32.const 11
                      i32.eq
                      i32.const 1
                      global.get 1
                      select
                      if  ;; label = @10
                        local.get 0
                        local.get 1
                        call $_*syscall/js.Value_.String
                        i32.const 11
                        global.get 1
                        i32.const 1
                        i32.eq
                        br_if 9 (;@1;)
                        drop
                      end
                      global.get 1
                      i32.eqz
                      if  ;; label = @10
                        local.get 2
                        i32.load offset=76
                        local.set 3
                        local.get 2
                        i32.load offset=72
                        br 4 (;@6;)
                      end
                    end
                    local.get 0
                    local.get 2
                    i32.const 80
                    i32.add
                    global.get 1
                    select
                    local.set 0
                    local.get 4
                    i32.const 12
                    i32.eq
                    i32.const 1
                    global.get 1
                    select
                    if  ;; label = @9
                      local.get 0
                      local.get 1
                      call $_*syscall/js.Type_.String
                      i32.const 12
                      global.get 1
                      i32.const 1
                      i32.eq
                      br_if 8 (;@1;)
                      drop
                    end
                    global.get 1
                    i32.eqz
                    if  ;; label = @9
                      local.get 2
                      i32.load offset=84
                      local.set 3
                      local.get 2
                      i32.load offset=80
                      br 3 (;@6;)
                    end
                  end
                  global.get 1
                  i32.eqz
                  if  ;; label = @8
                    local.get 1
                    i64.load
                    local.set 7
                    local.get 1
                    i32.load offset=8
                    local.set 1
                    local.get 2
                    i32.const 88
                    i32.add
                    local.set 0
                  end
                  local.get 4
                  i32.const 13
                  i32.eq
                  i32.const 1
                  global.get 1
                  select
                  if  ;; label = @8
                    local.get 0
                    local.get 7
                    local.get 1
                    call $_syscall/js.Value_.String
                    i32.const 13
                    global.get 1
                    i32.const 1
                    i32.eq
                    br_if 7 (;@1;)
                    drop
                  end
                  global.get 1
                  i32.eqz
                  if  ;; label = @8
                    local.get 2
                    i32.load offset=92
                    local.set 3
                    local.get 2
                    i32.load offset=88
                    br 2 (;@6;)
                  end
                end
                local.get 0
                local.get 2
                i32.const 96
                i32.add
                global.get 1
                select
                local.set 0
                local.get 4
                i32.const 14
                i32.eq
                i32.const 1
                global.get 1
                select
                if  ;; label = @7
                  local.get 0
                  local.get 1
                  call $_syscall/js.Type_.String
                  i32.const 14
                  global.get 1
                  i32.const 1
                  i32.eq
                  br_if 6 (;@1;)
                  drop
                end
                global.get 1
                if (result i32)  ;; label = @7
                  local.get 0
                else
                  local.get 2
                  i32.load offset=100
                  local.set 3
                  local.get 2
                  i32.load offset=96
                end
              end
              local.set 0
              global.get 1
              i32.eqz
              if  ;; label = @6
                local.get 2
                i32.const 184
                i32.add
                local.tee 1
                local.get 0
                i32.store
              end
              local.get 4
              i32.const 15
              i32.eq
              i32.const 1
              global.get 1
              select
              if  ;; label = @6
                local.get 0
                local.get 3
                call $runtime.printstring
                i32.const 15
                global.get 1
                i32.const 1
                i32.eq
                br_if 5 (;@1;)
                drop
              end
              global.get 1
              i32.eqz
              br_if 3 (;@2;)
            end
            global.get 1
            i32.eqz
            if  ;; label = @5
              loop  ;; label = @6
                local.get 0
                i32.const 0
                i32.ge_s
                if  ;; label = @7
                  local.get 2
                  i32.const 120
                  i32.add
                  local.get 0
                  i32.add
                  local.get 7
                  local.get 7
                  i64.const 10
                  i64.div_u
                  local.tee 7
                  i64.const 10
                  i64.mul
                  i64.sub
                  i32.wrap_i64
                  i32.const 48
                  i32.or
                  local.tee 5
                  i32.store8
                  local.get 3
                  local.get 0
                  local.get 5
                  i32.const 255
                  i32.and
                  i32.const 48
                  i32.eq
                  select
                  local.set 3
                  local.get 0
                  i32.const 1
                  i32.sub
                  local.set 0
                  br 1 (;@6;)
                end
              end
              local.get 3
              i32.const 20
              local.get 3
              i32.const 20
              i32.gt_s
              select
              local.get 3
              i32.sub
              local.set 0
              local.get 2
              i32.const 120
              i32.add
              local.tee 5
              local.get 3
              i32.add
              local.set 3
            end
            loop  ;; label = @5
              block  ;; label = @6
                global.get 1
                i32.eqz
                if  ;; label = @7
                  local.get 0
                  i32.eqz
                  br_if 1 (;@6;)
                  local.get 3
                  i32.load8_u
                  local.set 5
                end
                local.get 4
                i32.const 16
                i32.eq
                i32.const 1
                global.get 1
                select
                if  ;; label = @7
                  local.get 5
                  call $runtime.putchar
                  i32.const 16
                  global.get 1
                  i32.const 1
                  i32.eq
                  br_if 6 (;@1;)
                  drop
                end
                global.get 1
                i32.eqz
                if  ;; label = @7
                  local.get 0
                  i32.const 1
                  i32.sub
                  local.set 0
                  local.get 3
                  i32.const 1
                  i32.add
                  local.set 3
                  br 2 (;@5;)
                end
              end
            end
            local.get 4
            i32.const 17
            i32.eq
            i32.const 1
            global.get 1
            select
            if  ;; label = @5
              i32.const 58
              call $runtime.putchar
              i32.const 17
              global.get 1
              i32.const 1
              i32.eq
              br_if 4 (;@1;)
              drop
            end
            global.get 1
            i32.const 1
            local.get 1
            select
            if  ;; label = @5
              local.get 4
              i32.const 18
              i32.eq
              i32.const 1
              global.get 1
              select
              if  ;; label = @6
                i32.const 65651
                i32.const 3
                call $runtime.printstring
                i32.const 18
                global.get 1
                i32.const 1
                i32.eq
                br_if 5 (;@1;)
                drop
              end
              global.get 1
              i32.eqz
              br_if 2 (;@3;)
            end
            local.get 4
            i32.const 19
            i32.eq
            i32.const 1
            global.get 1
            select
            if  ;; label = @5
              i32.const 48
              call $runtime.putchar
              i32.const 19
              global.get 1
              i32.const 1
              i32.eq
              br_if 4 (;@1;)
              drop
            end
            local.get 4
            i32.const 20
            i32.eq
            i32.const 1
            global.get 1
            select
            if  ;; label = @5
              i32.const 120
              call $runtime.putchar
              i32.const 20
              global.get 1
              i32.const 1
              i32.eq
              br_if 4 (;@1;)
              drop
            end
            local.get 0
            i32.const 8
            global.get 1
            select
            local.set 0
            loop  ;; label = @5
              global.get 1
              i32.eqz
              if  ;; label = @6
                local.get 0
                i32.eqz
                br_if 3 (;@3;)
                local.get 1
                i32.const 28
                i32.shr_u
                local.tee 3
                i32.const 48
                i32.or
                local.get 3
                i32.const 87
                i32.add
                local.get 3
                i32.const 10
                i32.lt_u
                select
                local.set 3
              end
              local.get 4
              i32.const 21
              i32.eq
              i32.const 1
              global.get 1
              select
              if  ;; label = @6
                local.get 3
                call $runtime.putchar
                i32.const 21
                global.get 1
                i32.const 1
                i32.eq
                br_if 5 (;@1;)
                drop
              end
              global.get 1
              i32.eqz
              if  ;; label = @6
                local.get 0
                i32.const 1
                i32.sub
                local.set 0
                local.get 1
                i32.const 4
                i32.shl
                local.set 1
                br 1 (;@5;)
              end
            end
          end
          local.get 4
          i32.const 22
          i32.eq
          i32.const 1
          global.get 1
          select
          if  ;; label = @4
            call $runtime.nilPanic
            i32.const 22
            global.get 1
            i32.const 1
            i32.eq
            br_if 3 (;@1;)
            drop
          end
          global.get 1
          i32.eqz
          if  ;; label = @4
            unreachable
          end
        end
        local.get 4
        i32.const 23
        i32.eq
        i32.const 1
        global.get 1
        select
        if  ;; label = @3
          i32.const 41
          call $runtime.putchar
          i32.const 23
          global.get 1
          i32.const 1
          i32.eq
          br_if 2 (;@1;)
          drop
        end
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        i32.const 66380
        local.get 8
        i32.store
        local.get 2
        i32.const 208
        i32.add
        global.set $__stack_pointer
      end
      return
    end
    local.set 4
    global.get 2
    i32.load
    local.get 4
    i32.store
    global.get 2
    global.get 2
    i32.load
    i32.const 4
    i32.add
    i32.store
    global.get 2
    i32.load
    local.tee 4
    local.get 0
    i32.store
    local.get 4
    local.get 1
    i32.store offset=4
    local.get 4
    local.get 2
    i32.store offset=8
    local.get 4
    local.get 3
    i32.store offset=12
    local.get 4
    local.get 5
    i32.store offset=16
    local.get 4
    local.get 7
    i64.store offset=20 align=4
    local.get 4
    local.get 8
    i32.store offset=28
    local.get 4
    local.get 6
    i32.store offset=32
    global.get 2
    global.get 2
    i32.load
    i32.const 36
    i32.add
    i32.store)
  (func $runtime.stringConcat (type 16) (param i32 i32 i32 i32 i32)
    (local i32 i32 i32 i32 i32 i32)
    global.get 1
    i32.const 2
    i32.eq
    if  ;; label = @1
      global.get 2
      global.get 2
      i32.load
      i32.const 40
      i32.sub
      i32.store
      global.get 2
      i32.load
      local.tee 7
      i32.load
      local.set 0
      local.get 7
      i32.load offset=4
      local.set 1
      local.get 7
      i32.load offset=8
      local.set 2
      local.get 7
      i32.load offset=12
      local.set 3
      local.get 7
      i32.load offset=16
      local.set 4
      local.get 7
      i32.load offset=20
      local.set 5
      local.get 7
      i32.load offset=24
      local.set 8
      local.get 7
      i32.load offset=28
      local.set 9
      local.get 7
      i32.load offset=32
      local.set 10
      local.get 7
      i32.load offset=36
      local.set 7
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
        i32.const -64
        i32.add
        local.tee 5
        global.set $__stack_pointer
        local.get 5
        i64.const 0
        i64.store offset=56
        local.get 5
        i64.const 0
        i64.store offset=48
        local.get 5
        i32.const 6
        i32.store offset=36
        local.get 5
        local.get 1
        i32.store offset=8
        local.get 5
        local.get 2
        i32.store offset=12
        local.get 5
        i32.const 40
        i32.add
        local.get 5
        i32.const 8
        i32.add
        i32.store
        local.get 5
        local.get 3
        i32.store offset=16
        local.get 5
        local.get 4
        i32.store offset=20
        local.get 5
        i32.const 44
        i32.add
        local.get 5
        i32.const 16
        i32.add
        local.tee 9
        i32.store
        i32.const 66380
        i32.load
        local.set 10
        i32.const 66380
        local.get 5
        i32.const 32
        i32.add
        local.tee 8
        i32.store
        local.get 5
        local.get 10
        i32.store offset=32
      end
      block  ;; label = @2
        global.get 1
        i32.eqz
        if  ;; label = @3
          local.get 2
          i32.eqz
          if  ;; label = @4
            local.get 3
            local.set 8
            local.get 4
            local.set 9
            br 2 (;@2;)
          end
          local.get 4
          i32.eqz
          if  ;; label = @4
            local.get 1
            local.set 8
            local.get 2
            local.set 9
            br 2 (;@2;)
          end
          local.get 5
          i32.const 56
          i32.add
          local.set 7
          local.get 2
          local.get 4
          i32.add
          local.set 9
        end
        local.get 6
        i32.const 0
        global.get 1
        select
        i32.eqz
        if  ;; label = @3
          local.get 9
          call $runtime.alloc
          local.set 6
          i32.const 0
          global.get 1
          i32.const 1
          i32.eq
          br_if 2 (;@1;)
          drop
          local.get 6
          local.set 8
        end
        global.get 1
        i32.eqz
        if  ;; label = @3
          local.get 7
          local.get 8
          i32.store
          local.get 5
          i32.const 60
          i32.add
          local.get 8
          i32.store
          local.get 5
          i32.const 48
          i32.add
          local.get 8
          i32.store
          local.get 8
          local.get 1
          local.get 2
          memory.copy
          local.get 2
          local.get 8
          i32.add
          local.get 3
          local.get 4
          memory.copy
          local.get 5
          local.get 9
          i32.store offset=28
          local.get 5
          i32.const 52
          i32.add
          local.get 5
          i32.const 24
          i32.add
          i32.store
          local.get 5
          local.get 8
          i32.store offset=24
        end
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        i32.const 66380
        local.get 10
        i32.store
        local.get 0
        local.get 9
        i32.store offset=4
        local.get 0
        local.get 8
        i32.store
        local.get 5
        i32.const -64
        i32.sub
        global.set $__stack_pointer
      end
      return
    end
    local.set 6
    global.get 2
    i32.load
    local.get 6
    i32.store
    global.get 2
    global.get 2
    i32.load
    i32.const 4
    i32.add
    i32.store
    global.get 2
    i32.load
    local.tee 6
    local.get 0
    i32.store
    local.get 6
    local.get 1
    i32.store offset=4
    local.get 6
    local.get 2
    i32.store offset=8
    local.get 6
    local.get 3
    i32.store offset=12
    local.get 6
    local.get 4
    i32.store offset=16
    local.get 6
    local.get 5
    i32.store offset=20
    local.get 6
    local.get 8
    i32.store offset=24
    local.get 6
    local.get 9
    i32.store offset=28
    local.get 6
    local.get 10
    i32.store offset=32
    local.get 6
    local.get 7
    i32.store offset=36
    global.get 2
    global.get 2
    i32.load
    i32.const 40
    i32.add
    i32.store)
  (func $_syscall/js.Type_.String (type 2) (param i32 i32)
    (local i32 i32 i32)
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
        local.set 4
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        i32.const 9
        local.set 2
        i32.const 66064
        local.set 3
        block  ;; label = @3
          block  ;; label = @4
            block  ;; label = @5
              block  ;; label = @6
                block  ;; label = @7
                  block  ;; label = @8
                    block  ;; label = @9
                      block  ;; label = @10
                        block  ;; label = @11
                          local.get 1
                          br_table 7 (;@4;) 0 (;@11;) 1 (;@10;) 2 (;@9;) 3 (;@8;) 4 (;@7;) 5 (;@6;) 6 (;@5;) 8 (;@3;)
                        end
                        i32.const 4
                        local.set 2
                        i32.const 66073
                        local.set 3
                        br 6 (;@4;)
                      end
                      i32.const 7
                      local.set 2
                      i32.const 66077
                      local.set 3
                      br 5 (;@4;)
                    end
                    i32.const 6
                    local.set 2
                    i32.const 66084
                    local.set 3
                    br 4 (;@4;)
                  end
                  i32.const 6
                  local.set 2
                  i32.const 66090
                  local.set 3
                  br 3 (;@4;)
                end
                i32.const 6
                local.set 2
                i32.const 66096
                local.set 3
                br 2 (;@4;)
              end
              i32.const 6
              local.set 2
              i32.const 66102
              local.set 3
              br 1 (;@4;)
            end
            i32.const 8
            local.set 2
            i32.const 66108
            local.set 3
          end
          local.get 0
          local.get 3
          i32.store
          local.get 0
          local.get 2
          i32.store offset=4
          return
        end
      end
      local.get 4
      i32.const 0
      global.get 1
      select
      i32.eqz
      if  ;; label = @2
        i32.const 34
        i32.const 66128
        call $runtime._panic
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
  (func $_syscall/js.Error_.Error (type 5) (param i32 i64 i32)
    (local i32 i32 i32 i32 i32 i32)
    global.get 1
    i32.const 2
    i32.eq
    if  ;; label = @1
      global.get 2
      global.get 2
      i32.load
      i32.const 36
      i32.sub
      i32.store
      global.get 2
      i32.load
      local.tee 5
      i32.load
      local.set 0
      local.get 5
      i64.load offset=4 align=4
      local.set 1
      local.get 5
      i32.load offset=12
      local.set 2
      local.get 5
      i32.load offset=16
      local.set 3
      local.get 5
      i32.load offset=20
      local.set 6
      local.get 5
      i32.load offset=24
      local.set 7
      local.get 5
      i32.load offset=28
      local.set 8
      local.get 5
      i32.load offset=32
      local.set 5
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
        local.set 4
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        global.get $__stack_pointer
        i32.const 80
        i32.sub
        local.tee 3
        global.set $__stack_pointer
        local.get 3
        i32.const 72
        i32.add
        local.tee 5
        i32.const 0
        i32.store
        local.get 3
        i32.const -64
        i32.sub
        local.tee 6
        i64.const 0
        i64.store
        local.get 3
        i32.const 60
        i32.add
        local.get 2
        i32.store
        local.get 3
        i32.const 5
        i32.store offset=52
        local.get 3
        i32.const 0
        i32.store offset=40
        local.get 3
        i64.const 0
        i64.store offset=32
        local.get 3
        i32.const 56
        i32.add
        local.get 3
        i32.const 32
        i32.add
        i32.store
        i32.const 66380
        i32.load
        local.set 7
        i32.const 66380
        local.get 3
        i32.const 48
        i32.add
        i32.store
        local.get 3
        local.get 7
        i32.store offset=48
        local.get 3
        i32.const 16
        i32.add
        local.set 8
      end
      local.get 4
      i32.const 0
      global.get 1
      select
      i32.eqz
      if  ;; label = @2
        local.get 8
        local.get 1
        local.get 2
        i32.const 66008
        i32.const 7
        call $_syscall/js.Value_.Get
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
        local.get 6
        local.get 3
        i32.load offset=24
        local.tee 2
        i32.store
        local.get 3
        i32.const 8
        i32.add
        local.set 6
        local.get 3
        i64.load offset=16
        local.set 1
      end
      local.get 4
      i32.const 1
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        local.get 6
        local.get 1
        local.get 2
        call $_syscall/js.Value_.String
        i32.const 1
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        local.get 3
        i32.const 68
        i32.add
        local.get 3
        i32.load offset=8
        local.tee 2
        i32.store
        local.get 3
        i32.load offset=12
        local.set 6
      end
      local.get 4
      i32.const 2
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        local.get 3
        i32.const 66015
        i32.const 18
        local.get 2
        local.get 6
        call $runtime.stringConcat
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
        local.get 7
        i32.store
        local.get 5
        local.get 3
        i32.load
        local.tee 2
        i32.store
        local.get 3
        i32.load offset=4
        local.set 7
        local.get 0
        local.get 2
        i32.store
        local.get 0
        local.get 7
        i32.store offset=4
        local.get 3
        i32.const 80
        i32.add
        global.set $__stack_pointer
      end
      return
    end
    local.set 4
    global.get 2
    i32.load
    local.get 4
    i32.store
    global.get 2
    global.get 2
    i32.load
    i32.const 4
    i32.add
    i32.store
    global.get 2
    i32.load
    local.tee 4
    local.get 0
    i32.store
    local.get 4
    local.get 1
    i64.store offset=4 align=4
    local.get 4
    local.get 2
    i32.store offset=12
    local.get 4
    local.get 3
    i32.store offset=16
    local.get 4
    local.get 6
    i32.store offset=20
    local.get 4
    local.get 7
    i32.store offset=24
    local.get 4
    local.get 8
    i32.store offset=28
    local.get 4
    local.get 5
    i32.store offset=32
    global.get 2
    global.get 2
    i32.load
    i32.const 36
    i32.add
    i32.store)
  (func $_struct_syscall/js.Value__.String$invoke (type 2) (param i32 i32)
    (local i32 i32 i32 i64 i32)
    global.get 1
    i32.const 2
    i32.eq
    if  ;; label = @1
      global.get 2
      global.get 2
      i32.load
      i32.const 24
      i32.sub
      i32.store
      global.get 2
      i32.load
      local.tee 2
      i32.load
      local.set 0
      local.get 2
      i32.load offset=8
      local.set 3
      local.get 2
      i32.load offset=12
      local.set 4
      local.get 2
      i64.load offset=16 align=4
      local.set 5
      local.get 2
      i32.load offset=4
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
        i32.const 48
        i32.sub
        local.tee 2
        global.set $__stack_pointer
        local.get 2
        i32.const 36
        i32.add
        local.tee 3
        i64.const 0
        i64.store align=4
        local.get 3
        local.get 1
        i32.load offset=8
        local.tee 4
        i32.store
        local.get 2
        i32.const 3
        i32.store offset=28
        local.get 2
        i32.const 32
        i32.add
        local.get 2
        i32.const 8
        i32.add
        i32.store
        i32.const 66380
        i32.load
        local.set 3
        i32.const 66380
        local.get 2
        i32.const 24
        i32.add
        i32.store
        local.get 2
        local.get 3
        i32.store offset=24
        local.get 1
        i64.load
        local.set 5
        local.get 2
        i32.const 0
        i32.store offset=16
        local.get 2
        i64.const 0
        i64.store offset=8
      end
      local.get 6
      i32.const 0
      global.get 1
      select
      i32.eqz
      if  ;; label = @2
        local.get 2
        local.get 5
        local.get 4
        call $_syscall/js.Value_.String
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
        local.get 2
        i32.const 40
        i32.add
        local.get 2
        i32.load
        local.tee 1
        i32.store
        i32.const 66380
        local.get 3
        i32.store
        local.get 0
        local.get 2
        i32.load offset=4
        i32.store offset=4
        local.get 0
        local.get 1
        i32.store
        local.get 2
        i32.const 48
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
    local.get 2
    i32.store offset=4
    local.get 1
    local.get 3
    i32.store offset=8
    local.get 1
    local.get 4
    i32.store offset=12
    local.get 1
    local.get 5
    i64.store offset=16 align=4
    global.get 2
    global.get 2
    i32.load
    i32.const 24
    i32.add
    i32.store)
  (func $_*struct_syscall/js.Value__.String (type 2) (param i32 i32)
    (local i32 i32 i32 i32 i64 i32)
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
      local.tee 2
      i32.load
      local.set 0
      local.get 2
      i32.load offset=8
      local.set 3
      local.get 2
      i32.load offset=12
      local.set 4
      local.get 2
      i32.load offset=16
      local.set 5
      local.get 2
      i64.load offset=20 align=4
      local.set 6
      local.get 2
      i32.load offset=4
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
        local.set 7
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        global.get $__stack_pointer
        i32.const 32
        i32.sub
        local.tee 2
        global.set $__stack_pointer
        local.get 2
        i32.const 0
        i32.store offset=28
        local.get 2
        i64.const 2
        i64.store offset=20 align=4
        i32.const 66380
        i32.load
        local.set 3
        i32.const 66380
        local.get 2
        i32.const 16
        i32.add
        i32.store
        local.get 2
        local.get 3
        i32.store offset=16
        local.get 1
        i32.eqz
        local.set 4
      end
      block  ;; label = @2
        global.get 1
        i32.eqz
        if  ;; label = @3
          local.get 4
          br_if 1 (;@2;)
          local.get 2
          i32.const 24
          i32.add
          local.get 1
          i32.load offset=8
          local.tee 4
          i32.store
          local.get 1
          i64.load
          local.set 6
          local.get 2
          i32.const 8
          i32.add
          local.set 5
        end
        local.get 7
        i32.const 0
        global.get 1
        select
        i32.eqz
        if  ;; label = @3
          local.get 5
          local.get 6
          local.get 4
          call $_syscall/js.Value_.String
          i32.const 0
          global.get 1
          i32.const 1
          i32.eq
          br_if 2 (;@1;)
          drop
        end
        global.get 1
        i32.eqz
        if  ;; label = @3
          i32.const 66380
          local.get 3
          i32.store
          local.get 2
          i32.const 28
          i32.add
          local.get 2
          i32.load offset=8
          local.tee 1
          i32.store
          local.get 2
          i32.load offset=12
          local.set 3
          local.get 0
          local.get 1
          i32.store
          local.get 0
          local.get 3
          i32.store offset=4
          local.get 2
          i32.const 32
          i32.add
          global.set $__stack_pointer
          return
        end
      end
      local.get 7
      i32.const 1
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        call $runtime.nilPanic
        i32.const 1
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
    local.get 1
    local.get 4
    i32.store offset=12
    local.get 1
    local.get 5
    i32.store offset=16
    local.get 1
    local.get 6
    i64.store offset=20 align=4
    global.get 2
    global.get 2
    i32.load
    i32.const 28
    i32.add
    i32.store)
  (func $_*syscall/js.Value_.String (type 2) (param i32 i32)
    (local i32 i32 i32 i64 i32)
    global.get 1
    i32.const 2
    i32.eq
    if  ;; label = @1
      global.get 2
      global.get 2
      i32.load
      i32.const 24
      i32.sub
      i32.store
      global.get 2
      i32.load
      local.tee 2
      i32.load
      local.set 0
      local.get 2
      i32.load offset=8
      local.set 3
      local.get 2
      i32.load offset=12
      local.set 4
      local.get 2
      i64.load offset=16 align=4
      local.set 5
      local.get 2
      i32.load offset=4
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
        local.tee 2
        global.set $__stack_pointer
        local.get 2
        i32.const 16
        i32.add
        local.get 1
        i32.store
        local.get 2
        i64.const 0
        i64.store offset=20 align=4
        local.get 2
        i32.const 3
        i32.store offset=12
        i32.const 66380
        i32.load
        local.set 3
        i32.const 66380
        local.get 2
        i32.const 8
        i32.add
        i32.store
        local.get 2
        local.get 3
        i32.store offset=8
        local.get 1
        i32.eqz
        local.set 4
      end
      block  ;; label = @2
        global.get 1
        i32.eqz
        if  ;; label = @3
          local.get 4
          br_if 1 (;@2;)
          local.get 2
          i32.const 20
          i32.add
          local.get 1
          i32.load offset=8
          local.tee 4
          i32.store
          local.get 1
          i64.load
          local.set 5
        end
        local.get 6
        i32.const 0
        global.get 1
        select
        i32.eqz
        if  ;; label = @3
          local.get 2
          local.get 5
          local.get 4
          call $_syscall/js.Value_.String
          i32.const 0
          global.get 1
          i32.const 1
          i32.eq
          br_if 2 (;@1;)
          drop
        end
        global.get 1
        i32.eqz
        if  ;; label = @3
          i32.const 66380
          local.get 3
          i32.store
          local.get 2
          i32.const 24
          i32.add
          local.get 2
          i32.load
          local.tee 1
          i32.store
          local.get 2
          i32.load offset=4
          local.set 3
          local.get 0
          local.get 1
          i32.store
          local.get 0
          local.get 3
          i32.store offset=4
          local.get 2
          i32.const 32
          i32.add
          global.set $__stack_pointer
          return
        end
      end
      local.get 6
      i32.const 1
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        call $runtime.nilPanic
        i32.const 1
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
    local.get 1
    local.get 4
    i32.store offset=12
    local.get 1
    local.get 5
    i64.store offset=16 align=4
    global.get 2
    global.get 2
    i32.load
    i32.const 24
    i32.add
    i32.store)
  (func $_*syscall/js.Type_.String (type 2) (param i32 i32)
    (local i32 i32 i32 i32)
    global.get 1
    i32.const 2
    i32.eq
    if  ;; label = @1
      global.get 2
      global.get 2
      i32.load
      i32.const 20
      i32.sub
      i32.store
      global.get 2
      i32.load
      local.tee 2
      i32.load
      local.set 0
      local.get 2
      i32.load offset=4
      local.set 1
      local.get 2
      i32.load offset=8
      local.set 3
      local.get 2
      i32.load offset=12
      local.set 5
      local.get 2
      i32.load offset=16
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
        local.set 4
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
        i32.const 24
        i32.add
        local.get 1
        i32.store
        local.get 3
        i32.const 0
        i32.store offset=28
        local.get 3
        i32.const 2
        i32.store offset=20
        i32.const 66380
        i32.load
        local.set 5
        i32.const 66380
        local.get 3
        i32.const 16
        i32.add
        i32.store
        local.get 3
        local.get 5
        i32.store offset=16
        local.get 1
        i32.eqz
        local.set 2
      end
      block  ;; label = @2
        global.get 1
        i32.eqz
        if  ;; label = @3
          local.get 2
          br_if 1 (;@2;)
          local.get 3
          i32.const 8
          i32.add
          local.set 2
          local.get 1
          i32.load
          local.set 1
        end
        local.get 4
        i32.const 0
        global.get 1
        select
        i32.eqz
        if  ;; label = @3
          local.get 2
          local.get 1
          call $_syscall/js.Type_.String
          i32.const 0
          global.get 1
          i32.const 1
          i32.eq
          br_if 2 (;@1;)
          drop
        end
        global.get 1
        i32.eqz
        if  ;; label = @3
          i32.const 66380
          local.get 5
          i32.store
          local.get 3
          i32.const 28
          i32.add
          local.get 3
          i32.load offset=8
          local.tee 1
          i32.store
          local.get 3
          i32.load offset=12
          local.set 5
          local.get 0
          local.get 1
          i32.store
          local.get 0
          local.get 5
          i32.store offset=4
          local.get 3
          i32.const 32
          i32.add
          global.set $__stack_pointer
          return
        end
      end
      local.get 4
      i32.const 1
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        call $runtime.nilPanic
        i32.const 1
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
    local.set 4
    global.get 2
    i32.load
    local.get 4
    i32.store
    global.get 2
    global.get 2
    i32.load
    i32.const 4
    i32.add
    i32.store
    global.get 2
    i32.load
    local.tee 4
    local.get 0
    i32.store
    local.get 4
    local.get 1
    i32.store offset=4
    local.get 4
    local.get 3
    i32.store offset=8
    local.get 4
    local.get 5
    i32.store offset=12
    local.get 4
    local.get 2
    i32.store offset=16
    global.get 2
    global.get 2
    i32.load
    i32.const 20
    i32.add
    i32.store)
  (func $_syscall/js.Value_.String (type 5) (param i32 i64 i32)
    (local i32 i32 i32 i32 i32 i32 i32)
    global.get 1
    i32.const 2
    i32.eq
    if  ;; label = @1
      global.get 2
      global.get 2
      i32.load
      i32.const 36
      i32.sub
      i32.store
      global.get 2
      i32.load
      local.tee 6
      i32.load
      local.set 0
      local.get 6
      i64.load offset=4 align=4
      local.set 1
      local.get 6
      i32.load offset=12
      local.set 2
      local.get 6
      i32.load offset=16
      local.set 3
      local.get 6
      i32.load offset=20
      local.set 4
      local.get 6
      i32.load offset=24
      local.set 5
      local.get 6
      i32.load offset=28
      local.set 9
      local.get 6
      i32.load offset=32
      local.set 6
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
        local.set 8
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        global.get $__stack_pointer
        i32.const 128
        i32.sub
        local.tee 3
        global.set $__stack_pointer
        local.get 3
        i32.const 12
        i32.store offset=76
        local.get 3
        i32.const 84
        i32.add
        local.tee 9
        i64.const 0
        i64.store align=4
        local.get 3
        i64.const 0
        i64.store offset=108 align=4
        local.get 3
        i32.const 112
        i32.add
        local.get 2
        i32.store
        local.get 3
        i64.const 0
        i64.store offset=92 align=4
        local.get 3
        i32.const 96
        i32.add
        local.get 2
        i32.store
        local.get 3
        i32.const 88
        i32.add
        local.get 2
        i32.store
        local.get 9
        local.get 2
        i32.store
        local.get 3
        i32.const 0
        i32.store offset=124
        local.get 3
        i64.const 0
        i64.store offset=116 align=4
        local.get 3
        i64.const 0
        i64.store offset=100 align=4
        local.get 3
        local.get 1
        i64.store offset=56
        local.get 3
        local.get 2
        i32.store offset=64
        local.get 3
        i32.const 80
        i32.add
        local.get 3
        i32.const 56
        i32.add
        i32.store
        i32.const 66380
        i32.load
        local.set 9
        i32.const 66380
        local.get 3
        i32.const 72
        i32.add
        i32.store
        local.get 3
        local.get 9
        i32.store offset=72
        i32.const 11
        local.set 5
        i32.const 65924
        local.set 4
      end
      local.get 8
      i32.const 0
      global.get 1
      select
      i32.eqz
      if  ;; label = @2
        local.get 1
        local.get 2
        call $_syscall/js.Value_.Type
        local.set 7
        i32.const 0
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
        local.get 7
        local.set 6
      end
      block  ;; label = @2
        block  ;; label = @3
          block  ;; label = @4
            block  ;; label = @5
              block  ;; label = @6
                block  ;; label = @7
                  block  ;; label = @8
                    block  ;; label = @9
                      global.get 1
                      i32.eqz
                      if  ;; label = @10
                        block  ;; label = @11
                          local.get 6
                          br_table 8 (;@3;) 2 (;@9;) 3 (;@8;) 4 (;@7;) 0 (;@11;) 5 (;@6;) 6 (;@5;) 7 (;@4;) 9 (;@2;)
                        end
                      end
                      local.get 8
                      i32.const 1
                      i32.eq
                      i32.const 1
                      global.get 1
                      select
                      if  ;; label = @10
                        local.get 3
                        local.get 1
                        local.get 2
                        call $syscall/js.jsString
                        i32.const 1
                        global.get 1
                        i32.const 1
                        i32.eq
                        br_if 9 (;@1;)
                        drop
                      end
                      global.get 1
                      i32.eqz
                      if  ;; label = @10
                        local.get 3
                        i32.const 92
                        i32.add
                        local.get 3
                        i32.load
                        local.tee 4
                        i32.store
                        local.get 3
                        i32.load offset=4
                        local.set 5
                        br 7 (;@3;)
                      end
                    end
                    global.get 1
                    i32.eqz
                    if  ;; label = @9
                      i32.const 6
                      local.set 5
                      i32.const 65935
                      local.set 4
                      br 6 (;@3;)
                    end
                  end
                  local.get 4
                  local.get 3
                  i32.const 24
                  i32.add
                  global.get 1
                  select
                  local.set 4
                  local.get 8
                  i32.const 2
                  i32.eq
                  i32.const 1
                  global.get 1
                  select
                  if  ;; label = @8
                    local.get 4
                    local.get 1
                    local.get 2
                    call $syscall/js.jsString
                    i32.const 2
                    global.get 1
                    i32.const 1
                    i32.eq
                    br_if 7 (;@1;)
                    drop
                  end
                  global.get 1
                  i32.eqz
                  if  ;; label = @8
                    local.get 3
                    i32.const 100
                    i32.add
                    local.get 3
                    i32.load offset=24
                    local.tee 2
                    i32.store
                    local.get 3
                    i32.load offset=28
                    local.set 5
                    local.get 3
                    i32.const 16
                    i32.add
                    local.set 4
                  end
                  local.get 8
                  i32.const 3
                  i32.eq
                  i32.const 1
                  global.get 1
                  select
                  if  ;; label = @8
                    local.get 4
                    i32.const 65941
                    i32.const 10
                    local.get 2
                    local.get 5
                    call $runtime.stringConcat
                    i32.const 3
                    global.get 1
                    i32.const 1
                    i32.eq
                    br_if 7 (;@1;)
                    drop
                  end
                  global.get 1
                  i32.eqz
                  if  ;; label = @8
                    local.get 3
                    i32.const 104
                    i32.add
                    local.get 3
                    i32.load offset=16
                    local.tee 2
                    i32.store
                    local.get 3
                    i32.load offset=20
                    local.set 5
                    local.get 3
                    i32.const 8
                    i32.add
                    local.set 4
                  end
                  local.get 8
                  i32.const 4
                  i32.eq
                  i32.const 1
                  global.get 1
                  select
                  if  ;; label = @8
                    local.get 4
                    local.get 2
                    local.get 5
                    i32.const 65960
                    i32.const 1
                    call $runtime.stringConcat
                    i32.const 4
                    global.get 1
                    i32.const 1
                    i32.eq
                    br_if 7 (;@1;)
                    drop
                  end
                  global.get 1
                  i32.eqz
                  if  ;; label = @8
                    local.get 3
                    i32.const 108
                    i32.add
                    local.get 3
                    i32.load offset=8
                    local.tee 4
                    i32.store
                    local.get 3
                    i32.load offset=12
                    local.set 5
                    br 5 (;@3;)
                  end
                end
                local.get 4
                local.get 3
                i32.const 48
                i32.add
                global.get 1
                select
                local.set 4
                local.get 8
                i32.const 5
                i32.eq
                i32.const 1
                global.get 1
                select
                if  ;; label = @7
                  local.get 4
                  local.get 1
                  local.get 2
                  call $syscall/js.jsString
                  i32.const 5
                  global.get 1
                  i32.const 1
                  i32.eq
                  br_if 6 (;@1;)
                  drop
                end
                global.get 1
                i32.eqz
                if  ;; label = @7
                  local.get 3
                  i32.const 116
                  i32.add
                  local.get 3
                  i32.load offset=48
                  local.tee 2
                  i32.store
                  local.get 3
                  i32.load offset=52
                  local.set 5
                  local.get 3
                  i32.const 40
                  i32.add
                  local.set 4
                end
                local.get 8
                i32.const 6
                i32.eq
                i32.const 1
                global.get 1
                select
                if  ;; label = @7
                  local.get 4
                  i32.const 65951
                  i32.const 9
                  local.get 2
                  local.get 5
                  call $runtime.stringConcat
                  i32.const 6
                  global.get 1
                  i32.const 1
                  i32.eq
                  br_if 6 (;@1;)
                  drop
                end
                global.get 1
                i32.eqz
                if  ;; label = @7
                  local.get 3
                  i32.const 120
                  i32.add
                  local.get 3
                  i32.load offset=40
                  local.tee 2
                  i32.store
                  local.get 3
                  i32.load offset=44
                  local.set 5
                  local.get 3
                  i32.const 32
                  i32.add
                  local.set 4
                end
                local.get 8
                i32.const 7
                i32.eq
                i32.const 1
                global.get 1
                select
                if  ;; label = @7
                  local.get 4
                  local.get 2
                  local.get 5
                  i32.const 65960
                  i32.const 1
                  call $runtime.stringConcat
                  i32.const 7
                  global.get 1
                  i32.const 1
                  i32.eq
                  br_if 6 (;@1;)
                  drop
                end
                global.get 1
                i32.eqz
                if  ;; label = @7
                  local.get 3
                  i32.const 124
                  i32.add
                  local.get 3
                  i32.load offset=32
                  local.tee 4
                  i32.store
                  local.get 3
                  i32.load offset=36
                  local.set 5
                  br 4 (;@3;)
                end
              end
              global.get 1
              i32.eqz
              if  ;; label = @6
                i32.const 65961
                local.set 4
                i32.const 8
                local.set 5
                br 3 (;@3;)
              end
            end
            global.get 1
            i32.eqz
            if  ;; label = @5
              i32.const 65969
              local.set 4
              i32.const 8
              local.set 5
              br 2 (;@3;)
            end
          end
          global.get 1
          i32.eqz
          if  ;; label = @4
            i32.const 10
            local.set 5
            i32.const 65977
            local.set 4
          end
        end
        global.get 1
        i32.eqz
        if  ;; label = @3
          i32.const 66380
          local.get 9
          i32.store
          local.get 0
          local.get 5
          i32.store offset=4
          local.get 0
          local.get 4
          i32.store
          local.get 3
          i32.const 128
          i32.add
          global.set $__stack_pointer
          return
        end
      end
      local.get 8
      i32.const 8
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        i32.const 34
        i32.const 66128
        call $runtime._panic
        i32.const 8
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
    local.set 7
    global.get 2
    i32.load
    local.get 7
    i32.store
    global.get 2
    global.get 2
    i32.load
    i32.const 4
    i32.add
    i32.store
    global.get 2
    i32.load
    local.tee 7
    local.get 0
    i32.store
    local.get 7
    local.get 1
    i64.store offset=4 align=4
    local.get 7
    local.get 2
    i32.store offset=12
    local.get 7
    local.get 3
    i32.store offset=16
    local.get 7
    local.get 4
    i32.store offset=20
    local.get 7
    local.get 5
    i32.store offset=24
    local.get 7
    local.get 9
    i32.store offset=28
    local.get 7
    local.get 6
    i32.store offset=32
    global.get 2
    global.get 2
    i32.load
    i32.const 36
    i32.add
    i32.store)
  (func $runtime.memequal (type 4) (param i32 i32 i32 i32) (result i32)
    (local i32)
    i32.const 0
    local.set 4
    loop  ;; label = @1
      local.get 2
      local.get 4
      local.tee 3
      i32.ne
      if  ;; label = @2
        local.get 3
        i32.const 1
        i32.add
        local.set 4
        local.get 0
        local.get 3
        i32.add
        i32.load8_u
        local.get 1
        local.get 3
        i32.add
        i32.load8_u
        i32.eq
        br_if 1 (;@1;)
      end
    end
    local.get 2
    local.get 3
    i32.le_u)
  (func $_start (type 0)
    (local i32 i32 i32 i32 i32)
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
      local.tee 0
      i32.load
      local.set 1
      local.get 0
      i32.load offset=4
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
        local.set 3
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        global.get $__stack_pointer
        i32.const 16
        i32.sub
        local.tee 1
        global.set $__stack_pointer
        local.get 1
        i32.const 12
        i32.add
        local.tee 4
        i32.const 0
        i32.store
        local.get 1
        i64.const 2
        i64.store offset=4 align=4
        i32.const 66380
        i32.load
        local.set 2
        i32.const 66380
        local.get 1
        i32.store
        local.get 1
        local.get 2
        i32.store
        memory.size
        local.set 0
        i32.const 66388
        i32.const 1
        i32.store8
        i32.const 66228
        local.get 0
        i32.const 16
        i32.shl
        i32.store
        call $runtime.calculateHeapAddresses
        local.get 1
        i32.const 8
        i32.add
        i32.const 66364
        i32.load
        local.tee 0
        i32.store
        local.get 4
        local.get 0
        i32.store
        local.get 0
        i32.const 0
        i32.const 66228
        i32.load
        local.get 0
        i32.sub
        memory.fill
      end
      local.get 3
      i32.const 0
      global.get 1
      select
      i32.eqz
      if  ;; label = @2
        i32.const 1
        call $internal/task.start
        i32.const 0
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
      end
      local.get 3
      i32.const 1
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        call $runtime.scheduler
        i32.const 1
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
        local.get 2
        i32.store
        i32.const 66388
        i32.const 0
        i32.store8
        local.get 1
        i32.const 16
        i32.add
        global.set $__stack_pointer
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
    i32.store
    global.get 2
    i32.load
    local.tee 0
    local.get 1
    i32.store
    local.get 0
    local.get 2
    i32.store offset=4
    global.get 2
    global.get 2
    i32.load
    i32.const 8
    i32.add
    i32.store)
  (func $runtime.run$1$gowrapper (type 1) (param i32)
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
      local.get 1
      i32.const 0
      global.get 1
      select
      i32.eqz
      if  ;; label = @2
        call $runtime.run$1
        i32.const 0
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
      end
      local.get 1
      i32.const 1
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        call $runtime.deadlock
        i32.const 1
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
  (func $runtime.scheduler (type 0)
    (local i32 i32)
    global.get 1
    i32.const 2
    i32.eq
    if  ;; label = @1
      global.get 2
      global.get 2
      i32.load
      i32.const 4
      i32.sub
      i32.store
      global.get 2
      i32.load
      i32.load
      local.set 0
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
      loop  ;; label = @2
        block  ;; label = @3
          global.get 1
          i32.eqz
          if  ;; label = @4
            i32.const 66389
            i32.load8_u
            br_if 1 (;@3;)
            call $_*internal/task.Queue_.Pop
            local.tee 0
            i32.eqz
            br_if 1 (;@3;)
          end
          local.get 1
          i32.const 0
          global.get 1
          select
          i32.eqz
          if  ;; label = @4
            local.get 0
            call $_*internal/task.Task_.Resume
            i32.const 0
            global.get 1
            i32.const 1
            i32.eq
            br_if 3 (;@1;)
            drop
          end
          global.get 1
          i32.eqz
          br_if 1 (;@2;)
        end
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
    local.get 0
    i32.store
    global.get 2
    global.get 2
    i32.load
    i32.const 4
    i32.add
    i32.store)
  (func $runtime.run$1 (type 0)
    (local i32 i32 i32 i32 i32)
    global.get 1
    i32.const 2
    i32.eq
    if  ;; label = @1
      global.get 2
      global.get 2
      i32.load
      i32.const 16
      i32.sub
      i32.store
      global.get 2
      i32.load
      local.tee 2
      i32.load
      local.set 0
      local.get 2
      i32.load offset=4
      local.set 3
      local.get 2
      i32.load offset=8
      local.set 4
      local.get 2
      i32.load offset=12
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
        local.set 1
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        global.get $__stack_pointer
        i32.const 48
        i32.sub
        local.tee 0
        global.set $__stack_pointer
        local.get 0
        i32.const 44
        i32.add
        local.tee 2
        i32.const 0
        i32.store
        local.get 0
        i64.const 2
        i64.store offset=36 align=4
        i32.const 66380
        i32.load
        local.set 4
        i32.const 66380
        local.get 0
        i32.const 32
        i32.add
        i32.store
        local.get 0
        local.get 4
        i32.store offset=32
        i32.const 66228
        memory.size
        i32.const 16
        i32.shl
        i32.store
        local.get 0
        i32.const 16
        i32.add
        local.set 3
      end
      local.get 1
      i32.const 0
      global.get 1
      select
      i32.eqz
      if  ;; label = @2
        local.get 3
        i64.const 9221120241336057861
        i32.const 0
        i32.const 65720
        i32.const 6
        call $_syscall/js.Value_.Get
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
        local.get 0
        i32.const 40
        i32.add
        local.get 0
        i32.load offset=24
        local.tee 3
        i32.store
        i32.const 66392
        local.get 0
        i64.load offset=16
        i64.store
        i32.const 66400
        local.get 3
        i32.store
      end
      local.get 1
      i32.const 1
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        local.get 0
        i64.const 9221120241336057861
        i32.const 0
        i32.const 65726
        i32.const 5
        call $_syscall/js.Value_.Get
        i32.const 1
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        local.get 2
        local.get 0
        i32.load offset=8
        local.tee 3
        i32.store
        i32.const 66408
        local.get 0
        i64.load
        i64.store
        i32.const 66416
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
        i32.const 66160
        i32.const 5
        call $github.com/iansmith/parigot/abi.OutputString
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
        i32.const 66389
        i32.const 1
        i32.store8
        local.get 0
        i32.const 48
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
    local.get 4
    i32.store offset=8
    local.get 1
    local.get 2
    i32.store offset=12
    global.get 2
    global.get 2
    i32.load
    i32.const 16
    i32.add
    i32.store)
  (func $_syscall/js.Value_.Get (type 17) (param i32 i64 i32 i32 i32)
    (local i32 i32 i32 i32 i32 i32 i64 i64)
    global.get 1
    i32.const 2
    i32.eq
    if  ;; label = @1
      global.get 2
      global.get 2
      i32.load
      i32.const 48
      i32.sub
      i32.store
      global.get 2
      i32.load
      local.tee 6
      i32.load
      local.set 0
      local.get 6
      i32.load offset=12
      local.set 2
      local.get 6
      i32.load offset=16
      local.set 3
      local.get 6
      i32.load offset=20
      local.set 4
      local.get 6
      i32.load offset=24
      local.set 5
      local.get 6
      i32.load offset=28
      local.set 7
      local.get 6
      i32.load offset=32
      local.set 8
      local.get 6
      i32.load offset=36
      local.set 10
      local.get 6
      i64.load offset=40 align=4
      local.set 11
      local.get 6
      i64.load offset=4 align=4
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
        local.set 9
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        global.get $__stack_pointer
        i32.const 112
        i32.sub
        local.tee 5
        global.set $__stack_pointer
        local.get 5
        i32.const 88
        i32.add
        local.get 2
        i32.store
        local.get 5
        i32.const 76
        i32.add
        local.get 2
        i32.store
        local.get 5
        i32.const 0
        i32.store offset=108
        local.get 5
        i64.const 0
        i64.store offset=100 align=4
        local.get 5
        i64.const 0
        i64.store offset=92 align=4
        local.get 5
        i64.const 0
        i64.store offset=80
        local.get 5
        i32.const 10
        i32.store offset=68
        local.get 5
        i32.const 0
        i32.store offset=24
        local.get 5
        i64.const 0
        i64.store offset=16
        local.get 5
        i32.const 72
        i32.add
        local.get 5
        i32.const 16
        i32.add
        i32.store
        i32.const 66380
        i32.load
        local.set 7
        i32.const 66380
        local.get 5
        i32.const -64
        i32.sub
        local.tee 8
        i32.store
        local.get 5
        local.get 7
        i32.store offset=64
      end
      local.get 9
      i32.const 0
      global.get 1
      select
      i32.eqz
      if  ;; label = @2
        local.get 1
        local.get 2
        call $_syscall/js.Value_.Type
        local.set 6
        i32.const 0
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
        local.get 6
        local.set 8
      end
      local.get 10
      local.get 8
      i32.const -2
      i32.and
      i32.const 6
      i32.ne
      global.get 1
      select
      local.set 10
      block  ;; label = @2
        global.get 1
        i32.eqz
        if  ;; label = @3
          local.get 10
          br_if 1 (;@2;)
          local.get 5
          i32.const 0
          i32.store offset=40
          local.get 5
          i64.const 0
          i64.store offset=32
          local.get 5
          i32.const 80
          i32.add
          local.tee 8
          local.get 5
          i32.const 32
          i32.add
          i32.store
        end
        local.get 9
        i32.const 1
        i32.eq
        i32.const 1
        global.get 1
        select
        if  ;; label = @3
          local.get 1
          local.get 3
          local.get 4
          local.get 5
          call $syscall/js.valueGet
          local.set 12
          i32.const 1
          global.get 1
          i32.const 1
          i32.eq
          br_if 2 (;@1;)
          drop
          local.get 12
          local.set 11
        end
        local.get 9
        i32.const 2
        i32.eq
        i32.const 1
        global.get 1
        select
        if  ;; label = @3
          local.get 5
          local.get 11
          call $syscall/js.makeValue
          i32.const 2
          global.get 1
          i32.const 1
          i32.eq
          br_if 2 (;@1;)
          drop
        end
        global.get 1
        i32.eqz
        if  ;; label = @3
          local.get 5
          i32.const 56
          i32.add
          local.tee 3
          i64.const 0
          i64.store
          i32.const 66380
          local.get 7
          i32.store
          local.get 5
          i32.const 100
          i32.add
          local.get 5
          i32.load offset=8
          local.tee 7
          i32.store
          local.get 5
          i32.const 84
          i32.add
          local.get 7
          i32.store
          local.get 3
          local.get 2
          i32.store
          local.get 5
          i32.const 96
          i32.add
          local.get 5
          i32.const 48
          i32.add
          i32.store
          local.get 5
          i32.const 92
          i32.add
          local.get 5
          i32.const 48
          i32.add
          i32.store
          local.get 5
          local.get 1
          i64.store offset=48
          local.get 5
          i64.load
          local.set 1
          local.get 0
          local.get 7
          i32.store offset=8
          local.get 0
          local.get 1
          i64.store
          local.get 5
          i32.const 112
          i32.add
          global.set $__stack_pointer
          return
        end
      end
      local.get 0
      local.get 5
      i32.const 104
      i32.add
      global.get 1
      select
      local.set 0
      local.get 9
      i32.const 3
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        i32.const 12
        call $runtime.alloc
        local.set 6
        i32.const 3
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
        local.get 6
        local.set 2
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        local.get 0
        local.get 2
        i32.store
        local.get 5
        i32.const 108
        i32.add
        local.get 2
        i32.store
        local.get 2
        local.get 8
        i32.store offset=8
        local.get 2
        i32.const 9
        i32.store offset=4
        local.get 2
        i32.const 65872
        i32.store
      end
      local.get 9
      i32.const 4
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        i32.const 3045
        local.get 2
        call $runtime._panic
        i32.const 4
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
    local.set 6
    global.get 2
    i32.load
    local.get 6
    i32.store
    global.get 2
    global.get 2
    i32.load
    i32.const 4
    i32.add
    i32.store
    global.get 2
    i32.load
    local.tee 6
    local.get 0
    i32.store
    local.get 6
    local.get 1
    i64.store offset=4 align=4
    local.get 6
    local.get 2
    i32.store offset=12
    local.get 6
    local.get 3
    i32.store offset=16
    local.get 6
    local.get 4
    i32.store offset=20
    local.get 6
    local.get 5
    i32.store offset=24
    local.get 6
    local.get 7
    i32.store offset=28
    local.get 6
    local.get 8
    i32.store offset=32
    local.get 6
    local.get 10
    i32.store offset=36
    local.get 6
    local.get 11
    i64.store offset=40 align=4
    global.get 2
    global.get 2
    i32.load
    i32.const 48
    i32.add
    i32.store)
  (func $resume (type 0)
    (local i32 i32)
    global.get 1
    i32.const 2
    i32.eq
    if  ;; label = @1
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
        local.set 0
      end
      local.get 0
      i32.const 0
      global.get 1
      select
      i32.eqz
      if  ;; label = @2
        i32.const 2
        call $internal/task.start
        i32.const 0
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
      end
      global.get 1
      i32.const 1
      global.get 1
      if (result i32)  ;; label = @2
        local.get 1
      else
        i32.const 66388
        i32.load8_u
        i32.eqz
      end
      select
      if  ;; label = @2
        local.get 0
        i32.const 1
        i32.eq
        i32.const 1
        global.get 1
        select
        if  ;; label = @3
          call $runtime.minSched
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
          return
        end
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        i32.const 66388
        i32.const 1
        i32.store8
      end
      local.get 0
      i32.const 2
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        call $runtime.scheduler
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
        i32.const 66388
        i32.const 0
        i32.store8
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
    i32.store
    global.get 2
    i32.load
    local.get 1
    i32.store
    global.get 2
    global.get 2
    i32.load
    i32.const 4
    i32.add
    i32.store)
  (func $runtime.resume$1$gowrapper (type 1) (param i32)
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
      local.get 1
      i32.const 0
      global.get 1
      select
      i32.eqz
      if  ;; label = @2
        call $runtime.resume$1
        i32.const 0
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
      end
      local.get 1
      i32.const 1
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        call $runtime.deadlock
        i32.const 1
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
  (func $runtime.minSched (type 0)
    (local i32 i32)
    global.get 1
    i32.const 2
    i32.eq
    if  ;; label = @1
      global.get 2
      global.get 2
      i32.load
      i32.const 4
      i32.sub
      i32.store
      global.get 2
      i32.load
      i32.load
      local.set 0
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
      loop  ;; label = @2
        block  ;; label = @3
          global.get 1
          i32.eqz
          if  ;; label = @4
            i32.const 66389
            i32.load8_u
            br_if 1 (;@3;)
            call $_*internal/task.Queue_.Pop
            local.tee 0
            i32.eqz
            br_if 1 (;@3;)
          end
          local.get 1
          i32.const 0
          global.get 1
          select
          i32.eqz
          if  ;; label = @4
            local.get 0
            call $_*internal/task.Task_.Resume
            i32.const 0
            global.get 1
            i32.const 1
            i32.eq
            br_if 3 (;@1;)
            drop
          end
          global.get 1
          i32.eqz
          br_if 1 (;@2;)
        end
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
    local.get 0
    i32.store
    global.get 2
    global.get 2
    i32.load
    i32.const 4
    i32.add
    i32.store)
  (func $runtime.resume$1 (type 0)
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
        i32.const 0
        call $syscall/js.handleEvent
        i32.const 0
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
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
  (func $syscall/js.handleEvent (type 1) (param i32)
    (local i32 i32 i32 i32 i32 i32 i64 i32 i32 i32 i64 i64 i32 i32 i32 i32 i32 i32 i32 i32 f64 i32 i32 i32 i64 i64)
    global.get 1
    i32.const 2
    i32.eq
    if  ;; label = @1
      global.get 2
      global.get 2
      i32.load
      i32.const 108
      i32.sub
      i32.store
      global.get 2
      i32.load
      local.tee 3
      i32.load
      local.set 0
      local.get 3
      i32.load offset=4
      local.set 1
      local.get 3
      i32.load offset=12
      local.set 4
      local.get 3
      i32.load offset=16
      local.set 5
      local.get 3
      i64.load offset=20 align=4
      local.set 7
      local.get 3
      i32.load offset=28
      local.set 8
      local.get 3
      i32.load offset=32
      local.set 9
      local.get 3
      i32.load offset=36
      local.set 10
      local.get 3
      i64.load offset=40 align=4
      local.set 11
      local.get 3
      i64.load offset=48 align=4
      local.set 12
      local.get 3
      i32.load offset=56
      local.set 13
      local.get 3
      i32.load offset=60
      local.set 14
      local.get 3
      i32.load offset=64
      local.set 15
      local.get 3
      i32.load offset=68
      local.set 16
      local.get 3
      i32.load offset=72
      local.set 17
      local.get 3
      i32.load offset=76
      local.set 18
      local.get 3
      i32.load offset=80
      local.set 19
      local.get 3
      i32.load offset=84
      local.set 20
      local.get 3
      i32.load offset=88
      local.set 22
      local.get 3
      i32.load offset=92
      local.set 23
      local.get 3
      i32.load offset=96
      local.set 24
      local.get 3
      i64.load offset=100 align=4
      local.set 25
      local.get 3
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
        i32.const 736
        i32.sub
        local.tee 1
        global.set $__stack_pointer
        local.get 1
        i32.const 600
        i32.add
        local.tee 2
        i64.const 0
        i64.store
        local.get 1
        i32.const 488
        i32.add
        local.tee 0
        i64.const 0
        i64.store
        local.get 1
        i32.const 384
        i32.add
        local.tee 5
        i64.const 0
        i64.store
        local.get 1
        i32.const 368
        i32.add
        local.tee 8
        i64.const 0
        i64.store
        local.get 1
        i32.const 360
        i32.add
        local.tee 9
        i64.const 0
        i64.store
        local.get 1
        i64.const 0
        i64.store offset=592
        local.get 1
        i64.const 0
        i64.store offset=584
        local.get 1
        i64.const 0
        i64.store offset=576
        local.get 1
        i64.const 0
        i64.store offset=568
        local.get 1
        i64.const 0
        i64.store offset=560
        local.get 1
        i64.const 0
        i64.store offset=552
        local.get 1
        i64.const 0
        i64.store offset=544
        local.get 1
        i64.const 0
        i64.store offset=536
        local.get 1
        i64.const 0
        i64.store offset=528
        local.get 1
        i64.const 0
        i64.store offset=520
        local.get 1
        i64.const 0
        i64.store offset=512
        local.get 1
        i64.const 0
        i64.store offset=504
        local.get 1
        i64.const 0
        i64.store offset=496
        local.get 1
        i64.const 0
        i64.store offset=480
        local.get 1
        i64.const 0
        i64.store offset=472
        local.get 1
        i64.const 0
        i64.store offset=464
        local.get 1
        i64.const 0
        i64.store offset=456
        local.get 1
        i64.const 0
        i64.store offset=448
        local.get 1
        i64.const 0
        i64.store offset=440
        local.get 1
        i64.const 0
        i64.store offset=432
        local.get 1
        i64.const 0
        i64.store offset=424
        local.get 1
        i64.const 0
        i64.store offset=416
        local.get 1
        i64.const 0
        i64.store offset=408
        local.get 1
        i64.const 0
        i64.store offset=400
        local.get 1
        i64.const 0
        i64.store offset=392
        local.get 1
        i64.const 0
        i64.store offset=376
        local.get 1
        i64.const 403726925824
        i64.store offset=352
        local.get 1
        i64.const 0
        i64.store offset=728
        local.get 1
        i64.const 0
        i64.store offset=720
        local.get 1
        i64.const 0
        i64.store offset=712
        local.get 1
        i64.const 0
        i64.store offset=704
        local.get 1
        i64.const 0
        i64.store offset=696
        local.get 1
        i64.const 0
        i64.store offset=688
        local.get 1
        i64.const 0
        i64.store offset=680
        local.get 1
        i64.const 0
        i64.store offset=672
        local.get 1
        i64.const 0
        i64.store offset=664
        local.get 1
        i64.const 0
        i64.store offset=656
        local.get 1
        i64.const 0
        i64.store offset=648
        local.get 1
        i64.const 0
        i64.store offset=640
        local.get 1
        i64.const 0
        i64.store offset=632
        local.get 1
        i64.const 0
        i64.store offset=624
        local.get 1
        i64.const 0
        i64.store offset=616
        local.get 1
        i64.const 0
        i64.store offset=608
        local.get 9
        local.get 1
        i32.const 216
        i32.add
        i32.store
        local.get 1
        i64.const 0
        i64.store offset=216
        local.get 1
        i32.const 0
        i32.store offset=224
        i32.const 66380
        i32.load
        local.set 19
        i32.const 66380
        local.get 1
        i32.const 352
        i32.add
        i32.store
        local.get 1
        local.get 19
        i32.store offset=352
        local.get 1
        i32.const 200
        i32.add
        local.set 4
      end
      local.get 6
      i32.const 0
      global.get 1
      select
      i32.eqz
      if  ;; label = @2
        local.get 4
        i64.const 9221120241336057862
        i32.const 0
        i32.const 65731
        i32.const 13
        call $_syscall/js.Value_.Get
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
        local.get 2
        local.get 1
        i32.load offset=208
        local.tee 10
        i32.store
        local.get 1
        i32.const 500
        i32.add
        local.get 10
        i32.store
        local.get 0
        local.get 10
        i32.store
        local.get 5
        local.get 10
        i32.store
        local.get 8
        local.get 10
        i32.store
        local.get 1
        i32.const 364
        i32.add
        local.get 10
        i32.store
        local.get 1
        i32.const 372
        i32.add
        local.get 1
        i32.const 336
        i32.add
        local.tee 4
        i32.store
        local.get 1
        local.get 10
        i32.store offset=224
        local.get 1
        local.get 1
        i64.load offset=200
        local.tee 12
        i64.store offset=216
        local.get 12
        i64.const 9221120237041090562
        i64.eq
        local.set 0
      end
      block  ;; label = @2
        block  ;; label = @3
          block  ;; label = @4
            block  ;; label = @5
              block  ;; label = @6
                block  ;; label = @7
                  block  ;; label = @8
                    global.get 1
                    i32.eqz
                    if  ;; label = @9
                      local.get 0
                      br_if 1 (;@8;)
                      local.get 1
                      i32.const 376
                      i32.add
                      local.set 0
                    end
                    local.get 6
                    i32.const 1
                    i32.eq
                    i32.const 1
                    global.get 1
                    select
                    if  ;; label = @9
                      i32.const 16
                      call $runtime.alloc
                      local.set 3
                      i32.const 1
                      global.get 1
                      i32.const 1
                      i32.eq
                      br_if 8 (;@1;)
                      drop
                      local.get 3
                      local.set 2
                    end
                    global.get 1
                    i32.eqz
                    if  ;; label = @9
                      local.get 0
                      local.get 2
                      i32.store
                      local.get 1
                      i32.const 380
                      i32.add
                      local.tee 0
                      local.get 2
                      i32.store
                      local.get 2
                      i64.const 9221120237041090562
                      i64.store
                      local.get 2
                      i32.const 0
                      i32.store offset=8
                    end
                    local.get 6
                    i32.const 2
                    i32.eq
                    i32.const 1
                    global.get 1
                    select
                    if  ;; label = @9
                      i64.const 9221120241336057862
                      i32.const 0
                      i32.const 65731
                      i32.const 13
                      i32.const 63
                      local.get 2
                      call $_syscall/js.Value_.Set
                      i32.const 2
                      global.get 1
                      i32.const 1
                      i32.eq
                      br_if 8 (;@1;)
                      drop
                    end
                    local.get 0
                    local.get 1
                    i32.const 184
                    i32.add
                    global.get 1
                    select
                    local.set 0
                    local.get 6
                    i32.const 3
                    i32.eq
                    i32.const 1
                    global.get 1
                    select
                    if  ;; label = @9
                      local.get 0
                      local.get 12
                      local.get 10
                      i32.const 65744
                      i32.const 2
                      call $_syscall/js.Value_.Get
                      i32.const 3
                      global.get 1
                      i32.const 1
                      i32.eq
                      br_if 8 (;@1;)
                      drop
                    end
                    global.get 1
                    i32.eqz
                    if  ;; label = @9
                      local.get 1
                      i32.const 420
                      i32.add
                      local.get 1
                      i32.load offset=192
                      local.tee 2
                      i32.store
                      local.get 1
                      i32.const 404
                      i32.add
                      local.get 2
                      i32.store
                      local.get 1
                      i32.const 396
                      i32.add
                      local.get 2
                      i32.store
                      local.get 1
                      i32.const 388
                      i32.add
                      local.get 2
                      i32.store
                      local.get 1
                      i32.const 392
                      i32.add
                      local.get 1
                      i32.const 320
                      i32.add
                      i32.store
                      local.get 1
                      i32.const 412
                      i32.add
                      local.get 1
                      i32.const 336
                      i32.add
                      i32.store
                      local.get 1
                      i32.const 408
                      i32.add
                      local.get 1
                      i32.const 336
                      i32.add
                      i32.store
                      local.get 1
                      i32.const 400
                      i32.add
                      local.get 1
                      i32.const 336
                      i32.add
                      local.tee 4
                      i32.store
                      local.get 1
                      i64.load offset=184
                      local.set 7
                      local.get 1
                      i32.const 0
                      i32.store offset=328
                      local.get 1
                      local.get 7
                      i64.store offset=320
                      local.get 1
                      i32.const 344
                      i32.add
                      local.tee 0
                      i64.const 0
                      i64.store
                      local.get 0
                      local.get 2
                      i32.store
                      local.get 1
                      local.get 7
                      i64.store offset=336
                      local.get 7
                      local.get 2
                      call $_syscall/js.Value_.isNumber
                      i32.const 1
                      i32.and
                      i32.eqz
                      local.tee 0
                      br_if 2 (;@7;)
                      local.get 7
                      f64.reinterpret_i64
                      f64.const 0x0p+0 (;=0;)
                      local.get 7
                      i64.const 9221120237041090561
                      i64.ne
                      select
                      local.tee 21
                      i32.trunc_sat_f64_s
                      i32.const 0
                      i32.const 2147483647
                      i32.const -2147483648
                      local.get 21
                      f64.const -0x1p+31 (;=-2.14748e+09;)
                      f64.ge
                      local.tee 2
                      select
                      local.get 21
                      local.get 21
                      f64.ne
                      select
                      local.tee 0
                      local.get 21
                      f64.const 0x1.fffffffcp+30 (;=2.14748e+09;)
                      f64.le
                      local.tee 5
                      select
                      local.get 0
                      local.get 2
                      select
                      local.tee 0
                      i32.eqz
                      local.tee 4
                      br_if 3 (;@6;)
                      i32.const 66424
                      i32.load8_u
                      i32.eqz
                      local.set 4
                    end
                    block  ;; label = @9
                      block  ;; label = @10
                        global.get 1
                        i32.eqz
                        if  ;; label = @11
                          local.get 4
                          br_if 1 (;@10;)
                          local.get 1
                          i32.const 432
                          i32.add
                          i32.const 66224
                          i32.load
                          local.tee 2
                          i32.store
                          local.get 1
                          i32.const 428
                          i32.add
                          local.get 2
                          i32.store
                          local.get 2
                          i32.eqz
                          local.tee 4
                          br_if 6 (;@5;)
                          local.get 1
                          i32.const 436
                          i32.add
                          local.tee 4
                          i32.const 66428
                          i32.load
                          local.tee 5
                          i32.store
                          i32.const 66428
                          local.get 2
                          i32.store
                          local.get 2
                          local.get 5
                          i32.store
                        end
                        local.get 6
                        i32.const 4
                        i32.eq
                        i32.const 1
                        global.get 1
                        select
                        if  ;; label = @11
                          call $internal/task.Pause
                          i32.const 4
                          global.get 1
                          i32.const 1
                          i32.eq
                          br_if 10 (;@1;)
                          drop
                        end
                        global.get 1
                        i32.eqz
                        br_if 1 (;@9;)
                      end
                      global.get 1
                      i32.eqz
                      if  ;; label = @10
                        i32.const 66424
                        i32.const 1
                        i32.store8
                      end
                    end
                    global.get 1
                    i32.eqz
                    if  ;; label = @9
                      local.get 1
                      local.get 0
                      i32.store offset=320
                      local.get 1
                      i32.const 320
                      i32.add
                      i32.const 66204
                      i32.load8_u
                      i32.const 66196
                      i32.load
                      local.get 1
                      call $runtime.hash32
                      local.set 2
                      local.get 1
                      i32.const 440
                      i32.add
                      i32.const 66192
                      i32.load
                      local.tee 0
                      i32.store
                      i32.const -1
                      i32.const -1
                      i32.const 66206
                      i32.load8_u
                      local.tee 5
                      i32.shl
                      i32.const -1
                      i32.xor
                      local.get 5
                      i32.const 31
                      i32.gt_u
                      select
                      local.get 2
                      i32.and
                      local.tee 5
                      i32.const 66205
                      i32.load8_u
                      i32.const 66204
                      i32.load8_u
                      i32.add
                      i32.const 3
                      i32.shl
                      i32.const 12
                      i32.add
                      i32.mul
                      local.get 0
                      i32.add
                      local.set 0
                      local.get 2
                      i32.const 24
                      i32.shr_u
                      local.tee 2
                      i32.const 1
                      local.get 2
                      select
                      local.set 8
                      local.get 1
                      i32.const 448
                      i32.add
                      local.set 14
                      local.get 1
                      i32.const 444
                      i32.add
                      local.set 15
                      local.get 1
                      i32.const 464
                      i32.add
                      local.set 22
                      local.get 1
                      i32.const 456
                      i32.add
                      local.set 16
                      local.get 1
                      i32.const 460
                      i32.add
                      local.set 17
                      local.get 1
                      i32.const 452
                      i32.add
                      local.set 4
                    end
                    block  ;; label = @9
                      loop  ;; label = @10
                        block  ;; label = @11
                          global.get 1
                          i32.eqz
                          if  ;; label = @12
                            local.get 14
                            local.get 0
                            i32.store
                            local.get 4
                            local.get 0
                            i32.store
                            local.get 15
                            local.get 0
                            i32.store
                            local.get 0
                            i32.eqz
                            local.tee 2
                            br_if 1 (;@11;)
                            i32.const 0
                            local.set 2
                          end
                          loop  ;; label = @12
                            block  ;; label = @13
                              global.get 1
                              i32.eqz
                              if  ;; label = @14
                                local.get 2
                                i32.const 8
                                i32.eq
                                local.tee 5
                                br_if 1 (;@13;)
                                local.get 8
                                local.get 0
                                local.get 2
                                i32.add
                                i32.load8_u
                                i32.ne
                                local.set 5
                              end
                              block  ;; label = @14
                                global.get 1
                                i32.eqz
                                if  ;; label = @15
                                  local.get 5
                                  br_if 1 (;@14;)
                                  i32.const 66205
                                  i32.load8_u
                                  local.set 13
                                  i32.const 66204
                                  i32.load8_u
                                  local.set 9
                                  local.get 16
                                  i32.const 66208
                                  i32.load
                                  local.tee 18
                                  i32.store
                                  local.get 17
                                  i32.const 66212
                                  i32.load
                                  local.tee 5
                                  i32.store
                                  local.get 5
                                  i32.eqz
                                  br_if 10 (;@5;)
                                  local.get 1
                                  i32.const 320
                                  i32.add
                                  local.set 23
                                  i32.const 66204
                                  i32.load8_u
                                  local.set 24
                                  local.get 2
                                  local.get 9
                                  i32.mul
                                  local.get 0
                                  i32.add
                                  i32.const 12
                                  i32.add
                                  local.set 20
                                end
                                local.get 6
                                i32.const 5
                                i32.eq
                                i32.const 1
                                global.get 1
                                select
                                if  ;; label = @15
                                  local.get 23
                                  local.get 20
                                  local.get 24
                                  local.get 18
                                  local.get 5
                                  call_indirect (type 4)
                                  local.set 3
                                  i32.const 5
                                  global.get 1
                                  i32.const 1
                                  i32.eq
                                  br_if 14 (;@1;)
                                  drop
                                  local.get 3
                                  local.set 5
                                end
                                global.get 1
                                i32.eqz
                                if  ;; label = @15
                                  local.get 5
                                  i32.const 1
                                  i32.and
                                  i32.eqz
                                  local.tee 5
                                  br_if 1 (;@14;)
                                  local.get 1
                                  i32.const 336
                                  i32.add
                                  local.tee 4
                                  local.get 2
                                  local.get 13
                                  i32.mul
                                  local.get 9
                                  i32.const 3
                                  i32.shl
                                  i32.add
                                  local.get 0
                                  i32.add
                                  i32.const 12
                                  i32.add
                                  local.tee 2
                                  i32.const 66205
                                  i32.load8_u
                                  local.tee 5
                                  memory.copy
                                  br 6 (;@9;)
                                end
                              end
                              global.get 1
                              i32.eqz
                              if  ;; label = @14
                                local.get 2
                                i32.const 1
                                i32.add
                                local.set 2
                                br 2 (;@12;)
                              end
                            end
                          end
                          global.get 1
                          i32.eqz
                          if  ;; label = @12
                            local.get 22
                            local.get 0
                            i32.load offset=8
                            local.tee 0
                            i32.store
                            br 2 (;@10;)
                          end
                        end
                      end
                      global.get 1
                      i32.eqz
                      if  ;; label = @10
                        local.get 1
                        i32.const 336
                        i32.add
                        local.tee 4
                        i32.const 0
                        i32.const 66205
                        i32.load8_u
                        local.tee 2
                        memory.fill
                      end
                    end
                    global.get 1
                    i32.eqz
                    if  ;; label = @9
                      i32.const 66424
                      i32.load8_u
                      i32.eqz
                      local.set 4
                    end
                    block  ;; label = @9
                      block  ;; label = @10
                        block  ;; label = @11
                          global.get 1
                          i32.eqz
                          if  ;; label = @12
                            local.get 4
                            br_if 1 (;@11;)
                            local.get 1
                            i32.load offset=336
                            local.set 15
                            local.get 1
                            i32.load offset=340
                            local.set 14
                            local.get 1
                            i32.const 480
                            i32.add
                            i32.const 66428
                            i32.load
                            local.tee 2
                            i32.store
                            local.get 1
                            i32.const 476
                            i32.add
                            local.get 2
                            i32.store
                            local.get 1
                            i32.const 468
                            i32.add
                            local.get 2
                            i32.store
                            local.get 2
                            i32.eqz
                            local.tee 4
                            br_if 2 (;@10;)
                            i32.const 66428
                            local.get 2
                            i32.load
                            local.tee 5
                            i32.store
                            local.get 1
                            i32.const 472
                            i32.add
                            local.tee 4
                            local.get 5
                            i32.store
                            local.get 2
                            i32.const 0
                            i32.store
                          end
                          local.get 6
                          i32.const 6
                          i32.eq
                          i32.const 1
                          global.get 1
                          select
                          if  ;; label = @12
                            local.get 2
                            call $runtime.runqueuePushBack
                            i32.const 6
                            global.get 1
                            i32.const 1
                            i32.eq
                            br_if 11 (;@1;)
                            drop
                          end
                          global.get 1
                          i32.eqz
                          br_if 2 (;@9;)
                        end
                        local.get 6
                        i32.const 7
                        i32.eq
                        i32.const 1
                        global.get 1
                        select
                        if  ;; label = @11
                          i32.const 34
                          i32.const 65712
                          call $runtime._panic
                          i32.const 7
                          global.get 1
                          i32.const 1
                          i32.eq
                          br_if 10 (;@1;)
                          drop
                        end
                        global.get 1
                        i32.eqz
                        if  ;; label = @11
                          unreachable
                        end
                      end
                      global.get 1
                      i32.eqz
                      if  ;; label = @10
                        i32.const 66424
                        i32.const 0
                        i32.store8
                      end
                    end
                    local.get 0
                    local.get 0
                    i32.eqz
                    global.get 1
                    select
                    local.set 0
                    block  ;; label = @9
                      block  ;; label = @10
                        block  ;; label = @11
                          global.get 1
                          i32.eqz
                          if  ;; label = @12
                            local.get 0
                            br_if 1 (;@11;)
                            local.get 1
                            i32.const 0
                            i32.store offset=240
                            local.get 1
                            i64.const 0
                            i64.store offset=232
                            local.get 1
                            i32.const 484
                            i32.add
                            local.get 1
                            i32.const 232
                            i32.add
                            local.tee 4
                            i32.store
                            local.get 1
                            i32.const 168
                            i32.add
                            local.set 0
                          end
                          local.get 6
                          i32.const 8
                          i32.eq
                          i32.const 1
                          global.get 1
                          select
                          if  ;; label = @12
                            local.get 0
                            local.get 12
                            local.get 10
                            i32.const 65746
                            i32.const 4
                            call $_syscall/js.Value_.Get
                            i32.const 8
                            global.get 1
                            i32.const 1
                            i32.eq
                            br_if 11 (;@1;)
                            drop
                          end
                          global.get 1
                          i32.eqz
                          if  ;; label = @12
                            local.get 1
                            i32.const 592
                            i32.add
                            local.get 1
                            i32.load offset=176
                            local.tee 13
                            i32.store
                            local.get 1
                            i32.const 492
                            i32.add
                            local.get 13
                            i32.store
                            local.get 1
                            i32.const 0
                            i32.store offset=256
                            local.get 1
                            i64.const 0
                            i64.store offset=248
                            local.get 1
                            i32.const 496
                            i32.add
                            local.get 1
                            i32.const 248
                            i32.add
                            local.tee 4
                            i32.store
                            local.get 1
                            local.get 13
                            i32.store offset=240
                            local.get 1
                            local.get 1
                            i64.load offset=168
                            local.tee 25
                            i64.store offset=232
                            local.get 1
                            i32.const 152
                            i32.add
                            local.set 0
                          end
                          local.get 6
                          i32.const 9
                          i32.eq
                          i32.const 1
                          global.get 1
                          select
                          if  ;; label = @12
                            local.get 0
                            local.get 12
                            local.get 10
                            i32.const 65750
                            i32.const 4
                            call $_syscall/js.Value_.Get
                            i32.const 9
                            global.get 1
                            i32.const 1
                            i32.eq
                            br_if 11 (;@1;)
                            drop
                          end
                          global.get 1
                          i32.eqz
                          if  ;; label = @12
                            local.get 1
                            i32.const 572
                            i32.add
                            local.get 1
                            i32.load offset=160
                            local.tee 8
                            i32.store
                            local.get 1
                            i32.const 552
                            i32.add
                            local.get 8
                            i32.store
                            local.get 1
                            i32.const 544
                            i32.add
                            local.get 8
                            i32.store
                            local.get 1
                            i32.const 528
                            i32.add
                            local.get 8
                            i32.store
                            local.get 1
                            i32.const 516
                            i32.add
                            local.get 8
                            i32.store
                            local.get 1
                            i32.const 508
                            i32.add
                            local.get 8
                            i32.store
                            local.get 1
                            i32.const 504
                            i32.add
                            local.get 8
                            i32.store
                            local.get 1
                            i32.const 536
                            i32.add
                            local.get 1
                            i32.const 336
                            i32.add
                            i32.store
                            local.get 1
                            i32.const 532
                            i32.add
                            local.get 1
                            i32.const 336
                            i32.add
                            i32.store
                            local.get 1
                            local.get 8
                            i32.store offset=256
                            local.get 1
                            local.get 1
                            i64.load offset=152
                            local.tee 7
                            i64.store offset=248
                            local.get 1
                            i32.const 512
                            i32.add
                            local.tee 0
                            local.get 1
                            i32.const 320
                            i32.add
                            local.tee 4
                            i32.store
                            local.get 1
                            i32.const 0
                            i32.store offset=328
                            local.get 1
                            i64.const 0
                            i64.store offset=320
                          end
                          local.get 6
                          i32.const 10
                          i32.eq
                          i32.const 1
                          global.get 1
                          select
                          if  ;; label = @12
                            local.get 7
                            local.get 8
                            call $_syscall/js.Value_.Type
                            local.set 3
                            i32.const 10
                            global.get 1
                            i32.const 1
                            i32.eq
                            br_if 11 (;@1;)
                            drop
                            local.get 3
                            local.set 2
                          end
                          global.get 1
                          i32.eqz
                          if  ;; label = @12
                            local.get 2
                            i32.const -2
                            i32.and
                            i32.const 6
                            i32.ne
                            local.tee 0
                            br_if 8 (;@4;)
                          end
                          local.get 6
                          i32.const 11
                          i32.eq
                          i32.const 1
                          global.get 1
                          select
                          if  ;; label = @12
                            local.get 7
                            local.get 1
                            call $syscall/js.valueLength
                            local.set 3
                            i32.const 11
                            global.get 1
                            i32.const 1
                            i32.eq
                            br_if 11 (;@1;)
                            drop
                            local.get 3
                            local.set 9
                          end
                          global.get 1
                          i32.eqz
                          if  ;; label = @12
                            local.get 9
                            i32.const 268435455
                            i32.gt_u
                            local.tee 0
                            br_if 2 (;@10;)
                            local.get 9
                            i32.const 4
                            i32.shl
                            local.set 4
                            local.get 1
                            i32.const 540
                            i32.add
                            local.set 0
                          end
                          local.get 6
                          i32.const 12
                          i32.eq
                          i32.const 1
                          global.get 1
                          select
                          if  ;; label = @12
                            local.get 4
                            call $runtime.alloc
                            local.set 3
                            i32.const 12
                            global.get 1
                            i32.const 1
                            i32.eq
                            br_if 11 (;@1;)
                            drop
                            local.get 3
                            local.set 4
                          end
                          global.get 1
                          i32.eqz
                          if  ;; label = @12
                            local.get 0
                            local.get 4
                            i32.store
                            local.get 1
                            i32.const 580
                            i32.add
                            local.get 1
                            i32.const 336
                            i32.add
                            i32.store
                            local.get 1
                            i32.const 576
                            i32.add
                            local.get 1
                            i32.const 336
                            i32.add
                            i32.store
                            local.get 1
                            i32.const 564
                            i32.add
                            local.get 1
                            i32.const 320
                            i32.add
                            i32.store
                            local.get 1
                            i32.const 548
                            i32.add
                            local.get 1
                            i32.const 288
                            i32.add
                            i32.store
                            local.get 1
                            i32.const 588
                            i32.add
                            local.set 18
                            local.get 1
                            i32.const 584
                            i32.add
                            local.set 16
                            local.get 1
                            i32.const 568
                            i32.add
                            local.set 17
                            i32.const 0
                            local.set 2
                            local.get 4
                            local.set 0
                          end
                          loop  ;; label = @12
                            block  ;; label = @13
                              global.get 1
                              i32.eqz
                              if  ;; label = @14
                                local.get 2
                                local.get 9
                                i32.eq
                                local.tee 5
                                br_if 1 (;@13;)
                                local.get 1
                                i32.const 0
                                i32.store offset=296
                                local.get 1
                                i64.const 0
                                i64.store offset=288
                              end
                              local.get 6
                              i32.const 13
                              i32.eq
                              i32.const 1
                              global.get 1
                              select
                              if  ;; label = @14
                                local.get 7
                                local.get 8
                                call $_syscall/js.Value_.Type
                                local.set 3
                                i32.const 13
                                global.get 1
                                i32.const 1
                                i32.eq
                                br_if 13 (;@1;)
                                drop
                                local.get 3
                                local.set 5
                              end
                              global.get 1
                              i32.eqz
                              if  ;; label = @14
                                local.get 5
                                i32.const -2
                                i32.and
                                i32.const 6
                                i32.ne
                                br_if 11 (;@3;)
                                local.get 1
                                i32.const 0
                                i32.store offset=328
                                local.get 1
                                i64.const 0
                                i64.store offset=320
                                local.get 1
                                i32.const 136
                                i32.add
                                local.set 5
                              end
                              local.get 6
                              i32.const 14
                              i32.eq
                              i32.const 1
                              global.get 1
                              select
                              if  ;; label = @14
                                local.get 7
                                local.get 2
                                local.get 1
                                call $syscall/js.valueIndex
                                local.set 26
                                i32.const 14
                                global.get 1
                                i32.const 1
                                i32.eq
                                br_if 13 (;@1;)
                                drop
                                local.get 26
                                local.set 11
                              end
                              local.get 6
                              i32.const 15
                              i32.eq
                              i32.const 1
                              global.get 1
                              select
                              if  ;; label = @14
                                local.get 5
                                local.get 11
                                call $syscall/js.makeValue
                                i32.const 15
                                global.get 1
                                i32.const 1
                                i32.eq
                                br_if 13 (;@1;)
                                drop
                              end
                              global.get 1
                              i32.eqz
                              if  ;; label = @14
                                local.get 18
                                local.get 1
                                i32.load offset=144
                                local.tee 5
                                i32.store
                                local.get 16
                                local.get 5
                                i32.store
                                local.get 17
                                local.get 5
                                i32.store
                                local.get 1
                                i64.load offset=136
                                local.set 11
                                local.get 0
                                local.get 5
                                i32.store offset=8
                                local.get 0
                                local.get 11
                                i64.store
                                local.get 0
                                i32.const 16
                                i32.add
                                local.set 0
                                local.get 2
                                i32.const 1
                                i32.add
                                local.set 2
                                br 2 (;@12;)
                              end
                            end
                          end
                          global.get 1
                          i32.eqz
                          if  ;; label = @12
                            local.get 14
                            i32.eqz
                            local.tee 0
                            br_if 7 (;@5;)
                            local.get 1
                            i32.const 128
                            i32.add
                            local.set 0
                          end
                          local.get 6
                          i32.const 16
                          i32.eq
                          i32.const 1
                          global.get 1
                          select
                          if  ;; label = @12
                            local.get 0
                            local.get 25
                            local.get 13
                            local.get 4
                            local.get 9
                            local.get 9
                            local.get 15
                            local.get 14
                            call_indirect (type 18)
                            i32.const 16
                            global.get 1
                            i32.const 1
                            i32.eq
                            br_if 11 (;@1;)
                            drop
                          end
                          global.get 1
                          i32.eqz
                          if  ;; label = @12
                            local.get 1
                            i32.const 596
                            i32.add
                            local.get 1
                            i32.load offset=132
                            local.tee 2
                            i32.store
                            local.get 1
                            i32.load offset=128
                            local.set 0
                          end
                          local.get 6
                          i32.const 17
                          i32.eq
                          i32.const 1
                          global.get 1
                          select
                          if  ;; label = @12
                            local.get 12
                            local.get 10
                            i32.const 65754
                            i32.const 6
                            local.get 0
                            local.get 2
                            call $_syscall/js.Value_.Set
                            i32.const 17
                            global.get 1
                            i32.const 1
                            i32.eq
                            br_if 11 (;@1;)
                            drop
                          end
                          global.get 1
                          i32.eqz
                          br_if 3 (;@8;)
                        end
                        local.get 0
                        local.get 1
                        i32.const 112
                        i32.add
                        global.get 1
                        select
                        local.set 0
                        local.get 6
                        i32.const 18
                        i32.eq
                        i32.const 1
                        global.get 1
                        select
                        if  ;; label = @11
                          local.get 0
                          i64.const 9221120241336057861
                          i32.const 0
                          i32.const 65760
                          i32.const 7
                          call $_syscall/js.Value_.Get
                          i32.const 18
                          global.get 1
                          i32.const 1
                          i32.eq
                          br_if 10 (;@1;)
                          drop
                        end
                        global.get 1
                        i32.eqz
                        if  ;; label = @11
                          local.get 1
                          i32.const 668
                          i32.add
                          local.get 1
                          i32.load offset=120
                          local.tee 2
                          i32.store
                          local.get 1
                          i32.const 664
                          i32.add
                          local.get 2
                          i32.store
                          local.get 1
                          i32.const 644
                          i32.add
                          local.get 2
                          i32.store
                          local.get 1
                          i32.const 604
                          i32.add
                          local.get 2
                          i32.store
                          local.get 1
                          i64.const 0
                          i64.store offset=264
                          local.get 1
                          i32.const 608
                          i32.add
                          local.get 1
                          i32.const 264
                          i32.add
                          i32.store
                          local.get 1
                          i32.const 652
                          i32.add
                          local.get 1
                          i32.const 320
                          i32.add
                          i32.store
                          local.get 1
                          i32.const 648
                          i32.add
                          local.get 1
                          i32.const 320
                          i32.add
                          i32.store
                          local.get 1
                          i32.const 34
                          i32.store offset=264
                          local.get 1
                          i32.const 65792
                          i32.store offset=268
                          local.get 1
                          i32.const 660
                          i32.add
                          local.get 1
                          i32.const 304
                          i32.add
                          i32.store
                          local.get 1
                          i32.const 656
                          i32.add
                          local.get 1
                          i32.const 304
                          i32.add
                          i32.store
                          local.get 1
                          i32.const 700
                          i32.add
                          local.get 1
                          i32.const 288
                          i32.add
                          i32.store
                          local.get 1
                          i32.const 612
                          i32.add
                          local.get 1
                          i32.const 272
                          i32.add
                          i32.store
                          local.get 1
                          i64.load offset=112
                          local.set 7
                          local.get 1
                          i32.const 0
                          i32.store offset=280
                          local.get 1
                          i64.const 0
                          i64.store offset=272
                          local.get 1
                          i32.const 624
                          i32.add
                          local.tee 0
                          local.get 1
                          i32.const 336
                          i32.add
                          i32.store
                          local.get 1
                          i32.const 616
                          i32.add
                          local.set 4
                        end
                        local.get 6
                        i32.const 19
                        i32.eq
                        i32.const 1
                        global.get 1
                        select
                        if  ;; label = @11
                          i32.const 16
                          call $runtime.alloc
                          local.set 3
                          i32.const 19
                          global.get 1
                          i32.const 1
                          i32.eq
                          br_if 10 (;@1;)
                          drop
                          local.get 3
                          local.set 0
                        end
                        global.get 1
                        i32.eqz
                        if  ;; label = @11
                          local.get 4
                          local.get 0
                          i32.store
                          local.get 1
                          i32.const 636
                          i32.add
                          local.get 0
                          i32.store
                          local.get 1
                          i32.const 620
                          i32.add
                          local.set 4
                        end
                        local.get 6
                        i32.const 20
                        i32.eq
                        i32.const 1
                        global.get 1
                        select
                        if  ;; label = @11
                          i32.const 8
                          call $runtime.alloc
                          local.set 3
                          i32.const 20
                          global.get 1
                          i32.const 1
                          i32.eq
                          br_if 10 (;@1;)
                          drop
                          local.get 3
                          local.set 5
                        end
                        global.get 1
                        i32.eqz
                        if  ;; label = @11
                          local.get 4
                          local.get 5
                          i32.store
                          local.get 1
                          i32.const 640
                          i32.add
                          local.get 5
                          i32.store
                          local.get 1
                          i32.const 0
                          i32.store offset=344
                          local.get 1
                          i64.const 0
                          i64.store offset=336
                          local.get 1
                          i32.const 96
                          i32.add
                          local.set 4
                        end
                        local.get 6
                        i32.const 21
                        i32.eq
                        i32.const 1
                        global.get 1
                        select
                        if  ;; label = @11
                          local.get 4
                          i32.const 34
                          i32.const 65792
                          call $syscall/js.ValueOf
                          i32.const 21
                          global.get 1
                          i32.const 1
                          i32.eq
                          br_if 10 (;@1;)
                          drop
                        end
                        global.get 1
                        i32.eqz
                        if  ;; label = @11
                          local.get 1
                          i32.const 632
                          i32.add
                          local.get 1
                          i32.load offset=104
                          local.tee 8
                          i32.store
                          local.get 1
                          i32.const 628
                          i32.add
                          local.get 8
                          i32.store
                          local.get 0
                          local.get 1
                          i64.load offset=96
                          local.tee 11
                          i64.store
                          local.get 0
                          local.get 8
                          i32.store offset=8
                          local.get 5
                          local.get 11
                          i64.store
                          local.get 1
                          i32.const 80
                          i32.add
                          local.set 4
                        end
                        local.get 6
                        i32.const 22
                        i32.eq
                        i32.const 1
                        global.get 1
                        select
                        if  ;; label = @11
                          local.get 4
                          local.get 7
                          i32.const 65800
                          i32.const 5
                          local.get 5
                          i32.const 1
                          i32.const 1
                          local.get 1
                          call $syscall/js.valueCall
                          i32.const 22
                          global.get 1
                          i32.const 1
                          i32.eq
                          br_if 10 (;@1;)
                          drop
                        end
                        global.get 1
                        i32.eqz
                        if  ;; label = @11
                          local.get 1
                          i32.const 328
                          i32.add
                          local.tee 5
                          i64.const 0
                          i64.store
                          local.get 1
                          i32.const 312
                          i32.add
                          local.tee 4
                          i32.const 1
                          i32.store
                          local.get 5
                          local.get 2
                          i32.store
                          local.get 1
                          i64.const 0
                          i64.store offset=304
                          local.get 1
                          local.get 7
                          i64.store offset=320
                          local.get 1
                          local.get 0
                          i32.store offset=304
                          local.get 1
                          i32.const 1
                          i32.store offset=308
                          local.get 1
                          i64.load offset=80
                          local.set 11
                          local.get 1
                          i32.load8_u offset=88
                          local.tee 0
                          br_if 2 (;@9;)
                        end
                        local.get 6
                        i32.const 23
                        i32.eq
                        i32.const 1
                        global.get 1
                        select
                        if  ;; label = @11
                          local.get 7
                          local.get 2
                          call $_syscall/js.Value_.Type
                          local.set 3
                          i32.const 23
                          global.get 1
                          i32.const 1
                          i32.eq
                          br_if 10 (;@1;)
                          drop
                          local.get 3
                          local.set 0
                        end
                        local.get 4
                        local.get 0
                        i32.const -2
                        i32.and
                        i32.const 6
                        i32.ne
                        global.get 1
                        select
                        local.set 4
                        block  ;; label = @11
                          global.get 1
                          i32.eqz
                          if  ;; label = @12
                            local.get 4
                            br_if 1 (;@11;)
                            local.get 1
                            i32.const -64
                            i32.sub
                            local.set 0
                          end
                          local.get 6
                          i32.const 24
                          i32.eq
                          i32.const 1
                          global.get 1
                          select
                          if  ;; label = @12
                            local.get 0
                            local.get 7
                            local.get 2
                            i32.const 65800
                            i32.const 5
                            call $_syscall/js.Value_.Get
                            i32.const 24
                            global.get 1
                            i32.const 1
                            i32.eq
                            br_if 11 (;@1;)
                            drop
                          end
                          global.get 1
                          i32.eqz
                          if  ;; label = @12
                            local.get 1
                            i32.const 672
                            i32.add
                            local.tee 0
                            local.get 1
                            i32.load offset=72
                            local.tee 2
                            i32.store
                            local.get 1
                            i64.load offset=64
                            local.set 7
                          end
                          local.get 6
                          i32.const 25
                          i32.eq
                          i32.const 1
                          global.get 1
                          select
                          if  ;; label = @12
                            local.get 7
                            local.get 2
                            call $_syscall/js.Value_.Type
                            local.set 3
                            i32.const 25
                            global.get 1
                            i32.const 1
                            i32.eq
                            br_if 11 (;@1;)
                            drop
                            local.get 3
                            local.set 2
                          end
                          global.get 1
                          i32.eqz
                          if  ;; label = @12
                            local.get 2
                            i32.const 7
                            i32.ne
                            local.tee 0
                            br_if 10 (;@2;)
                            local.get 1
                            i32.const 0
                            i32.store offset=296
                            local.get 1
                            i64.const 0
                            i64.store offset=288
                            local.get 1
                            i32.const 16
                            i32.add
                            local.set 0
                          end
                          local.get 6
                          i32.const 26
                          i32.eq
                          i32.const 1
                          global.get 1
                          select
                          if  ;; label = @12
                            local.get 0
                            local.get 11
                            call $syscall/js.makeValue
                            i32.const 26
                            global.get 1
                            i32.const 1
                            i32.eq
                            br_if 11 (;@1;)
                            drop
                          end
                          global.get 1
                          i32.eqz
                          if  ;; label = @12
                            local.get 1
                            i32.const 708
                            i32.add
                            local.get 1
                            i32.load offset=24
                            local.tee 0
                            i32.store
                            local.get 1
                            i32.const 704
                            i32.add
                            local.get 0
                            i32.store
                            local.get 1
                            i64.load offset=16
                            local.set 7
                            local.get 1
                            i32.const 712
                            i32.add
                            local.set 4
                          end
                          local.get 6
                          i32.const 27
                          i32.eq
                          i32.const 1
                          global.get 1
                          select
                          if  ;; label = @12
                            i32.const 16
                            call $runtime.alloc
                            local.set 3
                            i32.const 27
                            global.get 1
                            i32.const 1
                            i32.eq
                            br_if 11 (;@1;)
                            drop
                            local.get 3
                            local.set 2
                          end
                          global.get 1
                          i32.eqz
                          if  ;; label = @12
                            local.get 4
                            local.get 2
                            i32.store
                            local.get 1
                            i32.const 716
                            i32.add
                            local.tee 4
                            local.get 2
                            i32.store
                            local.get 2
                            local.get 7
                            i64.store
                            local.get 2
                            local.get 0
                            i32.store offset=8
                          end
                          local.get 6
                          i32.const 28
                          i32.eq
                          i32.const 1
                          global.get 1
                          select
                          if  ;; label = @12
                            i32.const 127
                            local.get 2
                            call $runtime._panic
                            i32.const 28
                            global.get 1
                            i32.const 1
                            i32.eq
                            br_if 11 (;@1;)
                            drop
                          end
                          global.get 1
                          i32.eqz
                          if  ;; label = @12
                            unreachable
                          end
                        end
                        local.get 4
                        local.get 1
                        i32.const 720
                        i32.add
                        global.get 1
                        select
                        local.set 4
                        local.get 6
                        i32.const 29
                        i32.eq
                        i32.const 1
                        global.get 1
                        select
                        if  ;; label = @11
                          i32.const 12
                          call $runtime.alloc
                          local.set 3
                          i32.const 29
                          global.get 1
                          i32.const 1
                          i32.eq
                          br_if 10 (;@1;)
                          drop
                          local.get 3
                          local.set 2
                        end
                        global.get 1
                        i32.eqz
                        if  ;; label = @11
                          local.get 4
                          local.get 2
                          i32.store
                          local.get 1
                          i32.const 724
                          i32.add
                          local.tee 4
                          local.get 2
                          i32.store
                          local.get 2
                          local.get 0
                          i32.store offset=8
                          local.get 2
                          i32.const 10
                          i32.store offset=4
                          local.get 2
                          i32.const 65862
                          i32.store
                        end
                        local.get 6
                        i32.const 30
                        i32.eq
                        i32.const 1
                        global.get 1
                        select
                        if  ;; label = @11
                          i32.const 3045
                          local.get 2
                          call $runtime._panic
                          i32.const 30
                          global.get 1
                          i32.const 1
                          i32.eq
                          br_if 10 (;@1;)
                          drop
                        end
                        global.get 1
                        i32.eqz
                        if  ;; label = @11
                          unreachable
                        end
                      end
                      local.get 6
                      i32.const 31
                      i32.eq
                      i32.const 1
                      global.get 1
                      select
                      if  ;; label = @10
                        call $runtime.slicePanic
                        i32.const 31
                        global.get 1
                        i32.const 1
                        i32.eq
                        br_if 9 (;@1;)
                        drop
                      end
                      global.get 1
                      i32.eqz
                      if  ;; label = @10
                        unreachable
                      end
                    end
                    local.get 6
                    i32.const 32
                    i32.eq
                    i32.const 1
                    global.get 1
                    select
                    if  ;; label = @9
                      local.get 1
                      local.get 11
                      call $syscall/js.makeValue
                      i32.const 32
                      global.get 1
                      i32.const 1
                      i32.eq
                      br_if 8 (;@1;)
                      drop
                    end
                    global.get 1
                    i32.eqz
                    if  ;; label = @9
                      local.get 1
                      i32.const 732
                      i32.add
                      local.get 1
                      i32.load offset=8
                      local.tee 2
                      i32.store
                      local.get 1
                      i32.const 728
                      i32.add
                      local.tee 0
                      local.get 2
                      i32.store
                    end
                  end
                  global.get 1
                  i32.eqz
                  if  ;; label = @8
                    i32.const 66380
                    local.get 19
                    i32.store
                    local.get 1
                    i32.const 736
                    i32.add
                    global.set $__stack_pointer
                    return
                  end
                end
                local.get 4
                local.get 1
                i32.const 416
                i32.add
                global.get 1
                select
                local.set 4
                local.get 6
                i32.const 33
                i32.eq
                i32.const 1
                global.get 1
                select
                if  ;; label = @7
                  i32.const 12
                  call $runtime.alloc
                  local.set 3
                  i32.const 33
                  global.get 1
                  i32.const 1
                  i32.eq
                  br_if 6 (;@1;)
                  drop
                  local.get 3
                  local.set 0
                end
                global.get 1
                i32.eqz
                if  ;; label = @7
                  local.get 4
                  local.get 0
                  i32.store
                  local.get 1
                  i32.const 424
                  i32.add
                  local.tee 4
                  local.get 0
                  i32.store
                end
                local.get 6
                i32.const 34
                i32.eq
                i32.const 1
                global.get 1
                select
                if  ;; label = @7
                  local.get 7
                  local.get 2
                  call $_syscall/js.Value_.Type
                  local.set 3
                  i32.const 34
                  global.get 1
                  i32.const 1
                  i32.eq
                  br_if 6 (;@1;)
                  drop
                  local.get 3
                  local.set 4
                end
                global.get 1
                i32.eqz
                if  ;; label = @7
                  local.get 0
                  local.get 4
                  i32.store offset=8
                  local.get 0
                  i32.const 9
                  i32.store offset=4
                  local.get 0
                  i32.const 65892
                  i32.store
                end
                local.get 6
                i32.const 35
                i32.eq
                i32.const 1
                global.get 1
                select
                if  ;; label = @7
                  i32.const 3045
                  local.get 0
                  call $runtime._panic
                  i32.const 35
                  global.get 1
                  i32.const 1
                  i32.eq
                  br_if 6 (;@1;)
                  drop
                end
                global.get 1
                i32.eqz
                if  ;; label = @7
                  unreachable
                end
              end
              local.get 6
              i32.const 36
              i32.eq
              i32.const 1
              global.get 1
              select
              if  ;; label = @6
                call $runtime.deadlock
                i32.const 36
                global.get 1
                i32.const 1
                i32.eq
                br_if 5 (;@1;)
                drop
              end
              global.get 1
              i32.eqz
              if  ;; label = @6
                unreachable
              end
            end
            local.get 6
            i32.const 37
            i32.eq
            i32.const 1
            global.get 1
            select
            if  ;; label = @5
              call $runtime.nilPanic
              i32.const 37
              global.get 1
              i32.const 1
              i32.eq
              br_if 4 (;@1;)
              drop
            end
            global.get 1
            i32.eqz
            if  ;; label = @5
              unreachable
            end
          end
          local.get 4
          local.get 1
          i32.const 520
          i32.add
          global.get 1
          select
          local.set 4
          local.get 6
          i32.const 38
          i32.eq
          i32.const 1
          global.get 1
          select
          if  ;; label = @4
            i32.const 12
            call $runtime.alloc
            local.set 3
            i32.const 38
            global.get 1
            i32.const 1
            i32.eq
            br_if 3 (;@1;)
            drop
            local.get 3
            local.set 0
          end
          global.get 1
          i32.eqz
          if  ;; label = @4
            local.get 4
            local.get 0
            i32.store
            local.get 1
            i32.const 524
            i32.add
            local.tee 4
            local.get 0
            i32.store
            local.get 0
            local.get 2
            i32.store offset=8
            local.get 0
            i32.const 14
            i32.store offset=4
            local.get 0
            i32.const 65901
            i32.store
          end
          local.get 6
          i32.const 39
          i32.eq
          i32.const 1
          global.get 1
          select
          if  ;; label = @4
            i32.const 3045
            local.get 0
            call $runtime._panic
            i32.const 39
            global.get 1
            i32.const 1
            i32.eq
            br_if 3 (;@1;)
            drop
          end
          global.get 1
          i32.eqz
          if  ;; label = @4
            unreachable
          end
        end
        local.get 0
        local.get 1
        i32.const 556
        i32.add
        global.get 1
        select
        local.set 0
        local.get 6
        i32.const 40
        i32.eq
        i32.const 1
        global.get 1
        select
        if  ;; label = @3
          i32.const 12
          call $runtime.alloc
          local.set 3
          i32.const 40
          global.get 1
          i32.const 1
          i32.eq
          br_if 2 (;@1;)
          drop
          local.get 3
          local.set 2
        end
        global.get 1
        i32.eqz
        if  ;; label = @3
          local.get 0
          local.get 2
          i32.store
          local.get 1
          i32.const 560
          i32.add
          local.tee 0
          local.get 2
          i32.store
          local.get 2
          local.get 5
          i32.store offset=8
          local.get 2
          i32.const 11
          i32.store offset=4
          local.get 2
          i32.const 65881
          i32.store
        end
        local.get 6
        i32.const 41
        i32.eq
        i32.const 1
        global.get 1
        select
        if  ;; label = @3
          i32.const 3045
          local.get 2
          call $runtime._panic
          i32.const 41
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
      local.get 0
      local.get 1
      i32.const 56
      i32.add
      global.get 1
      select
      local.set 0
      local.get 6
      i32.const 42
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        local.get 0
        i32.const 65805
        i32.const 33
        i32.const 65800
        i32.const 5
        call $runtime.stringConcat
        i32.const 42
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        local.get 1
        i32.const 676
        i32.add
        local.get 1
        i32.load offset=56
        local.tee 0
        i32.store
        local.get 1
        i32.load offset=60
        local.set 5
        local.get 1
        i32.const 48
        i32.add
        local.set 4
      end
      local.get 6
      i32.const 43
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        local.get 4
        local.get 0
        local.get 5
        i32.const 65838
        i32.const 24
        call $runtime.stringConcat
        i32.const 43
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        local.get 1
        i32.const 680
        i32.add
        local.get 1
        i32.load offset=48
        local.tee 0
        i32.store
        local.get 1
        i32.load offset=52
        local.set 5
        local.get 1
        i32.const 40
        i32.add
        local.set 4
      end
      local.get 6
      i32.const 44
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        local.get 4
        local.get 2
        call $_syscall/js.Type_.String
        i32.const 44
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        local.get 1
        i32.const 684
        i32.add
        local.get 1
        i32.load offset=40
        local.tee 2
        i32.store
        local.get 1
        i32.load offset=44
        local.set 8
        local.get 1
        i32.const 32
        i32.add
        local.set 4
      end
      local.get 6
      i32.const 45
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        local.get 4
        local.get 0
        local.get 5
        local.get 2
        local.get 8
        call $runtime.stringConcat
        i32.const 45
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        local.get 1
        i32.const 688
        i32.add
        local.get 1
        i32.load offset=32
        local.tee 0
        i32.store
        local.get 1
        i32.load offset=36
        local.set 5
        local.get 1
        i32.const 692
        i32.add
        local.set 4
      end
      local.get 6
      i32.const 46
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        i32.const 8
        call $runtime.alloc
        local.set 3
        i32.const 46
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
        local.get 3
        local.set 2
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        local.get 4
        local.get 2
        i32.store
        local.get 1
        i32.const 696
        i32.add
        local.get 2
        i32.store
        local.get 2
        local.get 5
        i32.store offset=4
        local.get 2
        local.get 0
        i32.store
      end
      local.get 6
      i32.const 47
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        i32.const 34
        local.get 2
        call $runtime._panic
        i32.const 47
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
    local.set 3
    global.get 2
    i32.load
    local.get 3
    i32.store
    global.get 2
    global.get 2
    i32.load
    i32.const 4
    i32.add
    i32.store
    global.get 2
    i32.load
    local.tee 3
    local.get 0
    i32.store
    local.get 3
    local.get 1
    i32.store offset=4
    local.get 3
    local.get 2
    i32.store offset=8
    local.get 3
    local.get 4
    i32.store offset=12
    local.get 3
    local.get 5
    i32.store offset=16
    local.get 3
    local.get 7
    i64.store offset=20 align=4
    local.get 3
    local.get 8
    i32.store offset=28
    local.get 3
    local.get 9
    i32.store offset=32
    local.get 3
    local.get 10
    i32.store offset=36
    local.get 3
    local.get 11
    i64.store offset=40 align=4
    local.get 3
    local.get 12
    i64.store offset=48 align=4
    local.get 3
    local.get 13
    i32.store offset=56
    local.get 3
    local.get 14
    i32.store offset=60
    local.get 3
    local.get 15
    i32.store offset=64
    local.get 3
    local.get 16
    i32.store offset=68
    local.get 3
    local.get 17
    i32.store offset=72
    local.get 3
    local.get 18
    i32.store offset=76
    local.get 3
    local.get 19
    i32.store offset=80
    local.get 3
    local.get 20
    i32.store offset=84
    local.get 3
    local.get 22
    i32.store offset=88
    local.get 3
    local.get 23
    i32.store offset=92
    local.get 3
    local.get 24
    i32.store offset=96
    local.get 3
    local.get 25
    i64.store offset=100 align=4
    global.get 2
    global.get 2
    i32.load
    i32.const 108
    i32.add
    i32.store)
  (func $go_scheduler (type 0)
    (local i32 i32)
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
        local.set 0
      end
      global.get 1
      i32.const 1
      global.get 1
      if (result i32)  ;; label = @2
        local.get 1
      else
        i32.const 66388
        i32.load8_u
        i32.eqz
      end
      select
      if  ;; label = @2
        local.get 0
        i32.const 0
        global.get 1
        select
        i32.eqz
        if  ;; label = @3
          call $runtime.minSched
          i32.const 0
          global.get 1
          i32.const 1
          i32.eq
          br_if 2 (;@1;)
          drop
        end
        global.get 1
        i32.eqz
        if  ;; label = @3
          return
        end
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        i32.const 66388
        i32.const 1
        i32.store8
      end
      local.get 0
      i32.const 1
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        call $runtime.scheduler
        i32.const 1
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        i32.const 66388
        i32.const 0
        i32.store8
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
  (func $_syscall/js.Value_.Type (type 7) (param i64 i32) (result i32)
    (local i32 i32 i32 i32 i32)
    global.get 1
    i32.const 2
    i32.eq
    if  ;; label = @1
      global.get 2
      global.get 2
      i32.load
      i32.const 20
      i32.sub
      i32.store
      global.get 2
      i32.load
      local.tee 3
      i64.load align=4
      local.set 0
      local.get 3
      i32.load offset=8
      local.set 2
      local.get 3
      i32.load offset=12
      local.set 4
      local.get 3
      i32.load offset=16
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
        local.set 5
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        global.get $__stack_pointer
        i32.const 32
        i32.sub
        local.tee 2
        global.set $__stack_pointer
        local.get 2
        i32.const 0
        i32.store offset=28
        local.get 2
        i32.const 2
        i32.store offset=20
        local.get 2
        local.get 0
        i64.store
        local.get 2
        local.get 1
        i32.store offset=8
        i32.const 66380
        i32.load
        local.set 3
        i32.const 66380
        local.get 2
        i32.const 16
        i32.add
        i32.store
        local.get 2
        local.get 3
        i32.store offset=16
        local.get 2
        i32.const 24
        i32.add
        local.get 2
        i32.store
        local.get 0
        i64.eqz
        local.set 6
        i32.const 0
        local.set 4
      end
      block  ;; label = @2
        global.get 1
        i32.eqz
        if  ;; label = @3
          local.get 6
          br_if 1 (;@2;)
          local.get 0
          i64.const 9221120237041090562
          i64.eq
          if  ;; label = @4
            i32.const 1
            local.set 4
            br 2 (;@2;)
          end
          i32.const 2
          local.set 4
          local.get 0
          i64.const 9221120237041090563
          i64.sub
          i64.const 2
          i64.lt_u
          br_if 1 (;@2;)
          local.get 2
          i32.const 28
          i32.add
          local.get 1
          i32.store
          i32.const 3
          local.set 4
          local.get 0
          local.get 1
          call $_syscall/js.Value_.isNumber
          i32.const 1
          i32.and
          br_if 1 (;@2;)
          local.get 0
          i64.const 32
          i64.shr_u
          i64.const 7
          i64.and
          i64.const 1
          i64.sub
          local.tee 0
          i64.const 4
          i64.lt_u
          local.set 1
        end
        global.get 1
        i32.const 1
        local.get 1
        select
        if  ;; label = @3
          local.get 5
          i32.const 0
          global.get 1
          select
          i32.eqz
          if  ;; label = @4
            i32.const 34
            i32.const 66000
            call $runtime._panic
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
            unreachable
          end
        end
        global.get 1
        i32.eqz
        if  ;; label = @3
          local.get 0
          i32.wrap_i64
          i32.const 2
          i32.shl
          i32.const 66168
          i32.add
          i32.load
          local.set 4
        end
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        i32.const 66380
        local.get 3
        i32.store
        local.get 2
        i32.const 32
        i32.add
        global.set $__stack_pointer
        local.get 4
        return
      end
      unreachable
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
    i64.store align=4
    local.get 1
    local.get 2
    i32.store offset=8
    local.get 1
    local.get 4
    i32.store offset=12
    local.get 1
    local.get 3
    i32.store offset=16
    global.get 2
    global.get 2
    i32.load
    i32.const 20
    i32.add
    i32.store
    i32.const 0)
  (func $syscall/js.makeValue (type 19) (param i32 i64)
    (local i32 i32 i32 i32 i32 i32)
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
      local.tee 3
      i32.load
      local.set 0
      local.get 3
      i32.load offset=12
      local.set 2
      local.get 3
      i32.load offset=16
      local.set 4
      local.get 3
      i32.load offset=20
      local.set 5
      local.get 3
      i32.load offset=24
      local.set 6
      local.get 3
      i64.load offset=4 align=4
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
        local.set 7
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        global.get $__stack_pointer
        i32.const 48
        i32.sub
        local.tee 2
        global.set $__stack_pointer
        local.get 2
        i64.const 0
        i64.store offset=36 align=4
        local.get 2
        i64.const 0
        i64.store offset=28 align=4
        local.get 2
        i64.const 5
        i64.store offset=20 align=4
        i32.const 66380
        i32.load
        local.set 6
        i32.const 66380
        local.get 2
        i32.const 16
        i32.add
        i32.store
        local.get 2
        local.get 6
        i32.store offset=16
        local.get 1
        i64.const 9221120237041090560
        i64.and
        i64.const 9221120237041090560
        i64.ne
        local.set 5
        i32.const 0
        local.set 4
      end
      block  ;; label = @2
        global.get 1
        i32.eqz
        if  ;; label = @3
          local.get 5
          br_if 1 (;@2;)
          i32.const 0
          local.set 4
          local.get 1
          i64.const 30064771072
          i64.and
          i64.eqz
          br_if 1 (;@2;)
          local.get 2
          i32.const 16
          i32.add
          local.tee 4
          i32.const 8
          i32.add
          local.set 5
        end
        local.get 7
        i32.const 0
        global.get 1
        select
        i32.eqz
        if  ;; label = @3
          i32.const 8
          call $runtime.alloc
          local.set 3
          i32.const 0
          global.get 1
          i32.const 1
          i32.eq
          br_if 2 (;@1;)
          drop
          local.get 3
          local.set 4
        end
        global.get 1
        i32.eqz
        if  ;; label = @3
          local.get 5
          local.get 4
          i32.store
          local.get 2
          i32.const 28
          i32.add
          local.get 4
          i32.store
          local.get 4
          local.get 1
          i64.store
        end
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        i32.const 66380
        local.get 6
        i32.store
        local.get 2
        i32.const 32
        i32.add
        local.get 4
        i32.store
        local.get 2
        i32.const 40
        i32.add
        local.get 4
        i32.store
        local.get 2
        i32.const 0
        i32.store offset=8
        local.get 2
        i64.const 0
        i64.store
        local.get 2
        i32.const 36
        i32.add
        local.get 2
        i32.store
        local.get 0
        local.get 4
        i32.store offset=8
        local.get 0
        local.get 1
        i64.store
        local.get 2
        i32.const 48
        i32.add
        global.set $__stack_pointer
      end
      return
    end
    local.set 3
    global.get 2
    i32.load
    local.get 3
    i32.store
    global.get 2
    global.get 2
    i32.load
    i32.const 4
    i32.add
    i32.store
    global.get 2
    i32.load
    local.tee 3
    local.get 0
    i32.store
    local.get 3
    local.get 1
    i64.store offset=4 align=4
    local.get 3
    local.get 2
    i32.store offset=12
    local.get 3
    local.get 4
    i32.store offset=16
    local.get 3
    local.get 5
    i32.store offset=20
    local.get 3
    local.get 6
    i32.store offset=24
    global.get 2
    global.get 2
    i32.load
    i32.const 28
    i32.add
    i32.store)
  (func $_syscall/js.Value_.isNumber (type 7) (param i64 i32) (result i32)
    global.get $__stack_pointer
    i32.const 16
    i32.sub
    local.tee 1
    i32.const 0
    i32.store offset=8
    local.get 1
    i64.const 0
    i64.store
    i32.const 1
    local.set 1
    block  ;; label = @1
      local.get 0
      i64.const 9221120237041090560
      i64.sub
      i64.const 2
      i64.ge_u
      if (result i32)  ;; label = @2
        local.get 0
        i64.const 0
        i64.ne
        br_if 1 (;@1;)
        i32.const 0
      else
        local.get 1
      end
      return
    end
    local.get 0
    i64.const 9221120237041090560
    i64.and
    i64.const 9221120237041090560
    i64.ne)
  (func $syscall/js.jsString (type 5) (param i32 i64 i32)
    (local i32 i32 i32 i32 i32 i32 i32)
    global.get 1
    i32.const 2
    i32.eq
    if  ;; label = @1
      global.get 2
      global.get 2
      i32.load
      i32.const 36
      i32.sub
      i32.store
      global.get 2
      i32.load
      local.tee 4
      i32.load
      local.set 0
      local.get 4
      i32.load offset=12
      local.set 2
      local.get 4
      i32.load offset=16
      local.set 3
      local.get 4
      i32.load offset=20
      local.set 5
      local.get 4
      i32.load offset=24
      local.set 6
      local.get 4
      i32.load offset=28
      local.set 7
      local.get 4
      i32.load offset=32
      local.set 9
      local.get 4
      i64.load offset=4 align=4
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
        local.set 8
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        global.get $__stack_pointer
        i32.const 144
        i32.sub
        local.tee 3
        global.set $__stack_pointer
        local.get 3
        i32.const 100
        i32.add
        local.tee 5
        i64.const 0
        i64.store align=4
        local.get 3
        i32.const 92
        i32.add
        local.tee 6
        i64.const 0
        i64.store align=4
        local.get 3
        i64.const 13
        i64.store offset=84 align=4
        local.get 3
        i64.const 0
        i64.store offset=132 align=4
        local.get 3
        i64.const 0
        i64.store offset=124 align=4
        local.get 3
        i64.const 0
        i64.store offset=116 align=4
        local.get 3
        i64.const 0
        i64.store offset=108 align=4
        local.get 3
        i32.const 0
        i32.store offset=32
        local.get 3
        i64.const 0
        i64.store offset=24
        i32.const 66380
        i32.load
        local.set 9
        i32.const 66380
        local.get 3
        i32.const 80
        i32.add
        i32.store
        local.get 3
        local.get 9
        i32.store offset=80
        local.get 3
        i32.const 88
        i32.add
        local.get 3
        i32.const 24
        i32.add
        i32.store
        local.get 3
        i32.const 8
        i32.add
        local.set 7
      end
      local.get 8
      i32.const 0
      global.get 1
      select
      i32.eqz
      if  ;; label = @2
        local.get 7
        local.get 1
        local.get 3
        call $syscall/js.valuePrepareString
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
        local.get 3
        i32.const 48
        i32.add
        local.tee 7
        i64.const 0
        i64.store
        local.get 6
        local.get 2
        i32.store
        local.get 7
        local.get 2
        i32.store
        local.get 5
        local.get 3
        i32.const 40
        i32.add
        i32.store
        local.get 3
        i32.const 96
        i32.add
        local.get 3
        i32.const 40
        i32.add
        i32.store
        local.get 3
        local.get 1
        i64.store offset=40
        local.get 3
        i32.load offset=16
        local.tee 2
        i32.const 0
        i32.lt_s
        local.set 5
      end
      block  ;; label = @2
        global.get 1
        i32.eqz
        if  ;; label = @3
          local.get 5
          br_if 1 (;@2;)
          local.get 3
          i32.const 116
          i32.add
          local.set 5
          local.get 3
          i64.load offset=8
          local.set 1
        end
        local.get 8
        i32.const 1
        i32.eq
        i32.const 1
        global.get 1
        select
        if  ;; label = @3
          local.get 2
          call $runtime.alloc
          local.set 4
          i32.const 1
          global.get 1
          i32.const 1
          i32.eq
          br_if 2 (;@1;)
          drop
          local.get 4
          local.set 6
        end
        global.get 1
        i32.eqz
        if  ;; label = @3
          local.get 5
          local.get 6
          i32.store
          local.get 3
          i32.const 120
          i32.add
          local.get 6
          i32.store
          local.get 3
          i32.const 104
          i32.add
          local.tee 5
          local.get 6
          i32.store
        end
        local.get 8
        i32.const 2
        i32.eq
        i32.const 1
        global.get 1
        select
        if  ;; label = @3
          local.get 1
          local.get 6
          local.get 2
          local.get 2
          local.get 3
          call $syscall/js.valueLoadString
          i32.const 2
          global.get 1
          i32.const 1
          i32.eq
          br_if 2 (;@1;)
          drop
        end
        local.get 8
        i32.const 3
        i32.eq
        i32.const 1
        global.get 1
        select
        if  ;; label = @3
          local.get 1
          local.get 3
          call $syscall/js.finalizeRef
          i32.const 3
          global.get 1
          i32.const 1
          i32.eq
          br_if 2 (;@1;)
          drop
        end
        global.get 1
        i32.eqz
        if  ;; label = @3
          local.get 3
          i32.const 124
          i32.add
          local.get 3
          i32.const 72
          i32.add
          i32.store
          local.get 3
          i32.const 108
          i32.add
          local.tee 5
          local.get 3
          i32.const 56
          i32.add
          i32.store
          local.get 3
          i32.const 0
          i32.store offset=64
          local.get 3
          i64.const 0
          i64.store offset=56
          local.get 3
          i32.const 132
          i32.add
          local.set 7
        end
        local.get 8
        i32.const 4
        i32.eq
        i32.const 1
        global.get 1
        select
        if  ;; label = @3
          local.get 2
          call $runtime.alloc
          local.set 4
          i32.const 4
          global.get 1
          i32.const 1
          i32.eq
          br_if 2 (;@1;)
          drop
          local.get 4
          local.set 5
        end
        global.get 1
        i32.eqz
        if  ;; label = @3
          local.get 7
          local.get 5
          i32.store
          local.get 3
          i32.const 136
          i32.add
          local.get 5
          i32.store
          local.get 3
          i32.const 128
          i32.add
          local.get 5
          i32.store
          local.get 3
          i32.const 112
          i32.add
          local.get 5
          i32.store
          local.get 5
          local.get 6
          local.get 2
          memory.copy
          i32.const 66380
          local.get 9
          i32.store
          local.get 0
          local.get 2
          i32.store offset=4
          local.get 0
          local.get 5
          i32.store
          local.get 3
          i32.const 144
          i32.add
          global.set $__stack_pointer
          return
        end
      end
      local.get 8
      i32.const 5
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        call $runtime.slicePanic
        i32.const 5
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
    local.set 4
    global.get 2
    i32.load
    local.get 4
    i32.store
    global.get 2
    global.get 2
    i32.load
    i32.const 4
    i32.add
    i32.store
    global.get 2
    i32.load
    local.tee 4
    local.get 0
    i32.store
    local.get 4
    local.get 1
    i64.store offset=4 align=4
    local.get 4
    local.get 2
    i32.store offset=12
    local.get 4
    local.get 3
    i32.store offset=16
    local.get 4
    local.get 5
    i32.store offset=20
    local.get 4
    local.get 6
    i32.store offset=24
    local.get 4
    local.get 7
    i32.store offset=28
    local.get 4
    local.get 9
    i32.store offset=32
    global.get 2
    global.get 2
    i32.load
    i32.const 36
    i32.add
    i32.store)
  (func $syscall/js.ValueOf (type 20) (param i32 i32 i32)
    (local i32 i32 i32 i64 i32 i64)
    global.get 1
    i32.const 2
    i32.eq
    if  ;; label = @1
      global.get 2
      global.get 2
      i32.load
      i32.const 32
      i32.sub
      i32.store
      global.get 2
      i32.load
      local.tee 5
      i32.load
      local.set 0
      local.get 5
      i32.load offset=4
      local.set 1
      local.get 5
      i32.load offset=8
      local.set 2
      local.get 5
      i32.load offset=12
      local.set 3
      local.get 5
      i32.load offset=16
      local.set 7
      local.get 5
      i64.load offset=20 align=4
      local.set 6
      local.get 5
      i32.load offset=28
      local.set 5
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
        local.set 4
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        global.get $__stack_pointer
        i32.const 80
        i32.sub
        local.tee 3
        global.set $__stack_pointer
        local.get 3
        i32.const 0
        i32.store offset=76
        local.get 3
        i64.const 2
        i64.store offset=68 align=4
        i32.const 66380
        i32.load
        local.set 5
        i32.const 66380
        local.get 3
        i32.const -64
        i32.sub
        i32.store
        local.get 3
        local.get 5
        i32.store offset=64
        i32.const 0
        local.set 7
        i64.const 9221120237041090562
        local.set 6
      end
      block  ;; label = @2
        block  ;; label = @3
          global.get 1
          i32.eqz
          if  ;; label = @4
            block  ;; label = @5
              block  ;; label = @6
                block  ;; label = @7
                  block  ;; label = @8
                    local.get 1
                    br_table 5 (;@3;) 6 (;@2;) 6 (;@2;) 6 (;@2;) 1 (;@7;) 0 (;@8;)
                  end
                  local.get 1
                  i32.const 22
                  i32.eq
                  br_if 1 (;@6;)
                  local.get 1
                  i32.const 34
                  i32.eq
                  br_if 2 (;@5;)
                  local.get 1
                  i32.const 63
                  i32.ne
                  br_if 5 (;@2;)
                  local.get 3
                  local.get 2
                  i64.load
                  local.tee 6
                  i64.store offset=48
                  local.get 3
                  local.get 2
                  i32.load offset=8
                  local.tee 7
                  i32.store offset=56
                  local.get 3
                  i32.const 76
                  i32.add
                  local.get 3
                  i32.const 48
                  i32.add
                  i32.store
                  br 4 (;@3;)
                end
                local.get 3
                local.get 2
                f64.convert_i32_s
                call $syscall/js.floatValue
                local.get 3
                i32.load offset=8
                local.set 7
                local.get 3
                i64.load
                local.set 6
                br 3 (;@3;)
              end
              local.get 3
              i32.const 16
              i32.add
              local.get 2
              i64.load
              f64.convert_i64_u
              call $syscall/js.floatValue
              local.get 3
              i32.load offset=24
              local.set 7
              local.get 3
              i64.load offset=16
              local.set 6
              br 2 (;@3;)
            end
            local.get 2
            i32.load
            local.set 7
            local.get 2
            i32.load offset=4
            local.set 2
            local.get 3
            i32.const 32
            i32.add
            local.set 1
          end
          local.get 4
          i32.const 0
          global.get 1
          select
          i32.eqz
          if  ;; label = @4
            local.get 7
            local.get 2
            local.get 3
            call $syscall/js.stringVal
            local.set 8
            i32.const 0
            global.get 1
            i32.const 1
            i32.eq
            br_if 3 (;@1;)
            drop
            local.get 8
            local.set 6
          end
          local.get 4
          i32.const 1
          i32.eq
          i32.const 1
          global.get 1
          select
          if  ;; label = @4
            local.get 1
            local.get 6
            call $syscall/js.makeValue
            i32.const 1
            global.get 1
            i32.const 1
            i32.eq
            br_if 3 (;@1;)
            drop
          end
          global.get 1
          i32.eqz
          if  ;; label = @4
            local.get 3
            i32.load offset=40
            local.set 7
            local.get 3
            i64.load offset=32
            local.set 6
          end
        end
        global.get 1
        i32.eqz
        if  ;; label = @3
          i32.const 66380
          local.get 5
          i32.store
          local.get 3
          i32.const 72
          i32.add
          local.get 7
          i32.store
          local.get 0
          local.get 7
          i32.store offset=8
          local.get 0
          local.get 6
          i64.store
          local.get 3
          i32.const 80
          i32.add
          global.set $__stack_pointer
          return
        end
      end
      local.get 4
      i32.const 2
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        i32.const 34
        i32.const 66056
        call $runtime._panic
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
    local.set 4
    global.get 2
    i32.load
    local.get 4
    i32.store
    global.get 2
    global.get 2
    i32.load
    i32.const 4
    i32.add
    i32.store
    global.get 2
    i32.load
    local.tee 4
    local.get 0
    i32.store
    local.get 4
    local.get 1
    i32.store offset=4
    local.get 4
    local.get 2
    i32.store offset=8
    local.get 4
    local.get 3
    i32.store offset=12
    local.get 4
    local.get 7
    i32.store offset=16
    local.get 4
    local.get 6
    i64.store offset=20 align=4
    local.get 4
    local.get 5
    i32.store offset=28
    global.get 2
    global.get 2
    i32.load
    i32.const 32
    i32.add
    i32.store)
  (func $syscall/js.floatValue (type 21) (param i32 f64)
    (local i32 i32)
    global.get $__stack_pointer
    i32.const 32
    i32.sub
    local.tee 2
    local.get 1
    f64.store offset=24
    i32.const 0
    local.set 3
    local.get 0
    block (result i64)  ;; label = @1
      i64.const 9221120237041090561
      local.get 1
      f64.const 0x0p+0 (;=0;)
      f64.eq
      br_if 0 (;@1;)
      drop
      i64.const 9221120237041090560
      local.get 1
      local.get 1
      f64.ne
      br_if 0 (;@1;)
      drop
      i32.const 0
      local.set 3
      local.get 2
      i32.const 0
      i32.store offset=16
      local.get 2
      local.get 1
      f64.store offset=8
      local.get 1
      i64.reinterpret_f64
    end
    i64.store
    local.get 0
    local.get 3
    i32.store offset=8)
  (func $_syscall/js.Value_.Set (type 22) (param i64 i32 i32 i32 i32 i32)
    (local i32 i32 i32 i32 i32 i32 i64)
    global.get 1
    i32.const 2
    i32.eq
    if  ;; label = @1
      global.get 2
      global.get 2
      i32.load
      i32.const 52
      i32.sub
      i32.store
      global.get 2
      i32.load
      local.tee 7
      i64.load align=4
      local.set 0
      local.get 7
      i32.load offset=12
      local.set 2
      local.get 7
      i32.load offset=16
      local.set 3
      local.get 7
      i32.load offset=20
      local.set 4
      local.get 7
      i32.load offset=24
      local.set 5
      local.get 7
      i32.load offset=28
      local.set 6
      local.get 7
      i32.load offset=32
      local.set 8
      local.get 7
      i32.load offset=36
      local.set 9
      local.get 7
      i32.load offset=40
      local.set 11
      local.get 7
      i64.load offset=44 align=4
      local.set 12
      local.get 7
      i32.load offset=8
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
        local.set 10
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        global.get $__stack_pointer
        i32.const 144
        i32.sub
        local.tee 6
        global.set $__stack_pointer
        local.get 6
        i32.const 100
        i32.add
        local.tee 9
        i64.const 0
        i64.store align=4
        local.get 6
        i64.const 0
        i64.store offset=108 align=4
        local.get 6
        i32.const 112
        i32.add
        local.get 1
        i32.store
        local.get 9
        local.get 1
        i32.store
        local.get 6
        i64.const 12
        i64.store offset=92 align=4
        local.get 6
        i32.const 0
        i32.store offset=140
        local.get 6
        i64.const 0
        i64.store offset=132 align=4
        local.get 6
        i64.const 0
        i64.store offset=124 align=4
        local.get 6
        i64.const 0
        i64.store offset=116 align=4
        local.get 6
        i32.const 0
        i32.store offset=32
        local.get 6
        i64.const 0
        i64.store offset=24
        i32.const 66380
        i32.load
        local.set 9
        i32.const 66380
        local.get 6
        i32.const 88
        i32.add
        i32.store
        local.get 6
        local.get 9
        i32.store offset=88
        local.get 6
        i32.const 96
        i32.add
        local.tee 8
        local.get 6
        i32.const 24
        i32.add
        local.tee 11
        i32.store
      end
      local.get 10
      i32.const 0
      global.get 1
      select
      i32.eqz
      if  ;; label = @2
        local.get 0
        local.get 1
        call $_syscall/js.Value_.Type
        local.set 7
        i32.const 0
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
        local.get 7
        local.set 8
      end
      local.get 11
      local.get 8
      i32.const -2
      i32.and
      i32.const 6
      i32.ne
      global.get 1
      select
      local.set 11
      block  ;; label = @2
        global.get 1
        i32.eqz
        if  ;; label = @3
          local.get 11
          br_if 1 (;@2;)
          local.get 6
          i32.const 0
          i32.store offset=48
          local.get 6
          i64.const 0
          i64.store offset=40
          local.get 6
          i32.const 104
          i32.add
          local.get 6
          i32.const 40
          i32.add
          i32.store
          local.get 6
          i32.const 8
          i32.add
          local.set 8
        end
        local.get 10
        i32.const 1
        i32.eq
        i32.const 1
        global.get 1
        select
        if  ;; label = @3
          local.get 8
          local.get 4
          local.get 5
          call $syscall/js.ValueOf
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
          local.get 6
          i32.const 124
          i32.add
          local.get 6
          i32.load offset=16
          local.tee 4
          i32.store
          local.get 6
          i32.const 108
          i32.add
          local.get 4
          i32.store
          local.get 6
          i64.load offset=8
          local.set 12
        end
        local.get 10
        i32.const 2
        i32.eq
        i32.const 1
        global.get 1
        select
        if  ;; label = @3
          local.get 0
          local.get 2
          local.get 3
          local.get 12
          local.get 6
          call $syscall/js.valueSet
          i32.const 2
          global.get 1
          i32.const 1
          i32.eq
          br_if 2 (;@1;)
          drop
        end
        global.get 1
        i32.eqz
        if  ;; label = @3
          local.get 6
          i32.const 80
          i32.add
          local.tee 2
          i64.const 0
          i64.store
          local.get 6
          i32.const -64
          i32.sub
          local.tee 3
          i64.const 0
          i64.store
          i32.const 66380
          local.get 9
          i32.store
          local.get 2
          local.get 1
          i32.store
          local.get 3
          local.get 4
          i32.store
          local.get 6
          i32.const 116
          i32.add
          local.get 6
          i32.const 72
          i32.add
          i32.store
          local.get 6
          i32.const 120
          i32.add
          local.get 6
          i32.const 72
          i32.add
          i32.store
          local.get 6
          local.get 0
          i64.store offset=72
          local.get 6
          i32.const 132
          i32.add
          local.get 6
          i32.const 56
          i32.add
          i32.store
          local.get 6
          i32.const 128
          i32.add
          local.get 6
          i32.const 56
          i32.add
          i32.store
          local.get 6
          local.get 12
          i64.store offset=56
          local.get 6
          i32.const 144
          i32.add
          global.set $__stack_pointer
          return
        end
      end
      local.get 2
      local.get 6
      i32.const 136
      i32.add
      global.get 1
      select
      local.set 2
      local.get 10
      i32.const 3
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        i32.const 12
        call $runtime.alloc
        local.set 7
        i32.const 3
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
        local.get 7
        local.set 1
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        local.get 2
        local.get 1
        i32.store
        local.get 6
        i32.const 140
        i32.add
        local.get 1
        i32.store
        local.get 1
        local.get 8
        i32.store offset=8
        local.get 1
        i32.const 9
        i32.store offset=4
        local.get 1
        i32.const 65915
        i32.store
      end
      local.get 10
      i32.const 4
      i32.eq
      i32.const 1
      global.get 1
      select
      if  ;; label = @2
        i32.const 3045
        local.get 1
        call $runtime._panic
        i32.const 4
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
    local.set 7
    global.get 2
    i32.load
    local.get 7
    i32.store
    global.get 2
    global.get 2
    i32.load
    i32.const 4
    i32.add
    i32.store
    global.get 2
    i32.load
    local.tee 7
    local.get 0
    i64.store align=4
    local.get 7
    local.get 1
    i32.store offset=8
    local.get 7
    local.get 2
    i32.store offset=12
    local.get 7
    local.get 3
    i32.store offset=16
    local.get 7
    local.get 4
    i32.store offset=20
    local.get 7
    local.get 5
    i32.store offset=24
    local.get 7
    local.get 6
    i32.store offset=28
    local.get 7
    local.get 8
    i32.store offset=32
    local.get 7
    local.get 9
    i32.store offset=36
    local.get 7
    local.get 11
    i32.store offset=40
    local.get 7
    local.get 12
    i64.store offset=44 align=4
    global.get 2
    global.get 2
    i32.load
    i32.const 52
    i32.add
    i32.store)
  (func $tinygo_unwind (type 1) (param i32)
    i32.const 66648
    i32.load8_u
    if  ;; label = @1
      call $asyncify_stop_rewind
      i32.const 66648
      i32.const 0
      i32.store8
    else
      local.get 0
      global.get $__stack_pointer
      i32.store offset=4
      local.get 0
      call $asyncify_start_unwind
    end)
  (func $tinygo_launch (type 1) (param i32)
    (local i32)
    global.get $__stack_pointer
    local.set 1
    local.get 0
    i32.load offset=12
    global.set $__stack_pointer
    local.get 0
    i32.load offset=4
    local.get 0
    i32.load
    call_indirect (type 1)
    call $asyncify_stop_unwind
    local.get 1
    global.set $__stack_pointer)
  (func $tinygo_rewind (type 1) (param i32)
    (local i32 i32 i32)
    global.get $__stack_pointer
    local.set 1
    local.get 0
    i32.load offset=12
    global.set $__stack_pointer
    local.get 0
    i32.load offset=4
    local.set 2
    local.get 0
    i32.load
    local.set 3
    i32.const 66648
    i32.const 1
    i32.store8
    local.get 0
    i32.const 8
    i32.add
    call $asyncify_start_rewind
    local.get 2
    local.get 3
    call_indirect (type 1)
    call $asyncify_stop_unwind
    local.get 1
    global.set $__stack_pointer)
  (func $dummy (type 0)
    nop)
  (func $__wasm_call_dtors (type 0)
    call $dummy
    call $dummy)
  (func $malloc.command_export (type 3) (param i32) (result i32)
    (local i32)
    global.get 1
    i32.const 2
    i32.eq
    if  ;; label = @1
      global.get 2
      global.get 2
      i32.load
      i32.const 4
      i32.sub
      i32.store
      global.get 2
      i32.load
      i32.load
      local.set 0
    end
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
        local.get 1
      end
      i32.const 0
      global.get 1
      select
      i32.eqz
      if  ;; label = @2
        local.get 0
        call $malloc
        local.set 1
        i32.const 0
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
        call $__wasm_call_dtors
        local.get 0
        return
      end
      unreachable
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
    local.get 0
    i32.store
    global.get 2
    global.get 2
    i32.load
    i32.const 4
    i32.add
    i32.store
    i32.const 0)
  (func $free.command_export (type 1) (param i32)
    local.get 0
    call $free
    call $__wasm_call_dtors)
  (func $calloc.command_export (type 6) (param i32 i32) (result i32)
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
        local.get 2
      end
      i32.const 0
      global.get 1
      select
      i32.eqz
      if  ;; label = @2
        local.get 0
        local.get 1
        call $calloc
        local.set 2
        i32.const 0
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
        local.get 2
        local.set 0
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        call $__wasm_call_dtors
        local.get 0
        return
      end
      unreachable
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
    i32.store
    i32.const 0)
  (func $realloc.command_export (type 6) (param i32 i32) (result i32)
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
        local.get 2
      end
      i32.const 0
      global.get 1
      select
      i32.eqz
      if  ;; label = @2
        local.get 0
        local.get 1
        call $realloc
        local.set 2
        i32.const 0
        global.get 1
        i32.const 1
        i32.eq
        br_if 1 (;@1;)
        drop
        local.get 2
        local.set 0
      end
      global.get 1
      i32.eqz
      if  ;; label = @2
        call $__wasm_call_dtors
        local.get 0
        return
      end
      unreachable
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
    i32.store
    i32.const 0)
  (func $_start.command_export (type 0)
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
        call $_start
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
        call $__wasm_call_dtors
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
  (func $resume.command_export (type 0)
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
        call $resume
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
        call $__wasm_call_dtors
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
  (func $go_scheduler.command_export (type 0)
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
        call $go_scheduler
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
        call $__wasm_call_dtors
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
  (func $asyncify_start_unwind (type 1) (param i32)
    i32.const 1
    global.set 1
    local.get 0
    global.set 2
    global.get 2
    i32.load
    global.get 2
    i32.load offset=4
    i32.gt_u
    if  ;; label = @1
      unreachable
    end)
  (func $asyncify_stop_unwind (type 0)
    i32.const 0
    global.set 1
    global.get 2
    i32.load
    global.get 2
    i32.load offset=4
    i32.gt_u
    if  ;; label = @1
      unreachable
    end)
  (func $asyncify_start_rewind (type 1) (param i32)
    i32.const 2
    global.set 1
    local.get 0
    global.set 2
    global.get 2
    i32.load
    global.get 2
    i32.load offset=4
    i32.gt_u
    if  ;; label = @1
      unreachable
    end)
  (func $asyncify_stop_rewind (type 0)
    i32.const 0
    global.set 1
    global.get 2
    i32.load
    global.get 2
    i32.load offset=4
    i32.gt_u
    if  ;; label = @1
      unreachable
    end)
  (func $asyncify_get_state (type 8) (result i32)
    global.get 1)
  (table (;0;) 5 5 funcref)
  (memory (;0;) 2)
  (global $__stack_pointer (mut i32) (i32.const 65536))
  (global (;1;) (mut i32) (i32.const 0))
  (global (;2;) (mut i32) (i32.const 0))
  (export "memory" (memory 0))
  (export "malloc" (func $malloc.command_export))
  (export "free" (func $free.command_export))
  (export "calloc" (func $calloc.command_export))
  (export "realloc" (func $realloc.command_export))
  (export "_start" (func $_start.command_export))
  (export "resume" (func $resume.command_export))
  (export "go_scheduler" (func $go_scheduler.command_export))
  (export "asyncify_start_unwind" (func $asyncify_start_unwind))
  (export "asyncify_stop_unwind" (func $asyncify_stop_unwind))
  (export "asyncify_start_rewind" (func $asyncify_start_rewind))
  (export "asyncify_stop_rewind" (func $asyncify_stop_rewind))
  (export "asyncify_get_state" (func $asyncify_get_state))
  (elem (;0;) (i32.const 1) func $runtime.run$1$gowrapper $runtime.resume$1$gowrapper $runtime.memequal $runtime.hash32)
  (data $.rodata (i32.const 65536) "stack overflowout of memorypanic: panic: runtime error: nil pointer dereferenceindex out of rangeslice out of rangenilunreachable\00\00\00\00\00\00\00v\00\01\00\0b\00\00\00sync: unlock of unlocked Mutex\00\00\90\00\01\00\1e\00\00\00ObjectArray_pendingEventidthisargsresultconsolecall to released function\e7\00\01\00\19\00\00\00errorsyscall/js: Value.Call: property  is not a function, got Value.CallValue.GetValue.IndexValue.IntValue.SetIndexValue.Set<undefined><null><boolean: <number: ><symbol><object><function>bad type flag\c3\01\01\00\0d\00\00\00messageJavaScript error: ValueOf: invalid value\00\f1\01\01\00\16\00\00\00undefinednullbooleannumberstringsymbolobjectfunctionbad type\00\00\00\00D\02\01\00\08\00\00\00syscall/js: call of  on bleah\00\00\00\06\00\00\00\04\00\00\00\05\00\00\00\07\00\00\00")
  (data $.data (i32.const 66184) "\bc\02\01\00\00\00\00\00\80\03\01\00\11\b3\db0\00\00\00\00\04\08\01\00\00\00\00\00\03\00\00\00\00\00\00\00\04\00\00\00"))

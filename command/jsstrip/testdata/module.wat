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
 )

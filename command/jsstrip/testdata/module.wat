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
    )


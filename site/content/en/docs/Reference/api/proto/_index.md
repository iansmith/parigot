---
title: Protobuf Schema
date: _2023-08-04
---

These are the protobuf definitions for [parigot
itself]({{< ref "#syscall_v1_syscall-proto" >}})
and the built in services. 
# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [file/v1/file.proto](#file_v1_file-proto)
    - [CloseRequest](#file-v1-CloseRequest)
    - [CloseResponse](#file-v1-CloseResponse)
    - [CreateRequest](#file-v1-CreateRequest)
    - [CreateResponse](#file-v1-CreateResponse)
    - [DeleteRequest](#file-v1-DeleteRequest)
    - [DeleteResponse](#file-v1-DeleteResponse)
    - [FileInfo](#file-v1-FileInfo)
    - [LoadTestDataRequest](#file-v1-LoadTestDataRequest)
    - [LoadTestDataResponse](#file-v1-LoadTestDataResponse)
    - [OpenRequest](#file-v1-OpenRequest)
    - [OpenResponse](#file-v1-OpenResponse)
    - [ReadRequest](#file-v1-ReadRequest)
    - [ReadResponse](#file-v1-ReadResponse)
    - [StatRequest](#file-v1-StatRequest)
    - [StatResponse](#file-v1-StatResponse)
    - [WriteRequest](#file-v1-WriteRequest)
    - [WriteResponse](#file-v1-WriteResponse)
  
    - [FileErr](#file-v1-FileErr)
  
    - [File](#file-v1-File)
  
- [syscall/v1/syscall.proto](#syscall_v1_syscall-proto)
    - [BindMethodRequest](#syscall-v1-BindMethodRequest)
    - [BindMethodResponse](#syscall-v1-BindMethodResponse)
    - [BlockUntilCallRequest](#syscall-v1-BlockUntilCallRequest)
    - [BlockUntilCallResponse](#syscall-v1-BlockUntilCallResponse)
    - [DependencyExistsRequest](#syscall-v1-DependencyExistsRequest)
    - [DependencyExistsResponse](#syscall-v1-DependencyExistsResponse)
    - [DispatchRequest](#syscall-v1-DispatchRequest)
    - [DispatchResponse](#syscall-v1-DispatchResponse)
    - [ExitPair](#syscall-v1-ExitPair)
    - [ExitRequest](#syscall-v1-ExitRequest)
    - [ExitResponse](#syscall-v1-ExitResponse)
    - [ExportRequest](#syscall-v1-ExportRequest)
    - [ExportResponse](#syscall-v1-ExportResponse)
    - [FullyQualifiedService](#syscall-v1-FullyQualifiedService)
    - [HostBinding](#syscall-v1-HostBinding)
    - [LaunchRequest](#syscall-v1-LaunchRequest)
    - [LaunchResponse](#syscall-v1-LaunchResponse)
    - [LocateRequest](#syscall-v1-LocateRequest)
    - [LocateResponse](#syscall-v1-LocateResponse)
    - [MethodBinding](#syscall-v1-MethodBinding)
    - [ReadOneRequest](#syscall-v1-ReadOneRequest)
    - [ReadOneResponse](#syscall-v1-ReadOneResponse)
    - [RegisterRequest](#syscall-v1-RegisterRequest)
    - [RegisterResponse](#syscall-v1-RegisterResponse)
    - [RequireRequest](#syscall-v1-RequireRequest)
    - [RequireResponse](#syscall-v1-RequireResponse)
    - [ResolvedCall](#syscall-v1-ResolvedCall)
    - [ReturnValueRequest](#syscall-v1-ReturnValueRequest)
    - [ReturnValueResponse](#syscall-v1-ReturnValueResponse)
    - [ServiceByIdRequest](#syscall-v1-ServiceByIdRequest)
    - [ServiceByIdResponse](#syscall-v1-ServiceByIdResponse)
    - [ServiceByNameRequest](#syscall-v1-ServiceByNameRequest)
    - [ServiceByNameResponse](#syscall-v1-ServiceByNameResponse)
    - [ServiceMethodCall](#syscall-v1-ServiceMethodCall)
    - [SynchronousExitRequest](#syscall-v1-SynchronousExitRequest)
    - [SynchronousExitResponse](#syscall-v1-SynchronousExitResponse)
  
    - [KernelErr](#syscall-v1-KernelErr)
    - [MethodDirection](#syscall-v1-MethodDirection)
  
- [queue/v1/queue.proto](#queue_v1_queue-proto)
    - [CreateQueueRequest](#queue-v1-CreateQueueRequest)
    - [CreateQueueResponse](#queue-v1-CreateQueueResponse)
    - [DeleteQueueRequest](#queue-v1-DeleteQueueRequest)
    - [DeleteQueueResponse](#queue-v1-DeleteQueueResponse)
    - [LengthRequest](#queue-v1-LengthRequest)
    - [LengthResponse](#queue-v1-LengthResponse)
    - [LocateRequest](#queue-v1-LocateRequest)
    - [LocateResponse](#queue-v1-LocateResponse)
    - [MarkDoneRequest](#queue-v1-MarkDoneRequest)
    - [MarkDoneResponse](#queue-v1-MarkDoneResponse)
    - [QueueMsg](#queue-v1-QueueMsg)
    - [ReceiveRequest](#queue-v1-ReceiveRequest)
    - [ReceiveResponse](#queue-v1-ReceiveResponse)
    - [SendRequest](#queue-v1-SendRequest)
    - [SendResponse](#queue-v1-SendResponse)
  
    - [QueueErr](#queue-v1-QueueErr)
  
    - [Queue](#queue-v1-Queue)
  
- [protosupport/v1/protosupport.proto](#protosupport_v1_protosupport-proto)
    - [IdRaw](#protosupport-v1-IdRaw)
  
    - [File-level Extensions](#protosupport_v1_protosupport-proto-extensions)
    - [File-level Extensions](#protosupport_v1_protosupport-proto-extensions)
    - [File-level Extensions](#protosupport_v1_protosupport-proto-extensions)
  
- [queue/v1/queue.proto](#queue_v1_queue-proto)
    - [CreateQueueRequest](#queue-v1-CreateQueueRequest)
    - [CreateQueueResponse](#queue-v1-CreateQueueResponse)
    - [DeleteQueueRequest](#queue-v1-DeleteQueueRequest)
    - [DeleteQueueResponse](#queue-v1-DeleteQueueResponse)
    - [LengthRequest](#queue-v1-LengthRequest)
    - [LengthResponse](#queue-v1-LengthResponse)
    - [LocateRequest](#queue-v1-LocateRequest)
    - [LocateResponse](#queue-v1-LocateResponse)
    - [MarkDoneRequest](#queue-v1-MarkDoneRequest)
    - [MarkDoneResponse](#queue-v1-MarkDoneResponse)
    - [QueueMsg](#queue-v1-QueueMsg)
    - [ReceiveRequest](#queue-v1-ReceiveRequest)
    - [ReceiveResponse](#queue-v1-ReceiveResponse)
    - [SendRequest](#queue-v1-SendRequest)
    - [SendResponse](#queue-v1-SendResponse)
  
    - [QueueErr](#queue-v1-QueueErr)
  
    - [Queue](#queue-v1-Queue)
  
- [test/v1/test.proto](#test_v1_test-proto)
    - [AddTestSuiteRequest](#test-v1-AddTestSuiteRequest)
    - [AddTestSuiteResponse](#test-v1-AddTestSuiteResponse)
    - [AddTestSuiteResponse.SucceededEntry](#test-v1-AddTestSuiteResponse-SucceededEntry)
    - [ComparisonResult](#test-v1-ComparisonResult)
    - [ExecRequest](#test-v1-ExecRequest)
    - [ExecResponse](#test-v1-ExecResponse)
    - [QueuePayload](#test-v1-QueuePayload)
    - [StartRequest](#test-v1-StartRequest)
    - [StartResponse](#test-v1-StartResponse)
    - [SuiteInfo](#test-v1-SuiteInfo)
    - [SuiteReportRequest](#test-v1-SuiteReportRequest)
    - [SuiteReportRequest.DetailEntry](#test-v1-SuiteReportRequest-DetailEntry)
    - [SuiteReportRequest.TestSkipEntry](#test-v1-SuiteReportRequest-TestSkipEntry)
    - [SuiteReportRequest.TestSuccessEntry](#test-v1-SuiteReportRequest-TestSuccessEntry)
    - [SuiteReportResponse](#test-v1-SuiteReportResponse)
  
    - [TestErr](#test-v1-TestErr)
  
    - [MethodCallSuite](#test-v1-MethodCallSuite)
    - [Test](#test-v1-Test)
    - [UnderTest](#test-v1-UnderTest)
  
- [Scalar Value Types](#scalar-value-types)



<a name="file_v1_file-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## file/v1/file.proto



<a name="file-v1-CloseRequest"></a>

### CloseRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |






<a name="file-v1-CloseResponse"></a>

### CloseResponse
CloseResponse is not empty because it can return an error. However, there is no
action that the receiver of this response can take other than perhaps issuing a warning
to the system operators.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |






<a name="file-v1-CreateRequest"></a>

### CreateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| path | [string](#string) |  |  |
| content | [string](#string) |  |  |






<a name="file-v1-CreateResponse"></a>

### CreateResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| path | [string](#string) |  |  |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| truncated | [bool](#bool) |  |  |






<a name="file-v1-DeleteRequest"></a>

### DeleteRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| path | [string](#string) |  |  |






<a name="file-v1-DeleteResponse"></a>

### DeleteResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| path | [string](#string) |  |  |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |






<a name="file-v1-FileInfo"></a>

### FileInfo
Define the FileInfo struct


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| path | [string](#string) |  |  |
| is_dir | [bool](#bool) |  |  |
| size | [int32](#int32) |  |  |
| create_time | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | creation time |
| mod_time | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | modification time |






<a name="file-v1-LoadTestDataRequest"></a>

### LoadTestDataRequest
LoadTestDataRequest loads the contents of given directory from the _host_ file system into the /app directory
of the test filesystem (in memory).   This is only intended to be use for test code.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| dir_path | [string](#string) |  | path is a path to a directory on the _host_ filesystem that is to be loaded in /app |
| mount_location | [string](#string) |  | where this new file will exist in the in-memory filesystem... this path will be cleaned lexically and then joined to /app. Note that it is possible create paths with this parameter that cannot be opened because of restrictions on the path in open. |
| return_on_fail | [bool](#bool) |  | returnOnFail should be set to true if you do NOT want the normal behavior of using panic on error. If this value is set to true, the paths that cause an error on import are return in the TestDataResponse. |






<a name="file-v1-LoadTestDataResponse"></a>

### LoadTestDataResponse
LoadTestDataResponse contains a list of paths that caused an error during loading. This value is only
returned if the LoadRequest has the returnOnFail set to true.  If LoadDataRequest.return_on_fail is
false since by definition the error_path will be empty.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| error_path | [string](#string) | repeated |  |






<a name="file-v1-OpenRequest"></a>

### OpenRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| path | [string](#string) |  |  |






<a name="file-v1-OpenResponse"></a>

### OpenResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| path | [string](#string) |  |  |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |






<a name="file-v1-ReadRequest"></a>

### ReadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| buf | [bytes](#bytes) |  | Reads up to len(buf) bytes into buf |






<a name="file-v1-ReadResponse"></a>

### ReadResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| num_read | [int32](#int32) |  | The number of bytes read (0 &lt;= num_read &lt;= len(buf)) |






<a name="file-v1-StatRequest"></a>

### StatRequest
StatRequest asks for the information about a file


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| path | [string](#string) |  |  |






<a name="file-v1-StatResponse"></a>

### StatResponse
Use the FileInfo struct in the StatResponse message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| file_info | [FileInfo](#file-v1-FileInfo) |  |  |






<a name="file-v1-WriteRequest"></a>

### WriteRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| buf | [bytes](#bytes) |  | Writes len(buf) bytes from buf |






<a name="file-v1-WriteResponse"></a>

### WriteResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| num_write | [int32](#int32) |  | The number of bytes written from buf (0 &lt;= num_write &lt;= len(buf)) |





 


<a name="file-v1-FileErr"></a>

### FileErr


| Name | Number | Description |
| ---- | ------ | ----------- |
| NoError | 0 | mandatory |
| DispatchError | 1 | mandatory |
| UnmarshalError | 2 | mandatory |
| MarshalError | 3 | mandatory |
| InvalidPathError | 4 | InvalidPathError: Provided path name is not valid based on following rules: 1. The separator should be &#34;/&#34; 2. It should start with specific prefix -&gt; &#39;/parigot/app/&#39; 3. It should not contain any &#34;.&#34; or &#34;..&#34; in the path 4. It should not exceed a specific value for the number (max is 20) of parts in the path 5. It should avoid certain special characters, including: Asterisk (*)					Question mark (?)		Greater than (&gt;) Less than (&lt;)				Pipe symbol (|)			Ampersand (&amp;) Semicolon (;)				Dollar sign ($)			Backtick (`) Double quotation marks (&#34;)	Single quotation mark (&#39;)

Invalid example:

	&#39;/parigot/app/..&#39; -&gt; &#39;..&#39; is not allowed 	&#39;/parigot/app/./&#39; -&gt; &#39;.&#39; is not allowed 	&#39;/parigot/app/foo\bar&#39; -&gt; &#39;\&#39; is not allowed 	&#39;//parigot/app/foo&#39;, &#39;/parigot/app&#39; -&gt; prefix should be &#39;/parigot/app/&#39; |
| AlreadyInUseError | 5 | File status related errors

The file is already being used. |
| NotExistError | 6 | The file/path does not exist |
| FileClosedError | 7 | The file status is CLOSED, cannot be accessed by a read or write request |
| EOFError | 8 | The file is at the end of the file |
| ReadError | 9 | Some error happened during reading a file |
| WriteError | 10 | Some error happened during writing a file |
| OpenError | 11 | Some error happened during opening a file |
| DeleteError | 12 | Some error happened during deleting a file |
| CreateError | 13 | Some error happened during creating a file |
| NoDataFoundError | 14 | No data file found in the directory |
| LargeBufError | 15 | The buffer for the file is too large |
| InternalError | 16 | There are internal issues with the file service |


 

 


<a name="file-v1-File"></a>

### File


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Open | [OpenRequest](#file-v1-OpenRequest) | [OpenResponse](#file-v1-OpenResponse) | Open handles the READ-only operation on a file |
| Create | [CreateRequest](#file-v1-CreateRequest) | [CreateResponse](#file-v1-CreateResponse) | Create handles the WRITE-only operation on a file. It creates or truncates the name file in the path. If the file already exists, it is truncated. If the file does not exist, it is created. |
| Close | [CloseRequest](#file-v1-CloseRequest) | [CloseResponse](#file-v1-CloseResponse) | Close changes the status of a file to &#34;close&#34; |
| LoadTestData | [LoadTestDataRequest](#file-v1-LoadTestDataRequest) | [LoadTestDataResponse](#file-v1-LoadTestDataResponse) | Load does NOT check that the file(s) referred to are reasonable in length, do not contain symlinks, are readable, etc. Don&#39;t allow this call in prod. |
| Read | [ReadRequest](#file-v1-ReadRequest) | [ReadResponse](#file-v1-ReadResponse) |  |
| Write | [WriteRequest](#file-v1-WriteRequest) | [WriteResponse](#file-v1-WriteResponse) |  |
| Delete | [DeleteRequest](#file-v1-DeleteRequest) | [DeleteResponse](#file-v1-DeleteResponse) | Delete free a file from memory (datacache) or delete it from the disk |
| Stat | [StatRequest](#file-v1-StatRequest) | [StatResponse](#file-v1-StatResponse) |  |

 



<a name="syscall_v1_syscall-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## syscall/v1/syscall.proto



<a name="syscall-v1-BindMethodRequest"></a>

### BindMethodRequest
BindMethodRequest is used to tell parigot that the given service_id (located
at host_id) has an implementation for the given method name.  This will create
the mapping to a method_id, which is in the response.  The direction parameter
is either METHOD_DIRECTION_IN, OUT, or BOTH.  IN means thath the method has
no output parameter (the result is ignored), OUT means the method has no input
parameters, and BOTH means that both input and output parameters are used.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| host_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| service_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| method_name | [string](#string) |  |  |
| direction | [MethodDirection](#syscall-v1-MethodDirection) |  |  |






<a name="syscall-v1-BindMethodResponse"></a>

### BindMethodResponse
BindMethodResponse is the method_id of the service and method name provided
in the request.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| method_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |






<a name="syscall-v1-BlockUntilCallRequest"></a>

### BlockUntilCallRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| can_timeout | [bool](#bool) |  |  |






<a name="syscall-v1-BlockUntilCallResponse"></a>

### BlockUntilCallResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| param | [google.protobuf.Any](#google-protobuf-Any) |  |  |
| method | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| call | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| timed_out | [bool](#bool) |  |  |






<a name="syscall-v1-DependencyExistsRequest"></a>

### DependencyExistsRequest
DependencyExistsRequest is used to check if there a dependency
path from source to destination.  Callers should use either
a target service or a target service name, not both. 
The semantics are slightly different.  When you ask about the
name of a service, it is a question about what the service
has declared with require calls.   When you ask about a specific
service you are asking if a dependency path between the two services
exists and thus the dest service must be started before the
source.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| source_service_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| dest_service_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| service_name | [FullyQualifiedService](#syscall-v1-FullyQualifiedService) |  |  |






<a name="syscall-v1-DependencyExistsResponse"></a>

### DependencyExistsResponse
DependencyExistsResponse has the exists field set to 
true if there exists a sequence of dependencies that
join source and dest (from the request).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| exists | [bool](#bool) |  |  |






<a name="syscall-v1-DispatchRequest"></a>

### DispatchRequest
DispatchRequest is a request by a client to invoke a particular method with the parameters provided.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| method_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| param | [google.protobuf.Any](#google-protobuf-Any) |  | inside is another Request object, but we don&#39;t know its type |
| call_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  | reserved for internal use |
| host_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  | reserved for internal use |






<a name="syscall-v1-DispatchResponse"></a>

### DispatchResponse
DispatchResponse sent by the server back to a client.  This what is returned
as the intermediate value to the caller, because the caller
cannot block.  This call_id value can be used on the client side
to map to additional info about the call.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| call_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  | reserved for internal use |






<a name="syscall-v1-ExitPair"></a>

### ExitPair
ExitPair is a structure that is a service that is requesting
an exit and the exit code desired.  The service can be empty
if the caller wants the entire suite of services to be exited.
Code will be in the &#34;allowed&#34; range of 0 to 192.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| code | [int32](#int32) |  |  |






<a name="syscall-v1-ExitRequest"></a>

### ExitRequest
ExitRequest is how you can request for your wasm program, or the whole system
to exit. This will not terminate the process immediately as there may be other 
services running that need to be notified.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pair | [ExitPair](#syscall-v1-ExitPair) |  | For the code in the ExitPair, the valid values here are 0...192 and values&gt;192 or &lt;0 will be set to 192. The valid values for the service are a service id (typically the service making this request) or an zero valued service, indicating that the entire system should be brought down. |
| call_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  | reserved for internal use |
| host_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  | reserved for internal use |
| method_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  | reserved for internal use |






<a name="syscall-v1-ExitResponse"></a>

### ExitResponse
ExitResponse is needed because the exit request does
not cause the shutdown immediately. It causes the 
exit machinery to be invoked at some (soonish) point
in the future.  Note that due to concurrent calls to Exit() the exit
code received may not be the same as the one sent via ExitRequest!
Note that _only_ the caller of Exit() receives this response; if other
processes need to be shutdown, that is handled via SynchExit.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pair | [ExitPair](#syscall-v1-ExitPair) |  |  |






<a name="syscall-v1-ExportRequest"></a>

### ExportRequest
ExportRequest informs the kernel that the given
service id implements the named services on the
host given.  Note that the services provided must be
distinct.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| service | [FullyQualifiedService](#syscall-v1-FullyQualifiedService) | repeated |  |
| host_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |






<a name="syscall-v1-ExportResponse"></a>

### ExportResponse
Nothing to return.






<a name="syscall-v1-FullyQualifiedService"></a>

### FullyQualifiedService
FullyQualified service is the complete (protobuf) name of a service as a
a string.  This is typically something like &#34;foo.v1&#34; for the package and
the service name is &#34;Foo&#34;.  These are the names used by the export, require,
and locate calls.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| package_path | [string](#string) |  |  |
| service | [string](#string) |  |  |






<a name="syscall-v1-HostBinding"></a>

### HostBinding
HostBinding is the mapping between a service and a host. Note that a given
host may be bound to many services, but a single service is always bound
to exactly one host.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| host_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |






<a name="syscall-v1-LaunchRequest"></a>

### LaunchRequest
LaunchRequest is used to block a service until its depnedencies are ready.
It returns a future to the guest that can be used to take action once
launch is completed.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| call_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  | reserved for internal use |
| host_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  | reserved for internal use |
| method_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  | reserved for internal use |






<a name="syscall-v1-LaunchResponse"></a>

### LaunchResponse
LaunchResponse has nothing in it because the action will be handled by
a future created as a result of LaunchRequest.






<a name="syscall-v1-LocateRequest"></a>

### LocateRequest
LocateRequest is a read from the kernel of the service id associated with a package, service pair.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| package_name | [string](#string) |  |  |
| service_name | [string](#string) |  |  |
| called_by | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  | called_by is only needed for true clients. If you are doing a call to locate with a service that you did not and could not have known beforehand you should leave this empty. |






<a name="syscall-v1-LocateResponse"></a>

### LocateResponse
LocateResponse hands back the service Id of the package_name and service_name supplied in the request.
A service id can be thought of as a (network hostname,port) pair that defines which
service&#39;s &#34;location&#34;.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| host_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| service_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| binding | [MethodBinding](#syscall-v1-MethodBinding) | repeated |  |






<a name="syscall-v1-MethodBinding"></a>

### MethodBinding



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| method_name | [string](#string) |  |  |
| method_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |






<a name="syscall-v1-ReadOneRequest"></a>

### ReadOneRequest
ReadOneRequest gives a set of service/method pairs
that should be considered for a read.  The ReadOne
operation returns a single service/method pair that
has received a call. If the timeout expires, only 
the timeout bool is returned. If the timeout value
is 0, then an instanteous sample is returned.  If
the timeout value is negative, it means wait forever.
In addition to potential calls on the service who
requests this read, it is also possible that the
return value represents a completed call from a previous
point in the execution of the calling program.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| call | [ServiceMethodCall](#syscall-v1-ServiceMethodCall) | repeated |  |
| timeout_in_millis | [int32](#int32) |  |  |
| host_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |






<a name="syscall-v1-ReadOneResponse"></a>

### ReadOneResponse
ReadOneResponse is returned when the control
is turned over to parigot for a period of time via
a call to ReadOne.
ReadOneResponse returns timeout = true if a timeout 
has occurred. If timeout is true, all the other fields 
should be ignored. There are two types of results and
these are mutually exclusive.

If resolved is not nil, then this a notification that a
call made by this program have completed.  The
resolved field holds information about the completed call, 
and that data needs to be matched with the appropriate call ids 
and the promises resolved.

If resolved is nil, then the call is a call
on a service and method exposed by this server.  
In that case the pair indicates the method and service
being invoked, and the param and call id should be
to create a matching ReturnValueRequest.

Note that if the method denoted by the pair does not
take input, the value of param should be ignored and
it may be nil.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| timeout | [bool](#bool) |  |  |
| call | [ServiceMethodCall](#syscall-v1-ServiceMethodCall) |  |  |
| call_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| param | [google.protobuf.Any](#google-protobuf-Any) |  |  |
| resolved | [ResolvedCall](#syscall-v1-ResolvedCall) |  |  |
| exit | [bool](#bool) |  |  |






<a name="syscall-v1-RegisterRequest"></a>

### RegisterRequest
Register informs the kernel you are one of the known services
that can be accessed.  Clients use this so they can participate
in the dependency graph for startup order.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| fqs | [FullyQualifiedService](#syscall-v1-FullyQualifiedService) |  |  |
| host_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |






<a name="syscall-v1-RegisterResponse"></a>

### RegisterResponse
RegisterResponse indicates if the registering caller has created
the new service or not.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| existed_previously | [bool](#bool) |  |  |






<a name="syscall-v1-RequireRequest"></a>

### RequireRequest
Require establishes that the source given is going to import the service
given by dest.  It is not required that the source locate the dest, although
if one does call locate, a check is done to insure require was called previously.
This check is done to prevent a common programming mistake.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| dest | [FullyQualifiedService](#syscall-v1-FullyQualifiedService) | repeated |  |
| source | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |






<a name="syscall-v1-RequireResponse"></a>

### RequireResponse
RequireResponse is currently empty.






<a name="syscall-v1-ResolvedCall"></a>

### ResolvedCall
ResolvedCall is used to hold the output of a service/method call while we
are waiting for the future to be resolved.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| host_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| call_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| result | [google.protobuf.Any](#google-protobuf-Any) |  |  |
| result_error | [int32](#int32) |  |  |






<a name="syscall-v1-ReturnValueRequest"></a>

### ReturnValueRequest
ReturnValueRequest is used to return the result of a
function back to the caller. It is the result information
of a call to a service/method function.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| host_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| call_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| result | [google.protobuf.Any](#google-protobuf-Any) |  |  |
| result_error | [int32](#int32) |  |  |






<a name="syscall-v1-ReturnValueResponse"></a>

### ReturnValueResponse
ReturnValueResponse is currently empty.






<a name="syscall-v1-ServiceByIdRequest"></a>

### ServiceByIdRequest
ServiceByIdRequest looks up the given service by
its string representation.  This is probably only
useful for passing service objects over the wire.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_id | [string](#string) |  |  |






<a name="syscall-v1-ServiceByIdResponse"></a>

### ServiceByIdResponse
ServiceByIdResponse returns host binding for the
service or nothing.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| binding | [HostBinding](#syscall-v1-HostBinding) |  |  |






<a name="syscall-v1-ServiceByNameRequest"></a>

### ServiceByNameRequest
ServiceByName looks up the given service and returns all
the host bindings associated with it.   This does
change the internal data structures, only reports on them.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| fqs | [FullyQualifiedService](#syscall-v1-FullyQualifiedService) |  |  |






<a name="syscall-v1-ServiceByNameResponse"></a>

### ServiceByNameResponse
ServiceByNameResponse returns the list, possibly empty,
that has all the host bindings for the named service.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| binding | [HostBinding](#syscall-v1-HostBinding) | repeated |  |






<a name="syscall-v1-ServiceMethodCall"></a>

### ServiceMethodCall
ServiceMethodCall is the structure that holds &#34;what&#39;s been called&#34; in a service.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| method_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |






<a name="syscall-v1-SynchronousExitRequest"></a>

### SynchronousExitRequest
SynchronousExit is sent to a program (a service) that is being told
by the parigot system to run its cleanup (AtExit) handlers because it
is going down.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pair | [ExitPair](#syscall-v1-ExitPair) |  |  |






<a name="syscall-v1-SynchronousExitResponse"></a>

### SynchronousExitResponse
Synchronous exit response is sent to the at exit handlers for a service or
program.  There is no way to stop the shutdown once this is received, it
can be used only to clean up resources that need to be released.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pair | [ExitPair](#syscall-v1-ExitPair) |  |  |





 


<a name="syscall-v1-KernelErr"></a>

### KernelErr


| Name | Number | Description |
| ---- | ------ | ----------- |
| NoError | 0 |  |
| LocateError | 1 | LocateError is return when the kernel cannot find the requested service, given by a package name and service name pair. |
| UnmarshalFailed | 2 | UnmarshalFailed is used to indicate that in unmarshaling a request or result, the protobuf layer returned an error. |
| IdDispatch | 3 | IdDispatch means that a dispatch call failed due to an MethodId or ServiceId was not found. |
| NamespaceExhausted | 4 | NamespaceExhausted is returned when the kernel can no along accept additional packages, services, or methods. This is used primarily to thwart attempts at DOS attacks. |
| NotFound | 5 | NotFound means that a package, service, or method that was requested could not be found. |
| DataTooLarge | 6 | DataTooLarge means that the size of some part of method call was bigger than the buffer allocated to receive it. This could be a problem either on the call side or the return side. |
| MarshalFailed | 7 | Marshal means that a marshal of a protobuf has failed. |
| CallerUnavailable | 8 | CallerUnavailable means that the kernel could not find the original caller that requested the computation for which results have been provided. It is most likely because the caller was killed, exited or timed out. |
| ServiceAlreadyClosedOrExported | 9 | KernelServiceAlreadyClosedOrExported means that some process has already reported the service in question as closed or has already expressed that it is exporting (implementing this service). This is very likely a case where there are two servers that think they are or should be implementing the same service. |
| ServiceAlreadyRequired | 10 | ServiceAlreadyRequired means that this same process has already required the given service. |
| DependencyCycle | 11 | DependencyCycle means that no deterministic startup ordering exists for the set of exports and requires in use. In other words, you must refactor your program so that you do not have a cyle to make it come up cleanly. |
| NetworkFailed | 12 | NetworkFailed means that we successfully connected to the nameserver, but failed during the communication process itself. |
| NetworkConnectionLost | 13 | NetworkConnectionLost means that our internal connection to the remote nameserver was either still working but has lost &#34;sync&#34; in the protocol or the connection has become entirely broken. The kernel will close the connection to remote nameserver and reestablish it after this error. |
| DataTooSmall | 14 | DataTooSmall means that the kernel was speaking some protocol with a remote server, such as a remote nameserver, and data read from the remote said was smaller than the protocol dictated, e.g. it did not contain a checksum after a data block. |
| KernelConnectionFailed | 15 | ConnectionFailed means that the attempt to open a connection to a remote service has failed to connect. |
| NSRetryFailed | 16 | NSRetryFailed means that we tried twice to reach the nameserver with the given request, but both times could not do so. |
| DecodeError | 17 | DecodeError indicates that an attempt to extract a protobuf object from an encoded set of bytes has failed. Typically, this means that the encoder was not called. |
| ExecError | 18 | ExecError means that we received a response from the implenter of a particular service&#39;s function and the execution of that function failed. |
| KernelDependencyFailure | 19 | DependencyFailure means that the dependency infrastructure has failed. This is different than when a user creates bad set of depedencies (KernelDependencyCycle). This an internal to the kernel error. |
| AbortRequest | 20 | AbortRequest indicates that the program that receives this error should exit because the nameserver has asked it to do so. This means that some _other_ program has failed to start correctly, so this deployment cannot succeed. |
| EncodeError | 22 | EncodeError indicates that an attempt encode a protobuf with header and CRC has failed. |
| ClosedErr | 23 | ClosedErr indicates that that object is now closed. This is used as a signal when writing data between the guest and host. |
| GuestReadFailed | 24 | GuestReadFailed indicates that we did not successfully read from guest memory. This is usually caused by the proposed address to read from being out of bounds. |
| GuestWriteFailed | 25 | GuestWriteFailed indicates that we did not successfully write to guest memory. This is usually caused by the proposed address for writing to being out of bounds. |
| BadId | 26 | BadId indicates that you passed the zero value or the empty value of a an id type into a system call. This usually means that you did not properly initialize a protobuf. |
| NotReady | 27 | NotReady that the service that was trying to start was aborted because it returned false from Ready(). Usually this error indicates that the program has no way to continue running. |
| NotRequired | 28 | NotRequired that a service has tried to Locate() another service that that the first service did not Require() previously. |
| RunTimeout | 29 | RunTimeout means that the programs timeout has expired when waiting for all the required dependencies to be fulfilled. |
| ReadOneTimeout | 30 | ReadOneTimeout means that the program was trying to request a service/method pair to invoke, but the request timed out. |
| BadCallId | 31 | BadCallId is returned when trying to match up the results and the call of a function resulting in a promise. It is returned if either there is no such cid registered yet or the cid is already in use. |



<a name="syscall-v1-MethodDirection"></a>

### MethodDirection


| Name | Number | Description |
| ---- | ------ | ----------- |
| METHOD_DIRECTION_UNSPECIFIED | 0 |  |
| METHOD_DIRECTION_IN | 1 |  |
| METHOD_DIRECTION_OUT | 2 |  |
| METHOD_DIRECTION_BOTH | 3 |  |


 

 

 



<a name="queue_v1_queue-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## queue/v1/queue.proto



<a name="queue-v1-CreateQueueRequest"></a>

### CreateQueueRequest
Create creates a queue or returns an error.  Note that this is usually used
only once to set up the operating environment.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| queue_name | [string](#string) |  |  |






<a name="queue-v1-CreateQueueResponse"></a>

### CreateQueueResponse
CreateQueueResponse returns the queue just created.
Errors are passed back out of band.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |






<a name="queue-v1-DeleteQueueRequest"></a>

### DeleteQueueRequest
Delete queue deletes a queue and returns the queue id deleted, or sends
an error out of band.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |






<a name="queue-v1-DeleteQueueResponse"></a>

### DeleteQueueResponse
DeleteQueueResponse returns the (now invalid) queue id of what
was just deleted.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |






<a name="queue-v1-LengthRequest"></a>

### LengthRequest
Length requests and approximation of the number of elements in the queue


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |






<a name="queue-v1-LengthResponse"></a>

### LengthResponse
LengthResponse returns the queue id identifying the queue we 
computed the length for.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| length | [int64](#int64) |  |  |






<a name="queue-v1-LocateRequest"></a>

### LocateRequest
LocateRequest is request to access a given queue.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| queue_name | [string](#string) |  |  |






<a name="queue-v1-LocateResponse"></a>

### LocateResponse
LocateResponse returns the queue id corresponding to the name
provided.  It returns errors out of band.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |






<a name="queue-v1-MarkDoneRequest"></a>

### MarkDoneRequest
MarkDone request indicates that the caller has finished processing
each message in


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| msg | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) | repeated |  |






<a name="queue-v1-MarkDoneResponse"></a>

### MarkDoneResponse
MarkDone returns the list of unmodified (not marked done) messages 
remaining. In the normal case, this will be empty.  If there was an error
trying to mark items as done, it returns the error and 
puts the unmarked elements in the list unmodified


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| unmodified | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) | repeated |  |






<a name="queue-v1-QueueMsg"></a>

### QueueMsg
QueueMsg represents an object returned by a call to Receive.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| msg_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| receive_count | [int32](#int32) |  | ReceiveCount is an approximation of the number of times this messages has been delivered before this delivery. |
| received | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | ReceiveTime is an approximation to the first time the message was received. If the message has never been received before, this will be the zero value. |
| sender | [google.protobuf.Any](#google-protobuf-Any) |  | sender may be any type (or nil) at the discretion of sender |
| sent | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | when the message was sent |
| payload | [google.protobuf.Any](#google-protobuf-Any) |  | payload must be a serialized protobuf object |






<a name="queue-v1-ReceiveRequest"></a>

### ReceiveRequest
Receive pulls the available messages from the queue and returns
them.  Note that if multiple copies of the caller exist, the 
caller must be prepared to receive the same message multiple
times.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| message_limit | [int32](#int32) |  | it is expected that you can process all received messages inside the time limit

1 is usually the right choice here |






<a name="queue-v1-ReceiveResponse"></a>

### ReceiveResponse
Receive response hands the caller a list of messages to 
process. If you need to return an error, do so out of band.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| message | [QueueMsg](#queue-v1-QueueMsg) | repeated |  |






<a name="queue-v1-SendRequest"></a>

### SendRequest
Send requests enqueues the queue messages provided.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| msg | [QueueMsg](#queue-v1-QueueMsg) | repeated |  |






<a name="queue-v1-SendResponse"></a>

### SendResponse
If the queue msg id is an error then we are using the error_detail_msg to
return the value.  Note that the message id you provide here will
changed once we send you the success notification using your id.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| succeed | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) | repeated |  |
| fail | [QueueMsg](#queue-v1-QueueMsg) | repeated |  |
| failed_on | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |





 


<a name="queue-v1-QueueErr"></a>

### QueueErr


| Name | Number | Description |
| ---- | ------ | ----------- |
| NoError | 0 | mandatory |
| DispatchError | 1 | mandatory |
| UnmarshalError | 2 | mandatory |
| MarshalError | 3 | mandatory |
| InvalidName | 4 | InvalidName means that the given queue name is a not a valid identifier. Identifiers must contain only ascii alphanumeric characters and the symbols &#34;.&#34;, &#34;,&#34;,&#34;_&#34; and &#34;-&#34;. The first letter of a queue name must be an alphabetic character. |
| InternalError | 5 | InternalError means that the queue&#39;s implementation (not the values) passed to it) is the problem. This is roughly a 500 not a 401. This is usually caused by a problem with the internal database used to store the queue items. |
| NoPayload | 6 | NoPayload is an error that means that an attempt was made to create a message a nil payload. Payloads are mandatory and senders are optional. |
| NotFound | 7 | NotFound means that the Queue name requested could not be found. This the queue equivalent of 404. |
| AlreadyExists | 8 | AlreadyExists means that the Queue name is already in use. |
| UnmarshalFailed | 9 | Unmarshal error means that we could not use the protobuf unmarshal successfully for a payload or sender. |


 

 


<a name="queue-v1-Queue"></a>

### Queue
Queue supports a reliable source of messages.  Messages
may be delivered out of order or delivered multiple times.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateQueue | [CreateQueueRequest](#queue-v1-CreateQueueRequest) | [CreateQueueResponse](#queue-v1-CreateQueueResponse) | CreateQueue creates a new named queue. This is useful primarily in preparing for a deployment, not during normal execution. See LocateQueue to find an already existing queue. |
| Locate | [LocateRequest](#queue-v1-LocateRequest) | [LocateResponse](#queue-v1-LocateResponse) | Locate finds the named queue and returns the id. |
| DeleteQueue | [DeleteQueueRequest](#queue-v1-DeleteQueueRequest) | [DeleteQueueResponse](#queue-v1-DeleteQueueResponse) | DeleteQueue deletes a named queue. This request will return a specific error code if the queue does not exist. |
| Receive | [ReceiveRequest](#queue-v1-ReceiveRequest) | [ReceiveResponse](#queue-v1-ReceiveResponse) | Receive a queued message. Just receiving a message does not imply that it is fully processed. You need to call delete or the message will be redelivered at a future point. Messages are not guaranteed to be received in the order sent. If there are no messages ready, the response will be returned with a nil message. |
| MarkDone | [MarkDoneRequest](#queue-v1-MarkDoneRequest) | [MarkDoneResponse](#queue-v1-MarkDoneResponse) | Mark a message as done and delete. This should only be done _after_ the processing is completed. If you are worried about idempotentency in your processing,you will need to keep a record of which message Ids you have processed. |
| Length | [LengthRequest](#queue-v1-LengthRequest) | [LengthResponse](#queue-v1-LengthResponse) | Length returns the approximate number of items in the queue. |
| Send | [SendRequest](#queue-v1-SendRequest) | [SendResponse](#queue-v1-SendResponse) | Send a message for later delivery. |

 



<a name="protosupport_v1_protosupport-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## protosupport/v1/protosupport.proto



<a name="protosupport-v1-IdRaw"></a>

### IdRaw



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| high | [uint64](#uint64) |  |  |
| low | [uint64](#uint64) |  |  |





 

 


<a name="protosupport_v1_protosupport-proto-extensions"></a>

### File-level Extensions
| Extension | Type | Base | Number | Description |
| --------- | ---- | ---- | ------ | ----------- |
| parigot_error | bool | .google.protobuf.EnumOptions | 543211 |  |
| host_func_name | string | .google.protobuf.MethodOptions | 543212 |  |
| error_id_name | string | .google.protobuf.ServiceOptions | 543213 |  |

 

 



<a name="queue_v1_queue-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## queue/v1/queue.proto



<a name="queue-v1-CreateQueueRequest"></a>

### CreateQueueRequest
Create creates a queue or returns an error.  Note that this is usually used
only once to set up the operating environment.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| queue_name | [string](#string) |  |  |






<a name="queue-v1-CreateQueueResponse"></a>

### CreateQueueResponse
CreateQueueResponse returns the queue just created.
Errors are passed back out of band.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |






<a name="queue-v1-DeleteQueueRequest"></a>

### DeleteQueueRequest
Delete queue deletes a queue and returns the queue id deleted, or sends
an error out of band.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |






<a name="queue-v1-DeleteQueueResponse"></a>

### DeleteQueueResponse
DeleteQueueResponse returns the (now invalid) queue id of what
was just deleted.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |






<a name="queue-v1-LengthRequest"></a>

### LengthRequest
Length requests and approximation of the number of elements in the queue


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |






<a name="queue-v1-LengthResponse"></a>

### LengthResponse
LengthResponse returns the queue id identifying the queue we 
computed the length for.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| length | [int64](#int64) |  |  |






<a name="queue-v1-LocateRequest"></a>

### LocateRequest
LocateRequest is request to access a given queue.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| queue_name | [string](#string) |  |  |






<a name="queue-v1-LocateResponse"></a>

### LocateResponse
LocateResponse returns the queue id corresponding to the name
provided.  It returns errors out of band.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |






<a name="queue-v1-MarkDoneRequest"></a>

### MarkDoneRequest
MarkDone request indicates that the caller has finished processing
each message in


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| msg | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) | repeated |  |






<a name="queue-v1-MarkDoneResponse"></a>

### MarkDoneResponse
MarkDone returns the list of unmodified (not marked done) messages 
remaining. In the normal case, this will be empty.  If there was an error
trying to mark items as done, it returns the error and 
puts the unmarked elements in the list unmodified


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| unmodified | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) | repeated |  |






<a name="queue-v1-QueueMsg"></a>

### QueueMsg
QueueMsg represents an object returned by a call to Receive.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| msg_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| receive_count | [int32](#int32) |  | ReceiveCount is an approximation of the number of times this messages has been delivered before this delivery. |
| received | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | ReceiveTime is an approximation to the first time the message was received. If the message has never been received before, this will be the zero value. |
| sender | [google.protobuf.Any](#google-protobuf-Any) |  | sender may be any type (or nil) at the discretion of sender |
| sent | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | when the message was sent |
| payload | [google.protobuf.Any](#google-protobuf-Any) |  | payload must be a serialized protobuf object |






<a name="queue-v1-ReceiveRequest"></a>

### ReceiveRequest
Receive pulls the available messages from the queue and returns
them.  Note that if multiple copies of the caller exist, the 
caller must be prepared to receive the same message multiple
times.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| message_limit | [int32](#int32) |  | it is expected that you can process all received messages inside the time limit

1 is usually the right choice here |






<a name="queue-v1-ReceiveResponse"></a>

### ReceiveResponse
Receive response hands the caller a list of messages to 
process. If you need to return an error, do so out of band.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| message | [QueueMsg](#queue-v1-QueueMsg) | repeated |  |






<a name="queue-v1-SendRequest"></a>

### SendRequest
Send requests enqueues the queue messages provided.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| msg | [QueueMsg](#queue-v1-QueueMsg) | repeated |  |






<a name="queue-v1-SendResponse"></a>

### SendResponse
If the queue msg id is an error then we are using the error_detail_msg to
return the value.  Note that the message id you provide here will
changed once we send you the success notification using your id.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| succeed | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) | repeated |  |
| fail | [QueueMsg](#queue-v1-QueueMsg) | repeated |  |
| failed_on | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |





 


<a name="queue-v1-QueueErr"></a>

### QueueErr


| Name | Number | Description |
| ---- | ------ | ----------- |
| NoError | 0 | mandatory |
| DispatchError | 1 | mandatory |
| UnmarshalError | 2 | mandatory |
| MarshalError | 3 | mandatory |
| InvalidName | 4 | InvalidName means that the given queue name is a not a valid identifier. Identifiers must contain only ascii alphanumeric characters and the symbols &#34;.&#34;, &#34;,&#34;,&#34;_&#34; and &#34;-&#34;. The first letter of a queue name must be an alphabetic character. |
| InternalError | 5 | InternalError means that the queue&#39;s implementation (not the values) passed to it) is the problem. This is roughly a 500 not a 401. This is usually caused by a problem with the internal database used to store the queue items. |
| NoPayload | 6 | NoPayload is an error that means that an attempt was made to create a message a nil payload. Payloads are mandatory and senders are optional. |
| NotFound | 7 | NotFound means that the Queue name requested could not be found. This the queue equivalent of 404. |
| AlreadyExists | 8 | AlreadyExists means that the Queue name is already in use. |
| UnmarshalFailed | 9 | Unmarshal error means that we could not use the protobuf unmarshal successfully for a payload or sender. |


 

 


<a name="queue-v1-Queue"></a>

### Queue
Queue supports a reliable source of messages.  Messages
may be delivered out of order or delivered multiple times.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateQueue | [CreateQueueRequest](#queue-v1-CreateQueueRequest) | [CreateQueueResponse](#queue-v1-CreateQueueResponse) | CreateQueue creates a new named queue. This is useful primarily in preparing for a deployment, not during normal execution. See LocateQueue to find an already existing queue. |
| Locate | [LocateRequest](#queue-v1-LocateRequest) | [LocateResponse](#queue-v1-LocateResponse) | Locate finds the named queue and returns the id. |
| DeleteQueue | [DeleteQueueRequest](#queue-v1-DeleteQueueRequest) | [DeleteQueueResponse](#queue-v1-DeleteQueueResponse) | DeleteQueue deletes a named queue. This request will return a specific error code if the queue does not exist. |
| Receive | [ReceiveRequest](#queue-v1-ReceiveRequest) | [ReceiveResponse](#queue-v1-ReceiveResponse) | Receive a queued message. Just receiving a message does not imply that it is fully processed. You need to call delete or the message will be redelivered at a future point. Messages are not guaranteed to be received in the order sent. If there are no messages ready, the response will be returned with a nil message. |
| MarkDone | [MarkDoneRequest](#queue-v1-MarkDoneRequest) | [MarkDoneResponse](#queue-v1-MarkDoneResponse) | Mark a message as done and delete. This should only be done _after_ the processing is completed. If you are worried about idempotentency in your processing,you will need to keep a record of which message Ids you have processed. |
| Length | [LengthRequest](#queue-v1-LengthRequest) | [LengthResponse](#queue-v1-LengthResponse) | Length returns the approximate number of items in the queue. |
| Send | [SendRequest](#queue-v1-SendRequest) | [SendResponse](#queue-v1-SendResponse) | Send a message for later delivery. |

 



<a name="test_v1_test-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## test/v1/test.proto



<a name="test-v1-AddTestSuiteRequest"></a>

### AddTestSuiteRequest
AddTestSuiteRequest adds one or more test suites to the list of available 
suites for the TestService.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| suite | [SuiteInfo](#test-v1-SuiteInfo) | repeated |  |
| exec_package | [string](#string) |  |  |
| exec_service | [string](#string) |  |  |






<a name="test-v1-AddTestSuiteResponse"></a>

### AddTestSuiteResponse
AddTestSuiteResponse contains a map that takes a tuple, written as 
pkg.service.name, and maps it to a boolean to indicate if the given
tuple was added successfully.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| succeeded | [AddTestSuiteResponse.SucceededEntry](#test-v1-AddTestSuiteResponse-SucceededEntry) | repeated |  |






<a name="test-v1-AddTestSuiteResponse-SucceededEntry"></a>

### AddTestSuiteResponse.SucceededEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [bool](#bool) |  |  |






<a name="test-v1-ComparisonResult"></a>

### ComparisonResult
Comparison result describes a single comparison that was done during
a test.  This result is generally optional inside a ExecResult.  The
name field is not the package, service, or function name, it is a name
that can used to narrow down to a single comparison. name can be &#34;&#34;,
as can error_message.  The error_id can be nil. These can be zero
value because they are not crucial to the display of the results, although
it is highly recommended that if the success == false, then one of
error_message or error_id is set.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| error_message | [string](#string) |  |  |
| error_id | [protosupport.v1.IdRaw](#protosupport-v1-IdRaw) |  |  |
| success | [bool](#bool) |  |  |






<a name="test-v1-ExecRequest"></a>

### ExecRequest
ExecRequest is the type that flows _from_ the TestService to the 
package.service.func that is under test.  The package, service, and
name are in the request message because the callee might be doing 
trickery with names (see the map field in the SuiteInfo) and 
thus needs to know what to emulate, dispatch, etc.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| package | [string](#string) |  |  |
| service | [string](#string) |  |  |
| name | [string](#string) |  |  |






<a name="test-v1-ExecResponse"></a>

### ExecResponse
ExecResponse is what an object under test sends back to the TestService
describing the test outcome.  A single package/service/name can have
many comparisons.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| success | [bool](#bool) |  |  |
| skipped | [bool](#bool) |  |  |
| package | [string](#string) |  |  |
| service | [string](#string) |  |  |
| name | [string](#string) |  |  |
| detail | [ComparisonResult](#test-v1-ComparisonResult) | repeated |  |






<a name="test-v1-QueuePayload"></a>

### QueuePayload
QueuePayload is the payload that is sent to the TestService via sending
and receiving items from the queue.  Note that the TestSends these messages
during setup and retreives them in the background once the tests have started.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| func_name | [string](#string) |  |  |






<a name="test-v1-StartRequest"></a>

### StartRequest
StartRequest is what the client should use to start the tests running.
The provided data is exclusive, if filter_suite is provided filter_name
may not be, and vice versa.  Both the suite and name filters may be empty
to request running all tests.  filter_name and filter_suite must be legal
golang regular expressions.  Parallel is currently ignored.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| filter_suite | [string](#string) |  |  |
| filter_name | [string](#string) |  |  |
| parallel | [bool](#bool) |  |  |






<a name="test-v1-StartResponse"></a>

### StartResponse
StartResponse returns the number of tests that will be run, given
the filters provided in StartRequest.  If regex_failed means that one
of the regex fields (filter_suite or filter_name) was not a valid
golang regex.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| regex_failed | [bool](#bool) |  | so common we don&#39;t even call it an error |
| num_test | [int32](#int32) |  |  |






<a name="test-v1-SuiteInfo"></a>

### SuiteInfo
SuiteInfo is used to describe the set of test functions that a 
suite has. The map provided goes from the logical name of the
test (&#34;MyFunc&#34;) to the function the test service will actually 
request (&#34;MyTrickyDispatcher&#34;).   The key and value can be identical in
the simple case.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| package_path | [string](#string) |  |  |
| service | [string](#string) |  |  |
| function_name | [string](#string) | repeated |  |






<a name="test-v1-SuiteReportRequest"></a>

### SuiteReportRequest
SuiteReportRequest is passed to the suite _from_ the TestService and 
contains overall information about the suite&#39;s tests.  The maps have a key
that is the logical function name (the key in the map of SuiteRequest).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| success | [bool](#bool) |  |  |
| num_success | [int32](#int32) |  |  |
| num_failure | [int32](#int32) |  |  |
| num_skip | [int32](#int32) |  |  |
| package | [string](#string) |  |  |
| service | [string](#string) |  |  |
| test_success | [SuiteReportRequest.TestSuccessEntry](#test-v1-SuiteReportRequest-TestSuccessEntry) | repeated |  |
| test_skip | [SuiteReportRequest.TestSkipEntry](#test-v1-SuiteReportRequest-TestSkipEntry) | repeated |  |
| detail | [SuiteReportRequest.DetailEntry](#test-v1-SuiteReportRequest-DetailEntry) | repeated |  |






<a name="test-v1-SuiteReportRequest-DetailEntry"></a>

### SuiteReportRequest.DetailEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [ComparisonResult](#test-v1-ComparisonResult) |  |  |






<a name="test-v1-SuiteReportRequest-TestSkipEntry"></a>

### SuiteReportRequest.TestSkipEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [bool](#bool) |  |  |






<a name="test-v1-SuiteReportRequest-TestSuccessEntry"></a>

### SuiteReportRequest.TestSuccessEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [bool](#bool) |  |  |






<a name="test-v1-SuiteReportResponse"></a>

### SuiteReportResponse
SuiteReportResponse is empty because there is nothing valuable that
can be sent from the suite to the TestService.





 


<a name="test-v1-TestErr"></a>

### TestErr
Error codes

| Name | Number | Description |
| ---- | ------ | ----------- |
| NoError | 0 | mandatory |
| DispatchError | 1 | used by generated code |
| UnmarshalError | 2 | used by generated code |
| ServiceNotFound | 3 | ServiceNotFound means that the service that was supposed to be under test could not be found. |
| Exec | 4 | Exec means that the exec itself (not the thing being execed) has failed. |
| SendFailed | 5 | SendFailed means that the Test code itself could not create the necessary queue entries. |
| Internal | 6 | Internal means that the Test code itself (not the code under test) has had a problem. |
| RegexpFailed | 7 | RegexpFailed means that the regexp provided by the caller did not compile and is not a valid go regexp. |
| Marshal | 8 | Marshal is used to when we cannot marshal arguments into to a protobuf. |
| Queue | 9 | Queue means that the there was internal error with the queue that is used by the Test service. |
| DynamicLocate | 10 | DynamicLocate is returned when we cannot discover the protobuf package and service name pair provided. |


 

 


<a name="test-v1-MethodCallSuite"></a>

### MethodCallSuite


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Exec | [ExecRequest](#test-v1-ExecRequest) | [ExecResponse](#test-v1-ExecResponse) |  |
| SuiteReport | [SuiteReportRequest](#test-v1-SuiteReportRequest) | [SuiteReportResponse](#test-v1-SuiteReportResponse) |  |


<a name="test-v1-Test"></a>

### Test


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| AddTestSuite | [AddTestSuiteRequest](#test-v1-AddTestSuiteRequest) | [AddTestSuiteResponse](#test-v1-AddTestSuiteResponse) | AddTestSuite adds all the elements in the request as suites to the TestService list. If you add a Suite more than once, the function lists are merged--the functions that were already present are retained and the response will show these tests as failures. |
| Start | [StartRequest](#test-v1-StartRequest) | [StartResponse](#test-v1-StartResponse) | Start starts the TestManager running all the known tests, implicitly this includes all suites. If you give a suite filter the entire suite&#39;s tests are dropped if the suite name doesn&#39;t match the filter. For finer granularity you can supple a name filter which walks all the known tests and discards any test that doesn&#39;t match the filter. |


<a name="test-v1-UnderTest"></a>

### UnderTest
UnderTest is the service that services should implement to be &#34;under test&#34; and
testable via Test.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Exec | [ExecRequest](#test-v1-ExecRequest) | [ExecResponse](#test-v1-ExecResponse) | For the default test setup, this is the method that dispatches requests from the TestService to the appropriate test function. |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |


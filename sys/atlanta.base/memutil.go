package sys

import "github.com/iansmith/parigot/lib"

func (s *SysCall) ReadString(structPtr int64, dataOffset uintptr, lenOffset uintptr) string {
	return s.mem.LoadStringWithLen(int32(structPtr)+int32(dataOffset), int32(structPtr)+int32(lenOffset))
}
func (s *SysCall) ReadInt64(structPtr int64, dataOffset uintptr) int64 {
	return s.mem.GetInt64(int32(structPtr) + int32(dataOffset))
}
func (s *SysCall) WriteInt64(structPtr int64, dataOffset uintptr, value int64) {
	s.mem.SetInt64(int32(structPtr)+int32(dataOffset), value)
}

func (s *SysCall) Write64BitPair(structPtr int64, dataOffset uintptr, id lib.Id) {
	derefed := s.mem.GetInt32(int32(structPtr + int64(dataOffset)))
	// write the error info back to client
	s.mem.SetInt64(derefed, int64(id.Low()))
	s.mem.SetInt64(derefed+8, int64(id.High()))
}
func (s *SysCall) Read64BitPair(structPtr int64, dataOffset uintptr) (int64, int64) {
	// s.sysPrint("Read64BitPair1 ", "%x", int32(structPtr+int64(dataOffset)))
	// s.sysPrint("Read64BitPair2 ", "%x", int32(structPtr+int64(dataOffset)+8))
	low := s.mem.GetInt64(int32(structPtr + int64(dataOffset)))
	high := s.mem.GetInt64(int32(structPtr + int64(dataOffset) + 8))
	return low, high
}
func (s *SysCall) ReadSlice(structPtr int64, dataOffset uintptr, lenOffset uintptr) []byte {
	return s.mem.LoadSliceWithLenAddr(int32(structPtr)+int32(dataOffset),
		int32(structPtr)+int32(lenOffset))
}
func (s *SysCall) CopyToPtr(structPtr int64, dataOffset uintptr, content []byte) {
	s.mem.CopyToPtr(int32(structPtr)+int32(dataOffset), content)
}

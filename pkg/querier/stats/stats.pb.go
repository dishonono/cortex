// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: stats.proto

package stats

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_sortkeys "github.com/gogo/protobuf/sortkeys"
	_ "github.com/gogo/protobuf/types"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	_ "github.com/golang/protobuf/ptypes/duration"
	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type Stats struct {
	// The sum of all wall time spent in the querier to execute the query.
	WallTime time.Duration `protobuf:"bytes,1,opt,name=wall_time,json=wallTime,proto3,stdduration" json:"wall_time"`
	// The number of series fetched for the query
	FetchedSeriesCount uint64 `protobuf:"varint,2,opt,name=fetched_series_count,json=fetchedSeriesCount,proto3" json:"fetched_series_count,omitempty"`
	// The number of bytes of the chunks fetched for the query
	FetchedChunkBytes uint64 `protobuf:"varint,3,opt,name=fetched_chunk_bytes,json=fetchedChunkBytes,proto3" json:"fetched_chunk_bytes,omitempty"`
	// The number of bytes of data fetched for the query
	FetchedDataBytes uint64 `protobuf:"varint,4,opt,name=fetched_data_bytes,json=fetchedDataBytes,proto3" json:"fetched_data_bytes,omitempty"`
	// Extra fields to be reported on the stats log
	ExtraFields map[string]string `protobuf:"bytes,5,rep,name=extra_fields,json=extraFields,proto3" json:"extra_fields,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// The number of chunks fetched for the query
	FetchedChunksCount uint64 `protobuf:"varint,6,opt,name=fetched_chunks_count,json=fetchedChunksCount,proto3" json:"fetched_chunks_count,omitempty"`
	// The number of samples fetched for the query
	FetchedSamplesCount uint64 `protobuf:"varint,7,opt,name=fetched_samples_count,json=fetchedSamplesCount,proto3" json:"fetched_samples_count,omitempty"`
	// The limit hit when executing the query
	LimitHit string `protobuf:"bytes,8,opt,name=limit_hit,json=limitHit,proto3" json:"limit_hit,omitempty"`
	// The total number of split queries sent after going through all the middlewares.
	// It includes the number of requests that might be discarded by the queue.
	SplitQueries uint64 `protobuf:"varint,9,opt,name=split_queries,json=splitQueries,proto3" json:"split_queries,omitempty"`
	// The sum of wall time spent in the querier to fetch and merge data from storage.
	QueryStorageWallTime time.Duration `protobuf:"bytes,10,opt,name=query_storage_wall_time,json=queryStorageWallTime,proto3,stdduration" json:"query_storage_wall_time"`
	// The total number of postings touched in store gateway for a specific query.
	// Only successful requests from querier to store gateway are included.
	StoreGatewayTouchedPostingsCount uint64 `protobuf:"varint,11,opt,name=store_gateway_touched_postings_count,json=storeGatewayTouchedPostingsCount,proto3" json:"store_gateway_touched_postings_count,omitempty"`
	// The total size of postings touched in store gateway for a specific query, in bytes.
	// Only successful requests from querier to store gateway are included.
	StoreGatewayTouchedPostingBytes uint64 `protobuf:"varint,12,opt,name=store_gateway_touched_posting_bytes,json=storeGatewayTouchedPostingBytes,proto3" json:"store_gateway_touched_posting_bytes,omitempty"`
}

func (m *Stats) Reset()      { *m = Stats{} }
func (*Stats) ProtoMessage() {}
func (*Stats) Descriptor() ([]byte, []int) {
	return fileDescriptor_b4756a0aec8b9d44, []int{0}
}
func (m *Stats) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Stats) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Stats.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Stats) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Stats.Merge(m, src)
}
func (m *Stats) XXX_Size() int {
	return m.Size()
}
func (m *Stats) XXX_DiscardUnknown() {
	xxx_messageInfo_Stats.DiscardUnknown(m)
}

var xxx_messageInfo_Stats proto.InternalMessageInfo

func (m *Stats) GetWallTime() time.Duration {
	if m != nil {
		return m.WallTime
	}
	return 0
}

func (m *Stats) GetFetchedSeriesCount() uint64 {
	if m != nil {
		return m.FetchedSeriesCount
	}
	return 0
}

func (m *Stats) GetFetchedChunkBytes() uint64 {
	if m != nil {
		return m.FetchedChunkBytes
	}
	return 0
}

func (m *Stats) GetFetchedDataBytes() uint64 {
	if m != nil {
		return m.FetchedDataBytes
	}
	return 0
}

func (m *Stats) GetExtraFields() map[string]string {
	if m != nil {
		return m.ExtraFields
	}
	return nil
}

func (m *Stats) GetFetchedChunksCount() uint64 {
	if m != nil {
		return m.FetchedChunksCount
	}
	return 0
}

func (m *Stats) GetFetchedSamplesCount() uint64 {
	if m != nil {
		return m.FetchedSamplesCount
	}
	return 0
}

func (m *Stats) GetLimitHit() string {
	if m != nil {
		return m.LimitHit
	}
	return ""
}

func (m *Stats) GetSplitQueries() uint64 {
	if m != nil {
		return m.SplitQueries
	}
	return 0
}

func (m *Stats) GetQueryStorageWallTime() time.Duration {
	if m != nil {
		return m.QueryStorageWallTime
	}
	return 0
}

func (m *Stats) GetStoreGatewayTouchedPostingsCount() uint64 {
	if m != nil {
		return m.StoreGatewayTouchedPostingsCount
	}
	return 0
}

func (m *Stats) GetStoreGatewayTouchedPostingBytes() uint64 {
	if m != nil {
		return m.StoreGatewayTouchedPostingBytes
	}
	return 0
}

func init() {
	proto.RegisterType((*Stats)(nil), "stats.Stats")
	proto.RegisterMapType((map[string]string)(nil), "stats.Stats.ExtraFieldsEntry")
}

func init() { proto.RegisterFile("stats.proto", fileDescriptor_b4756a0aec8b9d44) }

var fileDescriptor_b4756a0aec8b9d44 = []byte{
	// 544 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x93, 0xcf, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0xbd, 0xcd, 0x1f, 0xe2, 0x4d, 0x90, 0x82, 0x09, 0xc2, 0x8d, 0xc4, 0x26, 0x50, 0x0e,
	0x39, 0x20, 0x07, 0x85, 0x0b, 0x02, 0x09, 0x55, 0x69, 0x0b, 0x1c, 0x10, 0x82, 0xa4, 0x12, 0x52,
	0x2f, 0xab, 0x4d, 0xb2, 0x71, 0xac, 0x3a, 0xde, 0x60, 0xaf, 0x29, 0xbe, 0x21, 0xf1, 0x02, 0x1c,
	0x79, 0x04, 0x1e, 0x25, 0xc7, 0x1c, 0x7b, 0x2a, 0xc4, 0xb9, 0x70, 0xec, 0x23, 0xa0, 0x9d, 0xb5,
	0x5b, 0xa8, 0x04, 0xe2, 0xe6, 0x9d, 0xdf, 0x37, 0x9f, 0xf6, 0x9b, 0x59, 0xe3, 0x6a, 0x24, 0x99,
	0x8c, 0x9c, 0x45, 0x28, 0xa4, 0xb0, 0x4a, 0x70, 0x68, 0x36, 0x5c, 0xe1, 0x0a, 0xa8, 0x74, 0xd5,
	0x97, 0x86, 0x4d, 0xe2, 0x0a, 0xe1, 0xfa, 0xbc, 0x0b, 0xa7, 0x51, 0x3c, 0xed, 0x4e, 0xe2, 0x90,
	0x49, 0x4f, 0x04, 0x19, 0xdf, 0xbe, 0xca, 0x59, 0x90, 0x68, 0x74, 0xef, 0x73, 0x19, 0x97, 0x86,
	0xca, 0xda, 0xda, 0xc5, 0xe6, 0x09, 0xf3, 0x7d, 0x2a, 0xbd, 0x39, 0xb7, 0x51, 0x1b, 0x75, 0xaa,
	0xbd, 0x6d, 0x47, 0x37, 0x3a, 0x79, 0xa3, 0xb3, 0x9f, 0x19, 0xf7, 0x2b, 0xcb, 0xb3, 0x96, 0xf1,
	0xf5, 0x7b, 0x0b, 0x0d, 0x2a, 0xaa, 0xeb, 0xd0, 0x9b, 0x73, 0xeb, 0x21, 0x6e, 0x4c, 0xb9, 0x1c,
	0xcf, 0xf8, 0x84, 0x46, 0x3c, 0xf4, 0x78, 0x44, 0xc7, 0x22, 0x0e, 0xa4, 0xbd, 0xd5, 0x46, 0x9d,
	0xe2, 0xc0, 0xca, 0xd8, 0x10, 0xd0, 0x9e, 0x22, 0x96, 0x83, 0x6f, 0xe6, 0x1d, 0xe3, 0x59, 0x1c,
	0x1c, 0xd3, 0x51, 0x22, 0x79, 0x64, 0x17, 0xa0, 0xe1, 0x46, 0x86, 0xf6, 0x14, 0xe9, 0x2b, 0x60,
	0x3d, 0xc0, 0xb9, 0x0b, 0x9d, 0x30, 0xc9, 0x32, 0x79, 0x11, 0xe4, 0xf5, 0x8c, 0xec, 0x33, 0xc9,
	0xb4, 0x7a, 0x17, 0xd7, 0xf8, 0x47, 0x19, 0x32, 0x3a, 0xf5, 0xb8, 0x3f, 0x89, 0xec, 0x52, 0xbb,
	0xd0, 0xa9, 0xf6, 0xee, 0x38, 0x7a, 0xae, 0x90, 0xda, 0x39, 0x50, 0x82, 0xe7, 0xc0, 0x0f, 0x02,
	0x19, 0x26, 0x83, 0x2a, 0xbf, 0xac, 0xfc, 0x9e, 0x08, 0xee, 0x97, 0x27, 0x2a, 0xff, 0x91, 0x08,
	0x2e, 0x98, 0x25, 0xea, 0xe1, 0x5b, 0x17, 0x33, 0x60, 0xf3, 0x85, 0x7f, 0x31, 0x84, 0x6b, 0xd0,
	0x92, 0xc7, 0x1d, 0x6a, 0xa6, 0x7b, 0xee, 0x62, 0xd3, 0xf7, 0xe6, 0x9e, 0xa4, 0x33, 0x4f, 0xda,
	0x95, 0x36, 0xea, 0x98, 0xfd, 0xe2, 0xf2, 0x4c, 0x8d, 0x16, 0xca, 0x2f, 0x3d, 0x69, 0xed, 0xe0,
	0xeb, 0xd1, 0xc2, 0xf7, 0x24, 0x7d, 0x1f, 0xc3, 0xf8, 0x6c, 0x13, 0xec, 0x6a, 0x50, 0x7c, 0xab,
	0x6b, 0xd6, 0x11, 0xbe, 0xad, 0x70, 0x42, 0x23, 0x29, 0x42, 0xe6, 0x72, 0x7a, 0xb9, 0x4f, 0xfc,
	0xff, 0xfb, 0x6c, 0x80, 0xc7, 0x50, 0x5b, 0xbc, 0xcb, 0x77, 0xfb, 0x1a, 0xdf, 0x57, 0xae, 0x9c,
	0xba, 0x4c, 0xf2, 0x13, 0x96, 0x50, 0x29, 0x62, 0x48, 0xb9, 0x10, 0x91, 0xf4, 0x02, 0x37, 0x8f,
	0x59, 0x85, 0x7b, 0xb5, 0x41, 0xfb, 0x42, 0x4b, 0x0f, 0xb5, 0xf2, 0x4d, 0x26, 0xd4, 0x99, 0x5f,
	0xe1, 0x9d, 0x7f, 0xfa, 0x65, 0xab, 0xad, 0x81, 0x5d, 0xeb, 0xef, 0x76, 0xb0, 0xe9, 0xe6, 0x33,
	0x5c, 0xbf, 0xba, 0x48, 0xab, 0x8e, 0x0b, 0xc7, 0x3c, 0x81, 0x97, 0x6c, 0x0e, 0xd4, 0xa7, 0xd5,
	0xc0, 0xa5, 0x0f, 0xcc, 0x8f, 0x39, 0x3c, 0x48, 0x73, 0xa0, 0x0f, 0x4f, 0xb6, 0x1e, 0xa3, 0xfe,
	0xd3, 0xd5, 0x9a, 0x18, 0xa7, 0x6b, 0x62, 0x9c, 0xaf, 0x09, 0xfa, 0x94, 0x12, 0xf4, 0x2d, 0x25,
	0x68, 0x99, 0x12, 0xb4, 0x4a, 0x09, 0xfa, 0x91, 0x12, 0xf4, 0x33, 0x25, 0xc6, 0x79, 0x4a, 0xd0,
	0x97, 0x0d, 0x31, 0x56, 0x1b, 0x62, 0x9c, 0x6e, 0x88, 0x71, 0xa4, 0xff, 0xc9, 0x51, 0x19, 0xa6,
	0xf9, 0xe8, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff, 0x16, 0xa4, 0x90, 0x1b, 0xb0, 0x03, 0x00, 0x00,
}

func (this *Stats) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Stats)
	if !ok {
		that2, ok := that.(Stats)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.WallTime != that1.WallTime {
		return false
	}
	if this.FetchedSeriesCount != that1.FetchedSeriesCount {
		return false
	}
	if this.FetchedChunkBytes != that1.FetchedChunkBytes {
		return false
	}
	if this.FetchedDataBytes != that1.FetchedDataBytes {
		return false
	}
	if len(this.ExtraFields) != len(that1.ExtraFields) {
		return false
	}
	for i := range this.ExtraFields {
		if this.ExtraFields[i] != that1.ExtraFields[i] {
			return false
		}
	}
	if this.FetchedChunksCount != that1.FetchedChunksCount {
		return false
	}
	if this.FetchedSamplesCount != that1.FetchedSamplesCount {
		return false
	}
	if this.LimitHit != that1.LimitHit {
		return false
	}
	if this.SplitQueries != that1.SplitQueries {
		return false
	}
	if this.QueryStorageWallTime != that1.QueryStorageWallTime {
		return false
	}
	if this.StoreGatewayTouchedPostingsCount != that1.StoreGatewayTouchedPostingsCount {
		return false
	}
	if this.StoreGatewayTouchedPostingBytes != that1.StoreGatewayTouchedPostingBytes {
		return false
	}
	return true
}
func (this *Stats) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 16)
	s = append(s, "&stats.Stats{")
	s = append(s, "WallTime: "+fmt.Sprintf("%#v", this.WallTime)+",\n")
	s = append(s, "FetchedSeriesCount: "+fmt.Sprintf("%#v", this.FetchedSeriesCount)+",\n")
	s = append(s, "FetchedChunkBytes: "+fmt.Sprintf("%#v", this.FetchedChunkBytes)+",\n")
	s = append(s, "FetchedDataBytes: "+fmt.Sprintf("%#v", this.FetchedDataBytes)+",\n")
	keysForExtraFields := make([]string, 0, len(this.ExtraFields))
	for k, _ := range this.ExtraFields {
		keysForExtraFields = append(keysForExtraFields, k)
	}
	github_com_gogo_protobuf_sortkeys.Strings(keysForExtraFields)
	mapStringForExtraFields := "map[string]string{"
	for _, k := range keysForExtraFields {
		mapStringForExtraFields += fmt.Sprintf("%#v: %#v,", k, this.ExtraFields[k])
	}
	mapStringForExtraFields += "}"
	if this.ExtraFields != nil {
		s = append(s, "ExtraFields: "+mapStringForExtraFields+",\n")
	}
	s = append(s, "FetchedChunksCount: "+fmt.Sprintf("%#v", this.FetchedChunksCount)+",\n")
	s = append(s, "FetchedSamplesCount: "+fmt.Sprintf("%#v", this.FetchedSamplesCount)+",\n")
	s = append(s, "LimitHit: "+fmt.Sprintf("%#v", this.LimitHit)+",\n")
	s = append(s, "SplitQueries: "+fmt.Sprintf("%#v", this.SplitQueries)+",\n")
	s = append(s, "QueryStorageWallTime: "+fmt.Sprintf("%#v", this.QueryStorageWallTime)+",\n")
	s = append(s, "StoreGatewayTouchedPostingsCount: "+fmt.Sprintf("%#v", this.StoreGatewayTouchedPostingsCount)+",\n")
	s = append(s, "StoreGatewayTouchedPostingBytes: "+fmt.Sprintf("%#v", this.StoreGatewayTouchedPostingBytes)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringStats(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *Stats) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Stats) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Stats) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.StoreGatewayTouchedPostingBytes != 0 {
		i = encodeVarintStats(dAtA, i, uint64(m.StoreGatewayTouchedPostingBytes))
		i--
		dAtA[i] = 0x60
	}
	if m.StoreGatewayTouchedPostingsCount != 0 {
		i = encodeVarintStats(dAtA, i, uint64(m.StoreGatewayTouchedPostingsCount))
		i--
		dAtA[i] = 0x58
	}
	n1, err1 := github_com_gogo_protobuf_types.StdDurationMarshalTo(m.QueryStorageWallTime, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdDuration(m.QueryStorageWallTime):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintStats(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x52
	if m.SplitQueries != 0 {
		i = encodeVarintStats(dAtA, i, uint64(m.SplitQueries))
		i--
		dAtA[i] = 0x48
	}
	if len(m.LimitHit) > 0 {
		i -= len(m.LimitHit)
		copy(dAtA[i:], m.LimitHit)
		i = encodeVarintStats(dAtA, i, uint64(len(m.LimitHit)))
		i--
		dAtA[i] = 0x42
	}
	if m.FetchedSamplesCount != 0 {
		i = encodeVarintStats(dAtA, i, uint64(m.FetchedSamplesCount))
		i--
		dAtA[i] = 0x38
	}
	if m.FetchedChunksCount != 0 {
		i = encodeVarintStats(dAtA, i, uint64(m.FetchedChunksCount))
		i--
		dAtA[i] = 0x30
	}
	if len(m.ExtraFields) > 0 {
		for k := range m.ExtraFields {
			v := m.ExtraFields[k]
			baseI := i
			i -= len(v)
			copy(dAtA[i:], v)
			i = encodeVarintStats(dAtA, i, uint64(len(v)))
			i--
			dAtA[i] = 0x12
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintStats(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintStats(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x2a
		}
	}
	if m.FetchedDataBytes != 0 {
		i = encodeVarintStats(dAtA, i, uint64(m.FetchedDataBytes))
		i--
		dAtA[i] = 0x20
	}
	if m.FetchedChunkBytes != 0 {
		i = encodeVarintStats(dAtA, i, uint64(m.FetchedChunkBytes))
		i--
		dAtA[i] = 0x18
	}
	if m.FetchedSeriesCount != 0 {
		i = encodeVarintStats(dAtA, i, uint64(m.FetchedSeriesCount))
		i--
		dAtA[i] = 0x10
	}
	n2, err2 := github_com_gogo_protobuf_types.StdDurationMarshalTo(m.WallTime, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdDuration(m.WallTime):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintStats(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintStats(dAtA []byte, offset int, v uint64) int {
	offset -= sovStats(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Stats) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = github_com_gogo_protobuf_types.SizeOfStdDuration(m.WallTime)
	n += 1 + l + sovStats(uint64(l))
	if m.FetchedSeriesCount != 0 {
		n += 1 + sovStats(uint64(m.FetchedSeriesCount))
	}
	if m.FetchedChunkBytes != 0 {
		n += 1 + sovStats(uint64(m.FetchedChunkBytes))
	}
	if m.FetchedDataBytes != 0 {
		n += 1 + sovStats(uint64(m.FetchedDataBytes))
	}
	if len(m.ExtraFields) > 0 {
		for k, v := range m.ExtraFields {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovStats(uint64(len(k))) + 1 + len(v) + sovStats(uint64(len(v)))
			n += mapEntrySize + 1 + sovStats(uint64(mapEntrySize))
		}
	}
	if m.FetchedChunksCount != 0 {
		n += 1 + sovStats(uint64(m.FetchedChunksCount))
	}
	if m.FetchedSamplesCount != 0 {
		n += 1 + sovStats(uint64(m.FetchedSamplesCount))
	}
	l = len(m.LimitHit)
	if l > 0 {
		n += 1 + l + sovStats(uint64(l))
	}
	if m.SplitQueries != 0 {
		n += 1 + sovStats(uint64(m.SplitQueries))
	}
	l = github_com_gogo_protobuf_types.SizeOfStdDuration(m.QueryStorageWallTime)
	n += 1 + l + sovStats(uint64(l))
	if m.StoreGatewayTouchedPostingsCount != 0 {
		n += 1 + sovStats(uint64(m.StoreGatewayTouchedPostingsCount))
	}
	if m.StoreGatewayTouchedPostingBytes != 0 {
		n += 1 + sovStats(uint64(m.StoreGatewayTouchedPostingBytes))
	}
	return n
}

func sovStats(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozStats(x uint64) (n int) {
	return sovStats(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Stats) String() string {
	if this == nil {
		return "nil"
	}
	keysForExtraFields := make([]string, 0, len(this.ExtraFields))
	for k, _ := range this.ExtraFields {
		keysForExtraFields = append(keysForExtraFields, k)
	}
	github_com_gogo_protobuf_sortkeys.Strings(keysForExtraFields)
	mapStringForExtraFields := "map[string]string{"
	for _, k := range keysForExtraFields {
		mapStringForExtraFields += fmt.Sprintf("%v: %v,", k, this.ExtraFields[k])
	}
	mapStringForExtraFields += "}"
	s := strings.Join([]string{`&Stats{`,
		`WallTime:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.WallTime), "Duration", "duration.Duration", 1), `&`, ``, 1) + `,`,
		`FetchedSeriesCount:` + fmt.Sprintf("%v", this.FetchedSeriesCount) + `,`,
		`FetchedChunkBytes:` + fmt.Sprintf("%v", this.FetchedChunkBytes) + `,`,
		`FetchedDataBytes:` + fmt.Sprintf("%v", this.FetchedDataBytes) + `,`,
		`ExtraFields:` + mapStringForExtraFields + `,`,
		`FetchedChunksCount:` + fmt.Sprintf("%v", this.FetchedChunksCount) + `,`,
		`FetchedSamplesCount:` + fmt.Sprintf("%v", this.FetchedSamplesCount) + `,`,
		`LimitHit:` + fmt.Sprintf("%v", this.LimitHit) + `,`,
		`SplitQueries:` + fmt.Sprintf("%v", this.SplitQueries) + `,`,
		`QueryStorageWallTime:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.QueryStorageWallTime), "Duration", "duration.Duration", 1), `&`, ``, 1) + `,`,
		`StoreGatewayTouchedPostingsCount:` + fmt.Sprintf("%v", this.StoreGatewayTouchedPostingsCount) + `,`,
		`StoreGatewayTouchedPostingBytes:` + fmt.Sprintf("%v", this.StoreGatewayTouchedPostingBytes) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringStats(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *Stats) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStats
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Stats: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Stats: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field WallTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStats
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStats
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthStats
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(&m.WallTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FetchedSeriesCount", wireType)
			}
			m.FetchedSeriesCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStats
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FetchedSeriesCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FetchedChunkBytes", wireType)
			}
			m.FetchedChunkBytes = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStats
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FetchedChunkBytes |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FetchedDataBytes", wireType)
			}
			m.FetchedDataBytes = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStats
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FetchedDataBytes |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExtraFields", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStats
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStats
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthStats
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ExtraFields == nil {
				m.ExtraFields = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowStats
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					wire |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				fieldNum := int32(wire >> 3)
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowStats
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthStats
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthStats
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowStats
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthStats
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue < 0 {
						return ErrInvalidLengthStats
					}
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipStats(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthStats
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.ExtraFields[mapkey] = mapvalue
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FetchedChunksCount", wireType)
			}
			m.FetchedChunksCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStats
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FetchedChunksCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FetchedSamplesCount", wireType)
			}
			m.FetchedSamplesCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStats
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FetchedSamplesCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LimitHit", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStats
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthStats
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStats
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.LimitHit = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SplitQueries", wireType)
			}
			m.SplitQueries = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStats
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SplitQueries |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field QueryStorageWallTime", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStats
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthStats
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthStats
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdDurationUnmarshal(&m.QueryStorageWallTime, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 11:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StoreGatewayTouchedPostingsCount", wireType)
			}
			m.StoreGatewayTouchedPostingsCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStats
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.StoreGatewayTouchedPostingsCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 12:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StoreGatewayTouchedPostingBytes", wireType)
			}
			m.StoreGatewayTouchedPostingBytes = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStats
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.StoreGatewayTouchedPostingBytes |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipStats(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthStats
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthStats
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipStats(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowStats
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowStats
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowStats
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthStats
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthStats
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowStats
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipStats(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthStats
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthStats = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowStats   = fmt.Errorf("proto: integer overflow")
)

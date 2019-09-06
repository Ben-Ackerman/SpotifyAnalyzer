// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api.proto

package api

import (
	fmt "fmt"

	proto "github.com/golang/protobuf/proto"

	math "math"

	context "golang.org/x/net/context"

	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Tracks struct {
	TrackInfo            []*Tracks_TrackInfo `protobuf:"bytes,1,rep,name=trackInfo,proto3" json:"trackInfo,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *Tracks) Reset()         { *m = Tracks{} }
func (m *Tracks) String() string { return proto.CompactTextString(m) }
func (*Tracks) ProtoMessage()    {}
func (*Tracks) Descriptor() ([]byte, []int) {
	return fileDescriptor_api_a309fa96c876db5b, []int{0}
}
func (m *Tracks) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Tracks.Unmarshal(m, b)
}
func (m *Tracks) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Tracks.Marshal(b, m, deterministic)
}
func (dst *Tracks) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Tracks.Merge(dst, src)
}
func (m *Tracks) XXX_Size() int {
	return xxx_messageInfo_Tracks.Size(m)
}
func (m *Tracks) XXX_DiscardUnknown() {
	xxx_messageInfo_Tracks.DiscardUnknown(m)
}

var xxx_messageInfo_Tracks proto.InternalMessageInfo

func (m *Tracks) GetTrackInfo() []*Tracks_TrackInfo {
	if m != nil {
		return m.TrackInfo
	}
	return nil
}

type Tracks_TrackInfo struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Artist               string   `protobuf:"bytes,2,opt,name=artist,proto3" json:"artist,omitempty"`
	GeniusURI            string   `protobuf:"bytes,3,opt,name=geniusURI,proto3" json:"geniusURI,omitempty"`
	Lyrics               string   `protobuf:"bytes,4,opt,name=lyrics,proto3" json:"lyrics,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Tracks_TrackInfo) Reset()         { *m = Tracks_TrackInfo{} }
func (m *Tracks_TrackInfo) String() string { return proto.CompactTextString(m) }
func (*Tracks_TrackInfo) ProtoMessage()    {}
func (*Tracks_TrackInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_api_a309fa96c876db5b, []int{0, 0}
}
func (m *Tracks_TrackInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Tracks_TrackInfo.Unmarshal(m, b)
}
func (m *Tracks_TrackInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Tracks_TrackInfo.Marshal(b, m, deterministic)
}
func (dst *Tracks_TrackInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Tracks_TrackInfo.Merge(dst, src)
}
func (m *Tracks_TrackInfo) XXX_Size() int {
	return xxx_messageInfo_Tracks_TrackInfo.Size(m)
}
func (m *Tracks_TrackInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_Tracks_TrackInfo.DiscardUnknown(m)
}

var xxx_messageInfo_Tracks_TrackInfo proto.InternalMessageInfo

func (m *Tracks_TrackInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Tracks_TrackInfo) GetArtist() string {
	if m != nil {
		return m.Artist
	}
	return ""
}

func (m *Tracks_TrackInfo) GetGeniusURI() string {
	if m != nil {
		return m.GeniusURI
	}
	return ""
}

func (m *Tracks_TrackInfo) GetLyrics() string {
	if m != nil {
		return m.Lyrics
	}
	return ""
}

func init() {
	proto.RegisterType((*Tracks)(nil), "api.Tracks")
	proto.RegisterType((*Tracks_TrackInfo)(nil), "api.Tracks.TrackInfo")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// LyricsClient is the client API for Lyrics service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LyricsClient interface {
	GetLyrics(ctx context.Context, in *Tracks, opts ...grpc.CallOption) (*Tracks, error)
}

type lyricsClient struct {
	cc *grpc.ClientConn
}

func NewLyricsClient(cc *grpc.ClientConn) LyricsClient {
	return &lyricsClient{cc}
}

func (c *lyricsClient) GetLyrics(ctx context.Context, in *Tracks, opts ...grpc.CallOption) (*Tracks, error) {
	out := new(Tracks)
	err := c.cc.Invoke(ctx, "/api.Lyrics/GetLyrics", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LyricsServer is the server API for Lyrics service.
type LyricsServer interface {
	GetLyrics(context.Context, *Tracks) (*Tracks, error)
}

func RegisterLyricsServer(s *grpc.Server, srv LyricsServer) {
	s.RegisterService(&_Lyrics_serviceDesc, srv)
}

func _Lyrics_GetLyrics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Tracks)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LyricsServer).GetLyrics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Lyrics/GetLyrics",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LyricsServer).GetLyrics(ctx, req.(*Tracks))
	}
	return interceptor(ctx, in, info, handler)
}

var _Lyrics_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Lyrics",
	HandlerType: (*LyricsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetLyrics",
			Handler:    _Lyrics_GetLyrics_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}

func init() { proto.RegisterFile("api.proto", fileDescriptor_api_a309fa96c876db5b) }

var fileDescriptor_api_a309fa96c876db5b = []byte{
	// 181 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4c, 0x2c, 0xc8, 0xd4,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4e, 0x2c, 0xc8, 0x54, 0x5a, 0xc3, 0xc8, 0xc5, 0x16,
	0x52, 0x94, 0x98, 0x9c, 0x5d, 0x2c, 0x64, 0xcc, 0xc5, 0x59, 0x02, 0x62, 0x79, 0xe6, 0xa5, 0xe5,
	0x4b, 0x30, 0x2a, 0x30, 0x6b, 0x70, 0x1b, 0x89, 0xea, 0x81, 0x94, 0x43, 0xe4, 0x21, 0x14, 0x48,
	0x32, 0x08, 0xa1, 0x4e, 0x2a, 0x97, 0x8b, 0x13, 0x2e, 0x2e, 0x24, 0xc4, 0xc5, 0x92, 0x97, 0x98,
	0x9b, 0x2a, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x19, 0x04, 0x66, 0x0b, 0x89, 0x71, 0xb1, 0x25, 0x16,
	0x95, 0x64, 0x16, 0x97, 0x48, 0x30, 0x81, 0x45, 0xa1, 0x3c, 0x21, 0x19, 0x2e, 0xce, 0xf4, 0xd4,
	0xbc, 0xcc, 0xd2, 0xe2, 0xd0, 0x20, 0x4f, 0x09, 0x66, 0xb0, 0x14, 0x42, 0x00, 0xa4, 0x2b, 0xa7,
	0xb2, 0x28, 0x33, 0xb9, 0x58, 0x82, 0x05, 0xa2, 0x0b, 0xc2, 0x33, 0x32, 0xe4, 0x62, 0xf3, 0x01,
	0xb3, 0x84, 0xd4, 0xb9, 0x38, 0xdd, 0x53, 0x4b, 0xa0, 0x1c, 0x6e, 0x24, 0x77, 0x4a, 0x21, 0x73,
	0x94, 0x18, 0x92, 0xd8, 0xc0, 0xbe, 0x35, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x62, 0x1e, 0x07,
	0xe9, 0xfa, 0x00, 0x00, 0x00,
}

package testhelloworld

import (
	"context"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/proto"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

var _ proto.Message = &pb.HelloReply{}

func newServer(t *testing.T) *bufconn.Listener {
	lis := bufconn.Listen(1024)
	t.Cleanup(func() {
		if err := lis.Close(); err != nil {
			t.Logf("failed to stop listener: %v", err)
		}
	})
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			t.Log(err)
		}
	}()
	return lis
}

func TestHelloWorld(t *testing.T) {
	lis := newServer(t)
	conn, err := grpc.Dial(
		"bufnet",
		grpc.WithInsecure(),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return lis.Dial()
		}),
	)
	if err != nil {
		t.Fatalf("did not connect %v", err)
	}
	t.Cleanup(func() { conn.Close() })

	c := pb.NewGreeterClient(conn)
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "x"})
	if err != nil {
		t.Fatal(err)
	}
	expResponse := &pb.HelloReply{
		Message: "Hello x",
	}

	// The below code compares a "golden" protobuf response struct with a
	// received reponse from an invoked rpc.
	// This is incorrect - testify will use reflect.DeepEquals but protobuf
	// Messages may set internal fields when passed into the library.
	// protocmpcheck should mark these calls as invalid.

	// These can and should probably be generated but :shrug:
	checkReflectDeepEqual(t, expResponse, r)
	checkAssertEqual(t, expResponse, r)
	checkRequireEqual(t, expResponse, r)
	checkAssertStructEqual(t, expResponse, r)
	checkRequireStructEqual(t, expResponse, r)
}

func checkReflectDeepEqual(t *testing.T, want, got *pb.HelloReply) {
	if !reflect.DeepEqual(want, got) { // want "comparing proto"
		t.Errorf("want: %+v; got: %+v", want, got)
	}
}

func checkRequireEqual(t *testing.T, want, got *pb.HelloReply) {
	require.Equal(t, want, got) // want "comparing proto"

	expSlice := []*pb.HelloReply{want}
	gotSlice := []*pb.HelloReply{got}
	require.Equal(t, expSlice, gotSlice)         // want "comparing proto"
	require.ElementsMatch(t, expSlice, gotSlice) // want "comparing proto"

	expInterfaceSlice := []interface{}{want}
	gotInterfaceSlice := []interface{}{got}
	require.Equal(t, expInterfaceSlice, gotInterfaceSlice)         // want "comparing proto"
	require.ElementsMatch(t, expInterfaceSlice, gotInterfaceSlice) // want "comparing proto"
}

func checkRequireStructEqual(t *testing.T, want, got *pb.HelloReply) {
	require := require.New(t)
	require.EqualValues(want, got) // want "comparing proto"
}

func checkAssertEqual(t *testing.T, exp, r *pb.HelloReply)       {}
func checkAssertStructEqual(t *testing.T, exp, r *pb.HelloReply) {}

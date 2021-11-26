package testhelloworld

import (
	"context"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func newServer(t *testing.T) string {
	tmp := t.TempDir()
	sockPath := filepath.Join(tmp, "listen.sock")
	lis, err := net.Listen("unix", sockPath)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := lis.Close(); err != nil {
			t.Logf("failed to stop listener: %v", err)
		}
	})
	t.Logf("listening at %s", sockPath)
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	go func() {
		err := s.Serve(lis)
		if err != nil {
			t.Log(err)
		}
	}()
	t.Cleanup(s.Stop)
	return sockPath
}

func TestHelloWorld(t *testing.T) {
	sockPath := newServer(t)
	conn, err := grpc.Dial(
		sockPath,
		grpc.WithInsecure(),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)
	if err != nil {
		t.Fatalf("did onot connect %v", err)
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
	checkAssertEqual(t, expResponse, r)
	checkRequireEqual(t, expResponse, r)
	checkAssertStructEqual(t, expResponse, r)
	checkRequireStructEqual(t, expResponse, r)
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

func checkAssertEqual(t *testing.T, exp, r *pb.HelloReply)        {}
func checkAssertStructEqual(t *testing.T, exp, r *pb.HelloReply)  {}
func checkRequireStructEqual(t *testing.T, exp, r *pb.HelloReply) {}

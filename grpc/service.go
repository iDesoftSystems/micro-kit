package service

import (
	"context"
	"net"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

// Options ...
type Options struct {
	Addr    string
	Network string
}

// Service ...
type Service struct {
	opts Options
	ctx  context.Context
	srv  *grpc.Server
}

// NewService returns a new gRRPC Wrap Service
func NewService(ctx context.Context, opts Options) *Service {
	return &Service{
		srv:  grpc.NewServer(),
		opts: opts,
		ctx:  ctx,
	}
}

// Server returns a gRPC server
func (s *Service) Server() *grpc.Server {
	return s.srv
}

// Run the gRPC server
func (s *Service) Run() error {
	conn, err := net.Listen(s.opts.Network, s.opts.Addr)
	if err != nil {
		return err
	}

	defer func() {
		if err := conn.Close(); err != nil {
			glog.Errorf("Failed to close %s %s: %v", s.opts.Network, s.opts.Addr, err)
		}
	}()

	go func() {
		defer s.srv.GracefulStop()
		<-s.ctx.Done()
	}()

	glog.Infof("Starting listening at %s", s.opts.Addr)
	return s.srv.Serve(conn)
}

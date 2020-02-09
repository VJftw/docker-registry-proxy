package cmd

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// NewGRPCServer returns a new GRPC server with the zap and recovery interceptors
func NewGRPCServer() *grpc.Server {
	return grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(Logger),
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
}

// NewGRPCConn returns a new GRPC connection
func NewGRPCConn(host string) (*grpc.ClientConn, error) {
	var grpcOpts []grpc.DialOption
	if viper.GetBool(flagGRPCInsecure) {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(host, grpcOpts...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

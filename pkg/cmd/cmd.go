package cmd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/VJftw/docker-registry-proxy/pkg/plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	flagNetworkAddress = "network_address"
	flagGRPCInsecure   = "grpc_insecure"
	// flagConfigPath     = "config_path"
)

var (
	// ErrUnsupportedScheme is returned by Execute if an invalid scheme is given as the network address
	ErrUnsupportedScheme = errors.New("unsupported net scheme")
)

// BuildVersion represents the built binary version
var BuildVersion = "dev"

var stop chan os.Signal

var (
	// Stopables represents extra closers to close when stopping
	Stopables = []Stoppable{}
	// WaitGroup is used to add processes to wait for and pass to goroutines
	WaitGroup sync.WaitGroup
)

// New returns a new root Cobra cmd with some default flags and graceful stopping
func New(programName string) *cobra.Command {
	stop = make(chan os.Signal)
	rootCmd := &cobra.Command{
		Use:   programName,
		Short: "",
		Long:  "",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// set up interrupts
			signal.Notify(stop, syscall.SIGTERM)
			signal.Notify(stop, syscall.SIGINT)

			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			plugin.Close()
			return nil
		},
		Version: BuildVersion,
	}

	viper.AutomaticEnv()
	viper.SetConfigName(programName)
	viper.SetConfigType("yaml")
	// rootCmd.Flags().String(flagConfigPath, ".", "the path to look for configuration yamls in")
	// viper.BindPFlag(flagConfigPath, rootCmd.Flags().Lookup(flagConfigPath))

	rootCmd.Flags().String(flagNetworkAddress, "tcp://:8080", "the network address to listen on e.g. unix:///var/run/test.sock or tcp://:8080")
	if err := viper.BindPFlag(flagNetworkAddress, rootCmd.Flags().Lookup(flagNetworkAddress)); err != nil {
		HandleErr(err)
	}

	// load plugins and add plugin flags
	err := plugin.LoadPluginsFromConfigSlice(viper.GetStringSlice("plugins"))
	if err != nil {
		HandleErr(err)
	}
	if err := LoadPlugins(rootCmd); err != nil {
		HandleErr(err)
	}
	return rootCmd
}

// Execute runs a given root cmd with the given grpc server
func Execute(
	root *cobra.Command,
	server Server,
	preFunc func() error,
) {
	root.RunE = func(cmd *cobra.Command, args []string) error {
		// if configPath := viper.GetString(flagConfigPath); configPath != "" {
		// 	viper.AddConfigPath(configPath)
		// }
		// if err := viper.ReadInConfig(); err != nil {
		// 	if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		// 		HandleErr(err)
		// 	}
		// }
		Logger.Info("loaded settings", zap.Any("settings", viper.AllSettings()))

		if err := ConfigurePlugins(); err != nil {
			return err
		}

		if preFunc != nil {
			if err := preFunc(); err != nil {
				return err
			}
		}

		for _, stopable := range Stopables {
			WaitGroup.Add(1)
			go stopable.Run(&WaitGroup)
		}

		networkAddress := viper.GetString(flagNetworkAddress)

		uri, err := url.Parse(networkAddress)
		if err != nil {
			return err
		}
		var address string
		switch uri.Scheme {
		case "tcp":
			address = uri.Host
		case "unix":
			address = uri.Path
		default:
			return fmt.Errorf("%s: %w", uri.Scheme, ErrUnsupportedScheme)
		}

		lis, err := net.Listen(uri.Scheme, address)
		if err != nil {
			return err
		}
		Logger.Info("listening on", zap.String("address", networkAddress))
		go func() {
			sig := <-stop
			Logger.Info("stopping", zap.String("signal", sig.String()))
			switch s := server.(type) {
			case *grpc.Server:
				s.GracefulStop()
			case *http.Server:
				if err := s.Shutdown(context.Background()); err != nil {
					Logger.Error("error shutting down server", zap.Error(err))
				}
			}
			lis.Close()
			for _, stopable := range Stopables {
				stopable.Stop()
			}
			Logger.Info("stopped")
			if err := Logger.Sync(); err != nil {
				Logger.Error("error syncing logs", zap.Error(err))
			}
		}()

		return server.Serve(lis)
	}

	HandleErr(root.Execute())
}

// HandleErr handles cmd errors
func HandleErr(err error) {
	if err != nil {
		Logger.Error("error", zap.Error(err))
		plugin.Close()
		os.Exit(1)
	}
}

// Server represents a server type
type Server interface {
	Serve(net.Listener) error
}

// Stoppable represents a stoppable type
type Stoppable interface {
	Stop()
	Run(*sync.WaitGroup)
}

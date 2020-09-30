package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	dockerregistryproxyv1 "github.com/VJftw/docker-registry-proxy/api/proto/v1"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type ConfigurationGRPCPlugin struct {
	plugin.Plugin
	Impl dockerregistryproxyv1.ConfigurationServer
}

func (p *ConfigurationGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	dockerregistryproxyv1.RegisterConfigurationServer(s, p.Impl)
	return nil
}

func (p *ConfigurationGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return dockerregistryproxyv1.NewConfigurationClient(c), nil
}

func MarshalConfigurationValue(t dockerregistryproxyv1.ConfigType, value interface{}) ([]byte, error) {
	switch t {
	case dockerregistryproxyv1.ConfigType_STRING:
		return json.Marshal(value.(string))
	case dockerregistryproxyv1.ConfigType_STRING_SLICE:
		return json.Marshal(value.([]string))
	}
	return nil, fmt.Errorf("unsupported type")
}

func UnmarshalConfigurationValue(t dockerregistryproxyv1.ConfigType, value []byte) (interface{}, error) {
	switch t {
	case dockerregistryproxyv1.ConfigType_STRING:
		var res string
		err := json.Unmarshal(value, &res)
		return res, err
	case dockerregistryproxyv1.ConfigType_STRING_SLICE:
		var res []string
		err := json.Unmarshal(value, &res)
		return res, err
	}
	return nil, fmt.Errorf("unsupported type")
}

func GetStringSliceValue(flag string, req *dockerregistryproxyv1.ConfigureRequest) []string {
	if v, ok := req.Attributes[flag]; ok {
		val, err := UnmarshalConfigurationValue(v.GetAttributeType(), v.GetValue())
		if err != nil {
			log.Printf("error: %s", err)
		}
		log.Printf("configured %s as '%v'", flag, val.([]string))
		return val.([]string)
	}
	return []string{}
}

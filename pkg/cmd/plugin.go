package cmd

import (
	"context"
	"fmt"
	"log"

	v1 "github.com/VJftw/docker-registry-proxy/pkg/genproto/v1"
	"github.com/VJftw/docker-registry-proxy/pkg/plugin"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ConfigurePlugins() error {
	for alias := range plugin.PluginClients {
		client, err := plugin.GetConfigurationClient(alias)
		if err != nil {
			HandleErr(err)
		}
		schema, _ := client.GetConfigurationSchema(context.Background(), &empty.Empty{})
		configureReq := &v1.ConfigureRequest{
			Attributes: map[string]*v1.ConfigurationAttributeValue{},
		}
		for attrName, attrConfig := range schema.GetAttributes() {
			viperAttrName := fmt.Sprintf("%s_%s", alias, attrName)
			marshalledValue, err := plugin.MarshalConfigurationValue(
				attrConfig.GetAttributeType(),
				viper.Get(viperAttrName),
			)
			if err != nil {
				return err
			}
			configureReq.Attributes[attrName] = &v1.ConfigurationAttributeValue{
				AttributeType: attrConfig.GetAttributeType(),
				Value:         marshalledValue,
			}
		}
		if _, err := client.Configure(context.Background(), configureReq); err != nil {
			log.Println(err)
		}
	}
	return nil
}

func LoadPlugins(rootCmd *cobra.Command) error {
	for alias := range plugin.PluginClients {
		client, err := plugin.GetConfigurationClient(alias)
		if err != nil {
			return err
		}
		schema, err := client.GetConfigurationSchema(context.Background(), &empty.Empty{})
		if err != nil {
			return err
		}
		for attrName, attrConfig := range schema.GetAttributes() {
			attrName = fmt.Sprintf("%s_%s", alias, attrName)
			switch attrConfig.GetAttributeType() {
			case v1.ConfigType_STRING:
				rootCmd.Flags().String(attrName, "", fmt.Sprintf("%s for %s plugin", attrConfig.GetDescription(), alias))
			case v1.ConfigType_STRING_SLICE:
				rootCmd.Flags().StringSlice(attrName, []string{}, fmt.Sprintf("%s for %s plugin", attrConfig.GetDescription(), alias))
			default:
				return fmt.Errorf("unsupported type: %s", attrConfig.GetAttributeType())
			}
			viper.BindPFlag(attrName, rootCmd.Flags().Lookup(attrName))
		}
	}
	return nil
}

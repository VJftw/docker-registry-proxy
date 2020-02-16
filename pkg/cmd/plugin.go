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
				GetViperAttr(attrConfig.GetAttributeType(), viperAttrName),
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
			if err := viper.BindPFlag(attrName, rootCmd.Flags().Lookup(attrName)); err != nil {
				return fmt.Errorf("could not bind pflag: %w", err)
			}
		}
	}
	return nil
}

func GetViperAttr(attrType v1.ConfigType, attrName string) interface{} {
	switch attrType {
	case v1.ConfigType_STRING:
		return viper.GetString(attrName)
	case v1.ConfigType_STRING_SLICE:
		return viper.GetStringSlice(attrName)
	case v1.ConfigType_BOOL:
		return viper.GetBool(attrName)
	}
	fmt.Println("unsupported type: %s", attrType)
	return nil
}

package cmd

import (
	"context"
	"fmt"
	"log"

	dockerregistryproxyv1 "github.com/VJftw/docker-registry-proxy/api/proto/v1"
	"github.com/VJftw/docker-registry-proxy/pkg/plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ConfigurePlugins() error {
	for alias := range plugin.PluginClients {
		client, err := plugin.GetConfigurationAPIClient(alias)
		if err != nil {
			HandleErr(err)
		}
		schema, _ := client.GetConfigurationSchema(context.Background(), &dockerregistryproxyv1.GetConfigurationSchemaRequest{})
		configureReq := &dockerregistryproxyv1.ConfigureRequest{
			Attributes: map[string]*dockerregistryproxyv1.ConfigurationAttributeValue{},
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
			configureReq.Attributes[attrName] = &dockerregistryproxyv1.ConfigurationAttributeValue{
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
		client, err := plugin.GetConfigurationAPIClient(alias)
		if err != nil {
			return err
		}
		schema, err := client.GetConfigurationSchema(context.Background(), &dockerregistryproxyv1.GetConfigurationSchemaRequest{})
		if err != nil {
			return err
		}
		for attrName, attrConfig := range schema.GetAttributes() {
			attrName = fmt.Sprintf("%s_%s", alias, attrName)
			switch attrConfig.GetAttributeType() {
			case dockerregistryproxyv1.ConfigType_CONFIG_TYPE_STRING:
				rootCmd.Flags().String(attrName, "", fmt.Sprintf("%s for %s plugin", attrConfig.GetDescription(), alias))
			case dockerregistryproxyv1.ConfigType_CONFIG_TYPE_STRING_SLICE:
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

func GetViperAttr(attrType dockerregistryproxyv1.ConfigType, attrName string) interface{} {
	switch attrType {
	case dockerregistryproxyv1.ConfigType_CONFIG_TYPE_STRING:
		return viper.GetString(attrName)
	case dockerregistryproxyv1.ConfigType_CONFIG_TYPE_STRING_SLICE:
		return viper.GetStringSlice(attrName)
	case dockerregistryproxyv1.ConfigType_CONFIG_TYPE_BOOL:
		return viper.GetBool(attrName)
	}
	fmt.Printf("unsupported type: %s\n", attrType)
	return nil
}

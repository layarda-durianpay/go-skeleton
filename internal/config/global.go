package config

import (
	"context"
	"sync"

	"github.com/durianpay/dpay-common/logger"
	consul "github.com/durianpay/dpay-consul"
)

var (
	gc     GlobalConfig
	onceGC sync.Once
)

// ProvideGlobalConfig will provide GlobalConfig in singleton
func ProvideGlobalConfig() GlobalConfig {
	onceGC.Do(func() {
		gc = NewGlobalConfig()
	})

	return gc
}

func NewGlobalConfig() GlobalConfig {
	var globalConfig globalServiceConfig

	staticGlobalEnv, err := consul.InitConsul("CONSUL_GLOBAL_CONFIG_PATH", false)
	if err != nil {
		logger.Errorw(context.Background(), "error initializing consul global variables", "error", err.Error())
		panic("error initializing consul")
	}

	dynamicGlobalEnv, err := consul.InitConsul("CONSUL_GLOBAL_CONFIG_PATH", true)
	if err != nil {
		logger.Errorw(context.Background(), "error initializing consul global variables", "error", err.Error())
		panic("error initializing consul")
	}

	variables := []consul.ConfigVariable{
		// dynamic var
		{
			Field:  &globalConfig.globalDynamicConfig,
			Key:    "DISBURSEMENT_DYNAMIC_CONFIG",
			Source: dynamicGlobalEnv,
		},

		// static var
		{
			Field:  &globalConfig.globalStaticConfig,
			Key:    "DISBURSEMENT_STATIC_CONFIG",
			Source: staticGlobalEnv,
		},
		{
			Field:  &globalConfig.debugPortForServer,
			Key:    "DEBUG_PORT_SERVER",
			Source: staticGlobalEnv,
		},
		{
			Field:  &globalConfig.grpcMaxConnectionAge,
			Key:    "GRPC_MAX_CONNECTION_AGE",
			Source: staticGlobalEnv,
		},
		{
			Field:  &globalConfig.merchantServiceGRPCAddr,
			Key:    "MERCHANT_SERVICE_GRPC_ADDR",
			Source: staticGlobalEnv,
		},
	}

	req := consul.InitVarRequest{
		ServiceStruct:    globalConfig,
		ServiceVariables: variables,
	}

	consul.InitConsulVariables(req)
	return &globalConfig
}

type globalServiceConfig struct {
	globalDynamicConfig string
	globalStaticConfig  string

	debugPortForServer int

	grpcMaxConnectionAge    int
	merchantServiceGRPCAddr string
}

func (c globalServiceConfig) GetGlobalDynamicConfig() string {
	return c.globalDynamicConfig
}

func (c globalServiceConfig) GetGlobalStaticConfig() string {
	return c.globalStaticConfig
}

func (c globalServiceConfig) GetMerchantServiceGRPCAddr() string {
	return c.merchantServiceGRPCAddr
}

func (c globalServiceConfig) GetDebugPortForServer() int {
	return c.debugPortForServer
}

func (c globalServiceConfig) GetGRPCMaxConnectionAge() int {
	return c.debugPortForServer
}

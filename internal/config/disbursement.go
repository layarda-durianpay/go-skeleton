package config

import (
	"context"
	"sync"

	"github.com/durianpay/dpay-common/logger"
	consul "github.com/durianpay/dpay-consul"
)

var (
	dc     DisbursementServiceConfig
	onceDC sync.Once
)

// ProvideDisbursementConfig will provide DisbursementServiceConfig in singleton
func ProvideDisbursementConfig() DisbursementServiceConfig {
	onceDC.Do(func() {
		dc = NewDisbursementServiceConfig()
	})

	return dc
}

func NewDisbursementServiceConfig() DisbursementServiceConfig {
	var disbursementConfig disbursementServiceConfig
	staticEnv, err := consul.InitConsul("CONSUL_DISBURSEMENT_SERVICE_CONFIG_PATH", false)
	if err != nil {
		logger.Errorw(context.Background(), "error initializing consul", "error", err.Error())
		panic("error initializing consul")
	}

	dynamicEnv, err := consul.InitConsul("CONSUL_DISBURSEMENT_SERVICE_CONFIG_PATH", true)
	if err != nil {
		logger.Errorw(context.Background(), "error initializing consul", "error", err.Error())
		panic("error initializing consul")
	}

	variables := []consul.ConfigVariable{
		{
			Field:  &disbursementConfig.disbursementDynamicConfig,
			Key:    "DISBURSEMENT_DYNAMIC_CONFIG",
			Source: dynamicEnv,
		},

		{
			Field:  &disbursementConfig.disbursementStaticConfig,
			Key:    "DISBURSEMENT_STATIC_CONFIG",
			Source: staticEnv,
		},
		{
			Field:  &disbursementConfig.enableConfigOpenTelemetry,
			Key:    "CONFIG_ENABLE_OPENTELEMETRY",
			Source: staticEnv,
		},
		{
			Field:  &disbursementConfig.enableConfigAllowTruncateAttributesOtel,
			Key:    "CONFIG_ENABLE_TRUNCATE_OPENTELEMTRY_ATTRIBUTE",
			Source: staticEnv,
		},
		{
			Field:  &disbursementConfig.startDebugServer,
			Key:    "START_DEBUG_SERVER",
			Source: staticEnv,
		},
		{
			Field:  &disbursementConfig.disbursementKafkaTopic,
			Key:    "DISBURSEMENT_KAFKA_TOPIC",
			Source: staticEnv,
		},
	}

	req := consul.InitVarRequest{
		ServiceStruct:    disbursementConfig,
		ServiceVariables: variables,
	}

	consul.InitConsulVariables(req)
	return &disbursementConfig
}

type disbursementServiceConfig struct {
	// sample
	disbursementDynamicConfig string
	disbursementStaticConfig  string
	startDebugServer          bool

	enableConfigOpenTelemetry               bool
	enableConfigAllowTruncateAttributesOtel bool

	// topic
	disbursementKafkaTopic string
}

func (c disbursementServiceConfig) GetDisbursementDynamicConfig() string {
	return c.disbursementDynamicConfig
}

func (c disbursementServiceConfig) GetDisbursementStaticConfig() string {
	return c.disbursementStaticConfig
}

func (c disbursementServiceConfig) GetEnableConfigOpenTelemetry() bool {
	return c.enableConfigOpenTelemetry
}

func (c disbursementServiceConfig) GetEnableConfigAllowTruncateAttributesOtel() bool {
	return c.enableConfigAllowTruncateAttributesOtel
}

func (c disbursementServiceConfig) GetStartDebugServer() bool {
	return c.startDebugServer
}

func (c disbursementServiceConfig) GetDisbursementKafkaTopic() string {
	return c.disbursementKafkaTopic
}

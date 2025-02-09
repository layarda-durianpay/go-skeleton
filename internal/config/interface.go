package config

type DisbursementServiceConfig interface {
	GetDisbursementDynamicConfig() string
	GetDisbursementStaticConfig() string
	GetEnableConfigOpenTelemetry() bool
	GetEnableConfigAllowTruncateAttributesOtel() bool
	GetStartDebugServer() bool
	GetDisbursementKafkaTopic() string
}

type GlobalConfig interface {
	GetGlobalDynamicConfig() string
	GetGlobalStaticConfig() string
	GetMerchantServiceGRPCAddr() string
	GetDebugPortForServer() int
	GetGRPCMaxConnectionAge() int
}

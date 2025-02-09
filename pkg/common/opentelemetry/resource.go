package opentelemetry

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type ResourceDeploymentOptions struct {
	Environment string
}

type ResourceOptions struct {
	ServiceName       string
	ServiceVersion    string
	ServiceInstanceID string
	Deployment        ResourceDeploymentOptions
}

func newResource(opt ResourceOptions) (*resource.Resource, error) {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(opt.ServiceName),
			semconv.ServiceVersion(opt.ServiceVersion),
			semconv.ServiceInstanceID(opt.ServiceInstanceID),
			attribute.String("deployment.environment", opt.Deployment.Environment),
		),
	)
	if err != nil {
		return nil, err
	}

	return r, nil
}

package deps

import (
	"mint-redeem-workflow/config"
	"mint-redeem-workflow/infra/brale"

	"github.com/uber-go/tally"
	apiv1 "github.com/uber/cadence-idl/go/proto/api/v1"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/client"
	"go.uber.org/cadence/compatibility"
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/transport/grpc"
)

type Dependencies struct {
	Config        config.ServiceConfig
	BraleClient   brale.BraleClient
	CadenceClient workflowserviceclient.Interface
}

const (
	clientName     = "mint-redeem"
	cadenceService = "cadence-frontend"
	hostPort       = "127.0.0.1:7833"
)

func NewDependencies() (*Dependencies, error) {
	cfg, err := config.NewServiceConfig()
	if err != nil {
		return nil, err
	}

	braleClient := brale.NewMockBraleClient()

	return &Dependencies{
		Config:      *cfg,
		BraleClient: braleClient,
	}, nil
}

func BuildCadenceServiceClient() (workflowserviceclient.Interface, error) {
	dispatcher := yarpc.NewDispatcher(yarpc.Config{
		Name: clientName,
		Outbounds: yarpc.Outbounds{
			cadenceService: {Unary: grpc.NewTransport().NewSingleOutbound(hostPort)},
		},
	})
	if err := dispatcher.Start(); err != nil {
		return nil, err
	}

	clientConfig := dispatcher.ClientConfig(cadenceService)

	return compatibility.NewThrift2ProtoAdapter(
		apiv1.NewDomainAPIYARPCClient(clientConfig),
		apiv1.NewWorkflowAPIYARPCClient(clientConfig),
		apiv1.NewWorkerAPIYARPCClient(clientConfig),
		apiv1.NewVisibilityAPIYARPCClient(clientConfig),
	), nil
}

func BuildCadenceClient() (client.Client, error) {
	service, err := BuildCadenceServiceClient()
	if err != nil {
		return nil, err
	}

	return client.NewClient(
		service, "test-domain2", &client.Options{MetricsScope: tally.NoopScope}), nil
}

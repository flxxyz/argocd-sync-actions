package argocd

import (
	"context"
	"io"

	argocdclient "github.com/argoproj/argo-cd/v2/pkg/apiclient"
	applicationpkg "github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	argoio "github.com/argoproj/argo-cd/v2/util/io"
)

// API is struct for argocd api.
type API struct {
	client     applicationpkg.ApplicationServiceClient
	connection io.Closer
}

// APIOptions is options for API.
type APIOptions struct {
	Address   string
	Token     string
	Insecure  bool
	PlainText bool
}

// NewAPI creates new API.
func NewAPI(options APIOptions) API {
	clientOptions := argocdclient.ClientOptions{
		PlainText:  options.PlainText,
		ServerAddr: options.Address,
		AuthToken:  options.Token,
		Insecure:   options.Insecure,
		GRPCWeb:    true,
	}

	connection, client := argocdclient.NewClientOrDie(&clientOptions).NewApplicationClientOrDie()

	return API{client: client, connection: connection}
}

// Sync syncs given application.
func (a API) Sync(appName string) error {
	request := applicationpkg.ApplicationSyncRequest{
		Name:  &appName,
		Prune: true,
	}

	_, err := a.client.Sync(context.Background(), &request)
	if err != nil {
		return err
	}

	defer argoio.Close(a.connection)

	return nil
}

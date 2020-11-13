package lnd

import (
	"context"
	// LND Client
	"github.com/lightninglabs/lndclient"

	// Internals
	"gitlab.com/nolim1t/golang-httpd-test/common"

	// GRPC Handling stuff
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/macaroon.v2"
)

// Config struct
type LndStruct struct {
	adminClient lnrpc.LightningClient
}

// Start Client
func startClient(conf common.Lnd) (c Lnd, err error) {
	conf.TlsFile = common.CleanAndExpandPath(conf.TlsFile)
	conf.MacaroonFile = common.CleanAndExpandPath(conf.MacaroonFile)

	transportCredentials, err := credentials.NewClientTLSFromFile(conf.TlsFile, conf.Host)

	if err != nil {
		return c, err
	}
	hostname := fmt.Sprintf("%s:%d", conf.Host, conf.Port)

	adminClient := getClient(transportCredentials, conf.Host, conf.MacaroonFile)

	notifier, err := NewNotifier(adminClient)
	if err != nil {
		return c, err
	}

	c = Lnd{
		adminClient: adminClient,
		notifier:    notifier,
	}

	go c.checkConnectionStatus()

	return c, nil
}

// Helper Functions
type rpcCreds map[string]string

func (m rpcCreds) RequireTransportSecurity() bool { return true }
func (m rpcCreds) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return m, nil
}

func newCreds(bytes []byte) rpcCreds {
	creds := make(map[string]string)
	creds["macaroon"] = hex.EncodeToString(bytes)
	return creds
}

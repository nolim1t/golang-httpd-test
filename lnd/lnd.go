package lnd

import (
	"context"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"time"
	// LND Client
	//"github.com/lightninglabs/lndclient"
	// LN RPC
	"github.com/lightningnetwork/lnd/lnrpc"
	// Internals
	"gitlab.com/nolim1t/golang-httpd-test/common"

	// GRPC Handling stuff
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/macaroon.v2"
)

type (
	LndInfo struct {
		Version             string   `json:"version"`
		CommitHash          string   `json:"commit_hash"`
		PubKey              string   `json:"identity_pubkey"`
		Alias               string   `json:"alias"`
		Color               string   `json:"color"`
		Peers               uint32   `json:"num_peers"`
		BlockHeight         uint32   `json:"block_height"`
		BlockHash           string   `json:"block_hash"`
		BestHeaderTimestamp int64    `json:"best_header_timestamp"`
		SyncedChain         bool     `json:"synced_to_chain"`
		SyncedGraph         bool     `json:"synced_to_graph"`
		Uris                []string `json:"uris"`
	}
)

// Client struct
type LndClient interface {
	Info(context.Context) (LndInfo, error)
}

// Config struct
type Lnd struct {
	adminClient lnrpc.LightningClient
}

func (lnd Lnd) Info(ctx context.Context) (info LndInfo, err error) {
	i, err := lnd.adminClient.GetInfo(ctx, &lnrpc.GetInfoRequest{})
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	// Return Info (see https://github.com/lightningnetwork/lnd/blob/master/lnrpc/rpc.pb.go#L4279)
	fmt.Printf("Info() raw result: %s ", i)
	return LndInfo{
		Version:             i.GetVersion(),
		CommitHash:          i.CommitHash,
		PubKey:              i.GetIdentityPubkey(),
		Alias:               i.GetAlias(),
		Color:               i.Color,
		Peers:               i.NumPeers,
		BlockHeight:         i.BlockHeight,
		BlockHash:           i.BlockHash,
		BestHeaderTimestamp: i.BestHeaderTimestamp,
		SyncedChain:         i.SyncedToChain,
		Uris:                i.GetUris(),
	}, nil
}

// Check connection
func (lnd Lnd) checkConnectionStatus() {
	failures := 0
	for {
		failures++

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := lnd.Info(ctx)
		if err == nil {
			if failures > 1 {
				fmt.Sprintln("lnd connection re-established")
			}
			failures = 0
		}
		cancel()

		if failures > 0 {
			fmt.Sprintln("lnd unreachable")
		}
		time.Sleep(time.Minute)
	}
}

// Start function
func Start(conf common.LndConfig) (Lnd, error) {
	return startClient(conf)
}

// Get Client
func getClient(transportCredentials credentials.TransportCredentials, fullHostname, file string) lnrpc.LightningClient {
	macaroonBytes, err := ioutil.ReadFile(file)
	if err != nil {
		panic(fmt.Sprintln("Cannot read macaroon file", err))
	}
	mac := &macaroon.Macaroon{}
	if err = mac.UnmarshalBinary(macaroonBytes); err != nil {
		panic(fmt.Sprintln("Cannot unmarshal macaroon", err))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connection, err := grpc.DialContext(ctx, fullHostname, []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(transportCredentials),
		grpc.WithPerRPCCredentials(newCreds(macaroonBytes)),
	}...)
	if err != nil {
		panic(fmt.Errorf("unable to connect to %s: %w", fullHostname, err))
	}

	return lnrpc.NewLightningClient(connection)
}

// Start Client
func startClient(conf common.LndConfig) (c Lnd, err error) {
	conf.TlsFile = common.CleanAndExpandPath(conf.TlsFile)
	conf.MacaroonFile = common.CleanAndExpandPath(conf.MacaroonFile)

	transportCredentials, err := credentials.NewClientTLSFromFile(conf.TlsFile, conf.Host)

	if err != nil {
		return c, err
	}
	hostname := fmt.Sprintf("%s:%d", conf.Host, conf.Port)

	adminClient := getClient(transportCredentials, hostname, conf.MacaroonFile)

	c = Lnd{
		adminClient: adminClient,
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

package node

import (
	"context"

	ds "github.com/ipfs/go-datastore"
	logging "github.com/ipfs/go-log/v2"

	"github.com/rollkit/rollkit/block"
	coreda "github.com/rollkit/rollkit/core/da"
	coreexecutor "github.com/rollkit/rollkit/core/execution"
	coresequencer "github.com/rollkit/rollkit/core/sequencer"
	"github.com/rollkit/rollkit/pkg/config"
	"github.com/rollkit/rollkit/pkg/genesis"
	"github.com/rollkit/rollkit/pkg/p2p"
	"github.com/rollkit/rollkit/pkg/service"
	"github.com/rollkit/rollkit/pkg/signer"
)

// Node is the interface for an application node
type Node interface {
	service.Service

	IsRunning() bool
}

type NodeOptions struct {
	ManagerOptions block.ManagerOptions
}

// NewNode returns a new Full or Light Node based on the config
// This is the entry point for composing a node, when compiling a node, you need to provide an executor
// Example executors can be found in apps/
func NewNode(
	ctx context.Context,
	conf config.Config,
	exec coreexecutor.Executor,
	sequencer coresequencer.Sequencer,
	da coreda.DA,
	signer signer.Signer,
	p2pClient *p2p.Client,
	genesis genesis.Genesis,
	database ds.Batching,
	metricsProvider MetricsProvider,
	logger logging.EventLogger,
	nodeOptions NodeOptions,
) (Node, error) {
	if conf.Node.Light {
		return newLightNode(conf, genesis, p2pClient, database, logger)
	}

	if err := nodeOptions.ManagerOptions.Validate(); err != nil {
		nodeOptions.ManagerOptions = block.DefaultManagerOptions()
	}

	return newFullNode(
		ctx,
		conf,
		p2pClient,
		signer,
		genesis,
		database,
		exec,
		sequencer,
		da,
		metricsProvider,
		logger,
		nodeOptions,
	)
}

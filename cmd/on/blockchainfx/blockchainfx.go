/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package blockchainfx

import (
	"context"
	"os"

	"github.com/it-chain/iLogger"

	"github.com/it-chain/engine/blockchain/api"
	"github.com/it-chain/engine/blockchain/infra/adapter"
	"github.com/it-chain/engine/blockchain/infra/mem"
	"github.com/it-chain/engine/blockchain/infra/repo"

	"github.com/it-chain/engine/api_gateway"
	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/rabbitmq/pubsub"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/conf"
	"go.uber.org/fx"
)

const publisherID = "publisher.1"
const BbPath = "./db"

var Module = fx.Options(
	fx.Provide(
		NewBlockRepository,
		mem.NewBlockPool,
		NewBlockAdapter,
		NewQueryService,
		NewBlockApi,
		NewSyncApi,
		NewConnectionEventHandler,
		NewBlockProposeHandler,
	),
	fx.Invoke(
		CreateGenesisBlock,
		RegisterRpcHandlers,
		RegisterPubsubHandlers,
		RegisterTearDown,
	),
)

func NewBlockAdapter() *adapter.HttpBlockAdapter {
	return adapter.NewHttpBlockAdapter()
}

func NewQueryService(blockAdapter *adapter.HttpBlockAdapter, connectionQueryApi *api_gateway.ConnectionQueryApi) *adapter.QuerySerivce {
	return adapter.NewQueryService(blockAdapter, connectionQueryApi)
}

func NewBlockRepository() (*repo.BlockRepository, error) {

	return repo.NewBlockRepository(BbPath)
}

func NewBlockApi(blockRepository *repo.BlockRepository, blockPool *mem.BlockPool, service common.EventService) (*api.BlockApi, error) {

	return api.NewBlockApi(publisherID, blockRepository, service, blockPool)
}

func NewSyncApi(blockRepository *repo.BlockRepository, eventService common.EventService, queryService *adapter.QuerySerivce, blockPool *mem.BlockPool) (*api.SyncApi, error) {

	api, err := api.NewSyncApi(publisherID, blockRepository, eventService, queryService, blockPool)
	return &api, err
}

func NewBlockProposeHandler(blockApi *api.BlockApi, config *conf.Configuration) *adapter.BlockProposeCommandHandler {
	return adapter.NewBlockProposeCommandHandler(blockApi, config.Engine.Mode)
}

func NewConnectionEventHandler(syncApi *api.SyncApi) *adapter.ConnectionEventHandler {
	return adapter.NewConnectionEventHandler(syncApi)
}

func CreateGenesisBlock(blockApi *api.BlockApi, config *conf.Configuration) {
	if err := blockApi.CommitGenesisBlock(config.Blockchain.GenesisConfPath); err != nil {
		panic(err)
	}
}

func RegisterRpcHandlers(server *rpc.Server, handler *adapter.BlockProposeCommandHandler) {
	iLogger.Infof(nil, "[Main] Blockchain is starting")
	if err := server.Register("block.propose", handler.HandleProposeBlockCommand); err != nil {
		panic(err)
	}
}

func RegisterPubsubHandlers(subscriber *pubsub.TopicSubscriber, handler *adapter.ConnectionEventHandler) {

	if err := subscriber.SubscribeTopic("connection.saved", handler); err != nil {
		panic(err)
	}

}

func RegisterTearDown(lifecycle fx.Lifecycle) {

	lifecycle.Append(fx.Hook{
		OnStart: func(context context.Context) error {
			return nil
		},
		OnStop: func(context context.Context) error {
			return os.RemoveAll(BbPath)
		},
	})
}
// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sender

import (
	backend "github.com/open-falcon/falcon-plus/common/backend_pool"
	"github.com/open-falcon/falcon-plus/modules/transfer/g"
	nset "github.com/toolkits/container/set"
)

func initConnPools() {
	cfg := g.Config()

	// judge
	if cfg.Judge.Enabled {
		judgeInstances := nset.NewStringSet()
		for _, instance := range cfg.Judge.Cluster {
			judgeInstances.Add(instance)
		}
		JudgeConnPools = backend.CreateSafeRpcConnPools(cfg.Judge.MaxConns, cfg.Judge.MaxIdle,
			cfg.Judge.ConnTimeout, cfg.Judge.CallTimeout, judgeInstances.ToSlice())
	}

	// transfer
	if cfg.Transfer.Enabled {
		transferInstances := nset.NewStringSet()
		for _, instance := range cfg.Transfer.Cluster {
			transferInstances.Add(instance)
		}
		TransferConnPools = backend.CreateSafeRpcConnPools(cfg.Transfer.MaxConns, cfg.Transfer.MaxIdle,
			cfg.Transfer.ConnTimeout, cfg.Transfer.CallTimeout, transferInstances.ToSlice())
	}

	// tsdb
	if cfg.Tsdb.Enabled {
		TsdbConnPoolHelper = backend.NewTsdbConnPoolHelper(cfg.Tsdb.Address, cfg.Tsdb.MaxConns, cfg.Tsdb.MaxIdle, cfg.Tsdb.ConnTimeout, cfg.Tsdb.CallTimeout)
	}

	// graph
	if cfg.Graph.Enabled {
		graphInstances := nset.NewSafeSet()
		for _, nitem := range cfg.Graph.ClusterList {
			for _, addr := range nitem.Addrs {
				graphInstances.Add(addr)
			}
		}
		GraphConnPools = backend.CreateSafeRpcConnPools(cfg.Graph.MaxConns, cfg.Graph.MaxIdle,
			cfg.Graph.ConnTimeout, cfg.Graph.CallTimeout, graphInstances.ToSlice())

	}
}

func DestroyConnPools() {
	JudgeConnPools.Destroy()
	GraphConnPools.Destroy()
	TsdbConnPoolHelper.Destroy()
}

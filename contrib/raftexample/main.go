// Copyright 2015 The etcd Authors
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

package main

import (
	"flag"
	"strings"

	"go.etcd.io/etcd/raft/raftpb"
)

// zhou: Http PUT value -> kvStore.Propose -> proposeC -> raft -> commitC -> map[string] string
func main() {
	cluster := flag.String("cluster", "http://127.0.0.1:9021", "comma separated cluster peers")
	id := flag.Int("id", 1, "node ID")
	kvport := flag.Int("port", 9121, "key-value server port")
	join := flag.Bool("join", false, "join an existing cluster")
	flag.Parse()

	// zhou: KV Store -> Raft, propose new log
	proposeC := make(chan string)
	defer close(proposeC)

	// zhou: Http Server -> Raft, propose new cluster config
	confChangeC := make(chan raftpb.ConfChange)
	defer close(confChangeC)

	// raft provides a commit stream for the proposals from the http api
	var kvs *kvstore
	// zhou: kvs.getSnapshot() still sits in package "main", so low capital is acceptable.
	getSnapshot := func() ([]byte, error) { return kvs.getSnapshot() }

	// zhou: "commitC", Raft -> KV Store, committed log sequence
	//       "errorC", Raft throw errors
	//       "snapshotterReady", Raft -> KV Store, 
	commitC, errorC, snapshotterReady := newRaftNode(*id, strings.Split(*cluster, ","), *join, getSnapshot, proposeC, confChangeC)

	// zhou: acting as Key-Value database
	kvs = newKVStore(<-snapshotterReady, proposeC, commitC, errorC)

	// zhou: client interface
	// the key-value http handler will propose updates to raft
	serveHttpKVAPI(kvs, *kvport, confChangeC, errorC)
}

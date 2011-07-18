package main

import "sync"
import "os"
import "gosqlite.googlecode.com/hg/sqlite"

type NodeStore struct {
	nodes map[string]*Node
	db *sqlite.Conn
	sync.RWMutex
}

func NewNodeStore(filename string) (*NodeStore, os.Error) {
	conn, err  := sqlite.Open(filename)
	if err != nil {
		return nil, err
	}
	return &NodeStore{ nodes : make(map[string]*Node), db : conn }, nil
}

func (store *NodeStore) Set(node *Node) {
	store.Lock()
	defer store.Unlock()
	store.nodes[node.Name] = node
}

func (store *NodeStore) Get(name string) (*Node, bool) {
	store.RLock()
	defer store.RUnlock()
	node, present := store.nodes[name]
	return node, present
}

func (store *NodeStore) GetAll() map[string]*Node {
	copied := make(map[string]*Node)
	store.RLock()
	defer store.RUnlock()
	for key, value := range store.nodes {
		copied[key] = value
	}
	return copied
}

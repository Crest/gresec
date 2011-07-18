package main

import "http"
import "fmt"

const lenNamePath = len("/name/")

func setContentText(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
}

func GetNodeByName(store *NodeStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		setContentText(w)
		name := req.URL.Path[lenNamePath:]
		node, present := store.Get(name)
		if present {
			fmt.Fprintln(w, node)
		} else {
			http.NotFound(w, req)
		}
	}
}

func SetNodeByName(store *NodeStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		setContentText(w)
		var node Node
		if _, err := fmt.Fscanln(req.Body, &node); err != nil {
			http.Error(w, "Failed to parse node", 400)
		} else {
			store.Set(&node)
		}
	}
}

func GetAllNodes(store *NodeStore) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		setContentText(w)
		for _, node := range store.GetAll() {
			fmt.Fprintln(w, node)
		}
		var peerCN string
		if len(req.TLS.PeerCertificates) > 0 {
			peerCN = req.TLS.PeerCertificates[0].Subject.CommonName 
		} else {
			peerCN = ""
		}

		fmt.Printf("complete    = %v\n", req.TLS.HandshakeComplete)
		fmt.Printf("cipherSuite = %v\n", req.TLS.CipherSuite)
		fmt.Printf("protocol    = %v\n", req.TLS.NegotiatedProtocol)
		fmt.Printf("  mutual?   = %v\n", req.TLS.NegotiatedProtocolIsMutual)
		fmt.Printf("peerCN      = %v\n", peerCN)
	}
}

package main

import "bytes"
import "net"
import "io"
import "io/ioutil"
import "os"
import "time"
import "fmt"
import "http"
import "crypto/rand"
import "crypto/tls"
import "crypto/x509"

var errEOF = io.ErrUnexpectedEOF

func toBytes(n int64) [8]byte {
	m := uint64(n)
	var bytes [8]byte
	bytes[0] = byte(m >> 56)
	bytes[1] = byte(m >> 48)
	bytes[2] = byte(m >> 40)
	bytes[3] = byte(m >> 32)
	bytes[4] = byte(m >> 25)
	bytes[5] = byte(m >> 16)
	bytes[6] = byte(m >> 8)
	bytes[7] = byte(m >> 0)
	return bytes
}

func secondsToBytes() [8]byte {
	return toBytes(time.UTC().Seconds())
}

func readNodes(r io.Reader) (map[string]Node, os.Error) {
	nodes := make(map[string]Node)
	for {
		var node Node
		if _, err := fmt.Fscanln(r, &node); err != nil {
			if err == errEOF {
				return nodes, nil
			} else {
				return nodes, err
			}
		} else {
			nodes[node.Name] = node
		}
	}
	return nodes, nil
}

func listenAndServe(addr string, certFile string, keyFile string, caFile string) os.Error {
	config := &tls.Config{
		Rand:               rand.Reader,
		Time:               time.Seconds,
		NextProtos:         []string{"http/1.1"},
		AuthenticateClient: true,
		CipherSuites:       []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA},
	}

	rootCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		return err
	}
	config.RootCAs = x509.NewCertPool()
	if !config.RootCAs.AppendCertsFromPEM(rootCert) {
		return os.NewError("Failed to add root certificate.")
	}

	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	conn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	tlsListener := tls.NewListener(conn, config)
	return http.Serve(tlsListener, nil)
}


func main() {
	buf := bytes.NewBufferString("eq4 46.4.89.243 10.0.0.2 2001:470:9ce6:200::2")
	nodes, err := readNodes(buf)
	if err != nil {
		fmt.Println("ERR: " + err.String())
	}
	node := nodes["eq4"]
	store, err := NewNodeStore("nodes")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERR: %v\n", err)
		return
	}
	store.Set(&node)
	http.HandleFunc("/name/", GetNodeByName(store))
	http.HandleFunc("/set", SetNodeByName(store))
	http.HandleFunc("/all", GetAllNodes(store))

	err = listenAndServe(":8080", "cert.pem", "key.pem", "cacert.pem")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERR: %v\n", err)
	}
}

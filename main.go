package main

import (
	"crypto/tls"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func main() {

	// Define flags
	pflag.String("cert-path", "server.crt", "Path to the certificate file")
	pflag.String("key-path", "server.key", "Path to the key file")
	pflag.String("dest-url", "dest.url.com", "Destination URL (without http prefix)")
	pflag.String("listen-addr", "localhost", "Listen address")
	pflag.Int("dest-port", 5000, "Destination port")
	pflag.Int("listen-port", 9000, "Listen port")

	// Bind flags to viper
	viper.BindPFlags(pflag.CommandLine)

	// Parse command line flags
	pflag.Parse()

	// Get values from viper
	certPath := viper.GetString("cert-path")
	keyPath := viper.GetString("key-path")
	destURL := viper.GetString("dest-url")
	listenAddr := viper.GetString("listen-addr")
	destPort := viper.GetInt("dest-port")
	listenPort := viper.GetInt("listen-port")

	destFullUrl := fmt.Sprintf("http://%s:%d", destURL, destPort)
	ReverseHttpsProxy(listenPort, listenAddr, destFullUrl, certPath, keyPath)
}

func ReverseHttpsProxy(port int, addr, dst, crt, key string) {
	u, e := url.Parse(dst)
	if e != nil {
		log.Fatal("Bad destination.")
	}
	h := httputil.NewSingleHostReverseProxy(u)
	var InsecureTransport http.RoundTripper = &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		TLSHandshakeTimeout: 10 * time.Second,
	}
	h.Transport = InsecureTransport
	fmt.Printf("Serving at %s:%d\n", addr, port)
	err := http.ListenAndServeTLS(fmt.Sprintf("%s:%d", addr, port), crt, key, h)
	if err != nil {
		log.Println("Error:", err)
	}
}

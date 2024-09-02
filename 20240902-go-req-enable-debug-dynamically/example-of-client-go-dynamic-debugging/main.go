package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/imroc/req/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var reqClient = req.C()

func init() {
	go func() {
		r := gin.Default()
		r.GET("/debug", func(c *gin.Context) {
			query := c.Query("enable")
			switch query {
			case "true": // Turn on debug with `curl 127.0.0.1:80/debug?enable=true`
				reqClient.EnableDumpAll()
				reqClient.EnableDebugLog()
				fmt.Println("Debug is enabled")
			case "false": // Turn off debug with `curl 127.0.0.1:80/debug?enable=false`
				reqClient.DisableDumpAll()
				reqClient.DisableDebugLog()
				fmt.Println("Debug is disabled")
			}
		})
		r.Run("0.0.0.0:80")
	}()
}

func main() {
	clientset := newClientset()
	for {
		// Run list services in kube-system for testing.
		services, err := clientset.CoreV1().Services("kube-system").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d services in the kube-system cluster\n", len(services.Items))
		time.Sleep(10 * time.Second)
	}
}

// create kubernetes clientset
func newClientset() *kubernetes.Clientset {
	// Creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	replaceTransport(config, reqClient.GetTransport())

	// Create kubernetes client with config.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

// replace client-go's Transport with *req.Transport
func replaceTransport(config *rest.Config, t *req.Transport) {
	// Extract tls.Config from rest.Config
	tlsConfig, err := rest.TLSConfigFor(config)
	if err != nil {
		panic(err.Error())
	}
	// Set TLSClientConfig to req's Transport.
	t.TLSClientConfig = tlsConfig
	// Override with req's Transport.
	config.Transport = t
	// rest.Config.TLSClientConfig should be empty if
	// custom Transport been set.
	config.TLSClientConfig = rest.TLSClientConfig{}
}

package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/SecurityDo/ingext_api/client"
)

// IngextAppAPI defines the contract for interacting with the backend
/*
type IngextAppAPI interface {
	Init(cluster, namespace, kubeCtx string) error
	Call(functionName string, functionArgs json.RawMessage) error

	// Domain specific methods
	AddStreamSource() error
	AddStreamSink() error
	// ... add other methods here
}*/

// Client is the concrete implementation of IngextAppAPI
type Client struct {
	Logger    *slog.Logger
	Cluster   string
	Namespace string
	// Add k8s clients, rest clients, etc here

	// Embed the K8s helper
	k8sClient    *K8sClusterClient
	ingextClient *client.IngextClient // If you have a separate client for ingext
}

// Option 1: Constructor injection (Recommended)
func NewClient(logger *slog.Logger) *Client {
	// Fallback: If caller passes nil, use a "No-Op" or Default logger
	// so the code doesn't panic.
	if logger == nil {
		// Using os.Stderr by default is safe for libraries
		logger = slog.New(slog.NewTextHandler(os.Stderr, nil))
	}
	return &Client{
		Logger:    logger,
		k8sClient: NewK8sClient(), // Initialize the helper
	}
}

// Ensure Client implements the interface
//var _ IngextAppAPI = (*Client)(nil)

var IngextAppAPI *Client

func (c *Client) Init(cluster, namespace string, kubeContext string) error {
	c.Cluster = cluster
	c.Namespace = namespace

	c.Logger.Debug("initializing k8s client", "context", kubeContext)

	// 1. Connect to Kubernetes
	if err := c.k8sClient.Connect(kubeContext); err != nil {
		return err
	}
	c.Logger.Info("connected to kubernetes cluster", "context", kubeContext)

	// TODO: Perform actual login / connection logic here
	c.Logger.Debug("Connecting to cluster ...\n", "cluster", cluster, "namespace", namespace)

	token, err := c.k8sClient.GetAppSecret(namespace, "app-secret", "token")
	if err != nil {
		return fmt.Errorf("failed to get app secret token: %s", err)
	}

	configText, err := c.k8sClient.GetAppConfig(namespace, "ingext-community-config", "site_config.json")
	if err != nil {
		return fmt.Errorf("failed to get app config: %s", err)
	}
	var config struct {
		SiteURL string `json:"siteURL"`
	}

	if err := json.Unmarshal([]byte(configText), &config); err != nil {
		c.Logger.Error("failed to parse site config", "error", err, "config", configText)
		return fmt.Errorf("failed to parse site config:  %s", err)
	}

	ingextClient := client.NewIngextClient(config.SiteURL, token, false, c.Logger)

	c.ingextClient = ingextClient

	c.Logger.Info("initialized ingext client",
		"siteURL", config.SiteURL,
		//"token", token,
	)

	return nil
}

func (c *Client) Call(functionName string, functionArgs json.RawMessage) error {
	// TODO: Implementation
	return nil
}

package commands

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/SecurityDo/ingext_api/internal/api"

	"github.com/SecurityDo/ingext_api/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Global flags
	cfgFile     string
	cluster     string
	namespace   string
	logLevel    string
	showVersion bool

	siteConfig string
	site       string
)

const (
	appVersion      = "1.1.0"
	defaultLogLevel = "warn"
)

/*
Default behavior: Use cmd.PrintErrf (or cmd.PrintErrln) for everything (interactive prompts, status logs, errors).
Exception: Use cmd.Printf (or cmd.Println) only when you are printing the final machine-readable output.
*/

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ingext",
	Short: "A CLI tool for managing ingext resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			cmd.Printf("ingext version %s\n", appVersion)
			return nil
		}

		return cmd.Help()
	},

	SilenceUsage:  true, // Don't show help text on runtime errors
	SilenceErrors: true, // Optional: if you want to print errors yourself in main.go
	// PersistentPreRunE runs BEFORE the subcommand (e.g., 'ingext stream add')
	// but AFTER flags are parsed.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip initialization if we're only showing the version
		if showVersion {
			return nil
		}

		// 1. Skip initialization for commands that don't need it (like 'config' or 'help')
		if cmd.Parent() != nil && cmd.Parent().Name() == "config" {
			return nil
		}
		if cmd.Name() == "config" {
			return nil
		}
		//if cmd.Name() == "version" {
		//	return nil
		//}

		// 1. SKIP logic for Help and Autocompletion
		// Cobra adds a "help" command automatically.
		// "__complete" is used during shell tab-completion.
		if cmd.Name() == "help" || cmd.Name() == "__complete" {
			return nil
		}

		// 2. Load values from Viper (which now holds flags + config file values)
		clusterName := viper.GetString("cluster")
		namespace := viper.GetString("namespace")

		levelValue := viper.GetString("log-level")

		// 1. Configure the Handler options
		opts := &slog.HandlerOptions{}

		switch strings.ToLower(levelValue) {
		case "debug":
			opts.Level = slog.LevelDebug
		case "info":
			opts.Level = slog.LevelInfo
		case "warn", "":
			opts.Level = slog.LevelWarn
		case "error":
			opts.Level = slog.LevelError
		default:
			return fmt.Errorf("invalid log level %q: choose debug, info, warn, error", levelValue)
		}

		// 2. Create the Handler pointing to STDERR
		handler := slog.NewTextHandler(cmd.ErrOrStderr(), opts)

		// 3. Create the Logger
		logger := slog.New(handler)

		// 4. Inject into your Client
		AppAPI = api.NewClient(logger)

		// ai register / unregister: no cluster/site/env; URL and token come from flags on the ai command.
		if IsAiTokenCommand(cmd) {
			return nil
		}

		// 5. Initialize the Global API — three modes in priority order:
		//    a) INGEXT_SITE_URL + INGEXT_TOKEN env vars (direct connect, no k8s)
		//    b) site_credentials.json file
		//    c) Kubernetes cluster

		envSiteURL := os.Getenv("INGEXT_SITE_URL")
		envToken := os.Getenv("INGEXT_TOKEN")

		if envSiteURL != "" && envToken != "" {
			// Mode (a): direct connect via environment variables
			AppAPI.InitDirect(envSiteURL, envToken)
			logger.Info("initialized ingext client from env vars", "siteURL", envSiteURL)
		} else {
			siteConfigPath := viper.GetString("site-config")
			siteName := viper.GetString("site")
			if siteName == "" {
				siteName = viper.GetString("default-site")
			}
			// Resolve site-config path: if not set, default to site_credentials.json in cwd
			if siteConfigPath == "" {
				if cwd, err := os.Getwd(); err == nil {
					siteConfigPath = filepath.Join(cwd, "site_credentials.json")
				}
			}

			useSiteConfig := false
			if siteConfigPath != "" {
				if _, err := os.Stat(siteConfigPath); err == nil {
					useSiteConfig = true
				}
			}

			if useSiteConfig {
				// Mode (b): site_credentials.json
				if err := AppAPI.InitFromSiteConfig(siteConfigPath, siteName); err != nil {
					return fmt.Errorf("failed to initialize from site config: %w", err)
				}
			} else {
				// Mode (c): Kubernetes
				if clusterName == "" {
					return fmt.Errorf("cluster name is required. Set INGEXT_SITE_URL + INGEXT_TOKEN env vars, place site_credentials.json in the current directory, or use --cluster")
				}
				kubeCtx := viper.GetString("context")
				if kubeCtx == "" {
					logger.Warn("no kube-context specified in config, using current system default")
				}
				if err := AppAPI.Init(clusterName, namespace, kubeCtx); err != nil {
					return fmt.Errorf("failed to initialize app API: %w", err)
				}
			}
		}

		// 6. Enable HTTP request/response debug dumps when log level is debug
		if strings.ToLower(levelValue) == "debug" {
			AppAPI.SetDebug(true)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func init() {
	cobra.OnInitialize(config.InitConfig)

	// Define global flags
	RootCmd.PersistentFlags().StringVar(&siteConfig, "site-config", "", "path to site_credentials.json (default: ./site_credentials.json)")
	RootCmd.PersistentFlags().StringVar(&site, "site", "", "site hostname from tokenMap (e.g. demo.cloud.fluencysecurity.com); if empty and using site config, first site is used")

	RootCmd.PersistentFlags().StringVar(&cluster, "cluster", "", "k8s cluster name")
	RootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "ingext", "namespace of the ingext app")
	RootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", defaultLogLevel, "log level: debug, info, warn, error")
	RootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "show version")
	RootCmd.Version = appVersion
	// Bind global flags to viper so they can be accessed anywhere
	viper.BindPFlag("site-config", RootCmd.PersistentFlags().Lookup("site-config"))
	viper.BindPFlag("site", RootCmd.PersistentFlags().Lookup("site"))

	viper.BindPFlag("cluster", RootCmd.PersistentFlags().Lookup("cluster"))
	viper.BindPFlag("namespace", RootCmd.PersistentFlags().Lookup("namespace"))
	viper.BindPFlag("log-level", RootCmd.PersistentFlags().Lookup("log-level"))
}

package commands

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/SecurityDo/ingext_api/internal/api"

	"github.com/SecurityDo/ingext_api/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Global flags
	cfgFile   string
	cluster   string
	namespace string
)

/*
Default behavior: Use cmd.PrintErrf (or cmd.PrintErrln) for everything (interactive prompts, status logs, errors).
Exception: Use cmd.Printf (or cmd.Println) only when you are printing the final machine-readable output.
*/

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "ingext",
	Short: "A CLI tool for managing ingext resources",

	SilenceUsage:  true, // Don't show help text on runtime errors
	SilenceErrors: true, // Optional: if you want to print errors yourself in main.go
	// PersistentPreRunE runs BEFORE the subcommand (e.g., 'ingext stream add')
	// but AFTER flags are parsed.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 1. Skip initialization for commands that don't need it (like 'config' or 'help')
		if cmd.Parent() != nil && cmd.Parent().Name() == "config" {
			return nil
		}
		if cmd.Name() == "config" {
			return nil
		}

		// 1. SKIP logic for Help and Autocompletion
		// Cobra adds a "help" command automatically.
		// "__complete" is used during shell tab-completion.
		if cmd.Name() == "help" || cmd.Name() == "__complete" {
			return nil
		}

		// 2. Load values from Viper (which now holds flags + config file values)
		clusterName := viper.GetString("cluster")
		namespace := viper.GetString("namespace")
		kubeCtx := viper.GetString("context")
		if clusterName == "" {
			return fmt.Errorf("cluster name is required. Run 'ingext config' or use --cluster")
		}

		// 1. Configure the Handler options
		opts := &slog.HandlerOptions{
			Level: slog.LevelInfo, // Default level
		}
		if verbose {
			opts.Level = slog.LevelDebug
		}

		// 2. Create the Handler pointing to STDERR
		// cmd.ErrOrStderr() ensures we use the proper writer wrapper from Cobra
		handler := slog.NewTextHandler(cmd.ErrOrStderr(), opts)

		// 3. Create the Logger
		logger := slog.New(handler)

		// If context is empty in config, we can default to empty string
		// (which means client-go uses the "current-context" from ~/.kube/config)
		if kubeCtx == "" {

			// Optional: log a warning
			logger.Warn("no kube-context specified in config, using current system default")

		}

		// 4. Inject into your Client
		// Now your client logs will go to Stderr, respecting the --verbose flag
		AppAPI = api.NewClient(logger)

		// 3. Initialize the Global API
		if err := AppAPI.Init(clusterName, namespace, kubeCtx); err != nil {
			return fmt.Errorf("failed to initialize app API: %w", err)
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

var verbose bool

func init() {
	cobra.OnInitialize(config.InitConfig)

	// Define global flags
	RootCmd.PersistentFlags().StringVar(&cluster, "cluster", "", "k8s cluster name")
	RootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "ingext", "namespace of the ingext app")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")
	// Bind global flags to viper so they can be accessed anywhere
	viper.BindPFlag("cluster", RootCmd.PersistentFlags().Lookup("cluster"))
	viper.BindPFlag("namespace", RootCmd.PersistentFlags().Lookup("namespace"))
}

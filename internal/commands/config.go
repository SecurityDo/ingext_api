package commands

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/SecurityDo/ingext_api/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	confProvider string
	confContext  string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure the ingext tool",
	Long:  `Manage ingext configuration stored in ~/.ingext/config.yaml.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// Subcommand: ADD
var configAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add configuration values for a cluster profile",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Build composite profile key: cluster:namespace
		targetCluster := cluster
		targetNamespace := namespace

		if cluster == "" {
			cmd.PrintErrln("Error: --cluster is required.")
			return
		}
		if namespace == "" {
			cmd.PrintErrln("Error: --namespace is required.")
			return
		}

		// If user didn't provide --cluster, fall back to current-cluster
		//if targetCluster == "" {
		//	targetCluster = viper.GetString("current-cluster")
		//}

		//if targetCluster == "" {
		//	cmd.PrintErrln("Error: No cluster name specified. Use --cluster <name> to configure a profile.")
		//	return
		//}

		// Build composite key
		profileKey := targetCluster + ":" + targetNamespace

		// 2. Set "Current Cluster" to the composite key
		viper.Set("current-cluster", profileKey)

		// 3. Save values using Dot Notation (clusters.<profileKey>.<field>)
		prefix := fmt.Sprintf("clusters.%s.", profileKey)

		// We always save the provider (defaults to 'eks' via flag if not typed)
		viper.Set(prefix+"provider", confProvider)

		// If --context was not provided, inherit from an existing profile with the same cluster
		if confContext == "" {
			allClusters := viper.GetStringMap("clusters")
			for key, val := range allClusters {
				if strings.HasPrefix(key, targetCluster+":") && key != profileKey {
					if details, ok := val.(map[string]interface{}); ok {
						if ctx, ok := details["context"]; ok {
							confContext = fmt.Sprintf("%v", ctx)
							break
						}
					}
				}
			}
			if confContext == "" {
				cmd.PrintErrln("Error: --context is required when no existing profile exists for this cluster.")
				return
			}
		}
		viper.Set(prefix+"context", confContext)

		// 4. Write to disk
		if err := config.SaveConfig(); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Printf("Configuration saved for profile '%s'.\n", profileKey)
	},
}

// Subcommand: LIST
var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured clusters",
	Run: func(cmd *cobra.Command, args []string) {
		current := viper.GetString("current-cluster")
		// GetStringMap returns map[string]interface{}
		clusters := viper.GetStringMap("clusters")

		w := tabwriter.NewWriter(os.Stderr, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "CURRENT\tPROFILE\tCLUSTER\tNAMESPACE\tPROVIDER")
		fmt.Fprintln(w, "-------\t-------\t-------\t---------\t--------")

		// Sort keys for consistent output
		var keys []string
		for k := range clusters {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, name := range keys {
			// Extract details from the nested map
			details, ok := clusters[name].(map[string]interface{})
			if !ok {
				continue
			}

			isCurrent := ""
			if name == current {
				isCurrent = "*"
			}

			// Parse composite key to extract cluster and namespace
			clusterPart := name
			namespacePart := ""
			if idx := strings.Index(name, ":"); idx >= 0 {
				clusterPart = name[:idx]
				namespacePart = name[idx+1:]
			}

			// Safe getter for provider
			prov := ""
			if v, ok := details["provider"]; ok {
				prov = fmt.Sprintf("%v", v)
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", isCurrent, name, clusterPart, namespacePart, prov)
		}
		w.Flush()
	},
}

// Subcommand: DELETE
var configDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a cluster configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if cluster == "" {
			cmd.PrintErrln("Error: --cluster is required to delete a profile.")
			return
		}
		if namespace == "" {
			cmd.PrintErrln("Error: --namespace is required to delete a profile.")
			return
		}

		// Build composite key
		profileKey := cluster + ":" + namespace

		// 1. Get the raw map
		allClusters := viper.GetStringMap("clusters")

		if _, exists := allClusters[profileKey]; !exists {
			fmt.Printf("Profile '%s' not found.\n", profileKey)
			return
		}

		// 2. Delete the key
		delete(allClusters, profileKey)

		// 3. Set the map back to Viper
		viper.Set("clusters", allClusters)

		// 4. Handle edge case: If we deleted the "current" profile, pick another if available
		current := viper.GetString("current-cluster")
		if current == profileKey {
			if len(allClusters) == 0 {
				viper.Set("current-cluster", "")
				fmt.Println("Warning: You deleted the currently active profile. No profiles remain.")
			} else {
				newCurrent := pickFirstCluster(allClusters)
				viper.Set("current-cluster", newCurrent)
				fmt.Printf("Switched current profile to '%s'.\n", newCurrent)
			}
		} else if current == "" && len(allClusters) > 0 {
			newCurrent := pickFirstCluster(allClusters)
			viper.Set("current-cluster", newCurrent)
			fmt.Printf("Current profile was unset. Switched to '%s'.\n", newCurrent)
		}

		if err := config.SaveConfig(); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Printf("Profile '%s' deleted.\n", profileKey)
	},
}

// Subcommand: SET (switch active profile)
var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the current cluster and namespace",
	Run: func(cmd *cobra.Command, args []string) {
		if cluster == "" {
			cmd.PrintErrln("Error: --cluster is required.")
			return
		}
		if namespace == "" {
			cmd.PrintErrln("Error: --namespace is required.")
			return
		}

		// Build composite key
		profileKey := cluster + ":" + namespace

		// 2. Set "Current Cluster" to the composite key
		viper.Set("current-cluster", profileKey)

		// 3. Save values using Dot Notation (clusters.<profileKey>.<field>)
		prefix := fmt.Sprintf("clusters.%s.", profileKey)

		// We always save the provider (defaults to 'eks' via flag if not typed)
		viper.Set(prefix+"provider", confProvider)

		// Verify the profile exists
		allClusters := viper.GetStringMap("clusters")
		if _, exists := allClusters[profileKey]; !exists {
			fmt.Printf("Profile '%s' not found. Use 'ingext config add' to create it first.\n", profileKey)
			return
		}

		// Switch current-cluster to the composite key
		viper.Set("current-cluster", profileKey)

		// We always save the provider (defaults to 'eks' via flag if not typed)
		viper.Set(prefix+"provider", confProvider)

		// If --context was not provided, inherit from an existing profile with the same cluster
		if confContext == "" {
			allClusters := viper.GetStringMap("clusters")
			for key, val := range allClusters {
				if strings.HasPrefix(key, cluster+":") && key != profileKey {
					if details, ok := val.(map[string]interface{}); ok {
						if ctx, ok := details["context"]; ok {
							confContext = fmt.Sprintf("%v", ctx)
							break
						}
					}
				}
			}
			if confContext == "" {
				cmd.PrintErrln("Error: --context is required when no existing profile exists for this cluster.")
				return
			}
		}
		viper.Set(prefix+"context", confContext)

		// 4. Write to disk
		if err := config.SaveConfig(); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Printf("Switched to profile '%s'.\n", profileKey)
	},
}

// Subcommand: VIEW (Updated to show current-cluster logic)
var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View current configuration settings",
	Run: func(cmd *cobra.Command, args []string) {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		// Identify current context (composite key like "datalake:ingext")
		current := viper.GetString("current-cluster")

		// Parse composite key
		clusterPart := current
		namespacePart := ""
		if idx := strings.Index(current, ":"); idx >= 0 {
			clusterPart = current[:idx]
			namespacePart = current[idx+1:]
		}

		prefix := fmt.Sprintf("clusters.%s.", current)

		fmt.Fprintln(w, "SETTING\tVALUE")
		fmt.Fprintln(w, "-------\t-----")

		fmt.Fprintf(w, "Profile\t%s\n", current)
		fmt.Fprintf(w, "Cluster\t%s\n", clusterPart)
		fmt.Fprintf(w, "Namespace\t%s\n", namespacePart)
		fmt.Fprintf(w, "Provider\t%s\n", viper.GetString(prefix+"provider"))
		fmt.Fprintf(w, "Context\t%s\n", viper.GetString(prefix+"context"))

		fmt.Fprintln(w, "-------\t-----")
		fmt.Fprintf(w, "Config File\t%s\n", viper.ConfigFileUsed())

		w.Flush()
	},
}

func init() {
	RootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configAddCmd)
	configCmd.AddCommand(configViewCmd)

	// Add new subcommands
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configDeleteCmd)
	_ = configSetCmd.MarkFlagRequired("cluster")
	_ = configSetCmd.MarkFlagRequired("namespace")
	_ = configDeleteCmd.MarkFlagRequired("cluster")
	_ = configDeleteCmd.MarkFlagRequired("namespace")

	configSetCmd.Flags().StringVar(&confProvider, "provider", "eks", "Provider (eks|aks|gke)")
	configSetCmd.Flags().StringVar(&confContext, "context", "", "Kubeconfig context name")

	// Configuration for 'config' command
	// Default value "eks" is set here for the FLAG
	configAddCmd.Flags().StringVar(&confProvider, "provider", "eks", "Provider (eks|aks|gke)")
	configAddCmd.Flags().StringVar(&confContext, "context", "", "Kubeconfig context name")

	_ = configAddCmd.MarkFlagRequired("namespace")
	_ = configAddCmd.MarkFlagRequired("cluster")

	// Set the global default for Viper as well (in case user views config without setting it)
	viper.SetDefault("provider", "eks")
}

// pickFirstCluster returns the first cluster name in sorted order.
func pickFirstCluster(clusters map[string]interface{}) string {
	var names []string
	for name := range clusters {
		names = append(names, name)
	}
	sort.Strings(names)
	return names[0]
}

/*
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure the ingext tool",
	Long:  `Sets configuration values in ~/.ingext/config.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		// Set values in Viper
		// Note: cluster and namespace are already bound in root.go,
		// but we need to explicitly set them if the user provided flags here to save them to file.

		if cluster != "" {
			viper.Set("cluster", cluster)
		}
		if namespace != "" {
			viper.Set("namespace", namespace)
		}
		if confProvider != "" {
			viper.Set("provider", confProvider)
		}
		if confContext != "" {
			viper.Set("context", confContext)
		}

		if err := config.SaveConfig(); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}

		fmt.Println("Configuration saved.")
	},
}*/
/*
// 1. Define the 'view' subcommand
var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View current configuration settings",
	Long:  "Displays the current configuration loaded from ~/.ingext/config.yaml and environment variables.",
	Run: func(cmd *cobra.Command, args []string) {
		// Use tabwriter to create a clean, aligned table
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		fmt.Fprintln(w, "SETTING\tVALUE")
		fmt.Fprintln(w, "-------\t-----")

		// Retrieve values from Viper (checks flags, env vars, and config file)
		fmt.Fprintf(w, "Cluster\t%s\n", viper.GetString("cluster"))
		fmt.Fprintf(w, "Namespace\t%s\n", viper.GetString("namespace"))
		fmt.Fprintf(w, "Provider\t%s\n", viper.GetString("provider"))
		fmt.Fprintf(w, "Context\t%s\n", viper.GetString("context"))

		// Access config file location
		fmt.Fprintln(w, "-------\t-----")
		fmt.Fprintf(w, "Config File\t%s\n", viper.ConfigFileUsed())

		w.Flush()
	},
}

func init() {
	RootCmd.AddCommand(configCmd)

	// Local flags for the config command
	configCmd.Flags().StringVar(&confProvider, "provider", "eks", "Provider (eks|aks|gke)")
	configCmd.Flags().StringVar(&confContext, "context", "", "Kubeconfig context name")
	// Note: --cluster and --namespace are inherited from Root, but usually 'config' commands
	// might want to enforce them or treat them differently. For now, we rely on the Root persistent flags.

	// 2. Register 'view' as a child of 'config'
	configCmd.AddCommand(configViewCmd)
}*/

package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var syslogPorts []string

var syslogCmd = &cobra.Command{
	Use:   "syslog",
	Short: "Syslog management commands",
}

var syslogGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get syslog configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := AppAPI.GetSyslogConfig()
		if err != nil {
			return err
		}
		jsonBytes, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal response: %w", err)
		}
		fmt.Println(string(jsonBytes))
		return nil
	},
}

var syslogRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register syslog configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := AppAPI.RegisterSyslogConfig(syslogPorts)
		if err != nil {
			return err
		}
		jsonBytes, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal response: %w", err)
		}
		fmt.Println(string(jsonBytes))
		return nil
	},
}

var syslogUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update syslog configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := AppAPI.UpdateSyslogConfig(syslogPorts)
		if err != nil {
			return err
		}
		jsonBytes, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal response: %w", err)
		}
		fmt.Println(string(jsonBytes))
		return nil
	},
}

var syslogDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete syslog configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := AppAPI.DeleteSyslogConfig(); err != nil {
			return err
		}
		fmt.Println("Syslog configuration deleted successfully.")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(syslogCmd)
	syslogCmd.AddCommand(syslogGetCmd, syslogRegisterCmd, syslogUpdateCmd, syslogDeleteCmd)

	syslogRegisterCmd.Flags().StringArrayVar(&syslogPorts, "port", []string{}, "Port types to register (tcp, udp, tls, tls-rfc6587)")
	_ = syslogRegisterCmd.MarkFlagRequired("port")

	syslogUpdateCmd.Flags().StringArrayVar(&syslogPorts, "port", []string{}, "Port types to update (tcp, udp, tls, tls-rfc6587)")
	_ = syslogUpdateCmd.MarkFlagRequired("port")
}

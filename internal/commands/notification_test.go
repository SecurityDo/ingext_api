package commands

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// add-email must default --action to Generic_Email_Action. account_setup runs it
// without --action, and an empty action leaves the endpoint with no FPL action.
func TestNotificationAddEmailDefaultAction(t *testing.T) {
	if err := notificationAddEmailCmd.Flags().Parse([]string{"--name", "BehaviorSummary_Export", "--to", "alerts@example.com"}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}
	if notifEmailAction != "Generic_Email_Action" {
		t.Fatalf("add-email --action defaulted to %q, want %q", notifEmailAction, "Generic_Email_Action")
	}
}

// pflag writes a flag's default into its variable when the flag is registered,
// so two flags bound to one variable with different defaults leave it holding
// whichever registered last. That silently emptied add-email's --action once
// add-slack reused its variable.
func TestNotificationFlagsSharingAVariableAgreeOnDefault(t *testing.T) {
	type binding struct {
		flag string
		def  string
	}
	seen := make(map[uintptr]binding)

	var walk func(c *cobra.Command)
	walk = func(c *cobra.Command) {
		c.Flags().VisitAll(func(f *pflag.Flag) {
			v := reflect.ValueOf(f.Value)
			if v.Kind() != reflect.Ptr {
				return
			}
			name := c.CommandPath() + " --" + f.Name
			if prev, ok := seen[v.Pointer()]; ok && prev.def != f.DefValue {
				t.Errorf("%s (default %q) and %s (default %q) share a variable", prev.flag, prev.def, name, f.DefValue)
				return
			}
			seen[v.Pointer()] = binding{flag: name, def: f.DefValue}
		})
		for _, sub := range c.Commands() {
			walk(sub)
		}
	}
	walk(notificationCmd)
}

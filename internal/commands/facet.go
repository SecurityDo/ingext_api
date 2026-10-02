package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/SecurityDo/ingext_api/model"
	"github.com/spf13/cobra"
)

var (
	facetName string
	facetFile string
)

var facetCmd = &cobra.Command{
	Use:   "facet",
	Short: "Manage search facet profiles",
}

var facetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List search facet profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := AppAPI.FacetList()
		if err != nil {
			return err
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tID\tGROUP\tINDEXES\tFACETS\tSTORAGE KEY")
		for _, e := range entries {
			fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%d\t%s\n", e.Name, e.ID, e.Group, strings.Join(e.Indexes, ","), len(e.Facets), e.Id)
		}
		return w.Flush()
	},
}

var facetGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Print one search facet profile as JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		entry, err := AppAPI.FacetGet(facetName)
		if err != nil {
			return err
		}
		b, err := json.MarshalIndent(entry, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal response: %w", err)
		}
		fmt.Println(string(b))
		return nil
	},
}

var facetApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Add a facet profile from a JSON file, or replace the one of the same name",
	Long: `Add a search facet profile from a JSON file (the Scripts repo Facets/*.json
format), or replace the stored profile with the same name in place.

The profile is matched by name, as the release sync does, so a later release of
the same facet replaces this one rather than duplicating it.

  ingext --gridaccount sahomeloans facet apply --file Facets/Default-Netskope_DLP.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		b, err := os.ReadFile(facetFile)
		if err != nil {
			return err
		}
		var entry model.FacetConfigEntry
		if err := json.Unmarshal(b, &entry); err != nil {
			return fmt.Errorf("parse %s: %w", facetFile, err)
		}
		created, err := AppAPI.FacetApply(&entry)
		if err != nil {
			return err
		}
		if created {
			fmt.Printf("Facet %q added.\n", entry.Name)
		} else {
			fmt.Printf("Facet %q updated (%s).\n", entry.Name, entry.Id)
		}
		return nil
	},
}

var facetDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a search facet profile by name",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := AppAPI.FacetDelete(facetName); err != nil {
			return err
		}
		fmt.Printf("Facet %q deleted.\n", facetName)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(facetCmd)
	facetCmd.AddCommand(facetListCmd, facetGetCmd, facetApplyCmd, facetDeleteCmd)

	facetGetCmd.Flags().StringVar(&facetName, "name", "", "Facet profile name")
	_ = facetGetCmd.MarkFlagRequired("name")
	facetDeleteCmd.Flags().StringVar(&facetName, "name", "", "Facet profile name")
	_ = facetDeleteCmd.MarkFlagRequired("name")
	facetApplyCmd.Flags().StringVarP(&facetFile, "file", "f", "", "Facet profile JSON file")
	_ = facetApplyCmd.MarkFlagRequired("file")
}

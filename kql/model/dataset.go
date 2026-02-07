package model

import (
	"encoding/json"
)

// ============================================================================
// DataSet: The Top-Level Result
// ============================================================================

// DataSet represents the result of a KQL execution.
// It contains an ordered list of tables (e.g., "PrimaryResult", "Facet1", "Facet2").
type DataSet struct {
	Tables []*DataTable
}

// NewDataSet creates an empty DataSet.
func NewDataSet() *DataSet {
	return &DataSet{
		Tables: make([]*DataTable, 0),
	}
}

// AddTable appends a materialized table to the result set.
func (ds *DataSet) AddTable(dt *DataTable) {
	ds.Tables = append(ds.Tables, dt)
}

// GetTable retrieves a table by name (case-sensitive).
// Returns nil if not found.
func (ds *DataSet) GetTable(name string) *DataTable {
	for _, t := range ds.Tables {
		if t.Name == name {
			return t
		}
	}
	return nil
}

// ============================================================================
// Serialization
// ============================================================================

// MarshalJSON serializes the DataSet into the standard KQL JSON response format.
// Output: { "Tables": [ { ...table1... }, { ...table2... } ] }
func (ds *DataSet) MarshalJSON() ([]byte, error) {
	// Simple wrapper to match Kusto REST API V1 structure
	return json.Marshal(struct {
		Tables []*DataTable `json:"Tables"`
	}{
		Tables: ds.Tables,
	})
}

// UnmarshalJSON deserializes a DataSet from the standard KQL JSON format.
func (ds *DataSet) UnmarshalJSON(data []byte) error {
	var raw struct {
		Tables []*DataTable `json:"Tables"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	ds.Tables = raw.Tables
	return nil
}

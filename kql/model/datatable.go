package model

import (
	"bytes"
	"encoding/json"
)

// ============================================================================
// Result Structures
// ============================================================================

// DataTable represents a materialized result set.
// It is the standard output unit of a KQL query.
type DataTable struct {
	Name    string      // Table name (usually "Table_0" or "PrimaryResult")
	Columns []ColumnDef // Ordered list of column metadata
	Rows    []Row       // Materialized data
}

// ColumnDef describes a single column's metadata.
type ColumnDef struct {
	Name string // Column Name
	Type KType  // Data Type (inferred or explicit)
}

// InferColumnDefs builds ColumnDef metadata from a non-empty slice of rows.
func InferColumnDefs(rows []Row) []ColumnDef {
	if len(rows) == 0 {
		return nil
	}
	colNames := rows[0].Scheme.Names
	defs := make([]ColumnDef, len(colNames))
	for i, name := range colNames {
		kType := TypeDynamic
		for _, r := range rows {
			val := r.Values[i]
			if val != nil && !val.IsNull() {
				kType = val.Type()
				break
			}
		}
		defs[i] = ColumnDef{Name: name, Type: kType}
	}
	return defs
}

// ============================================================================
// Serialization (JSON)
// ============================================================================

// MarshalJSON implements custom JSON encoding to match standard KQL API format.
// Format: { "TableName": "...", "Columns": [...], "Rows": [[v1, v2], ...] }
func (dt *DataTable) MarshalJSON() ([]byte, error) {
	// Define a proxy struct to control the JSON layout
	type jsonCol struct {
		ColumnName string `json:"ColumnName"`
		DataType   string `json:"DataType"`
	}

	type jsonTable struct {
		TableName string          `json:"TableName"`
		Columns   []jsonCol       `json:"Columns"`
		Rows      [][]interface{} `json:"Rows"`
	}

	// 1. Convert Columns
	jCols := make([]jsonCol, len(dt.Columns))
	for i, c := range dt.Columns {
		jCols[i] = jsonCol{
			ColumnName: c.Name,
			DataType:   c.Type.String(),
		}
	}

	// 2. Convert Rows
	// We convert KValues to standard Go types (interface{}) for the JSON encoder
	jRows := make([][]interface{}, len(dt.Rows))
	for rIdx, row := range dt.Rows {
		rowVals := make([]interface{}, len(row.Values))
		for cIdx, val := range row.Values {
			rowVals[cIdx] = kValueToJSON(val)
		}
		jRows[rIdx] = rowVals
	}

	return json.Marshal(jsonTable{
		TableName: dt.Name,
		Columns:   jCols,
		Rows:      jRows,
	})
}

// UnmarshalJSON deserializes a DataTable from the standard KQL JSON format.
func (dt *DataTable) UnmarshalJSON(data []byte) error {
	type jsonCol struct {
		ColumnName string `json:"ColumnName"`
		DataType   string `json:"DataType"`
	}

	// Use json.Decoder with UseNumber to preserve integer precision
	var raw struct {
		TableName string            `json:"TableName"`
		Columns   []jsonCol         `json:"Columns"`
		Rows      []json.RawMessage `json:"Rows"`
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(&raw); err != nil {
		return err
	}

	dt.Name = raw.TableName

	// Build ColumnDefs
	dt.Columns = make([]ColumnDef, len(raw.Columns))
	for i, c := range raw.Columns {
		dt.Columns[i] = ColumnDef{
			Name: c.ColumnName,
			Type: ParseKType(c.DataType),
		}
	}

	// Build schema
	names := make([]string, len(dt.Columns))
	for i, c := range dt.Columns {
		names[i] = c.Name
	}
	schema := NewColumnInfo(names)

	// Parse rows
	dt.Rows = make([]Row, 0, len(raw.Rows))
	for _, rawRow := range raw.Rows {
		var cells []json.RawMessage
		if err := json.Unmarshal(rawRow, &cells); err != nil {
			return err
		}
		values := make([]KValue, len(dt.Columns))
		for i := range dt.Columns {
			if i >= len(cells) {
				values[i] = KNullValue
				continue
			}
			v, err := ParseCellValue(cells[i], dt.Columns[i].Type.String())
			if err != nil {
				return err
			}
			values[i] = v
		}
		dt.Rows = append(dt.Rows, NewRow(schema, values))
	}
	return nil
}

// kValueToJSON extracts raw values for JSON serialization.
func kValueToJSON(v KValue) interface{} {
	if v == nil || v.IsNull() {
		return nil
	}
	switch t := v.(type) {
	case *KString:
		return t.Val
	case *KInt:
		return t.Val
	case *KLong:
		return t.Val
	case *KReal:
		return t.Val
	case *KBool:
		return t.Val
	case *KDateTime:
		return t.Val // JSON encoder handles time.Time automatically
	case *KTimespan:
		return t.String()
	default:
		return v.String()
	}
}

package model

import "encoding/json"

// ============================================================================
// Metadata (Shared across all rows in a batch)
// ============================================================================

// ColumnInfo describes the schema of the row.
// It is immutable and should be shared among all rows produced by an iterator.
type ColumnInfo struct {
	Names   []string       // ordered column names
	nameIdx map[string]int // fast lookup by name
}

// NewColumnInfo creates a schema definition.
// Ideally created once in the Iterator constructor (e.g., in NewProjectIter).
func NewColumnInfo(names []string) *ColumnInfo {
	idx := make(map[string]int, len(names))
	for i, n := range names {
		idx[n] = i
	}
	return &ColumnInfo{
		Names:   names,
		nameIdx: idx,
	}
}

// Index returns the column index for a name, or -1 if not found.
func (ci *ColumnInfo) Index(name string) int {
	if i, ok := ci.nameIdx[name]; ok {
		return i
	}
	return -1
}

// Extend creates a new ColumnInfo by appending new columns.
// Used by the 'extend' operator to evolve the schema.
func (ci *ColumnInfo) Extend(newCols []string) *ColumnInfo {
	finalNames := make([]string, len(ci.Names)+len(newCols))
	copy(finalNames, ci.Names)
	copy(finalNames[len(ci.Names):], newCols)
	return NewColumnInfo(finalNames)
}

// ============================================================================
// Data (Unique per row)
// ============================================================================

// Row is now a struct containing the Schema pointer and the Data slice.
type Row struct {
	Scheme *ColumnInfo
	Values []KValue
}

// NewRow creates a row with the given schema and values.
func NewRow(scheme *ColumnInfo, values []KValue) Row {
	return Row{
		Scheme: scheme,
		Values: values,
	}
}

// Get retrieves a value by column name (O(1) map lookup).
// Used by ExprEvaluator when the index isn't compiled.
func (r Row) Get(name string) KValue {
	idx := r.Scheme.Index(name)
	if idx == -1 {
		return KNullValue
	}
	return r.Values[idx]
}

// GetAt retrieves a value by index (O(1) direct access).
// Used for high-performance access when indices are known.
func (r Row) GetAt(i int) KValue {
	if i < 0 || i >= len(r.Values) {
		return KNullValue
	}
	return r.Values[i]
}

// Clone creates a shallow copy of the row (values are shared if pointers).
// Necessary for operators that buffer rows (like Sort or Join).
func (r Row) Clone() Row {
	newVals := make([]KValue, len(r.Values))
	copy(newVals, r.Values)
	return Row{
		Scheme: r.Scheme,
		Values: newVals,
	}
}

// MarshalJSON produces a flat JSON object {"col1": val1, "col2": val2, ...}.
func (r Row) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{}, len(r.Values))
	for i, val := range r.Values {
		m[r.Scheme.Names[i]] = kValueToJSON(val)
	}
	return json.Marshal(m)
}

// ToBag converts the Row back to a KDynamicBag for legacy parts of the system
// or for "pack" operations.
func (r Row) ToBag() *KDynamicBag {
	bag := NewKDynamicBag()
	for i, val := range r.Values {
		bag.Set(r.Scheme.Names[i], val)
	}
	return bag
}

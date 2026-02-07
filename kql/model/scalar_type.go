package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// ============================================================================
// Core Definitions
// ============================================================================

// KType represents the enumeration of all KQL scalar data types.
type KType int

const (
	TypeUnknown  KType = iota
	TypeNull           // Generic null (dynamic)
	TypeBool           // bool
	TypeInt            // int32
	TypeLong           // int64
	TypeReal           // float64
	TypeString         // string
	TypeDateTime       // time.Time
	TypeTimespan       // time.Duration
	TypeGuid           // [16]byte
	TypeDecimal        // shopspring/decimal (128-bit fixed precision)
	TypeDynamic        // The complex container type
)

func (t KType) String() string {
	switch t {
	case TypeNull:
		return "null"
	case TypeBool:
		return "bool"
	case TypeInt:
		return "int"
	case TypeLong:
		return "long"
	case TypeReal:
		return "real"
	case TypeString:
		return "string"
	case TypeDateTime:
		return "datetime"
	case TypeTimespan:
		return "timespan"
	case TypeGuid:
		return "guid"
	case TypeDecimal:
		return "decimal"
	case TypeDynamic:
		return "dynamic"
	default:
		return "unknown"
	}
}

// KValue is the universal interface that all KQL values must implement.
type KValue interface {
	Type() KType
	String() string
	IsNull() bool
}

// ============================================================================
// The Generic Null Type
// ============================================================================

// KNull represents a generic, untyped null (usually from JSON 'null').
type KNull struct{}

var KNullValue = &KNull{}

func (k *KNull) Type() KType                  { return TypeNull }
func (k *KNull) IsNull() bool                 { return true }
func (k *KNull) String() string               { return "null" }
func (k *KNull) MarshalJSON() ([]byte, error) { return []byte("null"), nil }

// ============================================================================
// Primitive Scalars (Null-Safe)
// ============================================================================

type KBool struct {
	Val   bool
	Valid bool
}

func NewKBool(v bool) *KBool  { return &KBool{Val: v, Valid: true} }
func (k *KBool) Type() KType  { return TypeBool }
func (k *KBool) IsNull() bool { return !k.Valid }
func (k *KBool) String() string {
	if !k.Valid {
		return "bool(null)"
	}
	return strconv.FormatBool(k.Val)
}

type KInt struct {
	Val   int32
	Valid bool
}

func NewKInt(v int32) *KInt  { return &KInt{Val: v, Valid: true} }
func (k *KInt) Type() KType  { return TypeInt }
func (k *KInt) IsNull() bool { return !k.Valid }
func (k *KInt) String() string {
	if !k.Valid {
		return "int(null)"
	}
	return strconv.FormatInt(int64(k.Val), 10)
}

type KLong struct {
	Val   int64
	Valid bool
}

func NewKLong(v int64) *KLong { return &KLong{Val: v, Valid: true} }
func (k *KLong) Type() KType  { return TypeLong }
func (k *KLong) IsNull() bool { return !k.Valid }
func (k *KLong) String() string {
	if !k.Valid {
		return "long(null)"
	}
	return strconv.FormatInt(k.Val, 10)
}

type KReal struct {
	Val   float64
	Valid bool
}

func NewKReal(v float64) *KReal { return &KReal{Val: v, Valid: true} }
func (k *KReal) Type() KType    { return TypeReal }
func (k *KReal) IsNull() bool   { return !k.Valid }
func (k *KReal) String() string {
	if !k.Valid {
		return "real(null)"
	}
	return strconv.FormatFloat(k.Val, 'g', -1, 64)
}

type KString struct {
	Val   string
	Valid bool
}

func NewKString(v string) *KString { return &KString{Val: v, Valid: true} }
func (k *KString) Type() KType     { return TypeString }
func (k *KString) IsNull() bool    { return !k.Valid }
func (k *KString) String() string {
	if !k.Valid {
		return "string(null)"
	}
	return fmt.Sprintf("%q", k.Val)
}

// ============================================================================
// Complex Scalars
// ============================================================================

type KDateTime struct {
	Val   time.Time
	Valid bool
}

func NewKDateTime(t time.Time) *KDateTime { return &KDateTime{Val: t, Valid: true} }
func (k *KDateTime) Type() KType          { return TypeDateTime }
func (k *KDateTime) IsNull() bool         { return !k.Valid }
func (k *KDateTime) String() string {
	if !k.Valid {
		return "datetime(null)"
	}
	return fmt.Sprintf("datetime(%s)", k.Val.Format(time.RFC3339Nano))
}

func (k *KDateTime) MarshalJSON() ([]byte, error) {
	if !k.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(k.Val.Format(time.RFC3339Nano))
}

type KTimespan struct {
	Val   time.Duration
	Valid bool
}

func NewKTimespan(d time.Duration) *KTimespan { return &KTimespan{Val: d, Valid: true} }
func (k *KTimespan) Type() KType              { return TypeTimespan }
func (k *KTimespan) IsNull() bool             { return !k.Valid }
func (k *KTimespan) String() string {
	if !k.Valid {
		return "timespan(null)"
	}
	return fmt.Sprintf("timespan(%s)", formatTimespan(k.Val))
}

func (k *KTimespan) MarshalJSON() ([]byte, error) {
	if !k.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(formatTimespan(k.Val))
}

// formatTimespan formats a duration in KQL's d.hh:mm:ss.fffffff format.
func formatTimespan(d time.Duration) string {
	neg := d < 0
	if neg {
		d = -d
	}

	totalSeconds := int64(d / time.Second)
	days := totalSeconds / 86400
	remaining := totalSeconds % 86400
	hours := remaining / 3600
	remaining %= 3600
	minutes := remaining / 60
	seconds := remaining % 60
	// Sub-second portion in 100-nanosecond ticks (7 decimal digits)
	fraction := (int64(d) % int64(time.Second)) / 100

	var sb strings.Builder
	if neg {
		sb.WriteByte('-')
	}
	if days > 0 {
		fmt.Fprintf(&sb, "%d.", days)
	}
	fmt.Fprintf(&sb, "%02d:%02d:%02d", hours, minutes, seconds)
	if fraction > 0 {
		fmt.Fprintf(&sb, ".%07d", fraction)
	}
	return sb.String()
}

type KGuid struct {
	Val   [16]byte
	Valid bool
}

func NewKGuid(uuid [16]byte) *KGuid { return &KGuid{Val: uuid, Valid: true} }
func (k *KGuid) Type() KType        { return TypeGuid }
func (k *KGuid) IsNull() bool       { return !k.Valid }
func (k *KGuid) String() string {
	if !k.Valid {
		return "guid(null)"
	}
	return fmt.Sprintf("guid(%08x-%04x-%04x-%04x-%012x)", k.Val[0:4], k.Val[4:6], k.Val[6:8], k.Val[8:10], k.Val[10:])
}

// KDecimal represents a KQL decimal using shopspring/decimal.
type KDecimal struct {
	Val   decimal.Decimal
	Valid bool
}

func NewKDecimal(v decimal.Decimal) *KDecimal { return &KDecimal{Val: v, Valid: true} }

// NewKDecimalFromString is a helper to easily create decimals from strings
func NewKDecimalFromString(s string) (*KDecimal, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil, err
	}
	return &KDecimal{Val: d, Valid: true}, nil
}

func (k *KDecimal) Type() KType  { return TypeDecimal }
func (k *KDecimal) IsNull() bool { return !k.Valid }
func (k *KDecimal) String() string {
	if !k.Valid {
		return "decimal(null)"
	}
	return fmt.Sprintf("decimal(%s)", k.Val.String())
}

// ============================================================================
// Dynamic Types (Bag and Array)
// ============================================================================

type KeyValuePair struct {
	Key   string
	Value KValue
}

// KDynamicBag represents a KQL Property Bag (JSON Object).
// We use a slice of pairs to maintain insertion order, unlike a Go map.
type KDynamicBag struct {
	//Fields map[string]KValue
	Pairs []KeyValuePair
}

//func NewKDynamicBag() *KDynamicBag {
//	return &KDynamicBag{Fields: make(map[string]KValue)}
//}

func NewKDynamicBag() *KDynamicBag {
	// Pre-allocate a small capacity to avoid immediate resizing
	return &KDynamicBag{Pairs: make([]KeyValuePair, 0, 4)}
}

func (k *KDynamicBag) Type() KType  { return TypeDynamic }
func (k *KDynamicBag) IsNull() bool { return k.Pairs == nil }

func (k *KDynamicBag) Set(key string, val KValue) {
	/*
		if k.Fields == nil {
			k.Fields = make(map[string]KValue)
		}
		k.Fields[key] = val
	*/
	k.Pairs = append(k.Pairs, KeyValuePair{Key: key, Value: val})
}

func (k *KDynamicBag) Get(key string) KValue {
	if k.IsNull() {
		return KNullValue
	}
	// Iterate backwards to support "last-write-wins" semantics if duplicates exist
	for i := len(k.Pairs) - 1; i >= 0; i-- {
		if k.Pairs[i].Key == key {
			return k.Pairs[i].Value
		}
	}
	return KNullValue
}

func (k *KDynamicBag) String() string {
	if k.IsNull() {
		return "dynamic(null)"
	}
	var sb strings.Builder
	sb.WriteString("{")
	for i, pair := range k.Pairs {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%q: %s", pair.Key, pair.Value.String()))
	}
	sb.WriteString("}")
	return sb.String()
}

// MarshalJSON implementation that preserves order
func (k *KDynamicBag) MarshalJSON() ([]byte, error) {
	if k.IsNull() {
		return []byte("null"), nil
	}
	var sb strings.Builder
	sb.WriteString("{")
	for i, pair := range k.Pairs {
		if i > 0 {
			sb.WriteString(",")
		}
		// JSON Key
		sb.WriteString(fmt.Sprintf("%q:", pair.Key))

		// JSON Value
		valBytes, err := json.Marshal(kValueToGoInterface(pair.Value))
		if err != nil {
			return nil, err
		}
		sb.WriteString(string(valBytes))
	}
	sb.WriteString("}")
	return []byte(sb.String()), nil
}

type KDynamicArray struct {
	Elements []KValue
}

func NewKDynamicArray(cap int) *KDynamicArray {
	return &KDynamicArray{Elements: make([]KValue, 0, cap)}
}

func (k *KDynamicArray) Type() KType  { return TypeDynamic }
func (k *KDynamicArray) IsNull() bool { return k.Elements == nil }

func (k *KDynamicArray) Append(val KValue) {
	k.Elements = append(k.Elements, val)
}

func (k *KDynamicArray) String() string {
	if k.IsNull() {
		return "dynamic(null)"
	}
	var sb strings.Builder
	sb.WriteString("[")
	for i, val := range k.Elements {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(val.String())
	}
	sb.WriteString("]")
	return sb.String()
}

func (k *KDynamicArray) MarshalJSON() ([]byte, error) {
	out := make([]interface{}, len(k.Elements))
	for i, val := range k.Elements {
		out[i] = kValueToGoInterface(val)
	}
	return json.Marshal(out)
}

// ============================================================================
// Helpers
// ============================================================================

func kValueToGoInterface(k KValue) interface{} {
	if k == nil {
		return nil
	}
	// Check for our generic KNull type
	if _, ok := k.(*KNull); ok {
		return nil
	}
	if k.IsNull() {
		return nil
	}
	switch v := k.(type) {
	case *KBool:
		return v.Val
	case *KInt:
		return v.Val
	case *KLong:
		return v.Val
	case *KReal:
		return v.Val
	case *KDecimal:
		return v.Val // Return decimal.Decimal directly (standard practice)
	case *KString:
		return v.Val
	case *KDateTime:
		return v.Val
	case *KDynamicBag:
		return v
	case *KDynamicArray:
		return v
	default:
		return k.String()
	}
}

func ParseDynamicJSON(jsonStr string) (KValue, error) {
	var root interface{}
	// Standard unmarshal will treat numbers as float64.
	// To support decimals from JSON, you would need d := json.NewDecoder(...); d.UseNumber()
	// and then handle json.Number in ConvertInterfaceToKValue.
	if err := json.Unmarshal([]byte(jsonStr), &root); err != nil {
		return nil, err
	}
	return ConvertInterfaceToKValue(root), nil
}

// ParseKType returns the KType for a type name string.
func ParseKType(s string) KType {
	switch s {
	case "null":
		return TypeNull
	case "bool":
		return TypeBool
	case "int":
		return TypeInt
	case "long":
		return TypeLong
	case "real":
		return TypeReal
	case "string":
		return TypeString
	case "datetime":
		return TypeDateTime
	case "timespan":
		return TypeTimespan
	case "guid":
		return TypeGuid
	case "decimal":
		return TypeDecimal
	case "dynamic":
		return TypeDynamic
	default:
		return TypeUnknown
	}
}

// kvalueJSON is the type-tagged JSON envelope for KValue serialization.
type kvalueJSON struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value,omitempty"`
}

// MarshalKValue serializes a KValue to type-tagged JSON.
func MarshalKValue(v KValue) ([]byte, error) {
	if v == nil || v.IsNull() {
		return json.Marshal(kvalueJSON{Type: "null"})
	}
	var typeName string
	var raw interface{}
	switch x := v.(type) {
	case *KBool:
		typeName = "bool"
		raw = x.Val
	case *KInt:
		typeName = "int"
		raw = x.Val
	case *KLong:
		typeName = "long"
		raw = x.Val
	case *KReal:
		typeName = "real"
		raw = x.Val
	case *KString:
		typeName = "string"
		raw = x.Val
	case *KDateTime:
		typeName = "datetime"
		raw = x.Val.Format(time.RFC3339Nano)
	case *KTimespan:
		typeName = "timespan"
		raw = formatTimespan(x.Val)
	case *KDynamicBag:
		typeName = "bag"
		raw = x // uses KDynamicBag.MarshalJSON
	case *KDynamicArray:
		typeName = "array"
		raw = x // uses KDynamicArray.MarshalJSON
	default:
		typeName = "string"
		raw = v.String()
	}
	valBytes, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}
	return json.Marshal(kvalueJSON{Type: typeName, Value: valBytes})
}

// UnmarshalKValue deserializes a KValue from type-tagged JSON.
func UnmarshalKValue(data []byte) (KValue, error) {
	var env kvalueJSON
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, err
	}
	switch env.Type {
	case "null":
		return KNullValue, nil
	case "bool":
		var v bool
		if err := json.Unmarshal(env.Value, &v); err != nil {
			return nil, err
		}
		return NewKBool(v), nil
	case "int":
		var v int32
		if err := json.Unmarshal(env.Value, &v); err != nil {
			return nil, err
		}
		return NewKInt(v), nil
	case "long":
		var v int64
		if err := json.Unmarshal(env.Value, &v); err != nil {
			return nil, err
		}
		return NewKLong(v), nil
	case "real":
		var v float64
		if err := json.Unmarshal(env.Value, &v); err != nil {
			return nil, err
		}
		return NewKReal(v), nil
	case "string":
		var v string
		if err := json.Unmarshal(env.Value, &v); err != nil {
			return nil, err
		}
		return NewKString(v), nil
	case "datetime":
		var s string
		if err := json.Unmarshal(env.Value, &s); err != nil {
			return nil, err
		}
		t, err := time.Parse(time.RFC3339Nano, s)
		if err != nil {
			return nil, fmt.Errorf("UnmarshalKValue datetime: %w", err)
		}
		return NewKDateTime(t), nil
	case "timespan":
		var s string
		if err := json.Unmarshal(env.Value, &s); err != nil {
			return nil, err
		}
		d, err := parseTimespan(s)
		if err != nil {
			return nil, fmt.Errorf("UnmarshalKValue timespan: %w", err)
		}
		return NewKTimespan(d), nil
	case "bag":
		val, err := ParseDynamicJSON(string(env.Value))
		if err != nil {
			return nil, fmt.Errorf("UnmarshalKValue bag: %w", err)
		}
		return val, nil
	case "array":
		val, err := ParseDynamicJSON(string(env.Value))
		if err != nil {
			return nil, fmt.Errorf("UnmarshalKValue array: %w", err)
		}
		return val, nil
	default:
		return nil, fmt.Errorf("UnmarshalKValue: unknown type %q", env.Type)
	}
}

// parseTimespan parses a KQL timespan string in d.hh:mm:ss.fffffff format.
func parseTimespan(s string) (time.Duration, error) {
	neg := false
	if len(s) > 0 && s[0] == '-' {
		neg = true
		s = s[1:]
	}

	// Split days if present
	var days int64
	if i := strings.Index(s, "."); i >= 0 {
		// Check if the dot is before the first colon (days separator)
		if ci := strings.Index(s, ":"); ci < 0 || i < ci {
			_, err := fmt.Sscanf(s[:i], "%d", &days)
			if err != nil {
				return 0, fmt.Errorf("parseTimespan: invalid days: %w", err)
			}
			s = s[i+1:]
		}
	}

	// Parse hh:mm:ss part
	parts := strings.SplitN(s, ":", 3)
	if len(parts) != 3 {
		return 0, fmt.Errorf("parseTimespan: expected hh:mm:ss, got %q", s)
	}
	var hours, minutes int64
	_, err := fmt.Sscanf(parts[0], "%d", &hours)
	if err != nil {
		return 0, err
	}
	_, err = fmt.Sscanf(parts[1], "%d", &minutes)
	if err != nil {
		return 0, err
	}

	// Parse seconds and fractional part
	secParts := strings.SplitN(parts[2], ".", 2)
	var seconds int64
	_, err = fmt.Sscanf(secParts[0], "%d", &seconds)
	if err != nil {
		return 0, err
	}

	var fraction int64
	if len(secParts) == 2 {
		_, err = fmt.Sscanf(secParts[1], "%d", &fraction)
		if err != nil {
			return 0, err
		}
		// fraction is in 100-nanosecond ticks (7 digits)
		// Pad or truncate to 7 digits
		fStr := secParts[1]
		for len(fStr) < 7 {
			fStr += "0"
		}
		_, _ = fmt.Sscanf(fStr[:7], "%d", &fraction)
	}

	d := time.Duration(days)*24*time.Hour +
		time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second +
		time.Duration(fraction)*100 // 100ns ticks

	if neg {
		d = -d
	}
	return d, nil
}

// ParseCellValue converts a raw JSON value to a KValue given a type hint.
func ParseCellValue(raw json.RawMessage, dataType string) (KValue, error) {
	// Handle JSON null
	if string(raw) == "null" {
		return KNullValue, nil
	}
	switch dataType {
	case "bool":
		var v bool
		if err := json.Unmarshal(raw, &v); err != nil {
			return nil, err
		}
		return NewKBool(v), nil
	case "int":
		// Use json.Number to preserve precision
		var n json.Number
		if err := json.Unmarshal(raw, &n); err != nil {
			return nil, err
		}
		i, err := n.Int64()
		if err != nil {
			return nil, err
		}
		return NewKInt(int32(i)), nil
	case "long":
		var n json.Number
		if err := json.Unmarshal(raw, &n); err != nil {
			return nil, err
		}
		i, err := n.Int64()
		if err != nil {
			return nil, err
		}
		return NewKLong(i), nil
	case "real":
		var n json.Number
		if err := json.Unmarshal(raw, &n); err != nil {
			return nil, err
		}
		f, err := n.Float64()
		if err != nil {
			return nil, err
		}
		return NewKReal(f), nil
	case "string":
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return nil, err
		}
		return NewKString(s), nil
	case "datetime":
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return nil, err
		}
		t, err := time.Parse(time.RFC3339Nano, s)
		if err != nil {
			return nil, err
		}
		return NewKDateTime(t), nil
	case "timespan":
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return nil, err
		}
		d, err := parseTimespan(s)
		if err != nil {
			return nil, err
		}
		return NewKTimespan(d), nil
	case "dynamic":
		val, err := ParseDynamicJSON(string(raw))
		if err != nil {
			return nil, err
		}
		return val, nil
	default:
		// Fall back to generic JSON parsing
		val, err := ParseDynamicJSON(string(raw))
		if err != nil {
			return nil, err
		}
		return val, nil
	}
}

func ConvertInterfaceToKValue(v interface{}) KValue {
	switch val := v.(type) {
	case map[string]interface{}:
		bag := NewKDynamicBag()
		// Note: Iterating this map is random in Go.
		for k, child := range val {
			bag.Set(k, ConvertInterfaceToKValue(child))
		}
		return bag
	case []interface{}:
		arr := NewKDynamicArray(len(val))
		for _, child := range val {
			arr.Append(ConvertInterfaceToKValue(child))
		}
		return arr
	case string:
		return NewKString(val)
	case float64:
		if val == float64(int64(val)) {
			return NewKLong(int64(val))
		}
		return NewKReal(val)
	case bool:
		return NewKBool(val)
	case nil:
		return KNullValue
	default:
		return &KString{Valid: false}
	}
}

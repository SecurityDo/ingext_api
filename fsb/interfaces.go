package fsb

type IEntityInfoAPI interface {
	EntityCheck(entity, key string) bool
	LookupEntityRow(entity string, key string) map[string]interface{}
	LookupEntityKeyColumn(entityInfo string, keyColumn string, key string, valueColumn string) (bool, string)
}


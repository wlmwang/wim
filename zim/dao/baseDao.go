package dao

type baseDao struct {
	tableName string
}

func (b *baseDao) SetTableName(tableName string) {
	b.tableName = tableName
}

func (b *baseDao) GetTableName() string {
	return b.tableName
}

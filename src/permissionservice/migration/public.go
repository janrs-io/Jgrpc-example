package migration

// ID Primary ID
type ID struct {
	ID *uint32 `json:"id" gorm:"column:id;primaryKey;type:int(10);unique;autoIncrement;comment:primary id"`
}

// Timestamp create time & update time
type Timestamp struct {
	CreateTime *int32 `json:"create_time" gorm:"column:create_time;type:int(10);default:0;comment:create time'"`
	UpdateTime *int32 `json:"update_time" gorm:"column:update_time;type:int(10);default:0;comment:update time"`
}

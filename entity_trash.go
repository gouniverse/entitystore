package entitystore

import (
	"time"
)

// EntityTrash type
type EntityTrash struct {
	ID     string `gorm:"type:varchar(40);column:id;primary_key;"`
	Status string `gorm:"type:varchar(10);column:status;"`
	Type   string `gorm:"type:varchar(40);column:type;"`
	CreatedAt time.Time  `gorm:"type:datetime;column:created_at;DEFAULT NULL;"`
	UpdatedAt time.Time  `gorm:"type:datetime;column:updated_at;DEFAULT NULL;"`
	DeletedAt time.Time `gorm:"type:datetime;column:deleted_at;DEFAULT NULL;"`
	DeletedBy string `gorm:"type:varchar(40);column:deleted_by;DEFAULT NULL;"`
}

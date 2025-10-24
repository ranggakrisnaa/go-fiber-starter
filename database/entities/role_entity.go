package entities

import "github.com/google/uuid"

type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4();" json:"id"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Description string    `gorm:"type:varchar(255);not null" json:"description"`
	IsActive    bool      `gorm:"type:bool;not null" json:"is_active"`

	Timestamp
}

func (r Role) GetName() string {
	return r.Name
}

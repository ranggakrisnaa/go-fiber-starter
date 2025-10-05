package entities

import "github.com/google/uuid"

type UserRole struct {
	ID     uuid.UUID `json:"id" gorm:"primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	RoleID uuid.UUID `gorm:"type:uuid;not null" json:"role_id"`
	User   User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Role   Role      `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE"`

	Timestamp
}

package models

import "time"

// User maps to the canonical users table managed by SQL migrations.
type User struct {
	ID            string    `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Email         string    `gorm:"column:email;type:varchar(255);not null"`
	PasswordHash  *string   `gorm:"column:password_hash;type:varchar(255)"`
	GoogleSubject *string   `gorm:"column:google_subject;type:varchar(255)"`
	Name          string    `gorm:"column:name;type:varchar(255);not null"`
	AvatarURL     *string   `gorm:"column:avatar_url;type:text"`
	CreatedAt     time.Time `gorm:"column:created_at;type:timestamptz;not null"`
	UpdatedAt     time.Time `gorm:"column:updated_at;type:timestamptz;not null"`
}

func (User) TableName() string {
	return "users"
}

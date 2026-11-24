package models

import "time"

// RawResponse represents a raw HTTP response fetched from an external source (e.g. TSETMC) and archived verbatim.
type RawResponse struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Endpoint    string    `gorm:"type:text;not null" json:"endpoint"`
	URL         string    `gorm:"type:text;not null" json:"url"`
	InsCode     *string   `gorm:"type:text" json:"ins_code"`
	FetchedAt   time.Time `gorm:"type:timestamptz;not null" json:"fetched_at"`
	HttpStatus  int       `gorm:"type:integer;not null" json:"http_status"`
	ContentType *string   `gorm:"type:text" json:"content_type"`
	Body        string    `gorm:"type:text;not null" json:"body,omitempty"`
	BodySha256  string    `gorm:"type:text;not null" json:"body_sha256"`
	IsBlocked   bool      `gorm:"type:boolean;not null;default:false" json:"is_blocked"`
	CreatedAt   time.Time `json:"created_at"`
}

// TableName overrides the table name for GORM.
func (RawResponse) TableName() string {
	return "raw_responses"
}

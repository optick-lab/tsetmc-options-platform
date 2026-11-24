package models

import "time"

// UnderlyingAsset represents a spot asset on which option contracts are traded.
type UnderlyingAsset struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Symbol    string    `gorm:"type:text;unique;not null" json:"symbol"`
	Name      *string   `gorm:"type:text" json:"name,omitempty"`
	InsCode   *string   `gorm:"type:text;unique" json:"ins_code"`
	AssetType string    `gorm:"type:text;not null" json:"asset_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName overrides the table name for GORM.
func (UnderlyingAsset) TableName() string {
	return "underlying_assets"
}

package model

import (
	"time"

	"github.com/google/uuid"
)

func (d *Streamer) UUIDToBinary() *Streamer {
	d.BinaryId, _ = d.Id.MarshalBinary()
	return d
}

func (d *Streamer) BinaryToUUID() *Streamer {
	d.Id, _ = uuid.FromBytes(d.BinaryId)
	return d
}

type Streamer struct {
	Id         uuid.UUID `json:"userID" gorm:"-"`
	BinaryId   []byte    `json:"-" gorm:"column:userID"`
	OpenID     *string   `json:"openID" gorm:"column:openID"`
	AgencyID   *int      `json:"agencyID" gorm:"column:agencyID"`
	Name       *string   `json:"name" gorm:"column:name"`
	RegionCode *string   `json:"userRegion" gorm:"column:regionCode"`
	AccountID  *int      `json:"accountID" gorm:"column:accountID"`
	SyncTime   time.Time `json:"syncTime" gorm:"column:syncTime"`
}

type StreamerAgency struct {
	Id        int       `json:"id" gorm:"column:id"`
	Name      string    `json:"name" gorm:"column:name"`
	AccountID *int      `json:"accountID" gorm:"column:accountID"`
	SyncTime  time.Time `json:"syncTime" gorm:"column:syncTime"`
}

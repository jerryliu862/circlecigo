package model

import "time"

type UserToken struct {
	Token string `json:"token"`
}

type User struct {
	Id          int       `json:"uid" gorm:"column:id"`
	GoogleID    string    `json:"googleID" gorm:"column:googleID"`
	Name        string    `json:"name" gorm:"column:name"`
	Email       string    `json:"email" gorm:"column:email"`
	Status      int       `json:"status" gorm:"column:status" validate:"oneof=0 1 2"`
	CreateTime  time.Time `json:"createTime" gorm:"column:createTime"`
	CreatorUID  int       `json:"creatorUID" gorm:"column:creatorUID"`
	ModifyTime  time.Time `json:"modifyTime,omitempty" gorm:"column:modifyTime"`
	ModifierUID int       `json:"modifierUID,omitempty" gorm:"column:modifierUID"`
}

type UserAuth struct {
	Id         int       `json:"id,omitempty" gorm:"column:id"`
	Uid        int       `json:"uid,omitempty" gorm:"column:uid"`
	RegionCode string    `json:"regionCode" gorm:"column:regionCode"`
	AuthType   int       `json:"authType" gorm:"column:authType" validate:"oneof=1 2 3 4 99"`
	AuthLevel  int       `json:"authLevel" gorm:"column:authLevel" validate:"oneof=1 2 3"`
	SendNotify bool      `json:"sendNotify,omitempty" gorm:"column:sendNotify"`
	CreateTime time.Time `json:"-" gorm:"column:createTime"`
	CreatorUID int       `json:"-" gorm:"column:creatorUID"`
}

type UserDetail struct {
	User         `json:",inline"`
	CreatorName  string     `json:"creatorName" gorm:"column:name"`
	ModifierName string     `json:"modifierName" gorm:"column:name"`
	AuthList     []UserAuth `json:"authList" validate:"dive"`
}

type UserWithTotalCount struct {
	User
	TotalCount int `json:"-" gorm:"column:totalCount"`
}

type GoogleUser struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

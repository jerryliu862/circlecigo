package repository

import (
	"17live_wso_be/config"
	"17live_wso_be/internal/customError"
	"17live_wso_be/internal/model"
	"context"

	"gorm.io/gorm"
)

func (c *Client) GetUser(ctx context.Context, user model.User) ([]model.User, error) {
	var data []model.User

	if err := c.Database.Table(User).Where(user).Find(&data).Error; err != nil {
		log.Errorf("fail to get user: %v. %s", user, err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	return data, nil
}

func (c *Client) GetUserAuthLevel(ctx context.Context, uid int, region string, authType int) (int, error) {
	var authLevel int

	if err := c.Database.Table(UserAuth).Model(&model.UserAuth{}).Where("uid = ? AND regionCode = ? AND authType = ?", uid, region, authType).Select("authLevel").First(&authLevel).Error; err != nil {
		log.Errorf("fail to get user auth level: uid %d, region %s, authType %d. %s", uid, region, authType, err.Error())
		return authLevel, customError.New(customError.DatabaseError)
	}

	return authLevel, nil
}

func (c *Client) GetUserAuthRegion(ctx context.Context, uid int, authType int, authLevel int) ([]string, error) {
	var regions []string

	if err := c.Database.Table(UserAuth).Model(&model.UserAuth{}).Where("uid = ? AND authType = ? AND authLevel >= ?", uid, authType, authLevel).Select("regionCode").Find(&regions).Error; err != nil {
		log.Errorf("fail to get user auth region: uid %d, authType %d, authLevel %d. %s", uid, authType, authLevel, err.Error())
		return regions, customError.New(customError.DatabaseError)
	}

	return regions, nil
}

func (c *Client) GetAdminMailList(ctx context.Context) ([]string, error) {
	var mailList []string

	if err := c.Database.Raw("SELECT DISTINCT u.email FROM UserAuth ua LEFT JOIN User u on u.id = ua.uid WHERE u.status = ? AND ua.sendNotify = ? AND ua.regionCode = ? AND ua.authType = ? AND ua.authLevel >= ?", model.UserStatusActive, true, model.RegionAll, model.UserAuthTypeSystem, model.UserAuthLevelEdit).Find(&mailList).Error; err != nil {
		log.Errorf("fail to get admin mail list: %s", err.Error())
		return mailList, customError.New(customError.DatabaseError)
	}

	return mailList, nil
}

func (c *Client) CreateUser(ctx context.Context, data model.User) (model.User, error) {
	if err := c.Database.Table(User).Create(&data).Error; err != nil {
		log.Errorf("fail to create user: %v. %s", data, err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	log.Infof("user created: %v", data)

	return data, nil
}

func (c *Client) CreateUserWithAuth(ctx context.Context, data model.UserDetail) error {
	err := c.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(User).Create(&data.User).Error; err != nil {
			return err
		}

		if len(data.AuthList) != 0 {
			for i := range data.AuthList {
				data.AuthList[i].Uid = data.User.Id
				data.AuthList[i].CreatorUID = data.User.CreatorUID
				data.AuthList[i].CreateTime = data.User.CreateTime
			}

			if err := tx.Table(UserAuth).Create(&data.AuthList).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Errorf("fail to create user with auth: %v. %s", data, err.Error())
		return customError.New(customError.DatabaseError)
	}

	log.Infof("user with auth created: %v", data)

	return nil
}

func (c *Client) UpdateUser(ctx context.Context, data model.User) error {
	if err := c.Database.Table(User).Save(&data).Error; err != nil {
		log.Errorf("fail to update user: %v. %s", data, err.Error())
		return customError.New(customError.DatabaseError)
	}

	log.Infof("user updated: %v", data)

	return nil
}

func (c *Client) UpdateUserWithAuth(ctx context.Context, data model.UserDetail) error {
	err := c.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(User).Save(&data.User).Error; err != nil {
			return err
		}

		if err := tx.Table(UserAuth).Where("uid = ?", data.User.Id).Delete(&model.UserAuth{}).Error; err != nil {
			return err
		}

		if len(data.AuthList) != 0 {
			for i := range data.AuthList {
				data.AuthList[i].Uid = data.User.Id
				data.AuthList[i].CreatorUID = data.User.ModifierUID
				data.AuthList[i].CreateTime = data.User.ModifyTime
			}

			if err := tx.Table(UserAuth).Create(&data.AuthList).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Errorf("fail to update user with auth: %v. %s", data, err.Error())
		return customError.New(customError.DatabaseError)
	}

	log.Infof("user with auth updated: %v", data)

	return nil
}

func (c *Client) ListUser(ctx context.Context, page model.PageFilter) ([]model.UserWithTotalCount, int, error) {
	var data []model.UserWithTotalCount
	var total int

	if err := c.Database.Raw("SELECT *, (SELECT COUNT(*) FROM User WHERE email <> ?) AS totalCount FROM User WHERE email <> ? ORDER BY id LIMIT ? OFFSET ?", config.New().User.Admin, config.New().User.Admin, page.PageSize, page.PageSize*(page.PageNum-1)).Find(&data).Error; err != nil {
		log.Errorf("fail to list user with filter %v: %s", page, err.Error())
		return nil, total, customError.New(customError.DatabaseError)
	}

	if len(data) != 0 {
		total = data[0].TotalCount
	}

	return data, total, nil
}

func (c *Client) GetUserDetail(ctx context.Context, data model.UserDetail) (model.UserDetail, error) {
	err := c.Database.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(User).Model(&data.User).Where("id = ?", data.User.CreatorUID).Select("name").First(&data.CreatorName).Error; err != nil {
			return err
		}

		if err := tx.Table(User).Model(&data.User).Where("id = ?", data.User.ModifierUID).Select("name").First(&data.ModifierName).Error; err != nil {
			return err
		}

		if err := tx.Table(UserAuth).Select("regionCode", "authType", "authLevel", "createTime").Where("uid = ?", data.User.Id).Find(&data.AuthList).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Errorf("fail to get user datail: %v. %s", data, err.Error())
		return data, customError.New(customError.DatabaseError)
	}

	return data, nil
}

package service

import (
	"17live_wso_be/config"
	"17live_wso_be/internal/customError"
	"17live_wso_be/internal/model"
	"17live_wso_be/util"
	"context"
	"time"
)

func (c *Client) Login(ctx context.Context, token *model.UserToken) (*model.UserToken, error) {
	log.Infof("service user login: %v", token)

	gUser, err := c.GoogleClient.AuthGoogleUser(ctx, token.Token)
	if err != nil {
		return nil, customError.New(customError.UserGoogleAuthFail)
	}

	if gUser.Sub == "" || gUser.Email == "" {
		log.Warnf("empty google user sub or email: %v", gUser)
		return nil, customError.New(customError.UserGoogleAuthFail)
	}

	var user model.User

	u, err := c.GetUser(ctx, model.User{GoogleID: gUser.Sub})
	if err != nil {
		return nil, err
	}

	if len(u) == 0 {
		user, err = c.activateUser(ctx, gUser)
		if err != nil {
			return nil, err
		}
	} else {
		user = u[0]
		if user.Email != gUser.Email {
			log.Warnf("user email incompatible: expected %s, query got %s", gUser.Email, user.Email)
			return nil, customError.New(customError.DatabaseError)
		} else if user.Status != model.UserStatusActive {
			log.Warnf("user login with invalid status: %d, %s", user.Status, user.Email)
			return nil, customError.New(customError.UserNotActive)
		} else if user.Name != gUser.Name {
			user.Name = gUser.Name
			if err := c.UpdateUser(ctx, user.Id, user); err != nil {
				return nil, err
			}
		}
	}

	return c.GetUserToken(ctx, user.Id)
}

func (c *Client) GetUser(ctx context.Context, user model.User) ([]model.User, error) {
	log.Infof("service get user: %v", user)
	return c.RepositoryClient.GetUser(ctx, user)
}

func (c *Client) ListUser(ctx context.Context, page model.PageFilter) ([]model.UserWithTotalCount, int, error) {
	log.Infof("service list user with filter: %v", page)
	return c.RepositoryClient.ListUser(ctx, page)
}

func (c *Client) GetUserDetail(ctx context.Context, user model.User) (model.UserDetail, error) {
	log.Infof("service get user detail: %v", user)
	return c.RepositoryClient.GetUserDetail(ctx, model.UserDetail{User: user})
}

func (c *Client) GetUserAuthRegion(ctx context.Context, uid int, authType int, authLevel int) ([]string, error) {
	log.Infof("service get user auth region: uid %d, authType %d, authLevel %d", uid, authType, authLevel)

	if c.PermissionCheck(ctx, uid, model.RegionAll, model.UserAuthTypeSystem, model.UserAuthLevelEdit) {
		log.Infof("user %d has got admin permission, return all regions", uid)
		return c.RepositoryClient.GetAllRegionCode(ctx)
	}

	return c.RepositoryClient.GetUserAuthRegion(ctx, uid, authType, authLevel)
}

func (c *Client) GetAdminMailList(ctx context.Context) ([]string, error) {
	log.Infof("service get admin mail list")
	return c.RepositoryClient.GetAdminMailList(ctx)
}

func (c *Client) CreateUser(ctx context.Context, creator int, data model.User) (model.User, error) {
	log.Infof("service create user: %v", data)

	time := time.Now().UTC()
	data.CreateTime = time
	data.ModifyTime = time
	data.CreatorUID = creator
	data.ModifierUID = creator

	return c.RepositoryClient.CreateUser(ctx, data)
}

func (c *Client) CreateUserWithAuth(ctx context.Context, creator int, data model.UserDetail) error {
	log.Infof("service create user with auth: %v", data)

	time := time.Now().UTC()
	data.User.CreateTime = time
	data.User.ModifyTime = time
	data.User.CreatorUID = creator
	data.User.ModifierUID = creator

	return c.RepositoryClient.CreateUserWithAuth(ctx, data)
}

func (c *Client) UpdateUser(ctx context.Context, modifier int, data model.User) error {
	log.Infof("service update user: %v", data)

	data.ModifyTime = time.Now().UTC()
	data.ModifierUID = modifier

	return c.RepositoryClient.UpdateUser(ctx, data)
}

func (c *Client) UpdateUserWithAuth(ctx context.Context, modifier int, data model.UserDetail) error {
	log.Infof("service update user with auth: %v", data)

	data.User.ModifyTime = time.Now().UTC()
	data.User.ModifierUID = modifier

	return c.RepositoryClient.UpdateUserWithAuth(ctx, data)
}

func (c *Client) activateUser(ctx context.Context, gUser model.GoogleUser) (model.User, error) {
	log.Infof("start activate user: %s", gUser.Email)

	var user model.User

	u, err := c.GetUser(ctx, model.User{Email: gUser.Email})
	if err != nil {
		return user, err
	}

	if len(u) == 0 {
		if util.AdmitEmailDomain(config.New().User.Domains, gUser.Email) {
			admin, err := c.GetUser(ctx, model.User{Email: config.New().User.Admin})
			if err != nil {
				return user, err
			} else if len(admin) != 1 {
				log.Warnf("fail to get system admin user")
				return user, customError.New(customError.DatabaseError)
			}
			user.GoogleID = gUser.Sub
			user.Email = gUser.Email
			user.Name = gUser.Name
			user.Status = model.UserStatusActive
			return c.CreateUser(ctx, admin[0].Id, user)
		} else {
			log.Warnf("user with outside domain does not exist: %s", gUser.Email)
			return user, customError.New(customError.UserEmailDomainInvalid)
		}
	}

	user = u[0]
	if user.Status != model.UserStatusInit {
		log.Warnf("user status should be 0, but got: %d, %s", user.Status, user.Email)
		return user, customError.New(customError.DatabaseError)
	}
	user.GoogleID = gUser.Sub
	user.Name = gUser.Name
	user.Status = model.UserStatusActive

	return user, c.UpdateUser(ctx, user.Id, user)
}

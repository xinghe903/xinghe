package biz

import (
	"auth/internal/biz/auth"
	"auth/internal/biz/auth/token"
	"auth/internal/biz/po"
	"auth/internal/biz/repo"
	"auth/internal/conf"
	"context"
	"strconv"
	"time"

	hashid "github.com/xinghe903/xinghe/pkg/distribute/id"
	"golang.org/x/exp/rand"

	authpb "auth/api/auth/v1"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/xinghe903/xinghe/pkg/encrypt"
)

type AuthUsecase struct {
	log   *log.Helper
	uRepo repo.UserRepo
	aRepo repo.AuthRepo
	enc   *encrypt.EncryptAes
	au    auth.Auth
}

func NewAuthUsecase(c *conf.Config, logger log.Logger, u repo.UserRepo, snow *hashid.Sonyflake,
	aRepo repo.AuthRepo) *AuthUsecase {
	return &AuthUsecase{
		log:   log.NewHelper(logger),
		uRepo: u,
		enc:   encrypt.NewEncryptAes(c.EncryptKey),
		au:    token.NewToken(snow, aRepo),
		aRepo: aRepo,
	}
}

func (a *AuthUsecase) Register(ctx context.Context, info *po.User) (string, error) {
	users, err := a.uRepo.List(ctx, &po.PageQuery[po.User]{Condition: &po.User{Name: info.Name}}, "")
	if err != nil {
		a.log.WithContext(ctx).Errorf("list user: %v", err.Error())
		return "", authpb.ErrorCreateUser("创建用户失败 %s", info.Name)
	}
	if len(users.Data) > 0 {
		return "", authpb.ErrorUsernameRepeat("用户名重复 %s", info.Name)
	}
	info.Password, err = a.enc.Encrypt(info.Password)
	id, err := a.uRepo.Create(ctx, info)
	if err != nil {
		a.log.WithContext(ctx).Errorf("create user: %v", err.Error())
		return "", authpb.ErrorCreateUser("创建用户失败 %s", info.Name)
	}
	user, err := a.uRepo.Get(ctx, id)
	if err != nil {
		a.log.WithContext(ctx).Errorf("get user: %v", err.Error())
		return "", authpb.ErrorCreateUser("创建用户失败 %s", info.Name)
	}
	a.aRepo.Create(ctx, &po.Auth{
		Name:      user.Name,
		NickName:  user.NickName,
		Code:      user.InstanceId,
		Status:    po.StatusUserLogout,
		ExpiredAt: time.Now(),
	})
	return id, nil
}

func (a *AuthUsecase) Login(ctx context.Context, u *po.User) (string, error) {
	users, err := a.uRepo.List(ctx, &po.PageQuery[po.User]{Condition: &po.User{Name: u.Name}}, "")
	if err != nil || len(users.Data) == 0 {
		a.log.WithContext(ctx).Errorf("用户名错误 %s,  %v", u.Name, err)
		return "", authpb.ErrorUserOrPasswordInvalid("用户名或密码错误")
	}
	user := users.Data[0]
	pdText, err := a.enc.Decrypt(user.Password)
	if err != nil || pdText != u.Password {
		a.log.WithContext(ctx).Warnf("密码错误")
		return "", authpb.ErrorUserOrPasswordInvalid("用户名或密码错误")
	}

	// 检查用户登录状态
	var authUser *po.Auth
	if authUsers, err := a.aRepo.List(ctx, &po.PageQuery[po.Auth]{
		Condition: &po.Auth{Code: user.InstanceId},
	}); err != nil || authUsers.Total != 1 {
		message := strconv.FormatInt(authUsers.Total, 10)
		if err != nil {
			message = err.Error()
		}
		a.log.WithContext(ctx).Errorf("list auth: %v", message)
		return "", authpb.ErrorLoginError("账号异常")
	}
	// 用户可以覆盖登录
	// 清空过期token
	a.aRepo.ClearExpiredToken(ctx, time.Now())
	// 生成token
	token, err := a.generateToken(ctx, user.InstanceId)
	if err != nil {
		return "", err
	}
	// 更新登录状态
	authUser.Token = token
	authUser.Status = po.StatusUserLogin
	authUser.ExpiredAt = time.Now().Add(time.Minute * 60) // 60分钟过期
	a.aRepo.Update(ctx, authUser)
	return token, nil
}

func (a *AuthUsecase) generateToken(ctx context.Context, userId string) (string, error) {
	retry := 10
	for i := 0; i < retry; i++ {
		// 生成token
		token, err := a.au.GenerateToken(userId)
		if err != nil {
			a.log.WithContext(ctx).Errorf("generate token: %v", err.Error())
			return "", authpb.ErrorLoginError("生成token失败")
		}
		if authUsers, err := a.aRepo.List(ctx, &po.PageQuery[po.Auth]{
			Condition: &po.Auth{Token: token},
		}); err != nil {
			a.log.WithContext(ctx).Errorf("query token: %v", err.Error())
			return "", authpb.ErrorLoginError("生成token失败")
		} else if authUsers.Total > 0 {
			continue
		}
		return token, nil
	}
	a.log.WithContext(ctx).Errorf("generate token failed retry count: %d", retry)
	return "", authpb.ErrorLoginError("生成token失败")
}

func (a *AuthUsecase) Logout(ctx context.Context, token string) error {
	if t := ctx.Value("access_token"); t != "" {
		token, _ = t.(string)
	}
	var authUser *po.Auth
	if authUsers, err := a.aRepo.List(ctx, &po.PageQuery[po.Auth]{
		Condition: &po.Auth{Token: token},
	}); err != nil || authUsers.Total != 1 {
		message := strconv.FormatInt(authUsers.Total, 10)
		if err != nil {
			message = err.Error()
		}
		a.log.WithContext(ctx).Errorf("logout auth: %v", message)
		return authpb.ErrorIllegalToken("非法token %s", token)
	} else {
		authUser = authUsers.Data[0]
	}
	// 更新登录状态
	authUser.Status = po.StatusUserLogout
	a.aRepo.Update(ctx, authUser)
	return nil
}

func (a *AuthUsecase) Auth(ctx context.Context, token string) (*po.User, error) {
	if t := ctx.Value("access_token"); t != "" {
		token, _ = t.(string)
	}

	var authUser *po.Auth
	if authUsers, err := a.aRepo.List(ctx, &po.PageQuery[po.Auth]{
		Condition: &po.Auth{Token: token},
	}); err != nil || authUsers.Total != 1 {
		message := strconv.FormatInt(authUsers.Total, 10)
		if err != nil {
			message = err.Error()
		}
		a.log.WithContext(ctx).Warnf("auth: %v", message)
		return nil, authpb.ErrorIllegalToken("非法token %s", token)
	} else {
		authUser = authUsers.Data[0]
	}
	if authUser.ExpiredAt.After(time.Now()) && authUser.Status == po.StatusUserLogin {
		return &po.User{
			Name:       authUser.Name,
			NickName:   authUser.NickName,
			InstanceId: authUser.Code,
		}, nil
	}
	a.log.WithContext(ctx).Warnf("token is expired %s", token)
	authUser.Status = po.StatusUserLogout
	a.aRepo.Update(ctx, authUser)
	return nil, authpb.ErrorTokenExpired("token已过期")
}

func (a *AuthUsecase) GetUserById(ctx context.Context, id string) (*po.User, error) {
	user, err := a.uRepo.Get(ctx, id)
	if err != nil {
		a.log.WithContext(ctx).Errorf("get user: %v", err.Error())
		return nil, authpb.ErrorUserInfo("获取用户信息失败")
	}
	return user, nil
}

func (a *AuthUsecase) ListUser(ctx context.Context, cond *po.PageQuery[po.User], username string) (*po.SearchList[po.User], error) {
	cond.Sort = []map[string]string{{"updated_at": "desc"}}
	list, err := a.uRepo.List(ctx, cond, username)
	if err != nil {
		a.log.WithContext(ctx).Errorf("list user: %v", err.Error())
		return nil, authpb.ErrorUserInfo("获取用户列表失败")
	}
	return list, nil
}

func (a *AuthUsecase) CreateUser(ctx context.Context, req *po.User) (*po.User, error) {
	textPassword := RandStringRunes(8) // 生成默认密码  固定长度8位
	req.Password = textPassword
	id, err := a.Register(ctx, req)
	if err != nil {
		a.log.WithContext(ctx).Errorf("create user: %v", err.Error())
		return nil, err
	}
	rsp, err := a.GetUserById(ctx, id)
	rsp.Password = textPassword
	return rsp, nil
}

func (a *AuthUsecase) UpdateUser(ctx context.Context, req *po.User) error {
	err := a.uRepo.Update(ctx, &po.User{
		InstanceId: req.InstanceId,
		NickName:   req.NickName,
		Email:      req.Email,
		Phone:      req.Phone,
	})
	if err != nil {
		a.log.WithContext(ctx).Errorf("update user: %v", err.Error())
		return authpb.ErrorUpdateUser("更新用户信息失败")
	}
	err = a.aRepo.UpdateByCode(ctx, &po.Auth{Code: req.InstanceId, NickName: req.NickName})
	if err != nil {
		a.log.WithContext(ctx).Errorf("update user: %v", err.Error())
		return authpb.ErrorUpdateUser("更新用户信息失败")
	}
	return nil
}

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

var letterRunes = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

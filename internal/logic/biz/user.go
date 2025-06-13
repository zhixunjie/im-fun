package biz

import (
	"context"
	"errors"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zhixunjie/im-fun/internal/logic/data"
	"github.com/zhixunjie/im-fun/internal/logic/data/cache"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/generate/model"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/request"
	"github.com/zhixunjie/im-fun/internal/logic/data/ent/response"
	"github.com/zhixunjie/im-fun/pkg/gmodel"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type UserUseCase struct {
	userRepo  *data.UserRepo
	userCache *cache.UserCache
}

func NewUserUseCase(userRepo *data.UserRepo, userCache *cache.UserCache) *UserUseCase {
	return &UserUseCase{
		userRepo:  userRepo,
		userCache: userCache,
	}
}

func (b *UserUseCase) Login(ctx context.Context, req *request.LoginReq) (rsp *response.LoginRsp, err error) {
	rsp = new(response.LoginRsp)
	switch req.AccountType {
	case gmodel.AccountTypeDID:
		// do nothing
	case gmodel.AccountTypeMobile:
		// TODO:check mobile code
	default:
		err = fmt.Errorf("invalid account type")
		return
	}

	// get user info
	user, err := b.getUser(ctx, req)
	if err != nil {
		err = fmt.Errorf("getUser err: %w", err)
		return
	}
	// get token
	token, err := b.genToken(user.ID)
	if err != nil {
		err = fmt.Errorf("genToken err: %w", err)
		return
	}
	// set token cache
	err = b.userCache.SetToken(ctx, user.ID, token)
	if err != nil {
		err = fmt.Errorf("cacheSetToken err: %w", err)
		return
	}

	rsp.Data = &response.LoginData{
		User:  user,
		Token: token,
	}

	return
}

func (b *UserUseCase) getUser(ctx context.Context, req *request.LoginReq) (user *model.User, err error) {
	// try to get user
	user, tmpErr := b.userRepo.GetUserByAccount(uint32(req.AccountType), req.AccountID)
	if tmpErr != nil && !errors.Is(tmpErr, gorm.ErrRecordNotFound) {
		err = fmt.Errorf("GetUserByAccount err: %w", tmpErr)
		return

	}
	if user != nil {
		return
	}
	// if not found, create new user
	user, err = b.genUser(ctx, req)
	if err != nil {
		err = fmt.Errorf("genUser err: %w", err)
		return
	}

	return
}

func (b *UserUseCase) genUser(ctx context.Context, req *request.LoginReq) (user *model.User, err error) {
	var userType gmodel.UserType
	accountType := req.AccountType
	if accountType == gmodel.AccountTypeDID {
		userType = gmodel.UserTypeVisitor
	} else {
		userType = gmodel.UserTypeNormal
	}

	sexType := b.RandomSexType()
	user = &model.User{
		UserType:    uint32(userType),
		AccountType: uint32(req.AccountType),
		AccountID:   req.AccountID,
		Nickname:    b.RandomNickName(sexType),
		Avatar:      b.RandomAvatar(),
		Sex:         uint32(sexType),
		//CreatedAt:   0,
		//UpdatedAt:   0,
	}

	err = b.userRepo.CreateUser(user)
	if err != nil {
		err = fmt.Errorf("CreateUser err: %w", err)
		return
	}
	return
}

// RandomSexType 随机生成性别
func (b *UserUseCase) RandomSexType() gmodel.SexType {
	randNum := rand.Intn(3)
	return gmodel.SexType(randNum)
}

// RandomNickName 随机生成昵称
func (b *UserUseCase) RandomNickName(sexType gmodel.SexType) string {
	m := map[gmodel.SexType]int{
		gmodel.SexTypeMale:    randomdata.Male,
		gmodel.SexTypeFemale:  randomdata.Female,
		gmodel.SexTypeUnknown: randomdata.RandomGender,
	}
	return randomdata.FullName(m[sexType])
}

// RandomAvatar 随机生成头像
func (b *UserUseCase) RandomAvatar() string {
	randNum := rand.Intn(15)
	return fmt.Sprintf("img/avatar/%v.jpeg", randNum)
}

// RefreshToken TODO
func (b *UserUseCase) RefreshToken(ctx context.Context, req *request.RefreshTokenReq) (rsp *response.RefreshTokenRsp, err error) {
	rsp = new(response.RefreshTokenRsp)

	return
}

// genToken 生成 token
func (b *UserUseCase) genToken(uid uint64) (string, error) {
	// 初始化：claims
	claims := &gmodel.AuthClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "",                                                     // 签名颁发者
			Subject:   "",                                                     // 签名主体
			IssuedAt:  jwt.NewNumericDate(time.Now()),                         // 签发时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 过期时间
		},
		AuthUserInfo: &gmodel.AuthUserInfo{
			Uid: uid,
		},
	}

	// 定义初始参数：采用哪种算法？claim 的值
	key := []byte("Kfv5opY4i6bYuUG")
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	// 开始生成 JWT ：先生成第一部分和第二部分，然后生成第三部分
	tokenStr, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// checkToken 检查 token 是否有效
func (b *UserUseCase) checkToken(tokenStr string) (*gmodel.AuthClaims, error) {
	// parse token
	key := []byte("Kfv5opY4i6bYuUG")
	token, err := jwt.ParseWithClaims(tokenStr, &gmodel.AuthClaims{}, func(token *jwt.Token) (i any, err error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	// check token
	if claims, ok := token.Claims.(*gmodel.AuthClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

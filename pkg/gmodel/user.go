package gmodel

type UserType int

const (
	UserTypeUnknown UserType = 0
	UserTypeVisitor UserType = 1 // 游客
	UserTypeNormal  UserType = 2 // 用户
)

type AccountType uint32

const (
	AccountTypeDID       AccountType = 1 // 设备
	AccountTypeMobile    AccountType = 2 // 手机号码
	AccountTypeInstagram AccountType = 3 // Instagram
	AccountTypeFacebook  AccountType = 4 // Facebook
	AccountTypeGoogle    AccountType = 5 // Google
	AccountTypeTwitter   AccountType = 6 // Twitter
	AccountTypeApple     AccountType = 7 // Apple
)

type SexType int

const (
	SexTypeUnknown SexType = 0
	SexTypeMale    SexType = 1 // 男
	SexTypeFemale  SexType = 2 // 女
)

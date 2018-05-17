package modules

import (
	"errors"
	"fmt"

	"time"

	"github.com/go-redis/redis"
	"github.com/hfdend/cxr/cli"
	"github.com/hfdend/cxr/models"
	"github.com/hfdend/cxr/utils"
)

type passport int

var Passport passport

const (
	KEY_VerificationCode = "verification_code_"
	KEY_Register         = "register"
)

func (p passport) Register(phone, code, password string) (*models.User, error) {
	if err := p.CheckVerificationCode(phone, code, KEY_Register); err != nil {
		return nil, err
	}
	user := new(models.User)
	user.Phone = phone
	user.Password = utils.AesEncode(password)
	if n, err := user.Insert(); err != nil {
		return nil, err
	} else if n == 0 {
		return nil, errors.New("手机号已注册")
	}
	return user, nil
}

func (p passport) SendRegisterCode(phone string) (code string, err error) {
	var u *models.User
	if u, err = models.UserDefault.GetByPhone(phone); err != nil {
		return
	} else if u.ID != 0 {
		err = errors.New("此号码以及被注册")
		return
	}
	code = fmt.Sprintf("%0.4d", utils.RandInterval(0, 10000))
	// TODO 执行发送验证码
	code = "1234"
	if err = p.SaveVerificationCode(phone, code, KEY_Register, 10*time.Minute); err != nil {
		return
	}
	return
}

func (p passport) Login(phone, password string) (*models.Token, error) {
	user, err := models.UserDefault.GetByPhone(phone)
	if err != nil {
		return nil, err
	} else if user.ID == 0 {
		return nil, errors.New("该手机号暂未注册")
	}
	return nil, nil
}

func (passport) SaveVerificationCode(phone, code, typ string, expiration time.Duration) error {
	key := fmt.Sprintf("%s%s_%s", KEY_VerificationCode, typ, phone)
	return cli.Redis.Set(key, code, expiration).Err()
}

func (p passport) CheckVerificationCode(phone, code, typ string) error {
	sCode, err := p.GetVerificationCode(phone, typ)
	if err != nil {
		return err
	}
	if sCode != code {
		return errors.New("验证码错误")
	}
	return nil
}

func (passport) GetVerificationCode(phone, typ string) (code string, err error) {
	key := fmt.Sprintf("%s%s_%s", KEY_VerificationCode, typ, phone)
	if code, err = cli.Redis.Get(key).Result(); err != nil && err == redis.Nil {
		err = errors.New("验证码已过期")
	}
	return
}

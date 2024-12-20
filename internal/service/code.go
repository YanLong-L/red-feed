package service

import (
	"context"
	"math/rand"
	"red-feed/internal/repository"
	"red-feed/internal/service/sms"
)

const (
	tplId = ""
)

var (
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
)

type CodeService interface {
	Send(ctx context.Context, biz, phone string) error
	Verfiy(ctx context.Context, biz, phone, code string) (bool, error)
}

type codeService struct {
	smsSvc   sms.Service
	codeRepo repository.CachedCodeRepository
}

func NewCodeService(smsSvc sms.Service, codeRepo repository.CachedCodeRepository) CodeService {
	return &codeService{
		smsSvc:   smsSvc,
		codeRepo: codeRepo,
	}
}

func (cs *codeService) Send(ctx context.Context, biz, phone string) error {
	// 生成一个验证码
	code := cs.generatCode()
	// 缓存这个验证码
	err := cs.codeRepo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	// 缓存成功，发送验证码
	err = cs.smsSvc.Send(ctx, tplId, []string{code}, phone)
	if err != nil {
		return err
	}
	return nil
}

func (cs *codeService) Verfiy(ctx context.Context, biz, phone, code string) (bool, error) {
	return cs.codeRepo.Verify(ctx, biz, phone, code)
}
 
func (cs *codeService) generatCode() string {
	const letterBytes = "0123456789" // 只包含数字
	const length = 6                 // 验证码长度为6位
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))] // 从0-9中随机选择一个字符
	}
	return string(b) // 返回生成的字符串形式的验证码
}

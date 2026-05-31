package email

import (
	"fmt"
	"log/slog"
	"strings"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dm "github.com/alibabacloud-go/dm-20151123/v2/client"
	"github.com/alibabacloud-go/tea/dara"
	"github.com/alibabacloud-go/tea/tea"
)

type Config struct {
	AccessKeyID     string
	AccessKeySecret string
	Region          string
	Endpoint        string
	AccountName     string
	FromAlias       string
	TagName         string
}

type Sender struct {
	cfg    Config
	client *dm.Client
}

func NewSender(cfg Config) (*Sender, error) {
	s := &Sender{cfg: cfg}
	if !s.Enabled() {
		return s, nil
	}

	clientCfg := &openapi.Config{
		AccessKeyId:     tea.String(cfg.AccessKeyID),
		AccessKeySecret: tea.String(cfg.AccessKeySecret),
		RegionId:        tea.String(cfg.Region),
	}
	if endpoint := strings.TrimSpace(cfg.Endpoint); endpoint != "" {
		clientCfg.Endpoint = tea.String(endpoint)
	}

	client, err := dm.NewClient(clientCfg)
	if err != nil {
		return nil, fmt.Errorf("create directmail client: %w", err)
	}

	s.client = client
	return s, nil
}

func (s *Sender) Enabled() bool {
	return strings.TrimSpace(s.cfg.AccessKeyID) != "" &&
		strings.TrimSpace(s.cfg.AccessKeySecret) != "" &&
		strings.TrimSpace(s.cfg.Region) != "" &&
		strings.TrimSpace(s.cfg.AccountName) != ""
}

func (s *Sender) SendVerificationCode(to, code string) error {
	subject := "Motewallet 注册验证码"
	body := fmt.Sprintf(
		"您好，\n\n您的 Motewallet 注册验证码为：%s\n\n验证码 10 分钟内有效，请勿泄露给他人。\n\n如非本人操作，请忽略此邮件。",
		code,
	)

	if !s.Enabled() {
		slog.Info("email verification code (DirectMail not configured)",
			slog.String("to", to),
			slog.String("code", code),
		)
		return nil
	}

	req := &dm.SingleSendMailRequest{
		AccountName:    tea.String(s.cfg.AccountName),
		AddressType:    tea.Int32(1),
		ReplyToAddress: tea.Bool(false),
		ToAddress:      tea.String(to),
		Subject:        tea.String(subject),
		TextBody:       tea.String(body),
	}

	if alias := strings.TrimSpace(s.cfg.FromAlias); alias != "" {
		req.FromAlias = tea.String(alias)
	}
	if tag := strings.TrimSpace(s.cfg.TagName); tag != "" {
		req.TagName = tea.String(tag)
	}

	if _, err := s.client.SingleSendMailWithOptions(req, &dara.RuntimeOptions{
		ConnectTimeout: tea.Int(10000),
		ReadTimeout:    tea.Int(10000),
	}); err != nil {
		return fmt.Errorf("directmail send: %w", err)
	}

	return nil
}

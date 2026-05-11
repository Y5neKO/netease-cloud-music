package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/chaunsin/netease-cloud-music/api/weapi"
	"github.com/chaunsin/netease-cloud-music/pkg/cookiecloud"
	"github.com/chaunsin/netease-cloud-music/pkg/log"

	qrcode2 "github.com/skip2/go-qrcode"
)

type UserProfile struct {
	UserID    int64  `json:"userId"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatarUrl"`
	VipType   int64  `json:"vipType"`
	Status    int64  `json:"status"`
}

type QrcodeResult struct {
	Key   string `json:"key"`
	Image string `json:"image"`
}

type QrcodeStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type AuthService struct{}

func (s *AuthService) app() *App {
	return ensureApp()
}

func (s *AuthService) IsLoggedIn() (bool, error) {
	a := s.app()
	if a.api == nil {
		return false, fmt.Errorf("api not initialized")
	}
	return !a.api.NeedLogin(context.Background()), nil
}

func (s *AuthService) GetUserProfile() (*UserProfile, error) {
	a := s.app()
	resp, err := a.api.GetUserInfo(context.Background(), &weapi.GetUserInfoReq{})
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}
	if resp.Code != 200 || resp.Profile == nil {
		return nil, fmt.Errorf("未登录或登录已过期")
	}
	return &UserProfile{
		UserID:    resp.Profile.UserId,
		Nickname:  resp.Profile.Nickname,
		AvatarURL: resp.Profile.AvatarUrl,
		VipType:   resp.Profile.VipType,
		Status:    resp.Profile.AccountStatus,
	}, nil
}

func (s *AuthService) CreateQrcode() (*QrcodeResult, error) {
	a := s.app()
	ctx := context.Background()
	key, err := a.api.QrcodeCreateKey(ctx, &weapi.QrcodeCreateKeyReq{Type: 1})
	if err != nil {
		return nil, fmt.Errorf("创建二维码key失败: %w", err)
	}
	if key.UniKey == "" {
		return nil, fmt.Errorf("二维码key为空")
	}
	qr, err := a.api.QrcodeGenerate(ctx, &weapi.QrcodeGenerateReq{
		CodeKey:  key.UniKey,
		Level:    qrcode2.Medium,
		Platform: "web",
	})
	if err != nil {
		return nil, fmt.Errorf("生成二维码失败: %w", err)
	}
	return &QrcodeResult{
		Key:   key.UniKey,
		Image: base64.StdEncoding.EncodeToString(qr.Qrcode),
	}, nil
}

func (s *AuthService) CheckQrcodeStatus(key string) (*QrcodeStatus, error) {
	a := s.app()
	resp, err := a.api.QrcodeCheck(context.Background(), &weapi.QrcodeCheckReq{Type: 1, Key: key})
	if err != nil {
		return nil, fmt.Errorf("检查扫码状态失败: %w", err)
	}
	code := int(resp.Code)
	status := &QrcodeStatus{Code: code}
	switch code {
	case 800:
		status.Message = "二维码已过期"
	case 801:
		status.Message = "等待扫码"
	case 802:
		status.Message = "已扫码，等待确认"
	case 803:
		status.Message = "登录成功"
	default:
		status.Message = fmt.Sprintf("未知状态: %d", code)
	}
	return status, nil
}

func (s *AuthService) SendSmsCode(phone string, countryCode float64) error {
	a := s.app()
	resp, err := a.api.SendSMS(context.Background(), &weapi.SendSMSReq{
		Cellphone: phone,
		CtCode:    int64(countryCode),
	})
	if err != nil {
		return fmt.Errorf("发送验证码失败: %w", err)
	}
	if resp.Code != 200 {
		msg := resp.Msg
		if msg == "" {
			msg = resp.Message
		}
		return fmt.Errorf("发送验证码失败(code=%d): %s", resp.Code, msg)
	}
	return nil
}

func (s *AuthService) LoginByPhone(phone, password, captcha string, countryCode float64) (*UserProfile, error) {
	a := s.app()
	if password == "" && captcha == "" {
		return nil, fmt.Errorf("密码和验证码不能同时为空")
	}
	resp, err := a.api.LoginCellphone(context.Background(), &weapi.LoginCellphoneReq{
		Phone:       phone,
		Countrycode: int64(countryCode),
		Remember:    true,
		Password:    password,
		Captcha:     captcha,
	})
	if err != nil {
		return nil, fmt.Errorf("登录失败: %w", err)
	}
	if resp.Code != 200 {
		msg := resp.Msg
		if msg == "" {
			msg = resp.Message
		}
		return nil, fmt.Errorf("登录失败(code=%d): %s", resp.Code, msg)
	}
	return s.GetUserProfile()
}

func (s *AuthService) LoginByCookie(cookieStr string) (*UserProfile, error) {
	a := s.app()
	if cookieStr == "" {
		return nil, fmt.Errorf("cookie 不能为空")
	}
	cookies := parseCookieString(cookieStr)
	if !hasMusicU(cookies) {
		return nil, fmt.Errorf("cookie 中未找到 MUSIC_U")
	}
	u, _ := url.Parse("https://music.163.com")
	a.client.SetCookies(u, cookies)
	return s.GetUserProfile()
}

func (s *AuthService) LoginByCookieFile(filePath string) (*UserProfile, error) {
	a := s.app()
	if filePath == "" {
		return nil, fmt.Errorf("文件路径不能为空")
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}
	var cookies []*http.Cookie
	if json.Valid(data) {
		cookies, _ = parseCookieJSON(data)
	}
	if len(cookies) == 0 {
		cookies, _ = http.ParseCookie(string(data))
	}
	if len(cookies) == 0 {
		return nil, fmt.Errorf("无法解析 cookie 文件")
	}
	if !hasMusicU(cookies) {
		return nil, fmt.Errorf("cookie 中未找到 MUSIC_U")
	}
	u, _ := url.Parse("https://music.163.com")
	a.client.SetCookies(u, cookies)
	return s.GetUserProfile()
}

func (s *AuthService) LoginByCookieCloud(server, uuid, password string) (*UserProfile, error) {
	a := s.app()
	if server == "" || uuid == "" || password == "" {
		return nil, fmt.Errorf("服务器地址、UUID 和密码不能为空")
	}
	cfg := &cookiecloud.Config{ApiUrl: server, Timeout: 30e9, Retry: 3}
	cc, err := cookiecloud.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("创建 CookieCloud 客户端失败: %w", err)
	}
	resp, err := cc.Get(context.Background(), &cookiecloud.GetReq{Uuid: uuid, Password: password})
	if err != nil {
		return nil, fmt.Errorf("获取 CookieCloud 数据失败: %w", err)
	}
	u, _ := url.Parse("https://music.163.com")
	var found bool
	for domain, domainCookies := range resp.CookieData {
		if !strings.HasSuffix(domain, "music.163.com") {
			continue
		}
		var httpCookies []*http.Cookie
		for _, c := range domainCookies {
			if c.Name == "MUSIC_U" {
				found = true
			}
			httpCookies = append(httpCookies, &http.Cookie{
				Domain: domain, Expires: c.GetExpired(),
				Name: c.Name, Path: c.Path, Secure: c.Secure, Value: c.Value,
			})
		}
		if len(httpCookies) > 0 {
			a.client.SetCookies(u, httpCookies)
		}
	}
	if !found {
		return nil, fmt.Errorf("CookieCloud 中未找到网易云音乐 MUSIC_U")
	}
	return s.GetUserProfile()
}

// SaveConfig 持久化前端配置到 ~/.ncmctl/gui_config.json
func (s *AuthService) SaveConfig(key string, value string) error {
	path := configPath()
	data := make(map[string]string)
	if raw, err := os.ReadFile(path); err == nil {
		json.Unmarshal(raw, &data)
	}
	data[key] = value
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(path, raw, 0o644)
}

// LoadConfig 读取前端配置
func (s *AuthService) LoadConfig(key string) string {
	path := configPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	var m map[string]string
	if json.Unmarshal(data, &m) != nil {
		return ""
	}
	return m[key]
}

func configPath() string {
	home, _ := os.UserHomeDir()
	return home + "/.ncmctl/gui_config.json"
}

func (s *AuthService) Logout() error {
	a := s.app()
	_, err := a.api.Layout(context.Background(), &weapi.LayoutReq{})
	if err != nil {
		log.Warn("调用登出接口失败: %s", err)
	}
	home, _ := os.UserHomeDir()
	_ = os.Remove(home + "/.ncmctl/cookie.json")
	return nil
}

func parseCookieString(raw string) []*http.Cookie {
	cookies, err := http.ParseCookie(raw)
	if err != nil {
		return nil
	}
	return cookies
}

func parseCookieJSON(data []byte) ([]*http.Cookie, error) {
	var items []struct {
		Domain string `json:"domain"`
		Name   string `json:"name"`
		Path   string `json:"path"`
		Secure bool   `json:"secure"`
		Value  string `json:"value"`
	}
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	var cookies []*http.Cookie
	for _, item := range items {
		cookies = append(cookies, &http.Cookie{
			Domain: item.Domain, Name: item.Name,
			Path: item.Path, Secure: item.Secure, Value: item.Value,
		})
	}
	return cookies, nil
}

func hasMusicU(cookies []*http.Cookie) bool {
	for _, c := range cookies {
		if c.Name == "MUSIC_U" {
			return true
		}
	}
	return false
}

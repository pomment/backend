package common

import (
	"github.com/hoisie/mustache"
	"github.com/pomment/pomment/utils/recaptcha"
)

type PommentConfig struct {
	System      PommentConfigSystem    `toml:"system"`
	Admin       PommentConfigAdmin     `toml:"admin"`
	Email       PommentConfigEmail     `toml:"email"`
	ReCAPTCHA   PommentConfigReCAPTCHA `toml:"reCAPTCHA"`
	Push        PommentConfigPush      `toml:"push"`
	Redis       PommentConfigRedis     `toml:"redis"`
	Avatar      PommentConfigAvatar    `toml:"avatar"`
	WebTemplate PommentConfigWebTemplate
}

type PommentConfigSystem struct {
	Host              string `toml:"host"`
	Port              int    `toml:"port"`
	URL               string `toml:"url"`
	DevelopCORSPolicy bool   `toml:"developCORSPolicy"`
}

type PommentConfigAdmin struct {
	Salt string                   `toml:"salt"`
	User []PommentConfigAdminUser `toml:"user"`
}

type PommentConfigAdminUser struct {
	Name     string `toml:"name" json:"name"`
	Email    string `toml:"email" json:"email"`
	Password string `toml:"password" json:"password"`
}

type PommentConfigEmail struct {
	Enabled       bool   `toml:"enabled"`
	Mode          string `toml:"mode"`
	Title         string `toml:"title"`
	Sender        string `toml:"sender"`
	SMTPHost      string `toml:"smtpHost"`
	SMTPPort      int    `toml:"smtpPort"`
	SMTPUsername  string `toml:"smtpUsername"`
	SMTPPassword  string `toml:"smtpPassword"`
	MailgunAPIKey string `toml:"mailgunAPIKey"`
	MailgunDomain string `toml:"mailgunDomain"`
}

type PommentConfigReCAPTCHA struct {
	Enabled      bool    `toml:"enabled"`
	SecretKey    string  `toml:"secretKey"`
	MinimumScore float64 `toml:"minimumScore"`
	Object       recaptcha.ReCAPTCHA
}

type PommentConfigPush struct {
	Enabled        bool   `toml:"enabled"`
	Gateway        string `toml:"gateway"`
	SiteName       string `toml:"siteName"`
	SiteKey        string `toml:"siteKey"`
	GravatarServer string `toml:"gravatarServer"`
}

type PommentConfigRedis struct {
	Enabled  bool   `toml:"enabled"`
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
	Database int    `toml:"database"`
}

type PommentConfigAvatar struct {
	UseSha256 bool `toml:"useSha256"`
}

type PommentConfigWebTemplate struct {
	EmailTitle         mustache.Template
	EmailBody          mustache.Template
	UnsubscribeConfirm mustache.Template
	UnsubscribeSuccess mustache.Template
	UnsubscribeError   mustache.Template
}

type Post struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	EmailHashed  string  `json:"emailHashed"`
	Website      string  `json:"website"`
	Parent       string  `json:"parent"`
	Content      string  `json:"content"`
	Hidden       bool    `json:"hidden"`
	ByAdmin      bool    `json:"byAdmin"`
	ReceiveEmail bool    `json:"receiveEmail"`
	EditKey      string  `json:"editKey"`
	CreatedAt    int64   `json:"createdAt"`
	UpdatedAt    int64   `json:"updatedAt"`
	OrigContent  string  `json:"origContent"`
	Avatar       string  `json:"avatar"`
	Rating       float64 `json:"rating"`
}

type PostSimple struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	EmailHashed string `json:"emailHashed"`
	Website     string `json:"website"`
	Parent      string `json:"parent"`
	Content     string `json:"content"`
	Hidden      bool   `json:"hidden"`
	ByAdmin     bool   `json:"byAdmin"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   int64  `json:"updatedAt"`
	Avatar      string `json:"avatar"`
}

type Thread struct {
	Title        string `json:"title"`
	FirstPostAt  int64  `json:"firstPostAt"`
	LatestPostAt int64  `json:"latestPostAt"`
	Amount       int    `json:"amount"`
	ID           string `json:"id"`
	Locked       bool   `json:"locked"`
	URL          string `json:"url"`
}

type ThreadMapItem struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type SendMessageContext struct {
	Post            Post   `json:"post"`
	PostParent      Post   `json:"postParent"`
	Thread          Thread `json:"thread"`
	TitleTemplate   string `json:"titleTemplate"`
	ContentTemplate string `json:"contentTemplate"`
}

type SendEmailContext struct {
	Recipient string
	Subject   string
	Body      string
}

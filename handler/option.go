package handler

import (
	"flag"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Token struct {
	Key string `ymal:"key"`
}

type Endpoints struct {
	Internet  string `ymal:"internet"`
	LocalAddr string `ymal:"localaddr"`
}

type Telegram struct {
	Token  string `ymal:"token"`
	ChatId string `ymal:"chat_id"`
}

type Dingding struct {
	Token  string `ymal:"token"`
	Secret string `ymal:"secret"`
}

type Etcd struct {
	Endpoints string `ymal:"endpoints"`
	UserName  string `ymal:"user_name"`
	Passwd    string `ymal:"passwd"`
	Version   string `ymal:"version"`
}

func (x *Etcd) Join(prefix ...string) string {
	return x.Version + strings.Join(prefix, "/")
}

type Mysql struct {
	Host string `ymal:"host"`
}

type Mongo struct {
	Host string `ymal:"host"`
	DB   string `ymal:"db"`
}

type Option struct {
	TCP      Endpoints `yaml:"tcp"`
	HTTP     Endpoints `yaml:"http"`
	QUIC     Endpoints `yaml:"quic"`
	Etcd     Etcd      `yaml:"etcd"`
	Dingding Dingding  `yaml:"dingding"`
	Telegram Telegram  `yaml:"telegram"`
	Token    Token     `yaml:"token"`
	Mysql    Mysql     `yaml:"mysql"`
	Mongo    Mongo     `yaml:"mongo"`
}

func (x *Handler) Parse() error {
	var fileName string
	flag.StringVar(&fileName, "option", "balance.yaml", "balance.yaml")
	flag.Parse()
	ymalBytes, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(ymalBytes, &x.Option); err != nil {
		return err
	}
	return nil
}

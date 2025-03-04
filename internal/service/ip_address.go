package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	"go-chat/config"
	"go-chat/internal/pkg/client"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/repository/repo"
)

type IpAddressService struct {
	*repo.Source
	config     *config.Config
	httpClient *client.RequestClient
}

func NewIpAddressService(source *repo.Source, conf *config.Config, httpClient *client.RequestClient) *IpAddressService {
	return &IpAddressService{Source: source, config: conf, httpClient: httpClient}
}

type IpAddressResponse struct {
	Code   string `json:"resultcode"`
	Reason string `json:"reason"`
	Result struct {
		Country  string `json:"Country"`
		Province string `json:"Province"`
		City     string `json:"City"`
		Isp      string `json:"Isp"`
	} `json:"result"`
	ErrorCode int `json:"error_code"`
}

func (i *IpAddressService) FindAddress(ip string) (string, error) {
	if val, err := i.getCache(ip); err == nil {
		return val, nil
	}

	params := &url.Values{}
	params.Add("key", i.config.App.JuheKey)
	params.Add("ip", ip)

	resp, err := i.httpClient.Get("http://apis.juhe.cn/ip/ipNew", params)
	if err != nil {
		return "", err
	}

	data := &IpAddressResponse{}
	if err := json.Unmarshal(resp, data); err != nil {
		return "", err
	}

	if data.Code != "200" {
		return "", errors.New(data.Reason)
	}

	arr := []string{data.Result.Country, data.Result.Province, data.Result.City, data.Result.Isp}
	val := strings.Join(sliceutil.Unique(arr), " ")
	val = strings.TrimSpace(val)

	_ = i.setCache(ip, val)

	return val, nil
}

func (i *IpAddressService) getCache(ip string) (string, error) {
	return i.Redis().HGet(context.TODO(), "rds:hash:ip-address", ip).Result()
}

func (i *IpAddressService) setCache(ip string, value string) error {
	return i.Redis().HSet(context.TODO(), "rds:hash:ip-address", ip, value).Err()
}

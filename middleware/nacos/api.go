package nacos

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/flyerxp/globalStruct/config"
	config2 "github.com/flyerxp/lib/config"
	"github.com/flyerxp/lib/logger"
	"github.com/flyerxp/lib/utils/json"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Client struct {
	BaseOption config.MidNacos
	HttpPool   *sync.Pool
	Context    context.Context
	Token      *AccessToken
}
type AccessToken struct {
	AccessToken string        `json:"accessToken"`
	TokenTtl    time.Duration `json:"tokenTtl"`
	GlobalAdmin bool          `json:"globalAdmin"`
	Username    string        `json:"username"`
	Expiration  int64         `json:"expiration,omitempty"`
}

var redisClient redis.UniversalClient

func NewClient(name string, ctx context.Context) (*Client, error) {
	for _, v := range config2.GetConf().Nacos {
		if v.Name == name {
			return newClient(v, ctx), nil
		}
	}
	logger.AddError(zap.Error(errors.New("nacos conf no find " + name)))
	return nil, errors.New("nacos conf no find " + name)
}
func newClient(o config.MidNacos, ctx context.Context) *Client {
	if redisClient == nil {
		redisClient = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:        o.Redis.Address,
			MasterName:   o.Redis.Master,
			Username:     o.Redis.User,
			Password:     o.Redis.Pwd,
			PoolTimeout:  time.Second,
			MaxIdleConns: 30,
		})
	}
	return &Client{
		o,
		&sync.Pool{
			New: func() any {
				n := newHttpClient(fmt.Sprintf("%x", md5.Sum([]byte(o.Url))))
				return n
			},
		},
		ctx,
		new(AccessToken),
	}
}
func (n *Client) GetKey(url string) string {
	key := n.BaseOption.Url + "@@" + n.BaseOption.User + "@@" + n.BaseOption.Pwd + "@@" + url
	return fmt.Sprintf("N%x", md5.Sum([]byte(key)))
}
func (n *Client) getUrl(url string) string {
	return n.BaseOption.Url + url
}
func (n *Client) getDataFromCache(cacheKey string) (*redis.StringCmd, error) {
	rv := redisClient.Get(n.Context, cacheKey)
	if rv.Err() == redis.Nil || rv.Val() == "" {
		return nil, errors.New("no exists")
	}
	return rv, nil
}

func (n *Client) GetToken(ctx context.Context) (*AccessToken, error) {
	if n.Token != nil && n.Token.Expiration > time.Now().Unix() {
		return n.Token, nil
	}
	key := n.GetKey("/v1/auth/login")
	rv, err := n.getDataFromCache(key)
	// 从缓存中获取
	if err == nil && rv.Err() != redis.Nil {
		token := new(AccessToken)
		bt, e := rv.Bytes()
		jsonErr := json.Decode(bt, token)
		if jsonErr == nil && token.Expiration > time.Now().Unix() {
			n.Token = token
			return token, e
		}
	}
	hc := n.HttpPool.Get().(*httpClient)
	hc.ctx(ctx)
	bToken, bErr := hc.SendRequest("POST", n.getUrl("/v1/auth/login"), "username="+n.BaseOption.User+"&password="+n.BaseOption.Pwd, 0, 0)
	n.HttpPool.Put(hc)
	if bErr != nil {
		logger.AddWarn(zap.Error(bErr))
		return nil, errors.New("nacos request fail")
	} else {
		token := new(AccessToken)
		err = json.Decode(bToken, token)
		if err == nil {
			token.TokenTtl -= 10
			token.Expiration = time.Now().Unix() + int64(token.TokenTtl)
			jsonStr, jsonErr := json.Encode(token)
			if jsonErr == nil && token.TokenTtl > 10 {
				redisClient.Set(n.Context, key, string(jsonStr), time.Second*token.TokenTtl)
			} else {
				return nil, jsonErr
			}
			n.Token = token
			return token, err
		} else {
			return nil, err
		}
	}
}

func (n *Client) GetConfig(ctx context.Context, did string, gp string, ns string) ([]byte, error) {
	key := n.GetKey("/nacos/v1/cs/configs" + "@@" + did + "@@" + gp + "@@" + ns)
	token, err := n.GetToken(ctx)
	rv, rErr := n.getDataFromCache(key)
	if rErr == nil && rv.String() == "" {
		return rv.Bytes()
	}
	//接口报错，则从cache取
	if err != nil {
		return []byte{}, err
	} else {
		s := logger.StartTime("nacos request")
		hc := n.HttpPool.Get().(*httpClient)
		bYaml, bErr := hc.SendRequest("GET", n.getUrl("/v1/cs/configs?accessToken="+token.AccessToken+"&tenant="+ns+"&dataId="+did+"&group="+gp), "username="+n.BaseOption.User+"&password="+n.BaseOption.Pwd, 0, 0)
		n.HttpPool.Put(hc)
		s.Stop()
		if bErr == nil {
			sYaml := string(bYaml)
			if rErr == nil && rv.String() != sYaml {
				redisClient.Set(ctx, key, sYaml, time.Hour*48)
			}
			return bYaml, nil
		} else {
			if rErr != nil && rv.Val() != "" {
				return rv.Bytes()
			}
		}
		return []byte{}, bErr
	}
}

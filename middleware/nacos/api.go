package nacos

import (
	"context"
	"errors"
	config2 "github.com/flyerxp/lib/config"
	"github.com/flyerxp/lib/logger"
	"github.com/flyerxp/lib/utils/json"
	"github.com/flyerxp/lib/utils/stringL"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"sync"
	"time"
)

type Client struct {
	BaseOption config2.MidNacos
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

func GetEngine(name string, ctx context.Context) (*Client, error) {
	for _, v := range config2.GetConf().Nacos {
		if v.Name == name {
			return newClient(v, ctx), nil
		}
	}
	logger.AddError(zap.Error(errors.New("nacos conf no find " + name)))
	return nil, errors.New("nacos conf no find " + name)
}
func newClient(o config2.MidNacos, ctx context.Context) *Client {
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
	c := &Client{
		o,
		&sync.Pool{
			New: func() any {
				n := newHttpClient(stringL.GetMd5(o.Url))
				return n
			},
		},
		ctx,
		new(AccessToken),
	}
	return c
}
func (n *Client) GetKey(url string) string {
	key := n.BaseOption.Url + "@@" + url
	return stringL.GetMd5(key)
}
func (n *Client) getUrl(url string) string {
	return n.BaseOption.Url + url
}
func (n *Client) getDataFromCache(cacheKey string) (*redis.StringCmd, error) {
	rv := redisClient.Get(n.Context, cacheKey)
	return rv, nil
}
func (n *Client) DelToken(ctx context.Context) {
	key := n.GetKey("/v1/auth/login")
	redisClient.Del(ctx, key)
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
	s := logger.StartTime("nacos-get-token")
	hc := n.HttpPool.Get().(*httpClient)
	hc.ctx(ctx)
	bToken, bErr := hc.SendRequest("POST", n.getUrl("/v1/auth/login"), "username="+n.BaseOption.User+"&password="+n.BaseOption.Pwd, 0, 0)
	n.HttpPool.Put(hc)
	s.Stop()
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

func (n *Client) DeleteCache(ctx context.Context, did string, gp string, ns string) string {
	key := n.GetKey("/nacos/v1/cs/configs" + "@@" + did + "@@" + gp + "@@" + ns)
	redisClient.Del(ctx, key)
	return key
}
func (n *Client) GetConfig(ctx context.Context, did string, gp string, ns string) ([]byte, error) {
	start := time.Now()
	key := n.GetKey("/nacos/v1/cs/configs" + "@@" + did + "@@" + gp + "@@" + ns)
	rv, rErr := n.getDataFromCache(key)
	if rErr == nil && rv.Err() != redis.Nil {
		logger.AddNacosTime(int(time.Since(start).Milliseconds()))
		return rv.Bytes()
	}
	token, err := n.GetToken(ctx)
	//接口报错，返回空
	if err != nil {
		logger.AddNacosTime(int(time.Since(start).Milliseconds()))
		logger.AddError(zap.Error(err))
		return []byte{}, err
	} else {
		s := logger.StartTime("nacos-get-config")
		hc := n.HttpPool.Get().(*httpClient)
		bYaml, bErr := hc.SendRequest("GET", n.getUrl("/v1/cs/configs?accessToken="+token.AccessToken+"&tenant="+ns+"&dataId="+did+"&group="+gp), "", 0, 0)
		n.HttpPool.Put(hc)
		s.Stop()
		if bErr == nil {
			sYaml := string(bYaml)
			if rv.String() != sYaml {
				redisClient.Set(ctx, key, sYaml, time.Second*86400*2)
			}
			logger.AddNacosTime(int(time.Since(start).Milliseconds()))
			return bYaml, nil
		} else {
			if rErr != nil && rv.Val() != "" {
				logger.AddNacosTime(int(time.Since(start).Milliseconds()))
				return rv.Bytes()
			}
		}
		logger.AddNacosTime(int(time.Since(start).Milliseconds()))
		return []byte{}, bErr
	}
}

/*
type ServiceRequest struct {
	Ip          string            `json:"ip"`
	Port        int               `json:"port"`
	ClusterName string            `json:"cluster_name"`
	ServiceName string            `json:"service_name"`
	GroupName   string            `json:"group_name"`
	NamespaceId string            `json:"namespace_id"`
	Healthy     bool              `json:"healthy"`
	Weight      int               `json:"weight"`
	Enabled     bool              `json:"enabled"`
	Metadata    map[string]string `json:"metadata"`
	Ephemeral   bool              `json:"ephemeral"`
}

func (n *Client) PutService(ctx context.Context, request ServiceRequest) ([]byte, error) {
	start := time.Now()
	if request.Ip == "" || request.NamespaceId == "" || request.Port == 0 || request.ClusterName == "" || request.GroupName == "" || request.ServiceName == "" {
		return []byte{}, errors.New("参数错误")
	}
	//n.DelToken(ctx)
	token, errToken := n.GetToken(ctx)
	if errToken != nil {
		logger.AddNacosTime(int(time.Since(start).Milliseconds()))
		logger.AddError(zap.Error(errToken))
		return []byte{}, errToken
	}
	fmt.Println(token.AccessToken)
	//接口报错，返回空
	s := logger.StartTime("nacos-put-service")
	hc := n.HttpPool.Get().(*httpClient)
	v, _ := query.Values(request)
	postData := v.Encode() + "&accessToken=" + token.AccessToken + "&tenant=" + request.NamespaceId
	fmt.Println(postData)
	fmt.Println(n.getUrl("/v1/ns/instance"))
	r, bErr := hc.SendRequest("POST", n.getUrl("/v1/ns/instance?"+postData), postData, time.Second*10, 3)
	fmt.Println(bErr)
	n.HttpPool.Put(hc)
	s.Stop()
	logger.AddNacosTime(int(time.Since(start).Milliseconds()))
	return r, bErr
}
*/

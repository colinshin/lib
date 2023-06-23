package pulsarL

import (
	"context"
	"errors"
	"github.com/apache/pulsar-client-go/pulsar"
	config2 "github.com/flyerxp/lib/config"
	"github.com/flyerxp/lib/logger"
	"github.com/flyerxp/lib/middleware/nacos"
	yaml2 "github.com/flyerxp/lib/utils/yaml"
	cmap "github.com/orcaman/concurrent-map/v2"
	"go.uber.org/zap"
	_ "log"
	"strings"
)

// Pulsar 容器
type PulsarContainer struct {
	PulsarContainer cmap.ConcurrentMap[string, *PulsarClient]
	PulsarConf      cmap.ConcurrentMap[string, config2.MidPulsarConf]
}

// Pulsar 客户端
type PulsarClient struct {
	CurrPulsar pulsar.Client
}

var pulsarEngine *PulsarContainer

func init() {
	go initEngine(context.Background())
}
func initEngine(ctx context.Context) {
	pulsarEngine = new(PulsarContainer)
	var confList []config2.MidPulsarConf
	pulsarEngine.PulsarConf = cmap.New[config2.MidPulsarConf]()
	pulsarEngine.PulsarContainer = cmap.New[*PulsarClient]()
	conf := config2.GetConf()
	confList = conf.Pulsar
	//本地文件中获取
	for _, v := range confList {
		if v.Name != "" {
			pulsarEngine.PulsarConf.Set(v.Name, v)
		}
	}

	if conf.PulsarNacos.Name != "" {
		var yaml []byte
		pulsarList := new(config2.PulsarConf)

		ns, e := nacos.GetEngine(conf.PulsarNacos.Name, ctx)
		if e == nil {
			yaml, e = ns.GetConfig(ctx, conf.PulsarNacos.Did, conf.PulsarNacos.Group, conf.PulsarNacos.Ns)
			if e == nil {
				e = yaml2.DecodeByBytes(yaml, pulsarList)
				if e == nil {
					for _, v := range pulsarList.List {
						pulsarEngine.PulsarConf.Set(v.Name, v)
					}
				} else {
					logger.AddError(zap.Error(errors.New("yaml conver error")))
				}
			} else {
				logger.AddError(zap.Error(e))
			}
		} else {
			logger.AddError(zap.Error(e))
		}
	}
}
func GetEngine(name string, ctx context.Context) (*PulsarClient, error) {
	if pulsarEngine == nil {
		initEngine(ctx)
	}
	e, ok := pulsarEngine.PulsarContainer.Get(name)
	if ok {
		return e, nil
	}
	o, okC := pulsarEngine.PulsarConf.Get(name)
	if okC {
		objPulsar := newClient(o)
		pulsarEngine.PulsarContainer.Set(name, objPulsar)
		return objPulsar, nil
	}
	logger.AddError(zap.Error(errors.New("no find Pulsar config " + name)))
	return nil, errors.New("no find Pulsar config " + name)
}

// https://github.com/golang-migrate/migrate/blob/master/database/Pulsar/README.md

func newClient(o config2.MidPulsarConf) *PulsarClient {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: "pulsar://" + strings.Join(o.Address, ","),
	})
	if err != nil {
		logger.AddError(zap.Error(err))
	}

	return &PulsarClient{client}
}

func (m *PulsarClient) GetPulsar() pulsar.Client {
	return m.CurrPulsar
}
func (p *PulsarContainer) Reset() {
	for k := range pulsarEngine.PulsarContainer.Items() {
		if v, ok := pulsarEngine.PulsarContainer.Pop(k); ok {
			v.CurrPulsar.Close()
		}
	}
	pulsarEngine = nil
}
func Reset() {
	Flush()
	pulsarEngine = nil
	_, _ = topicDistributionF()
}

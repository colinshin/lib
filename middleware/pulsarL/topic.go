package pulsarL

import (
	"context"
	"errors"
	"fmt"
	"github.com/flyerxp/lib/config"
	"github.com/flyerxp/lib/logger"
	"github.com/flyerxp/lib/middleware/nacos"
	yaml "github.com/flyerxp/lib/utils/yaml"
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"
)

type TopicS struct {
	Code    int    `yaml:"code" json:"code"`
	CodeStr string `yaml:"code_str" json:"code_str"`
	Delay   int    `yaml:"delay" json:"delay"`
	Cluster string `yaml:"cluster" json:"cluster"`
}
type TopicConfS struct {
	TopicDistribution map[int]string
	Topic             map[string]TopicS
	IsInitEd          bool
	isLoading         bool //是否正在异步重载
}
type yamlTopic struct {
	TopicDistribution map[int]string `yaml:"topic_distribution" json:"topic_distribution"`
	Topic             []TopicS       `yaml:"topic" json:"topic"`
}

func getCluster(code int, t map[int]string) string {
	i := int(code / 1000000)
	if c, ok := t[i]; ok {
		return c
	}
	if c, ok := topicConf.TopicDistribution[i]; ok {
		return c
	} else {
		eStr := fmt.Sprintf("topic:%d no find cluster %d/1000000 in cluster %d", code, code, i)
		logger.AddWarn(zap.Error(errors.New(eStr)))
	}
	return ""
}

var topicConf TopicConfS

func init() {
	go initTopic()
}

func getTopic(code int) (*TopicS, bool) {
	if !topicConf.IsInitEd {
		initTopic()
	}
	codeStr := strconv.Itoa(code)
	if t, ok := topicConf.Topic[codeStr]; ok {
		return &t, true
	} else {
		//如果没找到，自动重新载入配置
		eStr := fmt.Sprintf("topic:%d no find, 5 second reset load", code)
		logger.AddWarn(zap.Error(errors.New(eStr)))
		topicConf.IsInitEd = false
		if !topicConf.isLoading {
			topicConf.isLoading = true
			time.AfterFunc(time.Second*5, func() {
				initTopic()
				topicConf.isLoading = false
			})
		}
		return nil, false
	}
}

func initTopic() {
	if !topicConf.IsInitEd {
		conf, err := topicDistributionF()
		if err == nil {
			topicConf = conf
		}
		topicConf.IsInitEd = true
	}
}

func topicDistributionF() (TopicConfS, error) {
	topicConfTmp := TopicConfS{}
	topicConfTmp.TopicDistribution = make(map[int]string)
	topicConfTmp.Topic = make(map[string]TopicS)
	conf := config.GetConf().TopicNacos
	for _, v := range conf {
		n, e := nacos.GetEngine(v.Name, context.Background())
		if e != nil {
			logger.AddError(zap.Error(e))
		} else {
			b, be := n.GetConfig(context.Background(), v.Did, v.Group, v.Ns)
			if be == nil {
				tmp := new(yamlTopic)
				e = yaml.DecodeByBytes(b, tmp)
				if e != nil {
					logger.AddError(zap.Error(e))
				} else {
					getTopicConfig(&topicConfTmp, tmp)
				}
			} else {
				logger.AddError(zap.Error(e))
				return topicConfTmp, e
			}
		}
	}
	if getConfFile() == nil {
		tmpConf := new(yamlTopic)
		err := yaml.DecodeByFile(config.GetConfFile("pulsar.yml"), tmpConf)
		if err != nil {
			logger.AddError(zap.String("pusal topic error", "pulsar.yml read error"), zap.Error(err))
			return topicConfTmp, nil
		}
		getTopicConfig(&topicConfTmp, tmpConf)
	}
	return topicConfTmp, nil
}
func getConfFile() error {
	_, errf := os.Stat(config.GetConfFile("pulsar.yml"))
	if errf != nil && os.IsNotExist(errf) {
		return errf
	} else if errf != nil {
		logger.AddError(zap.String("topic read err", "read pulsar.yml err"), zap.Error(errf))
		return errf
	}
	return nil
}
func getTopicConfig(topicConfTmp *TopicConfS, tmpConf *yamlTopic) {
	for ck, cv := range tmpConf.TopicDistribution {
		topicConfTmp.TopicDistribution[ck] = cv
	}
	for _, cvt := range tmpConf.Topic {
		if cvt.CodeStr != "" {
			topicConfTmp.Topic[cvt.CodeStr] = cvt
		} else if clusterT := getCluster(cvt.Code, topicConfTmp.TopicDistribution); clusterT != "" {
			cvt.Cluster = clusterT
			topicConfTmp.Topic[strconv.Itoa(cvt.Code)] = cvt
		}
	}
}

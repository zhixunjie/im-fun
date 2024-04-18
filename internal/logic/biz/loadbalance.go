package biz

import (
	"github.com/zhixunjie/im-fun/internal/logic/conf"
	"github.com/zhixunjie/im-fun/pkg/logging"
)

type LoadBalance struct {
	conf    *conf.Config
	regions map[string]string // province -> region
}

func NewLoadBalance(conf *conf.Config) *LoadBalance {
	o := &LoadBalance{
		conf: conf,
	}
	o.initRegions()

	return o
}

// 建立省份与Region的映射关系
func (o *LoadBalance) initRegions() {
	o.regions = make(map[string]string, 0)
	for region, provinces := range o.conf.Regions {
		for _, province := range provinces {
			o.regions[province] = region
		}
	}
	logging.Infof("initRegions success,regions=%v", o.regions)
}

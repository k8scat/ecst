package cmd

import "github.com/hsowan-me/vss/server"

const (
	aliyunRegionStatusAvailable string = "available"
	aliyunRegionStatusSoldOut   string = "soldOut"
)

func isRegionAvailable(config *server.Config, regionID string) (available bool, err error) {
	regions, err := server.ListRegions(config)
	if err != nil {
		return
	}
	for _, region := range regions {
		if region.RegionId == regionID && region.Status == aliyunRegionStatusAvailable {
			available = true
			return
		}
	}
	return
}

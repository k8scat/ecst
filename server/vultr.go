package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/wanhuasong/vss/utils"
	"log"
	"net/url"
	"strconv"
	"time"
)

const (
	apiEndpoint      = "https://api.vultr.com"
	apiServerCreate  = "/v1/server/create"
	apiServerList    = "/v1/server/list"
	apiServerDestroy = "/v1/server/destroy"
	apiRegionList    = "/v1/regions/list"
	apiOSList        = "/v1/os/list"

	serverStatusOK = "ok"
)

type VultrClient struct {
	ApiKey string
}

type VultrRegion struct {
	DCID           string
	Name           string
	Country        string
	Continent      string
	Stat           string
	DDosProtection bool
	BlockStorage   bool
	RegionCode     string
	VPSPlanIDs     []int64
}

type VultrOS struct {
	OSID    int64
	Name    string
	Arch    string
	Family  string
	Windows bool
}

func NewVultrClient(apiKey string) *VultrClient {
	return &VultrClient{
		ApiKey: apiKey,
	}
}

func (c *VultrClient) authHeader() map[string]string {
	return map[string]string{
		"API-Key": c.ApiKey,
	}
}

func (c *VultrClient) CreateInstance(dcid, vpsPlanID, osID, firewallGroupID string) (instance *Instance, err error) {
	request := utils.NewRequest(apiEndpoint)
	data := map[string]string{
		"DCID":            dcid,
		"VPSPLANID":       vpsPlanID,
		"OSID":            osID,
		"FIREWALLGROUPID": firewallGroupID,
	}
	var payload *bytes.Buffer
	var contentType string
	payload, contentType, err = utils.ParseFormPayload(data)
	if err != nil {
		return
	}
	headers := c.authHeader()
	headers["Content-Type"] = contentType
	var s string
	s, err = request.Post(apiServerCreate, payload, headers)
	if err != nil {
		return
	}
	subID := gjson.Get(s, "SUBID").String()
	for {
		instance, err = c.GetInstance(subID)
		if err != nil {
			return
		}
		if instance.Status != serverStatusOK {
			log.Println("waiting instance start...")
			time.Sleep(time.Second * time.Duration(20))
			continue
		}
		return
	}
}

func (c *VultrClient) GetInstance(subID string) (instance *Instance, err error) {
	request := utils.NewRequest(apiEndpoint)
	params := &url.Values{
		"SUBID": {subID},
	}
	var s string
	s, err = request.Get(apiServerList, params, c.authHeader())
	if err != nil {
		return
	}
	if !gjson.Valid(s) {
		err = fmt.Errorf("get instance failed: %s", s)
		return
	}
	instance = &Instance{
		ID:       subID,
		PublicIP: gjson.Get(s, "main_ip").String(),
		Password: gjson.Get(s, "default_password").String(),
		Status:   gjson.Get(s, "server_state").String(),
	}
	return
}

func (c *VultrClient) ListInstances() (instances []*Instance, err error) {
	request := utils.NewRequest(apiEndpoint)
	var s string
	s, err = request.Get(apiServerList, nil, c.authHeader())
	if err != nil {
		return
	}
	subIDs := make(map[string]interface{})
	err = json.Unmarshal([]byte(s), &subIDs)
	if err != nil {
		return
	}
	for subID := range subIDs {
		r := gjson.Get(s, subID)
		instances = append(instances, &Instance{
			ID:       subID,
			PublicIP: r.Get("main_ip").String(),
			Password: r.Get("default_password").String(),
			Status:   r.Get("server_state").String(),
		})
	}
	return
}

func (c *VultrClient) DestroyInstance(subID string) error {
	request := utils.NewRequest(apiEndpoint)
	data := map[string]string{
		"SUBID": subID,
	}
	payload, contentType, err := utils.ParseFormPayload(data)
	if err != nil {
		return err
	}
	headers := c.authHeader()
	headers["Content-Type"] = contentType
	_, err = request.Post(apiServerDestroy, payload, headers)
	return err
}

func (c *VultrClient) ListDCIDs() (regions []*VultrRegion, err error) {
	request := utils.NewRequest(apiEndpoint)
	var s string
	s, err = request.Get(apiRegionList, nil, nil)
	if err != nil {
		return
	}
	i := 1
	for {
		r := gjson.Get(s, strconv.Itoa(i))
		if !r.Exists() {
			break
		}
		regions = append(regions, &VultrRegion{
			DCID:           r.Get("DCID").String(),
			Name:           r.Get("name").String(),
			Country:        r.Get("country").String(),
			Continent:      r.Get("continent").String(),
			Stat:           r.Get("state").String(),
			DDosProtection: r.Get("ddos_protection").Bool(),
			BlockStorage:   r.Get("block_storage").Bool(),
			RegionCode:     r.Get("regioncode").String(),
		})
		i++
	}
	return
}

func (c *VultrClient) ListOSIDs() (osList []*VultrOS, err error) {
	request := utils.NewRequest(apiEndpoint)
	var s string
	s, err = request.Get(apiOSList, nil, nil)
	if err != nil {
		return
	}
	osIDs := make(map[string]interface{})
	err = json.Unmarshal([]byte(s), &osIDs)
	if err != nil {
		return
	}
	for osID := range osIDs {
		r := gjson.Get(s, osID)
		osList = append(osList, &VultrOS{
			OSID:    r.Get("OSID").Int(),
			Name:    r.Get("name").String(),
			Arch:    r.Get("arch").String(),
			Family:  r.Get("family").String(),
			Windows: r.Get("windows").Bool(),
		})
	}
	return
}

func ListVPSPlanIDs() {

}

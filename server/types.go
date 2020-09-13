package server

const (
	ProviderAliyun string = "aliyun"
	ProviderVultr  string = "vultr"
)

type Config struct {
	AliyunAccessKeyID  string `json:"aliyun_access_key_id"`
	AliyunAccessSecret string `json:"aliyun_access_secret"`
	VultrAPIKey        string `json:"vultr_api_key"`
}

type Option struct {
	Provider   string `json:"provider"`
	InstanceID string `json:"instance_id"`
	ListType   string `json:"list_type"`

	// Aliyun options
	// create
	RegionID        string `json:"region_id"`
	InstanceType    string `json:"instance_type"`
	ImageID         string `json:"image_id"`
	SecurityGroupID string `json:"security_group_id"`
	VSwitchID       string `json:"v_switch_id"`

	// Vultr options
	// create
	DCID            string `json:"dcid"`
	VPSPlanID       string `json:"vps_plan_id"`
	OSID            string `json:"osid"`
	FirewallGroupID string `json:"firewall_group_id"`

	ScriptFile string `json:"script_file"`
}

type Instance struct {
	ID       string
	Password string
	PublicIP string
	Status   string
}

package server

import (
	"errors"
	"fmt"
	aliyunErrors "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/wanhuasong/vss/utils"
	"log"
	"strings"
	"time"
)

const (
	defaultRegionID             string = "cn-hangzhou"
	aliyunInstanceStatusRunning string = "Running"
)

type AliyunClient struct {
	AccessKeyID  string
	AccessSecret string
}

func NewAliyunClient(accessKeyID, accessSecret string) *AliyunClient {
	return &AliyunClient{
		AccessKeyID:  accessKeyID,
		AccessSecret: accessSecret,
	}
}

func (c *AliyunClient) CreateInstance(regionID, imageID, instanceType, securityGroupID, vSwitchID string) (instance *Instance, err error) {
	client, err := ecs.NewClientWithAccessKey(regionID, c.AccessKeyID, c.AccessSecret)
	if err != nil {
		return
	}
	request := ecs.CreateRunInstancesRequest()
	request.Scheme = "https"
	request.InstanceChargeType = "PostPaid"
	password := utils.GeneratePassword()
	request.Password = password
	request.Amount = requests.NewInteger(1)
	request.InternetMaxBandwidthIn = requests.NewInteger(100)
	request.InternetMaxBandwidthOut = requests.NewInteger(100)
	request.InstanceType = instanceType
	request.SecurityGroupId = securityGroupID
	request.ImageId = imageID
	request.VSwitchId = vSwitchID
	var res *ecs.RunInstancesResponse
	res, err = client.RunInstances(request)
	if err != nil {
		return
	}
	if !res.IsSuccess() {
		err = fmt.Errorf("create instance failed: %d, %s", res.GetHttpStatus(), res.GetHttpContentString())
		return
	}
	instanceID := res.InstanceIdSets.InstanceIdSet[0]
	log.Println("instance starting...")
	for {
		instance, err = c.GetInstance(instanceID)
		if err != nil {
			return
		}
		if instance.Status != aliyunInstanceStatusRunning {
			time.Sleep(time.Second * time.Duration(2))
			continue
		}
		instance.Password = password
		return
	}
}

func (c *AliyunClient) GetInstance(instanceID string) (instance *Instance, err error) {
	client, err := ecs.NewClientWithAccessKey(defaultRegionID, c.AccessKeyID, c.AccessSecret)
	if err != nil {
		return
	}
	request := ecs.CreateDescribeInstanceAttributeRequest()
	request.Scheme = "https"
	request.InstanceId = instanceID
	var res *ecs.DescribeInstanceAttributeResponse
	res, err = client.DescribeInstanceAttribute(request)
	if err != nil {
		if serverErr, isServerErr := err.(*aliyunErrors.ServerError); isServerErr {
			if serverErr.ErrorCode() == "InvalidInstanceId.NotFound" {
				err = errors.New(serverErr.Message())
				return
			}
		}
		return
	}
	if !res.IsSuccess() {
		err = fmt.Errorf("get instance failed, %d, %s", res.GetHttpStatus(), res.GetHttpContentString())
		return
	}
	instance = &Instance{
		ID:       instanceID,
		Status:   res.Status,
		PublicIP: strings.Join(res.PublicIpAddress.IpAddress, ", "),
	}
	return
}

func (c *AliyunClient) DestroyInstance(instanceID string) error {
	if _, err := c.GetInstance(instanceID); err != nil {
		return err
	}
	client, err := ecs.NewClientWithAccessKey(defaultRegionID, c.AccessKeyID, c.AccessSecret)
	if err != nil {
		return err
	}
	request := ecs.CreateDeleteInstanceRequest()
	request.Scheme = "https"
	request.InstanceId = instanceID
	request.Force = requests.NewBoolean(true)
	var res *ecs.DeleteInstanceResponse
	res, err = client.DeleteInstance(request)
	if err != nil {
		return err
	}
	if !res.IsSuccess() {
		err = fmt.Errorf("destroy instance failed: %d, %s", res.GetHttpStatus(), res.GetHttpContentString())
		return err
	}
	return nil
}

func (c *AliyunClient) ListInstances(regionID string) (instances []*Instance, err error) {
	client, err := ecs.NewClientWithAccessKey(regionID, c.AccessKeyID, c.AccessSecret)
	if err != nil {
		return
	}
	request := ecs.CreateDescribeInstancesRequest()
	request.Scheme = "https"
	pageNumber := 1
	request.PageNumber = requests.NewInteger(pageNumber)
	request.PageSize = requests.NewInteger(100)
	var res *ecs.DescribeInstancesResponse
	for {
		res, err = client.DescribeInstances(request)
		if err != nil {
			return
		}
		if !res.IsSuccess() {
			err = fmt.Errorf("list instances failed: %d, %s", res.GetHttpStatus(), res.GetHttpContentString())
			return
		}
		for _, instance := range res.Instances.Instance {
			instances = append(instances, &Instance{
				ID:       instance.InstanceId,
				PublicIP: strings.Join(instance.PublicIpAddress.IpAddress, ", "),
				Status:   instance.Status,
			})
		}
		if len(instances) == res.TotalCount {
			return
		} else {
			pageNumber += 1
			request.PageNumber = requests.NewInteger(pageNumber)
		}
	}
}

func (c *AliyunClient) ListRegions() (regions []ecs.Region, err error) {
	client, err := ecs.NewClientWithAccessKey(defaultRegionID, c.AccessKeyID, c.AccessSecret)
	if err != nil {
		return
	}
	request := ecs.CreateDescribeRegionsRequest()
	request.Scheme = "https"
	var res *ecs.DescribeRegionsResponse
	res, err = client.DescribeRegions(request)
	if err != nil {
		return
	}
	if !res.IsSuccess() {
		err = fmt.Errorf("list regions failed: %d, %s", res.GetHttpStatus(), res.GetHttpContentString())
		return
	}
	regions = res.Regions.Region
	return
}

package server

import (
	"errors"
	"fmt"
	aliyunErrors "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hsowan-me/vss/utils"
	"log"
	"time"
)

// https://api.aliyun.com/?accounttraceid=093de352c1224f028a5d0dbff1385cdcjvjz#/?product=Ecs&version=2014-05-26&api=CreateInstance&params={}&tab=DEMO&lang=GO
func RunAliyunInstance(config *Config, option *Option) (instance *AliyunInstance, err error) {
	if err = ValidateConfig(config, ProviderAliyun); err != nil {
		return
	}
	client, err := ecs.NewClientWithAccessKey(option.RegionID, config.AliyunAccessKeyID, config.AliyunAccessSecret)
	if err != nil {
		return
	}
	request := ecs.CreateRunInstancesRequest()
	request.Scheme = "https"
	request.InstanceChargeType = "PostPaid"
	password := generatePassword()
	request.Password = password
	request.Amount = requests.NewInteger(1)
	request.InternetMaxBandwidthIn = requests.NewInteger(100)
	request.InternetMaxBandwidthOut = requests.NewInteger(100)
	request.InstanceType = option.InstanceType
	request.SecurityGroupId = option.SecurityGroupID
	request.ImageId = option.ImageID
	request.VSwitchId = option.VSwitchID
	var res *ecs.RunInstancesResponse
	res, err = client.RunInstances(request)
	if err != nil {
		return
	}
	if !res.IsSuccess() {
		err = fmt.Errorf("create instance failed: %d, %s", res.GetHttpStatus(), res.GetHttpContentString())
		return
	}
	option.InstanceID = res.InstanceIdSets.InstanceIdSet[0]
	for {
		instance, err = GetAliyunInstance(config, option)
		if err != nil {
			return
		}
		if instance.Status != aliyunInstanceStatusRunning {
			log.Println("instance starting...")
			time.Sleep(time.Duration(2) * time.Second)
			continue
		}
		instance.Password = password
		return
	}
}

// https://api.aliyun.com/?accounttraceid=093de352c1224f028a5d0dbff1385cdcjvjz#/?product=Ecs&version=2014-05-26&api=DescribeInstanceAttribute&params={}&tab=DEMO&lang=GO
func GetAliyunInstance(config *Config, option *Option) (instance *AliyunInstance, err error) {
	if err = ValidateConfig(config, ProviderAliyun); err != nil {
		return
	}
	client, err := ecs.NewClientWithAccessKey(defaultRegionID, config.AliyunAccessKeyID, config.AliyunAccessSecret)
	if err != nil {
		return
	}
	request := ecs.CreateDescribeInstanceAttributeRequest()
	request.Scheme = "https"
	request.InstanceId = option.InstanceID
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
	instance = &AliyunInstance{
		InstanceID: option.InstanceID,
		Status:     res.Status,
		Ips:        res.PublicIpAddress.IpAddress,
	}
	return
}

// https://api.aliyun.com/#/?product=Ecs&version=2014-05-26&api=DeleteInstance
func DestroyAliyunInstance(config *Config, option *Option) error {
	if err := ValidateConfig(config, ProviderAliyun); err != nil {
		return err
	}
	if _, err := GetAliyunInstance(config, option); err != nil {
		return err
	}
	client, err := ecs.NewClientWithAccessKey(defaultRegionID, config.AliyunAccessKeyID, config.AliyunAccessSecret)
	if err != nil {
		return err
	}
	request := ecs.CreateDeleteInstanceRequest()
	request.Scheme = "https"
	request.InstanceId = option.InstanceID
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

func ListAliyunInstances(config *Config, option *Option) (instances []ecs.Instance, err error) {
	if err = ValidateConfig(config, ProviderAliyun); err != nil {
		return
	}
	client, err := ecs.NewClientWithAccessKey(option.RegionID, config.AliyunAccessKeyID, config.AliyunAccessSecret)
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
		instances = append(instances, res.Instances.Instance...)
		if len(instances) == res.TotalCount {
			return
		} else {
			pageNumber += 1
			request.PageNumber = requests.NewInteger(pageNumber)
		}
	}
}

func generatePassword() string {
	source := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789()`~!@#$%^&*-_+=|{}[]:;'<>,.?"
	password := fmt.Sprintf("%s%s", utils.RandomString(source, 25), "Vss^1")
	return password
}

func ListRegions(config *Config) (regions []ecs.Region, err error) {
	if err = ValidateConfig(config, ProviderAliyun); err != nil {
		return
	}
	client, err := ecs.NewClientWithAccessKey(defaultRegionID, config.AliyunAccessKeyID, config.AliyunAccessSecret)
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

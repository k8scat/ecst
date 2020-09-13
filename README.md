# VSS

Cloud server tool.

[![Build Status](https://travis-ci.com/wanhuasong/vss.svg?branch=master)](https://travis-ci.com/wanhuasong/vss)

## Support

- [x] [Aliyun](https://www.aliyun.com/)
- [x] [Vultr](https://www.vultr.com/)

## Start

```
# config
vss config --aliyun_access_key_id ${aliyun_access_key_id} --aliyun_access_secret ${aliyun_access_secret}

# create
vss create --provider aliyun --region_id ${region_id} \
--image_id ${image_id} --instance_type ${instance_type} \
--security_group_id ${security_group_id} --v_switch_id ${v_switch_id}

# list
vss list --provider aliyun --region_id ${region_id}

# destroy
vss destroy --provider aliyun --instance_id ${instance_id}

```

## Build

```
./build.sh [ linux|darwin|windows ]

```

## References

- [Aliyun VPC Switches](https://vpcnext.console.aliyun.com/vpc/cn-hongkong/switches)

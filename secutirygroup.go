package aliyun

import (
	"encoding/json"
	"net/http"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/fimreal/goutils/ezap"
	"github.com/gin-gonic/gin"
)

func Allow(c *gin.Context) {
	// 检查传入 json 参数是否符合
	var sgrule SGRule
	if err := c.ShouldBind(&sgrule); err != nil {
		ezap.Error(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ezap.Debugf("请求添加安全组规则, 传入 IP 地址: %s, 安全组 ID: %s, 备注: %s, 策略: %s", sgrule.IP, sgrule.SGID, sgrule.Remark, sgrule.Policy)
	if !sgrule.verify() {
		ezap.Errorw("aliyun securitygroup cap", "desc", "传入参数不符合要求")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "传入参数不符合要求"})
		return
	}

	err := sgrule.authorize()
	if err != nil {
		ezap.Error("添加安全组规则出错: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ezap.Debug("成功添加安全组规则")
	c.JSON(http.StatusOK, gin.H{"result": "成功添加安全组规则"})
}

// authorize() 添加安全组规则
func (s *SGRule) authorize() error {
	c := NewClient()

	r := ecs.CreateAuthorizeSecurityGroupRequest()
	r.Scheme = "https"
	r.IpProtocol = "tcp"
	r.Priority = "1"
	r.Policy = "accept"
	if s.Policy == "drop" || s.Policy == "deny" {
		r.Policy = "drop"
	}
	r.NicType = "internet"
	r.PortRange = "1/65535"

	r.SecurityGroupId = s.SGID
	r.SourceCidrIp = s.IP
	r.Description = s.Remark

	res, err := c.AuthorizeSecurityGroup(r)
	if err != nil {
		return err
	}

	ezap.Debugf("Response: %#v, IP: %s  was successfully added to the security group rule", res, r.SourceCidrIp)
	return nil
}

// verify() 检查传入参数是否正确
// 检查 IP 是否为内网 IP，如果是则返回错误
// 检查安全组 ID 对应安全组是否存在，不存在返回错误
// 否则通过
func (s *SGRule) verify() bool {

	if IsLanIPv4(s.IP) {
		ezap.Errorf("传入 ip[%s] 为内网 ip", s.IP)
		return false
	}

	c := NewClient()
	r := ecs.CreateDescribeSecurityGroupsRequest()
	r.Scheme = "https"
	r.RegionId = ak.REGION_ID
	r.SecurityGroupId = s.SGID
	ezap.Debugf("查询阿里云安全组信息, 请求参数: %+v", *r)

	resJSON, err := c.DescribeSecurityGroups(r)
	if err != nil {
		ezap.Error(err.Error())
		return false
	}
	// 解析返回数据
	res := &ecs.DescribeSecurityGroupsResponse{}
	err = json.Unmarshal(resJSON.GetHttpContentBytes(), res)
	if err != nil {
		ezap.Error(err.Error())
		return false
	}

	ezap.Debug("开始检查安全组是否存在，", s.SGID)
	if len(res.SecurityGroups.SecurityGroup) == 0 {
		ezap.Debugf("传入安全组 id: %s, 找到对应安全组 id: %+v", s.SGID, res.SecurityGroups.SecurityGroup)
		ezap.Error(s.SGID, " is not exists ")
		return false
	}

	return true
}

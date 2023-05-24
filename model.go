package aliyun

type SGRule struct {
	IP     string `json:"ip" form:"ip" validate:"required,ip"`
	SGID   string `json:"sgid" form:"sgid" validate:"required"`
	Remark string `json:"remark" form:"remark" validate:"required"`
	Policy string `json:"policy" form:"policy"` // 默认 accept, 可选 drop
}

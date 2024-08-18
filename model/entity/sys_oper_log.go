package entity

import "time"

type SysOperLog struct {
	OperId        int64     `xorm:"pk autoincr" json:"operId"`
	Title         string    `json:"title" form:"title"`
	BusinessType  int       `json:"businessType" form:"businessType"` // 0=其它,1=新增,2=修改,3=删除,4=授权,5=导出,6=导入,7=强退,8=生成代码,9=清空数据
	Method        string    `json:"method"`
	RequestMethod string    `json:"requestMethod"`
	OperatorType  int       `json:"operatorType"` // 0=其它,1=后台用户,2=手机端用户
	OperName      string    `json:"operName" form:"operName"`
	DeptName      string    `json:"deptName"`
	OperUrl       string    `json:"operUrl"`
	OperIp        string    `json:"operIp" form:"operIp"`
	OperLocation  string    `json:"operLocation"`
	OperParam     string    `json:"operParam"`
	JsonResult    string    `json:"jsonResult"`
	Status        int       `json:"status" form:"status"` // 0=正常,1=异常
	ErrorMsg      string    `json:"errorMsg"`
	CostTime      int64     `json:"costTime"`
	OperTime      time.Time `json:"operTime"`

	// 前端提交的额外参数
	Params map[string]any `xorm:"-" json:"params,omitempty"`
}

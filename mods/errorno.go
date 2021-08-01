package mods

const (
	SUCCESS = 200
	ERROR   = 500

	ERROR_MISSING_DATA = 1001
	ERROR_VERIFY_FAIL  = 1002
	ERROR_VALUE_ERROR  = 1003
	ERROR_PERMISSION   = 1004
	ERROR_USECOUNT     = 1005
	ERROR_UNKNOWN      = 1006

	ERROR_DATABASE = 2001
)

var errorMsg = map[int]string{
	SUCCESS:            "OK",
	ERROR:              "FAIL",
	ERROR_MISSING_DATA: "请求缺失数据",
	ERROR_VERIFY_FAIL:  "验证失败",
	ERROR_VALUE_ERROR:  "请求方法异常",
	ERROR_PERMISSION:   "权限不足",
	ERROR_USECOUNT:     "当日请求次数已达上限",
	ERROR_UNKNOWN:      "未知错误",
	ERROR_DATABASE:     "数据库连接出错",
}

func GetErrMsg(errno int) string {
	return errorMsg[errno]
}

package generate

import (
	"fmt"
	"os"
	"reflect"
)

// 设置查询条件
func setSearch(data []map[string]mapValue) string {
	str := ""
	for _, mapVal := range data {
		for k, v := range mapVal {
			if v.searchInfo != "" {
				str = fmt.Sprintf("%s\t%s := c.PostForm(\"%s\")\n", str, k, k)
				str = fmt.Sprintf("%s\tif %s != \"\" {\n", str, k)
				switch v.searchInfo {
				case "=":
				case ">":
				case "<":
				case "!=":
				case "like":
					v := fmt.Sprintf("fmt.Sprintf(\"%s%s\", %s, \"%s\")", "%s", "%s", k, "%")
					str = fmt.Sprintf("%s\t\tconditions = model.SetCondition(conditions,\"%s\",%s,\"like\")\n", str, k, v)
				}
				str = fmt.Sprintf("%s\t}\n", str)
			}
		}
	}
	return str
}

// 创建分类列表查询方法
func setPaginate(file *os.File, kind reflect.Type) string{
	str := `
func Paginate(c *gin.Context) {
	var conditions []model.Condition
`
	var rlt []map[string]mapValue
	rlt = getField(rlt, kind)
	data := setSearch(rlt)
	str = fmt.Sprintf("%s%s\n", str, data)

	str = fmt.Sprintf("%s\tvar %s []%s\n", str, kind.Name(), kind)

	data = "	page, _ := strconv.Atoi(c.DefaultQuery(\"page\", \"1\"))"
	str = fmt.Sprintf("%s%s\n", str, data)
	data = "	pageSize, _ := strconv.Atoi(c.DefaultQuery(\"page_size\", \"10\"))"
	str = fmt.Sprintf("%s%s\n", str, data)
	str = fmt.Sprintf("%s\tlists := model.Paginate(&%s, model.PageInfo{Page: int64(page), PageSize: int64(pageSize)}, conditions)\n", str, kind.Name())

	str = fmt.Sprintf("%s\tutils.SuccessData(c, lists)\n}", str)
	return str
	//file.Write([]byte(str))
}

// 创建详细数据方法
func setOne(file *os.File, kind reflect.Type) string{
	str := `
func Info(c *gin.Context) {
`
	// 查询信息
	data := getInfo(kind)
	str = fmt.Sprintf("%s\n\t%s\n", str, data)
	name := kind.Name()
	str = fmt.Sprintf("%s	utils.SuccessData(c, %sInfo)\n}", str, name)
	return str
	//file.Write([]byte(str))
}

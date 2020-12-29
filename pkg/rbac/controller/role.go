package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wuyan94zl/api/pkg/orm"
	"github.com/wuyan94zl/api/pkg/rbac/model"
	"github.com/wuyan94zl/api/pkg/utils"
	"strconv"
	"time"
)

func RoleCreate(c *gin.Context) {
	// 验证参数
	data := make(map[string][]string)

	data["name"] = []string{"required"}
	data["description"] = []string{"required"}

	validate := utils.Validator(c.Request, data)
	if validate != nil {
		utils.SuccessErr(c, 403, validate)
		return
	}
	var Role model.Role
	Role.Name = c.PostForm("name")
	Role.Description = c.PostForm("description")

	orm.GetInstance().Create(&Role)
	utils.SuccessData(c, Role) // 返回创建成功的信息
}
func RoleUpdate(c *gin.Context) {
	// 验证参数
	data := make(map[string][]string)

	data["name"] = []string{"required"}
	data["description"] = []string{"required"}

	validate := utils.Validator(c.Request, data)
	if validate != nil {
		utils.SuccessErr(c, 403, validate)
		return
	}
	id, _ := strconv.Atoi(c.Query("id"))
	var Role model.Role
	orm.GetInstance().First(&Role,id)
	if Role.Id == 0 {
		utils.SuccessErr(c, -1000, "数据不存在")
		return
	}

	Role.Name = c.PostForm("name")
	Role.Description = c.PostForm("description")
	orm.GetInstance().Save(Role)
	utils.SuccessData(c, Role) // 返回创建成功的信息
}
func RoleDelete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	var Role model.Role

	orm.GetInstance().First(&Role,id)
	if Role.Id == 0 {
		utils.SuccessErr(c, -1000, "数据不存在")
		return
	}
	orm.GetInstance().Delete(Role)
	orm.GetInstance().Delete(Role.Menus)
	orm.GetInstance().Delete(Role.Permissions)
	utils.SuccessData(c, "删除成功")
}
func RoleInfo(c *gin.Context) {
	//id, _ := strconv.Atoi(c.Query("id"))
	//var Role model.Role
	//orm.GetInstance().First(&Role, id, "Menus", "Permissions")
	//utils.SuccessData(c, Role)
	var Role model.Role
	where := make(map[string]interface{})

	where["name"] = "测试"
	where["id"] = 3
	where["created_at"] = orm.Where{Way: "between",Value: []time.Time{time.Now(),time.Now()}}

	or := make(map[string]interface{})
	or["id"] = orm.Where{Way: "IN",Value: []int64{2,4,5}}
	or["name"] = orm.Where{Way: "like",Value: "测试%"}

	orm.GetInstance().Where(where).Or(or).Get(&Role)
	utils.SuccessData(c, Role)
}

func RolePermissionMenu(c *gin.Context) {
	permissionId := c.DefaultPostForm("permission_id", "")
	roleId := c.Query("id")
	role := model.Role{}
	orm.GetInstance().First(&role,roleId)
	role.DelPermissionMenu()
	if permissionId != "" {
		role.SetPermissionMenu(permissionId)
	}
	utils.SuccessData(c, "设置成功")
}

func RolePaginate(c *gin.Context) {
	var Role []model.Role
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "3"))
	paginate := orm.SetPageList(&Role, int64(page), int64(pageSize))
	fmt.Println(paginate)
	orm.GetInstance().Order("id desc").Paginate(paginate)
	utils.SuccessData(c, paginate)
}

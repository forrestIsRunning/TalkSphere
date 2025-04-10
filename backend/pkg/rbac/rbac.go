package rbac

import (
	"TalkSphere/models"
	"TalkSphere/pkg/encrypt"
	"TalkSphere/pkg/mysql"
	"TalkSphere/pkg/snowflake"
	"TalkSphere/setting"

	"fmt"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var Enforcer *casbin.Enforcer

func InitCasbin() {
	// 使用 GORM 适配器
	adapter, err := gormadapter.NewAdapterByDB(mysql.DB)
	if err != nil {
		zap.L().Fatal("failed to initialize casbin adapter", zap.Error(err))
	}

	// Load both model and policy files when creating enforcer
	enforcer, err := casbin.NewEnforcer("conf/rbac_model.conf", "conf/rbac_policy.csv")
	if err != nil {
		zap.L().Fatal("failed to create casbin enforcer", zap.Error(err))
	}

	// Set the adapter and save policy to DB
	enforcer.SetAdapter(adapter)
	if err := enforcer.SavePolicy(); err != nil {
		zap.L().Fatal("failed to save policy to DB", zap.Error(err))
	}

	// 启用自动保存
	enforcer.EnableAutoSave(true)

	Enforcer = enforcer
}

// CheckPermission 检查权限
func CheckPermission(sub string, obj string, act string) bool {
	ok, err := Enforcer.Enforce(sub, obj, act)
	if err != nil {
		zap.L().Error("casbin enforce error", zap.Error(err))
		return false
	}
	return ok
}

// AddRole 为用户添加角色
func AddRole(user string, role string) bool {
	ok, err := Enforcer.AddGroupingPolicy(user, role)
	if err != nil {
		zap.L().Error("add role error", zap.Error(err))
		return false
	}
	return ok
}

// RemoveRole 删除用户角色
func RemoveRole(user string, role string) bool {
	ok, err := Enforcer.RemoveGroupingPolicy(user, role)
	if err != nil {
		zap.L().Error("remove role error", zap.Error(err))
		return false
	}
	return ok
}

// GetUserRole 根据用户ID获取用户角色
func GetUserRole(userID string) (string, error) {
	var rule struct {
		V1 string `gorm:"column:v1"`
	}

	err := mysql.DB.Table("casbin_rule").
		Select("v1").
		Where("ptype = 'g' AND v0 = ?", userID).
		First(&rule).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "guest", nil // 未注册或未登录的用户视为游客
		}
		return "", err
	}

	return rule.V1, nil
}

// InitSuperAdmin 初始化超级管理员
func InitSuperAdmin() {
	// 检查是否已存在超级管理员
	var user models.User
	result := mysql.DB.Where("username = ?", "super_admin").First(&user)
	if result.RowsAffected > 0 {
		// 已存在则跳过
		return
	}

	// 创建超级管理员用户
	userID := snowflake.GenID()
	user = models.User{
		ID:           userID,
		Username:     "super_admin",
		Email:        setting.Conf.SuperAdmin.Email,
		PasswordHash: encrypt.EncryptPassword(setting.Conf.SuperAdmin.Password), // 使用环境变量或配置文件中的密码
		Bio:          "System Super Administrator",
		Status:       1,
	}

	if err := mysql.DB.Create(&user).Error; err != nil {
		zap.L().Fatal("failed to create super admin", zap.Error(err))
		return
	}

	// 设置超级管理员角色
	if ok := AddRole(fmt.Sprintf("%d", userID), "super_admin"); !ok {
		zap.L().Fatal("failed to set super admin role")
	}
}

package rbac

import (
	"github.com/TalkSphere/backend/models"
	"github.com/TalkSphere/backend/pkg/encrypt"
	"github.com/TalkSphere/backend/pkg/mysql"
	"github.com/TalkSphere/backend/pkg/snowflake"
	"github.com/TalkSphere/backend/setting"

	"strconv"

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

	modelPath := "conf/rbac_model.conf"
	policyPath := "conf/rbac_policy.csv"

	// 创建enforcer，使用模型配置文件和数据库适配器
	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		zap.L().Fatal("failed to create casbin enforcer", zap.Error(err))
	}

	// 检查数据库中是否有策略
	var count int64
	if err := mysql.DB.Table("casbin_rule").Count(&count).Error; err != nil {
		zap.L().Fatal("failed to check policy existence", zap.Error(err))
	}

	zap.L().Info("Current policy count in database", zap.Int64("count", count))

	// 如果数据库中没有策略，则从 CSV 文件加载
	if count == 0 {
		zap.L().Info("No policy found in database, loading from CSV file")

		// 创建一个临时的 enforcer 来加载 CSV 文件
		tmpEnforcer, err := casbin.NewEnforcer(modelPath, policyPath)
		if err != nil {
			zap.L().Fatal("failed to create temporary enforcer", zap.Error(err))
		}

		// 获取策略规则
		rules, err := tmpEnforcer.GetPolicy()
		if err != nil {
			zap.L().Fatal("failed to get policies from CSV", zap.Error(err))
		}
		zap.L().Info("Loaded policies from CSV", zap.Int("rules_count", len(rules)))

		groupingRules, err := tmpEnforcer.GetGroupingPolicy()
		if err != nil {
			zap.L().Fatal("failed to get grouping policies from CSV", zap.Error(err))
		}
		zap.L().Info("Loaded grouping policies from CSV", zap.Int("grouping_rules_count", len(groupingRules)))

		// 将规则添加到主 enforcer
		if len(rules) > 0 {
			_, err = enforcer.AddPolicies(rules)
			if err != nil {
				zap.L().Fatal("failed to add policies", zap.Error(err))
			}
			zap.L().Info("Added policies to enforcer", zap.Int("count", len(rules)))
		}

		// 添加角色继承规则
		if len(groupingRules) > 0 {
			_, err = enforcer.AddGroupingPolicies(groupingRules)
			if err != nil {
				zap.L().Fatal("failed to add grouping policies", zap.Error(err))
			}
			zap.L().Info("Added grouping policies to enforcer", zap.Int("count", len(groupingRules)))
		}

		// 保存到数据库
		if err := enforcer.SavePolicy(); err != nil {
			zap.L().Fatal("failed to save policy to DB", zap.Error(err))
		}
		zap.L().Info("Successfully saved policies to database")

		// 验证策略是否已保存
		var newCount int64
		if err := mysql.DB.Table("casbin_rule").Count(&newCount).Error; err != nil {
			zap.L().Error("failed to check new policy count", zap.Error(err))
		}
		zap.L().Info("New policy count in database", zap.Int64("count", newCount))
	} else {
		zap.L().Info("Loading policy from database")
		// 从数据库加载策略
		if err := enforcer.LoadPolicy(); err != nil {
			zap.L().Fatal("failed to load policy from database", zap.Error(err))
		}
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
	zap.L().Info("开始删除用户角色",
		zap.String("user_id", user),
		zap.String("role", role))

	ok, err := Enforcer.RemoveGroupingPolicy(user, role)
	if err != nil {
		zap.L().Error("remove role error", zap.Error(err))
		return false
	}

	zap.L().Info("删除用户角色完成",
		zap.String("user_id", user),
		zap.String("role", role),
		zap.Bool("success", ok))
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

// RemoveAllRoles 删除用户的所有角色
func RemoveAllRoles(userID string) bool {
	zap.L().Info("开始删除用户所有角色",
		zap.String("user_id", userID))

	// 先获取用户当前的所有角色
	roles, err := Enforcer.GetRolesForUser(userID)
	if err != nil {
		zap.L().Error("获取用户角色失败",
			zap.String("user_id", userID),
			zap.Error(err))
		return false
	}

	zap.L().Info("当前用户的所有角色",
		zap.String("user_id", userID),
		zap.Strings("roles", roles))

	ok, err := Enforcer.RemoveFilteredGroupingPolicy(0, userID)
	if err != nil {
		zap.L().Error("删除用户所有角色失败",
			zap.String("user_id", userID),
			zap.Error(err))
		return false
	}

	// 保存策略到数据库
	if err := Enforcer.SavePolicy(); err != nil {
		zap.L().Error("保存策略到数据库失败",
			zap.String("user_id", userID),
			zap.Error(err))
		return false
	}

	zap.L().Info("删除用户所有角色完成",
		zap.String("user_id", userID),
		zap.Bool("success", ok))
	return ok
}

// InitSuperAdmin 初始化超级管理员
func InitSuperAdmin() {
	// 检查是否已存在超级管理员
	var user models.User
	result := mysql.DB.Where("username = ?", "super_admin").First(&user)

	if result.RowsAffected == 0 {
		// 不存在则创建超级管理员用户
		userID := snowflake.GenID()
		user = models.User{
			ID:           userID,
			Username:     "super_admin",
			Email:        setting.Conf.SuperAdmin.Email,
			PasswordHash: encrypt.EncryptPassword(setting.Conf.SuperAdmin.Password),
			Bio:          "System Super Administrator",
			Status:       1,
		}

		if err := mysql.DB.Create(&user).Error; err != nil {
			zap.L().Fatal("failed to create super admin", zap.Error(err))
			return
		}
	}

	// 无论是新建还是已存在，都确保角色正确
	userIDStr := strconv.FormatInt(user.ID, 10)

	// 检查当前角色
	currentRole, err := GetUserRole(userIDStr)
	if err != nil {
		zap.L().Error("failed to get user role", zap.Error(err))
	}

	// 如果不是超级管理员，则设置角色
	if currentRole != "super_admin" {
		// 先清除可能存在的其他角色
		RemoveAllRoles(userIDStr)

		// 设置超级管理员角色
		if ok := AddRole(userIDStr, "super_admin"); !ok {
			zap.L().Fatal("failed to set super admin role")
		}

		zap.L().Info("super admin role set successfully",
			zap.String("user_id", userIDStr),
			zap.String("username", user.Username))
	}
}

// AddPermissionForUser 为用户添加特定权限
func AddPermissionForUser(userID string, obj string, act string) bool {
	ok, err := Enforcer.AddPolicy(userID, obj, act)
	if err != nil {
		zap.L().Error("add permission error", zap.Error(err))
		return false
	}
	return ok
}

// RemovePermissionForUser 移除用户的特定权限
func RemovePermissionForUser(userID string, obj string, act string) bool {
	ok, err := Enforcer.RemovePolicy(userID, obj, act)
	if err != nil {
		zap.L().Error("remove permission error", zap.Error(err))
		return false
	}
	return ok
}

// GetUserPermissions 获取用户所有权限
func GetUserPermissions(userID string) ([][]string, error) {
	return Enforcer.GetFilteredPolicy(0, userID)
}

// HasPermission 检查用户是否有特定权限
func HasPermission(userID string, obj string, act string) bool {
	// 首先检查用户是否是超级管理员
	role, err := GetUserRole(userID)
	if err == nil && role == "super_admin" {
		return true
	}

	// 检查具体权限
	return CheckPermission(userID, obj, act)
}

// BatchUpdatePermissions 批量更新用户权限
func BatchUpdatePermissions(userID string, permissions [][]string) error {
	// 首先移除该用户的所有直接权限
	_, err := Enforcer.RemoveFilteredPolicy(0, userID)
	if err != nil {
		return err
	}

	// 添加新的权限
	for _, permission := range permissions {
		if len(permission) == 2 {
			_, err = Enforcer.AddPolicy(userID, permission[0], permission[1])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// GetUserAllPermissions 获取用户的所有权限（包括角色继承的权限）
func GetUserAllPermissions(userID string) map[string][]string {
	permissions := make(map[string][]string)

	// 获取用户的角色
	role, err := GetUserRole(userID)
	if err != nil {
		zap.L().Error("get user role error", zap.Error(err))
		return permissions
	}

	// 获取直接分配给用户的权限
	directPermissions, err := GetUserPermissions(userID)
	if err != nil {
		zap.L().Error("GetUserPermissions error", zap.Error(err))
		return permissions
	}
	for _, p := range directPermissions {
		if len(p) >= 3 {
			obj, act := p[1], p[2]
			if _, exists := permissions[obj]; !exists {
				permissions[obj] = make([]string, 0)
			}
			permissions[obj] = append(permissions[obj], act)
		}
	}

	// 获取角色的权限
	rolePermissions, err := Enforcer.GetFilteredPolicy(0, role)
	if err != nil {
		zap.L().Error("GetFilteredPolicy error", zap.Error(err))
		return permissions
	}
	for _, p := range rolePermissions {
		if len(p) >= 3 {
			obj, act := p[1], p[2]
			if _, exists := permissions[obj]; !exists {
				permissions[obj] = make([]string, 0)
			}
			permissions[obj] = append(permissions[obj], act)
		}
	}

	return permissions
}

// UpdateUserPermissions 更新用户的所有权限
func UpdateUserPermissions(userID string, permissions []map[string]interface{}) error {
	// 首先移除用户的所有直接权限
	_, err := Enforcer.RemoveFilteredPolicy(0, userID)
	if err != nil {
		return err
	}

	// 添加新的权限
	for _, perm := range permissions {
		obj, objOk := perm["path"].(string)
		actions, actionsOk := perm["actions"].([]interface{})

		if !objOk || !actionsOk {
			continue
		}

		for _, act := range actions {
			if actStr, ok := act.(string); ok {
				_, err = Enforcer.AddPolicy(userID, obj, actStr)
				if err != nil {
					zap.L().Error("add policy error",
						zap.String("user_id", userID),
						zap.String("obj", obj),
						zap.String("act", actStr),
						zap.Error(err))
				}
			}
		}
	}

	return nil
}

// CheckUserPermission 检查用户是否有特定权限（考虑继承关系）
func CheckUserPermission(userID string, obj string, act string) bool {
	// 检查是否是超级管理员
	role, err := GetUserRole(userID)
	if err == nil && role == "super_admin" {
		return true
	}

	// 检查直接权限
	if ok := CheckPermission(userID, obj, act); ok {
		return true
	}

	// 检查角色权限
	if role != "" {
		return CheckPermission(role, obj, act)
	}

	return false
}

// GetRolePermissions 获取角色的所有权限
func GetRolePermissions(role string) ([][]string, error) {
	return Enforcer.GetFilteredPolicy(0, role)
}

// RemoveUserDirectPermissions 移除用户的直接权限（保留角色权限）
func RemoveUserDirectPermissions(userID string) bool {
	_, err := Enforcer.RemoveFilteredPolicy(0, userID)
	if err != nil {
		zap.L().Error("remove user direct permissions error", zap.Error(err))
		return false
	}
	return true
}

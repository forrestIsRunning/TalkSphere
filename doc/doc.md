# TalkSphere 项目文档

## 项目初始化

### 1. 克隆脚手架
```bash
git clone git@github.com:Forrest-Tao/My_favorable_scaffold.git backend
```

### 2. 初始化仓库
```bash
git add .
git commit -m "Initial commit"
git push -u origin main
```

## 环境搭建

### MySQL 数据库

#### 1. 使用 Docker 启动
```bash
# 在项目根目录下执行
docker-compose -p mysql -f backend/deploy/mysql-docker-compose.yaml up -d
```

#### 2. 数据库配置
| 配置项 | 值 |
|--------|-----|
| 数据库名 | TalkSphere |
| 用户名 | forrest |
| 密码 | 571400yst |
| 端口 | 3306 |

#### 3. 表结构
<details>
<summary>boards 表（板块）</summary>

```sql
CREATE TABLE `boards` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(191) NOT NULL,
  `description` varchar(191) DEFAULT NULL,
  `creator_id` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_boards_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
```
</details>

> 📝 完整表结构请查看 [backend/models](../backend/models) 目录

## 服务启动

### 1. 词云图服务
详细说明见 [backend/scripts/README.md](../backend/scripts/worldCloud/README.md)

### 2. genData
详细说明见 [backend/genData/README.md](../backend/scripts/genData/README.md)

### 3. Golang 后端
详细说明见 [backend/README.md](../backend/README.md)

### 4. 前端服务
详细说明见 [frontend/README.md](../frontend/README.md)

## 目录结构
```
TalkSphere/
├── backend/         # 后端服务
│   ├── deploy/     # 部署相关配置
│   ├── models/     # 数据库模型
│   └── scripts/    # 词云图服务
├── frontend/       # 前端服务
└── doc/           # 项目文档
``` 
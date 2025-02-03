### 引入自己的脚手架
```bash
git clone git@github.com:Forrest-Tao/My_favorable_scaffold.git backend
```

### 初始化项目
```bash
git add .
git commit -m "Initial commit"
git push -u origin main
```

### mysql环境搭建
```bash
docker-compose -p mysql -f /Users/Zhuanz/go/src/forrest/TalkSphere/backend/deploy/mysql-docker-compose.yaml up -d
```

### 创建tables
```bash
-- 设置数据库字符集
CREATE DATABASE IF NOT EXISTS talksphere
CHARACTER SET utf8mb4
COLLATE utf8mb4_unicode_ci;

USE talksphere;

-- 用户表
CREATE TABLE users (
                       id BIGINT PRIMARY KEY AUTO_INCREMENT,
                       username VARCHAR(50) UNIQUE NOT NULL,
                       email VARCHAR(100) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       avatar_url VARCHAR(255),
                       bio TEXT,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                       status TINYINT DEFAULT 1 COMMENT '1: active, 0: inactive',
                       last_login_at TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Casbin 规则表
CREATE TABLE casbin_rule (
                             id BIGINT PRIMARY KEY AUTO_INCREMENT,
                             ptype VARCHAR(255),
                             v0 VARCHAR(255),
                             v1 VARCHAR(255),
                             v2 VARCHAR(255),
                             v3 VARCHAR(255),
                             v4 VARCHAR(255),
                             v5 VARCHAR(255),
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 板块表
CREATE TABLE boards (
                        id BIGINT PRIMARY KEY AUTO_INCREMENT,
                        name VARCHAR(100) NOT NULL,
                        description TEXT,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                        status TINYINT DEFAULT 1 COMMENT '1: active, 0: inactive',
                        sort_order INT DEFAULT 0,
                        creator_id BIGINT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 帖子表
CREATE TABLE posts (
                       id BIGINT PRIMARY KEY AUTO_INCREMENT,
                       title VARCHAR(255) NOT NULL,
                       content TEXT NOT NULL,
                       board_id BIGINT,
                       author_id BIGINT,
                       view_count INT DEFAULT 0 COMMENT '观看次数',
                       like_count INT DEFAULT 0 COMMENT '点赞数',
                       favorite_count INT DEFAULT 0 COMMENT '收藏数',
                       comment_count INT DEFAULT 0 COMMENT '评论数',
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                       status TINYINT DEFAULT 1 COMMENT '1: published, 0: draft, -1: deleted',
                       FULLTEXT KEY idx_post_search (title, content)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 帖子图片表
CREATE TABLE post_images (
                             id BIGINT PRIMARY KEY AUTO_INCREMENT,
                             post_id BIGINT,
                             user_id BIGINT NOT NULL,
                             image_url VARCHAR(255) NOT NULL,
                             status TINYINT DEFAULT 1,
                             sort_order INT DEFAULT 0,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             INDEX idx_post_id (post_id),
                             INDEX idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 标签表
CREATE TABLE tags (
                      id BIGINT PRIMARY KEY AUTO_INCREMENT,
                      name VARCHAR(50) UNIQUE NOT NULL,
                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 帖子标签关联表
CREATE TABLE post_tags (
                           post_id BIGINT,
                           tag_id BIGINT,
                           PRIMARY KEY (post_id, tag_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 评论表
CREATE TABLE comments (
                          id BIGINT PRIMARY KEY AUTO_INCREMENT,
                          post_id BIGINT NOT NULL COMMENT '帖子ID',
                          user_id BIGINT NOT NULL COMMENT '评论作者ID',
                          content TEXT NOT NULL COMMENT '评论内容',
                          parent_id BIGINT DEFAULT NULL COMMENT '父评论ID，顶级评论为NULL',
                          root_id BIGINT DEFAULT NULL COMMENT '根评论ID，顶级评论为NULL',
                          like_count INT DEFAULT 0 COMMENT '点赞数',
                          reply_count INT DEFAULT 0 COMMENT '回复数',
                          score INT DEFAULT 0 COMMENT '评论得分(用于排序)',
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                          status TINYINT DEFAULT 1 COMMENT '状态：1正常 0隐藏 -1删除',
                          FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
                          FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                          FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE,
                          FOREIGN KEY (root_id) REFERENCES comments(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建索引
CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_parent_id ON comments(parent_id);
CREATE INDEX idx_comments_root_id ON comments(root_id);
CREATE INDEX idx_comments_created_at ON comments(created_at);
CREATE INDEX idx_comments_score ON comments(score);

-- 点赞表
CREATE TABLE likes (
                       id BIGINT PRIMARY KEY AUTO_INCREMENT,
                       user_id BIGINT,
                       target_id BIGINT COMMENT '点赞目标ID',
                       target_type TINYINT COMMENT '1: post, 2: comment',
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       UNIQUE KEY unique_like (user_id, target_id, target_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 收藏表
CREATE TABLE favorites (
                           id BIGINT PRIMARY KEY AUTO_INCREMENT,
                           user_id BIGINT,
                           post_id BIGINT,
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           UNIQUE KEY unique_favorite (user_id, post_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 用户活跃度统计表
CREATE TABLE user_activities (
                                 id BIGINT PRIMARY KEY AUTO_INCREMENT,
                                 user_id BIGINT,
                                 date DATE,
                                 post_count INT DEFAULT 0,
                                 comment_count INT DEFAULT 0,
                                 like_count INT DEFAULT 0,
                                 view_count INT DEFAULT 0,
                                 UNIQUE KEY unique_user_daily (user_id, date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建必要的索引
CREATE INDEX idx_posts_board ON posts(board_id);
CREATE INDEX idx_posts_author ON posts(author_id);
CREATE INDEX idx_comments_post ON comments(post_id);
CREATE INDEX idx_comments_user ON comments(user_id);
CREATE INDEX idx_comments_parent ON comments(parent_id);
CREATE INDEX idx_likes_user ON likes(user_id);
CREATE INDEX idx_likes_target ON likes(target_id, target_type);
CREATE INDEX idx_favorites_user ON favorites(user_id);
CREATE INDEX idx_favorites_post ON favorites(post_id);
```
### 与数据库交互
```bash
➜  ~ docker ps -a | grep mysql
ab9cdd12c62a   mysql:8.0              "docker-entrypoint.s…"   43 hours ago   Up 43 hours   0.0.0.0:3306->3306/tcp, 33060/tcp

#进入docker容器
➜  ~ docker exec -it ab9cdd12c62a sh

#show databases
mysql> show databases;
+--------------------+
| Database           |
+--------------------+
| information_schema |
| performance_schema |
| talksphere         |
+--------------------+
3 rows in set (0.01 sec)

#change database
mysql> use talksphere;
Database changed

mysql> select * from users \G
*************************** 1. row ***************************
        id: 1
created_at: 2025-01-03 19:10:27.900
updated_at: 2025-01-03 19:10:27.900
deleted_at: NULL
   user_id: 198927022297088
  username: yansaitao
  password: 3537313430307973742341bdf8249049c40e8e0ce7305e78e5471b
     email: yansaitao@qq.com
    avatar:
       bio:
*************************** 2. row ***************************
        id: 2
created_at: 2025-01-04 15:52:22.188
updated_at: 2025-01-04 15:52:22.188
deleted_at: NULL
   user_id: 481263710375936
  username: alice
  password: 313233bdf8249049c40e8e0ce7305e78e5471b
     email: email@gamil.com
    avatar:
       bio:
2 rows in set (0.00 sec)
```
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
);

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
);

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
);

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
);

-- 帖子图片表
CREATE TABLE post_images (
                             id BIGINT PRIMARY KEY AUTO_INCREMENT,
                             post_id BIGINT,
                             image_url VARCHAR(255) NOT NULL,
                             sort_order INT DEFAULT 0,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 标签表
CREATE TABLE tags (
                      id BIGINT PRIMARY KEY AUTO_INCREMENT,
                      name VARCHAR(50) UNIQUE NOT NULL,
                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 帖子标签关联表
CREATE TABLE post_tags (
                           post_id BIGINT,
                           tag_id BIGINT,
                           PRIMARY KEY (post_id, tag_id)
);

-- 评论表 (树形结构)
CREATE TABLE comments (
                          id BIGINT PRIMARY KEY AUTO_INCREMENT,
                          post_id BIGINT,
                          user_id BIGINT,
                          content TEXT NOT NULL,
                          parent_id BIGINT COMMENT '父评论ID，顶级评论为NULL',
                          like_count INT DEFAULT 0,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                          status TINYINT DEFAULT 1 COMMENT '1: visible, 0: hidden'
);

-- 点赞表 (可同时用于帖子和评论的点赞)
CREATE TABLE likes (
                       id BIGINT PRIMARY KEY AUTO_INCREMENT,
                       user_id BIGINT,
                       target_id BIGINT COMMENT '点赞目标ID',
                       target_type TINYINT COMMENT '1: post, 2: comment',
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       UNIQUE KEY unique_like (user_id, target_id, target_type)
);

-- 收藏表
CREATE TABLE favorites (
                           id BIGINT PRIMARY KEY AUTO_INCREMENT,
                           user_id BIGINT,
                           post_id BIGINT,
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           UNIQUE KEY unique_favorite (user_id, post_id)
);

-- 用户活跃度统计表 (用于数据分析)
CREATE TABLE user_activities (
                                 id BIGINT PRIMARY KEY AUTO_INCREMENT,
                                 user_id BIGINT,
                                 date DATE,
                                 post_count INT DEFAULT 0,
                                 comment_count INT DEFAULT 0,
                                 like_count INT DEFAULT 0,
                                 view_count INT DEFAULT 0,
                                 UNIQUE KEY unique_user_daily (user_id, date)
);

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

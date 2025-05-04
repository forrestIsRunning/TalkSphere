ALTER TABLE posts
ADD COLUMN excerpt VARCHAR(255) DEFAULT NULL COMMENT '帖子摘要' AFTER content; 

-- 修改现有的 excerpt 列的字符集
ALTER TABLE posts MODIFY excerpt VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 修改整个表的字符集（为了确保所有文本列都使用正确的字符集）
ALTER TABLE posts CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
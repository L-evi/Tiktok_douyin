# 创建一个数据库 用户名 密码
CREATE DATABASE IF NOT EXISTS `tiktok` DEFAULT CHARACTER SET utf8 COLLATE utf8mb4_general_ci;
CREATE USER 'tiktok'@'%' IDENTIFIED BY 'tiktok';
GRANT ALL PRIVILEGES ON `tiktok`.* TO 'tiktok'@'%';
FLUSH PRIVILEGES;

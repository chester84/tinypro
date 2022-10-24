CREATE DATABASE IF NOT EXISTS `tinypro` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'tinypro_prod'@'%' IDENTIFIED BY  'xxxxxxxx';
GRANT ALL PRIVILEGES ON tinypro.* TO 'tinypro_prod'@'%';
FLUSH PRIVILEGES;



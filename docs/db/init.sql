CREATE DATABASE IF NOT EXISTS `tinypro` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE tinypro;

INSERT IGNORE INTO `admin` (`id`, `email`, `nickname`, `password`, `status`, `register_time`, `last_login_time`) VALUES
(1, 'admin@tinypro.com', 'admin', 'c8269e9955fc424c54328820d50cda6c', 1, 1516710209872, 0);

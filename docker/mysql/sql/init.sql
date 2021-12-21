-- users
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL DEFAULT '',
  `email` varchar(50) NOT NULL DEFAULT '',
  `password` varchar(200) NOT NULL DEFAULT '',
  `salt` varchar(10) NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
INSERT INTO `users` (`id`, `name`, `email`, `password`, `salt`) VALUES
(1, 'ユーザー１', 'sample@example.com', 'password', 'salt');

-- tasks
DROP TABLE IF EXISTS `tasks`;
CREATE TABLE `tasks` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `user_id` int(10) NOT NULL,
  `title` varchar(50) NOT NULL DEFAULT '',
  `content` varchar(200) NOT NULL DEFAULT '',
  `due_date` datetime,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
INSERT INTO `tasks` (`title`, `user_id`, `content`, `due_date`) VALUES
('お米を炊く', 1, '買い出しに帰ってくる頃に炊けるようにセットしておく', '2021-12-05 20:30:09'),
('買い出しに行く', 1,'スーパーで、卵と鶏肉と三葉を買う', '2021-12-05 20:30:09'),
('晩御飯を作る', 1,'親子丼を作る', '2021-12-05 20:30:09'),
('お風呂に入る', 1,'肩まで浸かって10数えよう', '2021-12-05 20:30:09');

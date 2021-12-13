DROP TABLE IF EXISTS `users`;
CREATE TABLE `tasks` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(63) NOT NULL DEFAULT '',
  `content` varchar(255) NOT NULL DEFAULT '',
  `due_date` datetime,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `tasks` (`title`, `content`, `due_date`) VALUES
('お米を炊く', '買い出しに帰ってくる頃に炊けるようにセットしておく', '2021-12-05 20:30:09'),
('買い出しに行く', 'スーパーで、卵と鶏肉と三葉を買う', '2021-12-05 20:30:09'),
('晩御飯を作る', '親子丼を作る', '2021-12-05 20:30:09'),
('お風呂に入る', '肩まで浸かって10数えよう', '2021-12-05 20:30:09');

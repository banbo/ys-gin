CREATE TABLE `ys_gin_example`.`user` (
  `uid` varchar(30) NOT NULL,
  `name` varchar(30) NOT NULL,
  `age` int(10) NOT NULL,
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `user` VALUES ('270200547871555584', 'test1', 18);
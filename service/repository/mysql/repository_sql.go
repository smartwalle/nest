package mysql

const nestSQL = "" +
	"CREATE TABLE IF NOT EXISTS `%s` (" +
	"`id` int(11) NOT NULL AUTO_INCREMENT," +
	"`ctx` int(11) DEFAULT '0'," +
	"`name` varchar(128) DEFAULT NULL," +
	"`left_value` int(11) DEFAULT NULL," +
	"`right_value` int(11) DEFAULT NULL," +
	"`depth` int(11) DEFAULT NULL," +
	"`status` int(11) DEFAULT '1'," +
	"`created_on` datetime DEFAULT NULL," +
	"`updated_on` datetime DEFAULT NULL," +
	"PRIMARY KEY (`id`)," +
	"UNIQUE KEY `%s_id_uindex` (`id`)," +
	"KEY `%s_ctx_index` (`ctx`)," +
	"KEY `%s_id_ctx_index` (`id`,`ctx`)," +
	"KEY `%s_ctx_right_value_index` (`ctx`,`right_value`)," +
	"KEY `%s_ctx_left_value_index` (`ctx`,`left_value`)" +
	") ENGINE=InnoDB DEFAULT CHARSET=utf8;"

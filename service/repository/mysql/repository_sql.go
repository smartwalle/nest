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
	"UNIQUE KEY `nest_id_uindex` (`id`)," +
	"KEY `nest_ctx_index` (`ctx`)," +
	"KEY `nest_id_ctx_index` (`id`,`ctx`)" +
	") ENGINE=InnoDB DEFAULT CHARSET=utf8;"

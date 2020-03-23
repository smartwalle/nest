package sql

const initMySQLTable = "" +
	"CREATE TABLE IF NOT EXISTS `%s` (" +
	"`id` bigint(20) NOT NULL AUTO_INCREMENT," +
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

const initPostgreSQLTable = "" +
	"create table if not exists public.nest(" +
	"id bigint not null constraint nest_pk primary key," +
	"ctx integer default 0," +
	"name varchar(128)," +
	"left_value integer," +
	"right_value integer," +
	"depth integer," +
	"status integer default 1," +
	"created_on timestamp with time zone," +
	"updated_on timestamp with time zone);" +
	"create unique index if not exists nest_id_uindex on public.nest (id);" +
	"create index if not exists nest_ctx_index on public.nest (ctx);" +
	"create index if not exists nest_ctx_left_value_index on public.nest (ctx, left_value);" +
	"create index if not exists nest_ctx_right_value_index on public.nest (ctx, right_value);" +
	"create index if not exists nest_id_ctx_index on public.nest (id, ctx);"

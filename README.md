
```sql
CREATE TABLE `org_department` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ctx` int(11) DEFAULT 0 COMMENT,
  `name` varchar(128) DEFAULT NULL,
  `left_value` int(11) DEFAULT NULL,
  `right_value` int(11) DEFAULT NULL,
  `depth` int(11) DEFAULT NULL,
  `status` int(11) DEFAULT 1000,
  `created_on` datetime DEFAULT NULL,
  `updated_on` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `org_department_id_uindex` (`id`),
  KEY `org_department_ctx_index` (`ctx`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8
```
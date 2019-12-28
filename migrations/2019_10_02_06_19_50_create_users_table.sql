/* -- migrate_up -- */

create table users
(
	id            bigint primary key not null auto_increment,
	guid          varchar(40)        not null unique,
	name          varchar(100)            default '',
	bio           varchar(400)            default '',
	username      varchar(100)       not null unique,
	mobile_number varchar(11)        not null unique,
	password      varchar(100)       not null,
	active        bool                    default true,
	superuser     bool                    default false,
	created_date  timestamp               default current_timestamp,
	modified_date timestamp          null default null on update current_timestamp
);

/* -- migrate_down -- */
drop table users;
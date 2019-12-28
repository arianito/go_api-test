/* -- migrate_up -- */
create table photos
(
	id            bigint primary key not null auto_increment,
	guid          varchar(40)        not null unique,
	user_id       bigint             not null,
	photo_url     varchar(255)       not null,
	photo_ratio   float              not null,
	photo_mime    varchar(255)       not null,
	photo_size    bigint             not null,
	in_use        bool                    default false,
	created_date  timestamp               default current_timestamp,
	modified_date timestamp          null default null on update current_timestamp,
	foreign key (user_id) references users (id)
);
/* -- migrate_down -- */
drop table photos;
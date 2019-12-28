/* -- migrate_up -- */
create table posts
(
	id            bigint primary key not null auto_increment,
	guid          varchar(40)        not null unique,
	user_id       bigint             not null,
	title         varchar(100)            default '',
	photo_id      bigint             not null,
	created_date  timestamp               default current_timestamp,
	modified_date timestamp          null default null on update current_timestamp,
	foreign key (user_id) references users (id),
	foreign key (photo_id) references photos (id)
);
/* -- migrate_down -- */
drop table posts;
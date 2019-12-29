/* -- migrate_up -- */
create table post_comments
(
	id           bigint primary key not null auto_increment,
	user_id      bigint             not null,
	post_id      bigint             not null,
	message      varchar(400)       not null,
	created_date timestamp default current_timestamp,
	foreign key (user_id) references users (id),
	foreign key (post_id) references posts (id)
);
/* -- migrate_down -- */
drop table post_comments;
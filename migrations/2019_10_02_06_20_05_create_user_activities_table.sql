/* -- migrate_up -- */
create table user_activities
(
	id            bigint primary key not null auto_increment,
	user_id       bigint             not null,
	activity_type int                not null,
	user_b        bigint,
	post_b        bigint,
	created_date  timestamp default current_timestamp,
	foreign key (user_id) references users (id),
	foreign key (user_b) references users (id),
	foreign key (post_b) references posts (id)
);
/* -- migrate_down -- */
drop table user_activities;
create table card (
	id varchar(255) not null primary key,
	previous_step varchar(255) not null,
	current_step varchar(255) not null,
	status varchar(255) not null
);

create table actionlog (
	card_id varchar(255) not null,
	action varchar(255) not null,
	status varchar(255) not null,
	msg text,
	foreign key (card_id)
		references card(id)
		on delete cascade
);


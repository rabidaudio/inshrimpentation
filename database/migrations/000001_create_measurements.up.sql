create table measurements (
	id integer not null primary key,
	time timestamp not null,
	type varchar not null,
	value numeric not null
);

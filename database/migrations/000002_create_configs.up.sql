create table configs (
	id integer not null primary key,
	type string not null,
	interval integer not null,
	last_completed timestamp not null
);
CREATE UNIQUE INDEX idx_configs_type ON configs(type);

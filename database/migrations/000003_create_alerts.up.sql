CREATE TABLE alerts (
    id integer not null primary key,
    type string not null,
    time timestamp not null,
    suppressed boolean not null,
    message text not null
);
CREATE INDEX idx_alerts on alerts(suppressed, time);

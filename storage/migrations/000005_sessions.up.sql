create table sessions (
  token text primary key,
  data bytea not null,
  expiry timestamp not null
);

create index on sessions (expiry);

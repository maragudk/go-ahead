create table sessions (
  token string primary key,
  data bytes not null,
  expiry timestamp not null,
  index sessions_expiry_idx (expiry)
);

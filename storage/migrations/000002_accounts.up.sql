-- accounts have names and will be used primarily for billing.
create table accounts (
  id uuid primary key default gen_random_uuid(),
  name string not null,
  created timestamp not null default current_timestamp(),
  updated timestamp not null default current_timestamp()
);

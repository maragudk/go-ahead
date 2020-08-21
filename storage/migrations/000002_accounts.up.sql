-- accounts have names and will be used primarily for billing.
create table accounts (
  id uuid primary key default gen_random_uuid(),
  name string not null,
  created_at timestamp not null default current_timestamp(),
  updated_at timestamp not null default current_timestamp()
);

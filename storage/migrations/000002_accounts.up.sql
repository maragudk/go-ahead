create extension if not exists "uuid-ossp";

-- accounts have names and will be used primarily for billing.
create table accounts (
  id uuid primary key default uuid_generate_v4(),
  name text not null,
  created timestamp not null default now(),
  updated timestamp not null default now()
);

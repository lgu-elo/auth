alter table users alter column created_at set default now();
alter table users alter column updated_at set default now();
alter table users alter column "role" set default 'user';
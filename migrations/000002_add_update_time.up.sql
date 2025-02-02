alter table public.position
    add updatedAt date default now() not null;


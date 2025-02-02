create table if not exists public.position
(
    id    serial not null
        constraint position_pk
            primary key,
    name  text   not null,
    price integer
);


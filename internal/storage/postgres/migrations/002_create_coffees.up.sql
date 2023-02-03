create table if not exists "coffees" (
       "id" UUID primary key,
       "name" varchar(255) not null,
       "created_at" timestamp not null,
       "state" varchar(255) not null,
       "bitterness" varchar(255) not null,
       "acidity" float not null,
       "short_desc" varchar(255) not null
);

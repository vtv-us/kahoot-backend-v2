create table "slide" (
    "id" text primary key,
    "owner" text not null,
    "title" text not null,
    "content" text not null,
    "created_at" timestamp not null default now(),
    "updated_at" timestamp not null default now()
);

create table "question" (
    "id" text primary key,
    "slide_id" text not null,
    "raw_question" text not null,
    "answer_a" text not null,
    "answer_b" text not null,
    "answer_c" text not null,
    "answer_d" text not null,
    "correct_answer" text not null,
    "created_at" timestamp not null default now(),
    "updated_at" timestamp not null default now()
);

alter table "question" add foreign key ("slide_id") references "slide" ("id");

create index on "slide" using btree ("id");

alter table "slide" add foreign key ("owner") references "user" ("user_id");
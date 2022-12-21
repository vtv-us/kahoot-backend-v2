create table "user_question" (
    "question_id" text not null,
    "slide_id" text not null,
    "username" text not null,
    "content" text not null,
    "votes" integer not null default 0,
    "answered" boolean not null default false,
    "created_at" timestamp not null default now(),
    constraint "user_question_pkey" primary key ("question_id")
)

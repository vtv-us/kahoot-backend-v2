create table "collab" (
    "user_id" text not null,
    "slide_id" text not null,
    "created_at" timestamptz not null default (now()),
    constraint "collab_pkey" primary key ("user_id", "slide_id")
);

alter table "question" add column "type" text not null default 'multiple-choice';
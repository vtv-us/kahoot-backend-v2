create table "chat_msg" (
    "id" text not null,
    "slide_id" text not null,
    "username" text not null,
    "content" text not null,
    "created_at" timestamptz not null default (now()),
    constraint "chat_msg_pkey" primary key ("id")
);
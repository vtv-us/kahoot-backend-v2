create table "answer_history" (
    "id" text primary key,
    "slide_id" text not null,
    "raw_question" text not null,
    "raw_answer" text not null,
    "num_chosen" integer not null default 0,
    "created_at" timestamp not null default now(),
    "updated_at" timestamp not null default now()
);

create function "slide_delete_trigger" () returns trigger as $$
begin
    delete from "answer_history" where "slide_id" = old."id";
    return nil;
end;$$ LANGUAGE plpgsql;

create trigger "slide_delete_trigger" before delete on "slide" for each row execute procedure "slide_delete_trigger"();

create function "check_question_index_exists" () returns trigger as $$
begin
    if exists (select 1 from "question" where "slide_id" = new."slide_id" and "index" = new."index" and "id" != new."id") then
        raise exception 'question index already exists';
    end if;
    return new;
end;$$ LANGUAGE plpgsql;

create trigger "check_question_index_exists" before insert or update on "question" for each row execute procedure "check_question_index_exists"();

create function "check_answer_index_exists" () returns trigger as $$
begin
    if exists (select 1 from "answer" where "question_id" = new."question_id" and "index" = new."index" and "id" != new."id") then
        raise exception 'answer index already exists';
    end if;
    return new;
end;$$ LANGUAGE plpgsql;

create trigger "check_answer_index_exists" before insert or update on "answer" for each row execute procedure "check_answer_index_exists"();

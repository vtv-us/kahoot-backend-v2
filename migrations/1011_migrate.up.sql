create or replace function "delete_slide" () returns trigger as $$
begin
    delete from "question" where "slide_id" = old."id";
    return old;
end;$$ LANGUAGE plpgsql;

create trigger "delete_slide" before delete on "slide" for each row execute procedure "delete_slide"();

create or replace function "delete_question" () returns trigger as $$
begin
    delete from "answer" where "question_id" = old."id";
    return old;
end;$$ LANGUAGE plpgsql;

create trigger "delete_question" before delete on "question" for each row execute procedure "delete_question"();

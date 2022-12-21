create function "check_true_index_answer" () returns trigger as $$
begin
    if not exists (select 1 from "answer" where "question_id" = new."question_id" and ("index" = new."index" - 1 or new."index" = 1)) then
        raise exception 'wrong index';
    end if;
    return new;
end;$$ LANGUAGE plpgsql;

create trigger "check_true_index_answer" before insert or update on "answer" for each row execute procedure "check_true_index_answer"();

create function "check_true_index_question" () returns trigger as $$
begin
    if not exists (select 1 from "question" where "slide_id" = new."slide_id" and ("index" = new."index" - 1 or new."index" = 1)) then
        raise exception 'wrong index';
    end if;
    return new;
end;$$ LANGUAGE plpgsql;

create trigger "check_true_index_question" before insert or update on "question" for each row execute procedure "check_true_index_question"();

create function "auto_decrease_index_when_delete_question" () returns trigger as $$
begin
    if old."index" != (select max("index") from "question" where "slide_id" = old."slide_id") then
        update "question" set "index" = "index" - 1 where "slide_id" = old."slide_id" and "index" > old."index";
    end if;
    return nil;
end;$$ LANGUAGE plpgsql;

create trigger "auto_decrease_index_when_delete_question" before delete on "question" for each row execute procedure "auto_decrease_index_when_delete_question"();

create function "auto_decrease_index_when_delete_answer" () returns trigger as $$
begin
    if old."index" != (select max("index") from "answer" where "question_id" = old."question_id") then
        update "answer" set "index" = "index" - 1 where "question_id" = old."question_id" and "index" > old."index";
    end if;
    return nil;
end;$$ LANGUAGE plpgsql;

create trigger "auto_decrease_index_when_delete_answer" before delete on "answer" for each row execute procedure "auto_decrease_index_when_delete_answer"();


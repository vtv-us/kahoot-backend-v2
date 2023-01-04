create or replace function "auto_decrease_index_when_delete_question" () returns trigger as $$
begin
    if old."index" != (select max("index") from "question" where "slide_id" = old."slide_id") then
        update "question" set "index" = "index" - 1 where "slide_id" = old."slide_id" and "index" > old."index";
    end if;
    return old;
end;$$ LANGUAGE plpgsql;
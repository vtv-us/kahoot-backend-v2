create function "group_delete_trigger" () returns trigger as $$
begin
    delete from "user_group" where "group_id" = old."group_id";
    return old;
end;$$ LANGUAGE plpgsql;

create trigger "group_delete_trigger" before delete on "group" for each row execute procedure "group_delete_trigger"();
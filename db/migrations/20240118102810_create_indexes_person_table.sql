-- +goose Up
-- +goose StatementBegin
create index idx_person_name ON person (name);
create index idx_person_age ON person (age);
create index idx_person_gender ON person (gender);
create index idx_person_nationality ON person (nationality);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index idx_person_name;
drop index idx_person_age;
drop index idx_person_gender;
drop index idx_person_nationality;
-- +goose StatementEnd

ALTER TABLE "users"
    ADD company_id bigserial NOT NULL;

ALTER TABLE "users"
    ADD CONSTRAINT fk_users_companies FOREIGN KEY (company_id) REFERENCES companies (id);

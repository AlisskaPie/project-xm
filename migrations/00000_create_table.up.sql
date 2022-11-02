CREATE TYPE companyType AS ENUM ('Corporations', 'NonProfit', 'Cooperative', 'Sole Proprietorship');
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE company (
    id uuid DEFAULT uuid_generate_v4(),
    name VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    amount_of_employees INT NOT NULL DEFAULT 0,
    registered BOOLEAN NOT NULL DEFAULT FALSE,
    type companyType
);

CREATE UNIQUE INDEX id_idx ON company (id ASC);
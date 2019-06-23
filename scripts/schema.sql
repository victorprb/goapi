CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    uuid uuid DEFAULT uuid_generate_v1mc(),
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    cpf VARCHAR(11) UNIQUE NOT NULL,
    email VARCHAR(50) UNIQUE NOT NULL,
    created TIMESTAMP NOT NULL,
    active BOOLEAN NOT NULL DEFAULT true,
    PRIMARY KEY (uuid)
);

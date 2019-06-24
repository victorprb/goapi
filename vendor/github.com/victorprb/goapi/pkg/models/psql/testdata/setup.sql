CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    uuid uuid DEFAULT gen_random_uuid(),
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    cpf VARCHAR(11) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created TIMESTAMP NOT NULL,
    active BOOLEAN NOT NULL DEFAULT true,
    PRIMARY KEY (uuid)
);

INSERT INTO users (uuid, first_name, last_name, cpf, email, hashed_password, created, active) VALUES(
    '29c31a03-f111-4f6f-b49e-07b0117fbb42',
    'Cloud',
    'Strife',
    '12312312300',
    'cloud.strife@ffvii.com',
    '$2a$12$Egfcp2Ch5O9KYjPo/oBtouBcd0QNbWjwW52dFqJD8oV0mg2VUqwze',
    '2019-06-20 15:00:00',
    'true'
);

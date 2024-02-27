CREATE TABLE data_records
(
    id BIGSERIAL NOT NULL PRIMARY KEY,
    uploaded_at TIMESTAMP NOT NULL,
    type VARCHAR(255) NOT NULL,
    checksum VARCHAR(255),
    data TEXT,
    filepath VARCHAR(255),
    name VARCHAR(255) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    key VARCHAR
);
CREATE TYPE role AS ENUM ('admin', 'hr');

CREATE TABLE iF NOT EXISTS adminAndHr (
    id UUID DEFAULT gen_random_uuid(),
    username VARCHAR(50) NOT NULL,
    password VARCHAR(50) NOT NULL,
    role role NOT NULL,
    status varchar(50),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at BIGINT NOT NULL DEFAULT 0
    
    );

CREATE TABLE IF NOT EXISTS employees (
    id UUID Primary key DEFAULT gen_random_uuid(),
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    position VARCHAR(50) NOT NULL,
    department VARCHAR(50) NOT NULL,
    comment VARCHAR(50),
    created_by VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at BIGINT NOT NULL DEFAULT 0
    );

CREATE TABLE IF NOT EXISTS black_list (
    id UUID DEFAULT gen_random_uuid(),
    employee_id  UUID REFERENCES employee(id),
    reason VARCHAR(50) NOT NULL,
    added_by VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP ,
    deleted_at BIGINT NOT NULL DEFAULT 0
    );



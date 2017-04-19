use midash;

create table user (
    id INTEGER PRIMARY KEY AUTO_INCREMENT NOT NULL,
    salt varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    public_id varchar(255) NOT NULL,
    private_id varchar(255) NOT NULL,
    updated_at timestamp NOT NULL,
    created_at timestamp NOT NULL
);

create table profile (
    id INTEGER PRIMARY KEY AUTO_INCREMENT NOT NULL,
    user_id varchar(255) NOT NULL,
    public_id varchar(255) NOT NULL,
    firstName varchar(255) NOT NULL,
    lastName varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    updated_at timestamp NOT NULL,
    created_at timestamp NOT NULL
);

create table session (
    id INTEGER PRIMARY KEY AUTO_INCREMENT NOT NULL,
    user_id varchar(255) NOT NULL,
    public_id varchar(255) NOT NULL,
    expiration timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    created_at timestamp NOT NULL
);
use midash;

create table users (
    id INTEGER AUTO_INCREMENT NOT NULL,
    hash varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    public_id PRIMARY KEY varchar(255) NOT NULL,
    private_id varchar(255) NOT NULL,
    updated_at timestamp NOT NULL,
    created_at timestamp NOT NULL
);

create table profiles (
    id INTEGER AUTO_INCREMENT NOT NULL,
    address text NOT NULL,
    user_public_id varchar(255) NOT NULL,
    public_id PRIMARY KEY varchar(255) NOT NULL,
    firstName varchar(255) NOT NULL,
    lastName varchar(255) NOT NULL,
    updated_at timestamp NOT NULL,
    created_at timestamp NOT NULL,
    INDEX user_id (user_public_id)
);

create table sessions (
    id INTEGER AUTO_INCREMENT NOT NULL,
    user_public_id varchar(255) NOT NULL,
    public_id PRIMARY KEY varchar(255) NOT NULL,
    token varchar(255) NOT NULL,
    expiration timestamp NOT NULL,
    created_at timestamp NOT NULL,
    INDEX user_id (user_public_id)
);
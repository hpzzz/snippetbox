

CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255)NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);

CREATE TABLE snippets (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL,
    creator INTEGER NOT NULL
);

 ALTER TABLE snippets ADD FOREIGN KEY (creator) REFERENCES users(id);


CREATE INDEX isx_snippets_created ON snippets(created);

INSERT INTO users (name, email, hashed_password, created) VALUES (
    'Kari Hari',
    'Kari@example.com',
    '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
    '2020-07-23 18:12:00',
    '1'
);
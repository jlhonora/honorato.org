## How to create the database and tables

1. Open `psql`
2. Create a user: `CREATE USER pgtestuser WITH PASSWORD 'somepassword';`
3. Create a database: `CREATE DATABASE pgtestdb;`
4. Grant privileges: `GRANT ALL PRIVILEGES ON DATABASE pgtestdb to pgtestuser;`
5. Logout from psql with `\q` and test it: `psql -d pgtestdb -U pgtestuser`
6. Create table: 

```
CREATE TABLE users (
   id bigserial primary key,
   username text NOT NULL,
   password text NOT NULL,
   created_at timestamp);
```

7. Insert something:
```
INSERT INTO users VALUES (default, 'somevalue', 'someothervalue', now());
```

8. Query it:
```
SELECT * FROM users;
```

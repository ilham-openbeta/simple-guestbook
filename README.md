# Simple Guestbook

Simple Guestbook API using Gorilla Mux, GORM, and MySQL databases. Simple Go testing, swagger documentation, and client example included.

Check file static/docs/index.html or redoc.html to view API documentation.

Check file static/index.html for example.

Features :

- Pretty simple create & read the post from the database for the guestbook

### How to install

1. Clone this repository

2. Create the messages table

```sql
CREATE TABLE IF NOT EXISTS messages (
    id int(5) NOT NULL,
    name varchar(50) NOT NULL,
    message varchar(160) NOT NULL,
    contact varchar(50) NOT NULL,
    created_at datetime NOT NULL
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE `messages` ADD PRIMARY KEY (`id`);

ALTER TABLE `messages` MODIFY `id` int(5) NOT NULL AUTO_INCREMENT;
```

2. Create .env files, check env.example for example.

3. Fix package import
   
   ```bash
   go get
   ```

4. Test the code
   
   ```bash
   go test -v
   ```

5. Run the code
   
   ```bash
   go run .
   ```

6. Build the application
   
   ```bash
   go build .
   ```

### TODO

- Add account system to create, update, and delete the post message using JWT auth
- Arrange file structure.
- Generate the OpenApi document automatically. Currently, I write it manually on swagger.yaml file
- Fix user interfaces for client example, make it more beautiful
- Form validation


### LICENSE

MIT

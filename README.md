Create docker.env file in root directory and add following values:
```dotenv
POSTGRES_CONNECTION=postgres://<user>:<password>@db:5432/meem_db?sslmode=disable
POSTGRES_PASSWORD=<password>
POSTGRES_USER=<user>
GLOBAL_PASSWORD_SALT=<random string>
MONGODB_URI=mongodb+srv://<user>:<password>@<cluster>/test?tlsInsecure=true
SECRET_SIGNING_KEY=<random string>
```

DB scheme migration:

```psql -U postgres -h localhost -d meem_db -a -f "<path>/20210923213515_postgresql.up.sql"```


Run

```docker build -t messenger . ```

```docker run -d -p 8080:8080 --name messenger --restart always --env-file ./.env messenger```



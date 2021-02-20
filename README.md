# Url-Shortener

A URL shortener service (basically a glorified hashmap of urls)  written in Golang.

## Requirement

- Go version: 1.16
- Docker and Docker-Compose

## Start App locally

The URL shortener can use [Redis](https://redis.io/), [Postgres](https://www.postgresql.org/) or [Mysqsl](https://dev.mysql.com/) as database.

You can choose between databases Redis, Postgres or Mysql by modifying the `DB_ENGINE` variable in the `.env` file.
By default, Postgres is used: `DB_ENGINE=postgres`

See `docker-compose.yaml` definition for running and configuring databases images locally.

If you just started and wish to run locally go, you may need to download packages.
```bash
make mod
```

You can start the service locally by running:
```bash
make start
```

App will be available on `localhost:8000`. 

You can change the port with `PORT=8000` in `.env` file.

## How it works

### Create a new short URL
Let's create a new short url  that will redirect to `https://www.google.com/`.

First we get a new url code by sending a POST request to  `/create` with `"url": "https://www.google.com/` json data.

````bash
curl --request POST \
  --url http://localhost:8000/create \
  --header 'content-type: application/json' \
  --data '{
	"url": "`https://www.google.com`"
}	'
````

You will receive a json response with the new short url as `code`:
````json5
{
  "code": "7TQ_AlDGR", // New Code URL!!
}
````

That's it! 

Now you can make a get request to `localhost:8000/7TQ_AlDGR` to be redirected to `https://www.google.com`!

### Update short URL
As you can imagine, the url path `/7TQ_AlDGR` is not very friendly. That is why you have the possibility to update this url code!

You need to send a `PATCH` request to `/update` with the previous `code` that you created earlier and ask for a `new_code` that you define yourself!
Careful, the new URL should be more than **3 characters** and less than **20 characters** long,  if not it should not be called a short url right?!

```bash
curl --request PATCH \
  --url http://localhost:8000/update \
  --header 'content-type: application/json' \
  --cookie JSESSIONID=89293DD417BC8EBDAB2DCF7BE157A217 \
  --data '{
	"code": "7TQ_AlDGR",
	"new_code": "MyNewAwsomeURL"
}		'
```

The request will respond with the following json:
```bash
{
  "code": "MyNewAwsomeURL",
}
```

You now have replaced the new short URL `/7TQ_AlDGR` by your customized short URL `/MyNewAwsomeURL`!

Try it out at `localhost:8000/MyNewAwsomeURL`, it should redirect you to `https://www.google.com` .


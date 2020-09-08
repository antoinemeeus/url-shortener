# Url-Shortener

An URL shortener service (basically a glorified hashmap of urls)  written in Golang.

## Requirement

Go version: 1.14

## Start App locally

The URL shortenenr can use [Redis](https://redis.io/) or [Postgres](https://www.postgresql.org/) as database.

You can choose between Redis or Postgres by modifying the `DB_ENGINE` variable in the `.env` file.
By default Postgres is used: `DB_ENGINE=postgres`

See `docker-compose.yaml` definition for running and configuring databases images locally.

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
  "ID": 1,
  "created_at": "2020-09-08T21:05:00.389655Z",
  "updated_at": "2020-09-08T23:05:00.391674+02:00",
  "deleted_at": null,
  "code": "7TQ_AlDGR", // New Code URL!!
  "new_code": "",
  "url": "https://www.google.com"
}
````

That's it! 

Now you can make a get request to `localhost:8000/7TQ_AlDGR` to be redirected to `https://www.google.com`!

### Update short URL
As you can imagine, the url path `/7TQ_AlDGR` is not very friendly. That is why you can update this url code!

You need to send a `PUT` request to `/update` with the previous `code` that you created earlier and ask for a `new_code` that you define yourself!
Careful, the new URL should be less than **20 characters** long,  if not it should not be called a short url right?!

```bash
curl --request PUT \
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
  "ID": 8,
  "created_at": "2020-09-08T21:36:09.428789Z",
  "updated_at": "2020-09-08T23:36:09.438257+02:00",
  "deleted_at": null,
  "code": "MyNewAwsomeURL",
  "new_code": "MyNewAwsomeURL",
  "url": "https://google.com"
}
```

And now you replaced the new short URL `/7TQ_AlDGR` by customized short URL `/MyNewAwsomeURL`!

Try it out at `localhost:8000/MyNewAwsomeURL`, it should redirect you to `https://www.google.com` .


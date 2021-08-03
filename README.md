
# TC List

A barebones Go app for tracking lists

## Deploying to Heroku

```sh
$ heroku create
$ git push heroku main
$ heroku open
```

## API

GET `/items` lists all items

```json
[
    {
        "id": 3,
        "name": "Spinach",
        "action": ""
    },
    {
        "id": 4,
        "name": "Oreos",
        "action": ""
    },
    {
        "id": 2,
        "name": "Coke",
        "action": "BOUGHT"
    }
]
```

POST `/items` with name creates

```json
{
  "name": "Bananas"
}
```

POST `/items` with id and action updates

```json
{
  "id": 1,
  "action": "BOUGHT"
}
```

DELETE `/items` removes all items with an `action` specified


## Documentation

For more information about using Go on Heroku, see these Dev Center articles:

- [Go on Heroku](https://devcenter.heroku.com/categories/go)

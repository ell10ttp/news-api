## news-api

An implementation of a RESTful API to consume a public news feed.

In this implementation, we hold all sources and categories in memory. In an production environment or at larger scale this would be performed in a more stable database. BBC and Sky News are set up as default, with a category model that holds a map of the common news types to the values of the xml feed of their location. This is not consistent across different source providers and the rss feed link should be checked before. All links provided for BBC and Sky News are valid on startup.

# deploying the application

The application uses logging similar style to most SOCs, with common tags for easy extraction by a tool such as Elasticsearch, and includes some debugging provided. In version control, in the [deployments](https://github.com/ell10ttp/news-api/tree/main/deployments) directory, is already held a '.env' file with current environment variables. Deployment is immediately ready to go through.

Clone the directory:
`git clone git@github.com:ell10ttp/news-api.git`

Traverse in the [deployments](https://github.com/ell10ttp/news-api/tree/main/deployments) and run docker-compose.
`docker-compose up`

This will launch a container of the application.

The default server port is 5000. To change this, edit the SERVER_PORT env var in the .env file.

# using the application

Each Source has an Unique ID held in struct. This is what is used to reference the source in future calls.

To see available sources:

```
curl --request GET http://localhost:5000/source
```

To create a new source (e.g. adding CNN):
 - Name, Description, Url, Language and Country are all required fields
 - Adding new categories is not yet implemented.

```
curl --header "Content-Type: application/json" \
--data '{"Name":"CNN", "Description":"News from CCN delivered directly", "Url":"http://rss.cnn.com/rss/edition.rss", "Language": "en-US", "Country": "USA"}' \
--request POST http://localhost:5000/source
```

To retrieve a specific news feed (e.g. BBC News, default source id: 1)

```
curl --request GET http://localhost:5000/source/1
```

Articles are held in the `items` json array.
Thumbnails, if present, are returned as URLs for the browser/app to collect held in the `extensions.thumbnails`.

To check which categories are available on a source:

```
curl --request GET http://localhost:5000/source/1/category
```
This will return a JSON response like so:
`{"action":"retrieve available categories","successful":true,"numberOfCategories":6,"categories":["business","entertainment","politics","technology","uk","world"]}`


To retrieve a specific news feed with a category filter:
Use any of the categories as a url parameter on the retrieval of a news feed e.g:
```
curl --request GET http://localhost:5000/source/1?category=technology
```


To return news feed sorted by publish date, using a time.Time comparison:
```
curl --request GET http://localhost:5000/source/1?sort=true
```

The above two operations may be combined and are not mutually exclusive and can be combined:

```
curl --request GET http://localhost:5000/source/1?sort=true&category=technology
```

# Article Service

## Introduction
RESTful Go based JSON API built using the Gorilla framework. The API allows CRUD based operations on an articles database.
It provides following endpoints:
1. POST /articles 

This handles the receipt of some article data in json format, and store it within the postgres database.

2. GET /articles/{id} 

This return the JSON representation of the article in the following format:
```
{
  "id": "1",
  "title": "latest science shows that potato chips are better for you than sugar",
  "date" : "2016-09-22",
  "body" : "some text, potentially containing simple markup about how potato chips are great",
  "tags" : ["health", "fitness", "science"]
}
```

3. GET /tags/{tagName}/{date} 

This returns the list of article ids that have that tag name on the given date and some summary data about that tag for that day in the following format:
```
{
  "tag" : "health",
  "count" : 17,
    "articles" :
      [
        "1",
        "7"
      ],
    "related_tags" :
      [
        "science",
        "fitness"
      ]
}
```

## Getting Started

### Prerequisites

1. Install Docker Engine - [ Instructions to install Docker can be found [here](https://docs.docker.com/get-docker/)]
2. Install Docker Compose Plugin
3. Install Go - [Instructions to install Go can be found [here](https://go.dev/doc/install)]

### Running Application
1. Clone the Github repository git@github.com:sg83/go-microservice.git.
2. Start Docker.
3. Navigate to the article-api directory.
4. Run the following command to start the API and Postgres database containers:
```
docker compose up
```
The database container will start first, followed by the API service container. The API server will listen on port 8080. Once both containers are running, you can test the endpoints using a client such as Postman.

Alternatively, you can test the endpoints using curl. Here are some example commands:
```
curl localhost:8080/articles/1   

curl localhost:8080/articles -XPOST -d '{"Title": "Article3", "Body": "Some text about lifestyle and fitness", "Date": "2023-04-07", "Tags":["lifestyle", "fitness", "yoga"]}'

curl localhost:8080/tags/health/20230407 
```


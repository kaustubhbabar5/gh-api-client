# GitHub Users API
This is an application which provides an API to retrieve detailed information about github users

## Requirements
+ it should allow retrieval of up to a max of 10 users
+ it should return back some basic information about user such as
    + name
    + login
    + company
    + number of followers
    + number of public repos
+ the returned users should be sorted alphabetically by name

## Use Cases
1. usernames not found in cache or not able to communicate with cache
    -   Use github API to fetch information of each user and aggregate response
    -   backfill the cache
    -   return user information with http status code 200
2. usernames not found in cache and Github API fails for all users
    -   return http response with Http Status Code 500 (Internal Server Error)
3. partial usernames found in cache
    -   Hit the Github API for fetch information for remaining users
        -   if Github API is successful
            -   backfill the cache and aggregate response
            -   respond with 
        -   if not Github API fails
            -   return partial users and `could not able to connect to Github API` for remaining Users 


## Start Up instructions
-   make sure you have docker and docker-compose installed
-   make sure you have .env in the current directly
    -   ```
        HOST=0.0.0.0
        PORT=8080
        REDIS_URL=0.0.0.0:6379
        REDIS_PASSWORD=<some-random-password>
        GITHUB_AUTH_TOKEN=<github-auth-token>
        ```
    -   you can generate `github-auth-token` (Personal access tokens) form [here](https://github.com/settings/tokens)
    - use `make run` to start the server
    - Swagger [API docs](http://localhost:8080/swagger/index.html)


## Run tests
-   make sure you have .env setup
-   use `make test` to run tests

## Considerations
-   Fields which are not existent in the GitHub API response are omitted in the server's API response
-   Users in response are sorted by `Login` field alphabetically in ascending order, case insensitive
-   GET request with body is used. [Ref](https://www.elastic.co/guide/en/elasticsearch/guide/current/_empty_search.html#get_vs_post)
-   As Github does not have API for fetching user details for multiple users in a single call and  latencies for sequentially fetching 10 users details was over 3 seconds, fetched user details concurrently
-   Left out things
    -   Use case of rate limits
    -   use of contexts for timeouts and cancellation
    -   retry mechanism in case of not being able to communicate to github or rate limit reached
    -   instrumenting service for observability
    -   use of Github GraphQL API instead of REST, which allows fetching user details in bulk in single API call




    
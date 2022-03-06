# Rover Pic

## Run Server
In the project directory, inspect the config file `roverpic.toml`. Modify as required.
Replace `DEMO_KEY` in `APIKey` with your real API key (or not).

In the project directory, execute 
```
go run .
``` 
to run the server.
It looks for the config file in the current directory.

## Run Tests
This will run all tests:
```
go test ./...
```
Tests were however only written for the `roverapi` package.

## Send Requests
The endpoint for the assignment is implemented at `/download` and requires `date`
as a query parameter in the format `YYYY-MM-DD`.

Requests can be send to the server using curl, for example:
```
curl 'localhost:8080/download?date=2021-05-11'
```
returns the result:
```
{"api_success":true,"downloaded":340}
```

## Notes
- The number of concurrent downloads is limited for the server, which prevents
  unbounded resource usage. The number can be set in the config.
- Server errors are returned over http, and could be security risk. 
  No attempt was made to hide sensitive information in this assignment.
- Photos are re-downloaded for every request. If we introduce a requirement to skip
  ones already downloaded, we would also have to check for and delete any partially 
  downloaded files when an error occurs, or perhaps download to a `.partial` file.
- When a request over HTTP is cancelled, the currently downloading photos in that
  request are not cancelled. To do this we would need to introduce IO functions
  with context support.

# API pagination playground

An experimentation in API pagination and its pitfalls.

There are three components:

1) A Go backend API which handles requests
2) A Python client who constantly hammers the API and updates & creates records.
3) A Python client who tries to export all the data via API pagination. 

You will need [goose](https://github.com/pressly/goose#install), docker-compose, Go and Python3.

- `docker-compose up`
- `goose -dir migrations/ postgres "user=calendar dbname=postgres sslmode=disable password=calendar" reset`
- `goose -dir migrations/ postgres "user=calendar dbname=postgres sslmode=disable password=calendar" up`
- `go build`  

Then,

- Start the API: `./api-pagination-playground`
- Start client who updates and posts data: `python3 update_events.py`
- Try to export all records with the other client:
`python3 export_events.py`

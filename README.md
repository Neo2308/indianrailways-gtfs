## How to run
```shell
X_API_KEY= go run ./
```
TODO: Add instructions to get API key

## Examples of urls

- Show current map
  - http://localhost:8080/map
- Show map of train number
  - http://localhost:8080/map/22646
- Search by keyword and show the respective trains on map
  - http://localhost:8080/map/search/VANDE
- Show map of trains coming to a station in next 8 hours
  - http://localhost:8080/map/liveStation/BPL
- Show all data from cache
  - http://localhost:8080/show/all
- Download data in gtfs format
  - http://localhost:8080/save
  - Note:
    - only uses currently loaded data to generate gtfs
    - Best to run http://localhost:8080/show/all before this.



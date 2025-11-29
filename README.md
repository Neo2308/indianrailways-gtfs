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
  - http://localhost:8080/map/show/all
- Download data in gtfs format
  - http://localhost:8080/save
  - Note:
    - only uses currently loaded data to generate gtfs
    - Best to run http://localhost:8080/map/show/all before this.

## Checking for errors post fetch
- Search everywhere for `"rc": "4` this will show all 4xx errors in the train service profile info.
- Data for these trains can be re-fetched by deleting respective cached files and using:
  - http://localhost:8080/map/show/all

## Status Report
**Data generated**: Around 9th Nov 2025

**Station fixing report (from the code output):**
* 708 total errors
* 0 total warnings
* 140 unique error stations
* 0 unique warning stations

**GTFS review report** (from [transport.data.gouv.fr](https://transport.data.gouv.fr/))
https://transport.data.gouv.fr/validation/517479?token=3ffe848d-2f44-47a6-aa60-c375ac43639b
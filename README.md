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
###### API Errors
- Search everywhere for `"rc": "4` this will show all 4xx errors in the train service profile info.
- Data for these trains can be re-fetched by deleting respective cached files and using:
  - http://localhost:8080/map/show/all
###### Missing info
- Trigger http://localhost:8080/save to see if any errors/warnings are generated in the console.
- Check [fixes-data/adjustments.yaml](fixes-data/adjustments.yaml) and fix the sections for erroringStations, warningStations and move them to updates/overrides

## Status Report
**Data generated**: Around 9th Nov 2025

**Station fixing report (from the code output):**
* 0 total errors
* 0 total warnings
* 0 unique error stations
* 0 unique warning stations
* 6 ignored trains

**GTFS review report** (from [transport.data.gouv.fr](https://transport.data.gouv.fr/))
https://transport.data.gouv.fr/validation/517981?token=18b3d9ec-ab04-4ee2-aa52-2443d550ab60

## Users
- [Transitous](https://transitous.org)
  - Shoutout to [Lach-anonym](https://github.com/Lach-anonym) for adding the info here ([PR](https://github.com/public-transport/transitous/pull/1614))
- [MobilityDatabase](https://mobilitydatabase.org/feeds/gtfs/mdb-2867)
- [CatenaryMaps](https://maps.catenarymaps.org/#pos=4.39/20.69/84.4)
  - Shoutout to [Ethanc8](https://github.com/ethanc8) for adding the info here ([PR](https://github.com/catenarytransit/transitland-atlas/pull/58))
- [TransitRouter](https://transitrouter.vonter.in/#/railways/)
  - Shoutout to [Vonter](https://github.com/Vonter)
  - [Source Code](https://github.com/Vonter/transitrouter)

(Please open an issue/PR to add your project here using this data)
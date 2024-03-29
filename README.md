

# A guide to [_Soarch_ ]( https://adityamohan29.github.io/Soarch/)

### A flight search tool that gives you the best routes for optimal prices based on spatial querying. 

Soarch makes use of kiwi.com's open flight search APIs to find prices for **_direct flights_** between two points and uses spatial querying to formulate a route based on optimal prices. For best results, it's recommended to get the flight numbers and individual routes through the soarch interface and search for the flight for the particular route from the respective airline's website.

It displays the airport IATA code, airport names (when hovering on the IATA codes), flight numbers, deep link to Kiwi's booking page (hyperlinked to the flight codes), total duration and the total price of the flight.

Soarch makes use of POSTGIS (an extension of the POSTGRES database) to store geographic data and perform spatial queries. The tool includes the functionality to display the best routes with two stops or less ( which would suffice to connect any two airports in the world ). Golang was chosen as the primary language for the backend for API processing and serving and concurrent function execution (goroutines). The backend service is currently depolyed on heroku.com. **_Important note_** : The free tier dyno was chosen for this heroku application which causes the first request to take some amount of time. However, subsequent requests made within an hour of activity would be sped up.

Update: This project can no longer be hosted since heroku had deleted all free tier dyno databases 


### Methodology

This tool different spatial querying mechanisms to display the routes depending on the number of stops. 
#### 1. For One-Stop flights

This is a pretty straightforward dump of the Kiwi API Call between two airports. 


#### 2. For Two-Stop flights

For flights from point A to B, soarch first finds the neighbouring airports (using POSTGIS's nearest neighbour operation)  to point B ( let's call them _mid_airports_ ) and finds the direct flights to _mid_airports_ from point A. If there exists a direct flight to _mid_airports_ , then the API is called from those airports as the origin, till the destination provided by the user.


#### 2. For Three-Stop flights


In this case, for flights from point A to B, soarch first performs a nearest neighbour query in A and B and finds direct flights to airports close to A ( _mid_flights_A_ ) and B ( _mid_flights_B_ ) themseleves. Once these direct flights are established soarch finds the flights from _mid_flights_A_ to _mid_flights_B_ ).

At the end, the result from all these procedures are joined into a table and are ordered based on their price.


A few remarks:

 1. Currently, this web application works best on chrome in Windows/Linux environments. 
 2. The first request might take a bit longer than usual, and the requests which follow might end up being slightly faster.


Enjoy soarching!


#### References

Airport Data: https://ourairports.com/data/  \
Kiwi API: https://docs.kiwi.com/search-api/#/paths/~1flights/get  \
PostGIS Queries: http://postgis.net/

Update: Heroku has revoked my DB access for some reason, so the website won't be working until I fix it. Sorry for the inconvenience! 

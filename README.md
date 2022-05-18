## https://adityamohan29.github.io/Soarch/


# A guide to [_Soarch_ ]( https://adityamohan29.github.io/Soarch/)

### A flight search tool that gives you the best routes for optimal prices based on spatial querying. 

Soarch makes use of kiwi.com's open flight search APIs to find prices for **_direct flights_** between two points and uses spatial querying to formulate a route based on optimal prices. For best results, it's recommended to get the flight numbers and individual routes through the soarch interface and search for the flight for the particular route from the respective airline's website.

It displays the airport IATA code, airport names (when hovering on the IATA codes), flight numbers, deep link to Kiwi's booking page (hyperlinked to the flight codes), total duration and the total price of the flight.

Soarch makes use of POSTGIS (an extension of the POSTGRES database) to store geographic data and perform spatial queries. The tool includes the functionality to display the best routes with two stops or less ( which would suffice to connect any two airports in the world ). Golang was chosen as the primary language for the backend for API processing and serving and concurrent function execution (goroutines).


### Methodology

This tool different spatial querying mechanisms to display the routes depending on the number of stops. 
#### 1. For One-Stop flights

This is a pretty straightforward dump of the Kiwi API Call between two airports. 


#### 2. For Two-Stop flights

For flights from point A to B, soarch first finds the neighbouring airports (using POSTGIS's nearest neighbour operation)  to point B ( let's call them _mid_airports_ ) and finds the direct flights to _mid_airports_ from point A. If there exists a direct flight to _mid_airports_ , then the API is called from those airports as the origin, till the destination provided by the user.


#### 2. For Three-Stop flights


In this case, for flights from point A to B, soarch first performs a nearest neighbour query in A and B and finds direct flights to airports close to A ( _mid_flights_A_ ) and B ( _mid_flights_B_ ) themseleves. Once these direct flights are established soarch finds the flights from _mid_flights_A_ to _mid_flights_B_ ).

At the end, the result from all these procedures are joined into a table and are ordered based on their price.

Enjoy soarching!


#### References

Airport Data: https://ourairports.com/data/
Kiwi API: https://docs.kiwi.com/search-api/#/paths/~1flights/get 

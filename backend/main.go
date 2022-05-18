package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"time"

	"github.com/gorilla/mux"
)

var db *sql.DB

type Flight struct {
	A1            string
	A2            string
	A3            string
	A4            string
	Price         string
	Duration1     string
	Duration2     string
	Duration3     string
	Dtime1        string
	Atime1        string
	Dtime2        string
	Atime2        string
	Dtime3        string
	Atime3        string
	Dtime1UTC     string
	Atime1UTC     string
	Dtime2UTC     string
	Atime2UTC     string
	Dtime3UTC     string
	Atime3UTC     string
	Deeplink1     string
	Deeplink2     string
	Deeplink3     string
	Srciata       string
	Mid1iata      string
	Mid2iata      string
	Dstiata       string
	Airline1      string
	Flightno1     string
	Airline2      string
	Flightno2     string
	Airline3      string
	Flightno3     string
	TotalDuration string
}

func init() {

	var err error
	//connStr := "postgres://postgres:dummypwd4soarch@localhost/postgis_db?sslmode=disable"
	//connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s  sslmode=disable",
	//"postgres", 5432, "postgres", "dummypwd4soarch")
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	//db, err = sql.Open("postgres", connStr)

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("Connected to POST_GIS DB")

}

func main() {

	router := mux.NewRouter()
	truncateAllTables()
	router.HandleFunc("/find/{apfrom}&{apto}&{nearestAirportsSrc}&{nearestAirportsDst}&{dateOfFlight}", handleSoarch).Methods("GET")
	//err := http.ListenAndServe(":8090", router)
	err := http.ListenAndServe(GetPort(),router)
	if err != nil {
 		log.Fatal("ListenAndServe: ", err)
			}

}

func GetPort() string {
	 	var port = os.Getenv("PORT")
	 	// Set a default port if there is nothing in the environment
	 	if port == "" {
	 		port = "4747"
	 		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
		}
		return ":" + port


}


func writeJsonResponse(w http.ResponseWriter) {

	query := "select * from final_merged order by total limit 30;"

	var flights []Flight
	rows, err := db.Query(query)
	defer rows.Close()
	for rows.Next() {
		var src, dtime1, atime1, dtime2, atime2, dtime3, atime3, dtime1UTC, atime1UTC, dtime2UTC, atime2UTC, dtime3UTC, atime3UTC, mid1, mid2, dst, d1, d2, d3, deeplink1, deeplink2, deeplink3, srciata, mid1iata, mid2iata, dstiata, airline1, flightno1, airline2, flightno2, airline3, flightno3, totalduration string
		var total int
		err = rows.Scan(&src, &mid1, &mid2, &dst, &total, &d1, &d2, &d3, &dtime1, &atime1, &dtime2, &atime2, &dtime3, &atime3, &dtime1UTC, &atime1UTC, &dtime2UTC, &atime2UTC, &dtime3UTC, &atime3UTC, &deeplink1, &deeplink2, &deeplink3, &srciata, &mid1iata, &mid2iata, &dstiata, &airline1, &flightno1, &airline2, &flightno2, &airline3, &flightno3, &totalduration)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)

		}
		totalString := fmt.Sprint(total)

		flights = append(flights, Flight{

			A1:            src,
			A2:            mid1,
			A3:            mid2,
			A4:            dst,
			Price:         totalString,
			Duration1:     d1,
			Duration2:     d2,
			Duration3:     d3,
			Dtime1:        dtime1,
			Atime1:        atime1,
			Dtime2:        dtime2,
			Atime2:        atime2,
			Dtime3:        dtime3,
			Atime3:        atime3,
			Dtime1UTC:     dtime1UTC,
			Atime1UTC:     atime1UTC,
			Dtime2UTC:     dtime2UTC,
			Atime2UTC:     atime2UTC,
			Dtime3UTC:     dtime3UTC,
			Atime3UTC:     atime3UTC,
			Deeplink1:     deeplink1,
			Deeplink2:     deeplink2,
			Deeplink3:     deeplink3,
			Srciata:       srciata,
			Mid1iata:      mid1iata,
			Mid2iata:      mid2iata,
			Dstiata:       dstiata,
			Airline1:      airline1,
			Flightno1:     flightno1,
			Airline2:      airline2,
			Flightno2:     flightno2,
			Airline3:      airline3,
			Flightno3:     flightno3,
			TotalDuration: totalduration,
		})

	}

	postBody, _ := json.Marshal(flights)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Content-Type", "application/json")

	w.Write(postBody)

}

func truncateAllTables() {

	truncateTableQuery := "TRUNCATE TABLE T1, T2,PRICE_LIST_M1_T1,PRICE_LIST_M2_T1,PRICE_LIST_M1_T2,PRICE_LIST_M2_T2,PRICE_LIST_ONE,PRICE_LIST_three_nearDst,PRICE_LIST_three_nearSrc,PRICE_LIST_MID_FLIGHTS,price_list_three_todst,FINAL_PRICES, FINAL_PRICES_THREE,SRC_3,DST_3,FINAL_MERGED;DELETE FROM spatial_ref_sys  WHERE srid IN ( SELECT srid FROM spatial_ref_sys where srid!='4326' ORDER BY srid desc limit 4000);"
	db.Exec(truncateTableQuery)
	fmt.Println("Truncated tables")
}

func handleSoarch(w http.ResponseWriter, r *http.Request) {

	truncateAllTables()

	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	params := mux.Vars(r)

	apfrom := params["apfrom"]
	apto := params["apto"]
	date := params["dateOfFlight"]

	coordQuery := "SELECT latitude,longitude from airports where iata_code='" + apfrom + "';"
	rows, err := db.Query(coordQuery)
	defer rows.Close()

	var latitude, longitude, latdst, longdst string

	for rows.Next() {

		err = rows.Scan(&latitude, &longitude)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}

	}
	coordQuery = "SELECT latitude,longitude from airports where iata_code='" + apto + "';"
	rows, err = db.Query(coordQuery)
	defer rows.Close()

	for rows.Next() {

		err = rows.Scan(&latdst, &longdst)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return
		}

	}

	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)

	}

	nearestAirportsSrc := params["nearestAirportsSrc"]
	nearestAirportsDst := params["nearestAirportsDst"]
	srcAirports := populateAirportList(w, r, latitude, longitude, latdst, longdst, 0, nearestAirportsSrc, nearestAirportsDst)

	dstAirports := populateAirportList(w, r, latitude, longitude, latdst, longdst, 1, nearestAirportsSrc, nearestAirportsDst)

	flightsAPI(w, srcAirports, dstAirports, date)

}

func populateAirportList(w http.ResponseWriter, r *http.Request, latitude string, longitude string, latdst string, longdst string, t int, nearestAirportsSrc string, nearestAirportsDst string) []string {

	var queryStringSrc = ""
	var createTableQuery = ""
	var createTableQuery3 = ""
	if t == 0 {
		queryStringSrc = "SELECT Airport_Name, Country, Latitude, Longitude, IATA_CODE,ST_GeographyFromText(ST_AsText(coordinates))<-> ST_GeographyFromText('SRID=4326;POINT(" + longitude + " " + latitude + ")')::geography as dist FROM airports where PORT_TYPE='large_airport' ORDER by dist limit " + nearestAirportsSrc + ";"
		createTableQuery = " INSERT INTO T1  SELECT Airport_Name, Country, REGION, IATA_CODE,ST_AsText(coordinates), ST_DISTANCE(ST_GeographyFromText(ST_AsText(coordinates)),ST_GeographyFromText('POINT(" + longdst + " " + latdst + ")')) as dist FROM airports where ST_DISTANCE(ST_GeographyFromText(ST_AsText(coordinates)),ST_GeographyFromText('POINT(" + longdst + " " + latdst + ")')::geography) < ( (ST_Distance(ST_GeographyFromText('SRID=4326;POINT(" + longitude + " " + latitude + ")')::geography,ST_GeographyFromText('SRID=4326;POINT(" + longdst + " " + latdst + ")')::geography)))/1.5 order by dist;"
		createTableQuery3 = "INSERT INTO src_3  SELECT Airport_Name, Country, REGION, IATA_CODE,ST_AsText(coordinates), ST_DISTANCE(ST_GeographyFromText(ST_AsText(coordinates)),ST_GeographyFromText('POINT(" + longitude + " " + latitude + ")')) as dist FROM airports order by dist limit 200 ;"
	} else {
		queryStringSrc = "SELECT Airport_Name, Country, Latitude, Longitude, IATA_CODE,ST_GeographyFromText(ST_AsText(coordinates))<-> ST_GeographyFromText('SRID=4326;POINT(" + longdst + " " + latdst + ")')::geography as dist FROM airports where PORT_TYPE='large_airport' ORDER by dist limit " + nearestAirportsDst + ";"
		createTableQuery = "INSERT INTO T2  SELECT Airport_Name, Country, REGION, IATA_CODE,ST_AsText(coordinates), ST_DISTANCE(ST_GeographyFromText(ST_AsText(coordinates)),ST_GeographyFromText('POINT(" + longitude + " " + latitude + ")')) as dist FROM airports where ST_DISTANCE(ST_GeographyFromText(ST_AsText(coordinates)),ST_GeographyFromText('POINT(" + longitude + " " + latitude + ")')::geography) < ((ST_Distance(ST_GeographyFromText('SRID=4326;POINT(" + longitude + " " + latitude + ")')::geography,ST_GeographyFromText('SRID=4326;POINT(" + longdst + " " + latdst + ")')::geography)))/1.5 order by dist;"
		createTableQuery3 = "INSERT INTO dst_3  SELECT Airport_Name, Country, REGION, IATA_CODE,ST_AsText(coordinates), ST_DISTANCE(ST_GeographyFromText(ST_AsText(coordinates)),ST_GeographyFromText('POINT(" + longdst + " " + latdst + ")')) as dist FROM airports order by dist limit 200 ;"
	}
	rows, err := db.Query(queryStringSrc)
	db.Exec(createTableQuery)
	db.Exec(createTableQuery3)
	defer rows.Close()

	airportList := []string{}

	for rows.Next() {

		var airport_name string
		var country string
		var latitudeVal float32
		var longitudeVal float32
		var dist float32
		var iata string
		err = rows.Scan(&airport_name, &country, &latitudeVal, &longitudeVal, &iata, &dist)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)
			return airportList
		}

		airportList = append(airportList, iata)
	}

	return airportList
}

type Response struct {
	FlightData []flightData `json:"data"`
}

type flightData struct {
	FlyFrom     string `json:"flyFrom"`
	FlyTo       string `json:"flyTo"`
	CityFrom    string `json:"cityFrom"`
	CityTo      string `json:"cityTo"`
	CountryFrom string `json:"countryFrom"`
	CountryTo   string `json:"countryTo"`
	Price       int    `json:"price"`
	DeepLink    string `json:"deep_link"`
	DTime       int    `json:"dTime"`
	DTimeUTC    int    `json:"dTimeUTC"`
	ATime       int    `json:"aTime"`
	ATimeUTC    int    `json:"aTimeUTC"`
	FlyDuration int    `json:"fly_duration"`
	Route       []Route
}

type Route struct {
	Airline  string `json:"airline"`
	Flightno int    `json:"flight_no"`
}

func flightsAPI(w http.ResponseWriter, srcAirports []string, dstAirports []string, departDate string) {

	startTime := time.Now()
	var srcList string
	var dstList string

	for indsrc := range srcAirports {
		srcList = srcList + "," + srcAirports[indsrc]
	}
	srcList = srcList[1:]
	for inddst := range dstAirports {
		dstList = dstList + "," + dstAirports[inddst]
	}
	dstList = dstList[1:]

	go flightsOneStop(w, srcList, dstList, "ONE", false, true, departDate, 0, 1)

	go flightsTwoStops(w, srcList, dstList, "T1", departDate)
	go flightsTwoStops(w, srcList, dstList, "T2", departDate)

	flightsThreeStops(w, srcList, dstList, "", departDate)

	_, err := db.Exec("INSERT INTO final_prices  SELECT price_list_m2_t2.city_from as src,price_list_m2_t2.iata_from as src_iata, price_list_m2_t2.city_to as mid,price_list_m2_t2.iata_to as mid_iata, price_list_m1_t2.city_to as dst,price_list_m1_t2.iata_to as dst_iata, price_list_m2_t2.price+price_list_m1_t2.price as total, price_list_m2_t2.duration as d1, price_list_m2_t2.aTime as atime1,price_list_m2_t2.aTimeUTC as atimeutc1,price_list_m2_t2.dTime as dtime,price_list_m2_t2.dTimeUTC as dtimeutc1,price_list_m2_t2.deepLink as deeplink1, price_list_m1_t2.duration as d2,price_list_m1_t2.aTime as atime2,price_list_m1_t2.aTimeUTC as atimeutc2,price_list_m1_t2.dTime as dtime2,price_list_m1_t2.dTimeUTC as dtimeutc2,price_list_m1_t2.deepLink as deeplink2, price_list_m2_t2.airline as airline1, price_list_m2_t2.flightno as flightno1, price_list_m1_t2.airline as airline2, price_list_m1_t2.flightno as flightno2 from price_list_m1_t2, price_list_m2_t2 where price_list_m2_t2.iata_to=price_list_m1_t2.iata_from AND price_list_m2_t2.aTimeUTC<price_list_m1_t2.dTimeUTC order by total;")

	_, err = db.Exec("INSERT into final_prices  SELECT price_list_m1_t1.city_from as src,price_list_m1_t1.iata_from as src_iata, price_list_m1_t1.city_to as mid,price_list_m1_t1.iata_to as mid_iata, price_list_m2_t1.city_to as dst,price_list_m2_t1.iata_to as dst_iata, price_list_m2_t1.price+price_list_m1_t1.price as total, price_list_m1_t1.duration as d1, price_list_m1_t1.aTime as atime1,price_list_m1_t1.aTimeUTC as atimeutc1,price_list_m1_t1.dTime as dtime1,price_list_m1_t1.dTimeUTC as dtimeutc1,price_list_m1_t1.deepLink as deeplink1, price_list_m2_t1.duration as d2,price_list_m2_t1.aTime as atime2,price_list_m2_t1.aTimeUTC as atimeutc2,price_list_m2_t1.dTime as dtime2,price_list_m2_t1.dTimeUTC as dtimeutc2,price_list_m2_t1.deepLink as deeplink2,price_list_m1_t1.airline as airline1, price_list_m1_t1.flightno as flightno1, price_list_m2_t1.airline as airline2, price_list_m2_t1.flightno as flightno2 from price_list_m1_t1, price_list_m2_t1 where price_list_m1_t1.iata_to=price_list_m2_t1.iata_from AND price_list_m1_t1.aTimeUTC<price_list_m2_t1.dTimeUTC  order by total;")

	_, err = db.Exec("insert into final_merged(src,dst,d1,total,dtime1, dtime1utc, atime1, atime1utc, deeplink1, src_iata,dst_iata,airline1,flightno1 ) select p1.city_from as src, p1.city_to as dst, p1.duration as d1, p1.price as total,p1.dtime, p1.dtimeutc, p1.atime, p1.atimeutc, p1.deeplink, p1.iata_from, p1.iata_to , p1.airline, p1.flightno from price_list_one p1 order by total limit 30;insert into final_merged(src,mid1,dst,d1,d2,atime1,dtime1,atime1UTC,dtime1UTC,deeplink1,atime2,dtime2,atime2UTC,dtime2UTC,deeplink2, total, src_iata,mid1_iata,dst_iata,airline1,flightno1,airline2,flightno2) select src,mid,dst,d1,d2,atime1,dtime1,atime1UTC,dtime1UTC,deeplink1,atime2,dtime2,atime2UTC,dtime2UTC,deeplink2, total, src_iata, mid_iata, dst_iata,airline1,flightno1,airline2,flightno2 from final_prices order by total limit 30;insert into final_merged(src,mid1,mid2,dst,d1,d2,d3,atime1,dtime1,atime1UTC,dtime1UTC,deeplink1,atime2,dtime2,atime2UTC,dtime2UTC,deeplink2, atime3,dtime3,atime3UTC,dtime3UTC,deeplink3, total, src_iata,mid1_iata,mid2_iata,dst_iata,airline1,flightno1,airline2,flightno2,airline3,flightno3) select src,mid1,mid2,dst,d1,d2,d3,atime1,dtime1,atime1UTC,dtime1UTC,deeplink1,atime2,dtime2,atime2UTC,dtime2UTC,deeplink2, atime3,dtime3,atime3UTC,dtime3UTC,deeplink3,total, src_iata,mid1_iata,mid2_iata,dst_iata,airline1,flightno1,airline2,flightno2,airline3,flightno3 from final_prices_three order by total limit 30 ;UPDATE final_merged SET mid1 = 'invalid' WHERE mid1 IS NULL;UPDATE final_merged SET mid2 = 'invalid' WHERE mid2 IS NULL;UPDATE final_merged SET d1 = 'invalid' WHERE d1 IS NULL;UPDATE final_merged SET d2 = 'invalid',mid1_iata='invalid' WHERE d2 IS NULL;UPDATE final_merged SET d3 = 'invalid', mid2_iata='invalid' WHERE d3 IS NULL;UPDATE final_merged SET atime1 = -1 WHERE atime1 IS NULL;UPDATE final_merged SET dtime1 = -1 WHERE dtime1 IS NULL;UPDATE final_merged SET atime2 = -1 WHERE atime2 IS NULL;UPDATE final_merged SET dtime2 = -1 WHERE dtime2 IS NULL;UPDATE final_merged SET atime3 = -1 WHERE atime3 IS NULL;UPDATE final_merged SET dtime3 = -1 WHERE dtime3 IS NULL;UPDATE final_merged SET atime1UTC = -1 WHERE atime1UTC IS NULL;UPDATE final_merged SET dtime1UTC = -1 WHERE dtime1UTC IS NULL;UPDATE final_merged SET atime2UTC = -1 WHERE atime2UTC IS NULL;UPDATE final_merged SET dtime2UTC = -1 WHERE dtime2UTC IS NULL;UPDATE final_merged SET atime3UTC = -1 WHERE atime3UTC IS NULL;UPDATE final_merged SET dtime3UTC = -1 WHERE dtime3UTC IS NULL; UPDATE final_merged SET Deeplink1 = 'invalid' WHERE deeplink1 IS NULL; update final_merged set deeplink2='invalid' where deeplink2 IS NULL; update final_merged set deeplink3='invalid' where deeplink3 IS NULL; update final_merged set airline2='invalid' where airline2 IS NULL;update final_merged set flightno2='invalid' where flightno2 IS NULL;update final_merged set airline3='invalid' where airline3 IS NULL;update final_merged set flightno3='invalid' where flightno3 IS NULL;update final_merged set src = (select airport_name from airports where iata_code=src_iata);update final_merged set mid1 = (select airport_name from airports where iata_code=mid1_iata) where mid1!='invalid';update final_merged set mid2 = (select airport_name from airports where iata_code=mid2_iata) where mid2!='invalid';update final_merged set dst = (select airport_name from airports where iata_code=dst_iata) ")
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	durationRows, err := db.Query("SELECT atime1utc,dtime1utc,atime2utc,dtime2utc,atime3utc,dtime3utc FROM final_merged;")
	defer durationRows.Close()
	var a1, d1, a2, d2, a3, d3 int

	for durationRows.Next() {
		err = durationRows.Scan(&a1, &d1, &a2, &d2, &a3, &d3)

		var timeDiff int
		if a2 == -1 && a3 == -1 {
			timeDiff = int(time.Unix(int64(a1), 0).Sub(time.Unix(int64(d1), 0)).Minutes())
			timeDiffHour := timeDiff / 60
			timeDiffMins := timeDiff % 60

			duration := fmt.Sprint(timeDiffHour) + "hrs " + fmt.Sprint(timeDiffMins) + "mins"
			db.Exec("UPDATE final_merged SET time1 ='" + duration + "' where atime1utc='" + fmt.Sprint(a1) + "' and dtime1utc='" + fmt.Sprint(d1) + "';")

		} else if a3 == -1 {
			timeDiff = int(time.Unix(int64(a2), 0).Sub(time.Unix(int64(d1), 0)).Minutes())
			timeDiffHour := timeDiff / 60
			timeDiffMins := timeDiff % 60

			duration := fmt.Sprint(timeDiffHour) + "hrs " + fmt.Sprint(timeDiffMins) + "mins"
			db.Exec("UPDATE final_merged SET time1 ='" + duration + "' where atime2utc='" + fmt.Sprint(a2) + "' and dtime1utc='" + fmt.Sprint(d1) + "';")

		} else {
			timeDiff = int(time.Unix(int64(a3), 0).Sub(time.Unix(int64(d1), 0)).Minutes())
			timeDiffHour := timeDiff / 60
			timeDiffMins := timeDiff % 60

			duration := fmt.Sprint(timeDiffHour) + "hrs " + fmt.Sprint(timeDiffMins) + "mins"
			db.Exec("UPDATE final_merged SET time1 ='" + duration + "' where atime3utc='" + fmt.Sprint(a3) + "' and dtime1utc='" + fmt.Sprint(d1) + "';")
		}

	}

	writeJsonResponse(w)

	fmt.Printf("Total Time taken: %v", time.Since(startTime))

}

func flightsOneStop(w http.ResponseWriter, srcList string, dstList string, table string, oneForCity bool, display bool, dateString string, addDaysDA int, addDaysDB int) {

	fmt.Println(("FL ONE STOP "))
	str_oneForCity := fmt.Sprint(oneForCity)

	date, err := time.Parse("2006-01-02", dateString)

	dateFinal := date.AddDate(0, 0, addDaysDA).String()[0:10]
	departBeforeString := date.AddDate(0, 0, addDaysDB).String()[0:10]
	APIqueryString := "https://api.skypicker.com/flights?partner=soarchsoarch&depart_after=" + dateFinal + "&direct_flights=true&one_for_city=" + str_oneForCity + "&depart_before=" + departBeforeString + "&curr=USD&sort=price&limit=300&fly_from=airport:" + srcList + "&fly_to=airport:" + dstList
	fmt.Println(APIqueryString)
	response, err := http.Get(APIqueryString)

	responseData, err := ioutil.ReadAll(response.Body)
	var responseObject Response
	json.Unmarshal(responseData, &responseObject)

	for option := range responseObject.FlightData {

		dtime := time.Unix(int64(responseObject.FlightData[option].DTime), 0).UTC().Format(time.UnixDate)[4:16]
		atime := time.Unix(int64(responseObject.FlightData[option].ATime), 0).UTC().Format(time.UnixDate)[4:16]

		insertStmt := "insert into price_list_" + table + "(iata_from,iata_to,city_from,city_to,duration,price,dTime,dTimeUTC,aTime,aTimeUTC,deepLink,airline,flightno) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11, $12, $13);"
		_, err := db.Exec(insertStmt, responseObject.FlightData[option].FlyFrom, responseObject.FlightData[option].FlyTo, responseObject.FlightData[option].CityFrom, responseObject.FlightData[option].CityTo, responseObject.FlightData[option].FlyDuration, responseObject.FlightData[option].Price, dtime, responseObject.FlightData[option].DTimeUTC, atime, responseObject.FlightData[option].ATimeUTC, responseObject.FlightData[option].DeepLink, responseObject.FlightData[option].Route[0].Airline, responseObject.FlightData[option].Route[0].Flightno)

		if err != nil {
			fmt.Print(err.Error())
			return
		}

	}

	if err != nil {
		fmt.Print(err.Error())
		return
	}

}

func flightsTwoStops(w http.ResponseWriter, srcList string, dstList string, table string, dateString string) {

	fmt.Println("ERER")
	tableSizeRows, err := db.Query("SELECT COUNT(*) FROM " + table)
	defer tableSizeRows.Close()
	var tableSize int

	date, err := time.Parse("2006-01-02", dateString)
	departBeforeString := date.AddDate(0, 0, 2).String()[0:10]

	for tableSizeRows.Next() {

		err = tableSizeRows.Scan(&tableSize)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)

		}
	}

	for i := 0; i < tableSize; i += 500 {

		s := fmt.Sprint(i)
		queryString1 := "SELECT IATA_CODE, AIRPORT_NAME FROM " + table + " limit 500 offset " + s

		rows1, err := db.Query(queryString1)
		defer rows1.Close()

		airportList := []string{}

		for rows1.Next() {

			var airport_name string
			var iata string
			err = rows1.Scan(&iata, &airport_name)
			if err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(500), 500)

			}

			airportList = append(airportList, iata)

		}

		var midList1 string
		for indmid1 := range airportList {
			midList1 = midList1 + "," + airportList[indmid1]
		}
		if len(midList1) > 0 {
			midList1 = midList1[1:]
		}

		fmt.Println("API1")

		var APIqueryString string
		if table == "T1" {
			APIqueryString = "https://api.skypicker.com/flights?partner=soarchsoarch&depart_after=" + dateString + "&direct_flights=true&depart_before=" + departBeforeString + "&curr=USD&sort=price&asc=1&limit=200&fly_from=airport:" + srcList + "&fly_to=airport:" + midList1
		} else {
			APIqueryString = "https://api.skypicker.com/flights?partner=soarchsoarch&depart_after=" + dateString + "&direct_flights=true&depart_before=" + departBeforeString + "&curr=USD&sort=price&asc=1&limit=200&fly_from=airport:" + midList1 + "&fly_to=airport:" + dstList
		}
		fmt.Println(APIqueryString)
		response, err := http.Get(APIqueryString)

		responseData, err := ioutil.ReadAll(response.Body)
		var responseObject Response
		json.Unmarshal(responseData, &responseObject)

		for option := range responseObject.FlightData {

			dtime := time.Unix(int64(responseObject.FlightData[option].DTime), 0).UTC().Format(time.UnixDate)[4:16]
			atime := time.Unix(int64(responseObject.FlightData[option].ATime), 0).UTC().Format(time.UnixDate)[4:16]

			insertStmt := "insert into price_list_m1_" + table + "(iata_from,iata_to,city_from,city_to,duration,price,dTime,dTimeUTC,aTime,aTimeUTC,deepLink,airline,flightno) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11, $12, $13);"
			_, err := db.Exec(insertStmt, responseObject.FlightData[option].FlyFrom, responseObject.FlightData[option].FlyTo, responseObject.FlightData[option].CityFrom, responseObject.FlightData[option].CityTo, responseObject.FlightData[option].FlyDuration, responseObject.FlightData[option].Price, dtime, responseObject.FlightData[option].DTimeUTC, atime, responseObject.FlightData[option].ATimeUTC, responseObject.FlightData[option].DeepLink, responseObject.FlightData[option].Route[0].Airline, responseObject.FlightData[option].Route[0].Flightno)

			if err != nil {
				fmt.Print(err.Error())

			}

		}

		midList1 = ""
		var queryStringMidList string
		if table == "T1" {
			queryStringMidList = "select distinct iata_to from price_list_m1_t1"
		} else {
			queryStringMidList = "select distinct iata_from from price_list_m1_t2"
		}

		rows2, err := db.Query(queryStringMidList)
		defer rows2.Close()

		airportList = []string{}

		for rows2.Next() {

			var iata string
			err = rows2.Scan(&iata)
			if err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(500), 500)

			}

			airportList = append(airportList, iata)

		}

		for indmid1 := range airportList {
			midList1 = midList1 + "," + airportList[indmid1]
		}

		if len(midList1) > 0 {
			midList1 = midList1[1:]
		}

		departBeforeString = date.AddDate(0, 0, 3).String()[0:10]

		fmt.Println("API2")

		if table == "T1" {
			APIqueryString = "https://api.skypicker.com/flights?partner=soarchsoarch&depart_after=" + dateString + "&direct_flights=true&depart_before=" + departBeforeString + "&curr=USD&sort=price&asc=1&limit=200&fly_from=airport:" + midList1 + "&fly_to=airport:" + dstList
		} else {
			APIqueryString = "https://api.skypicker.com/flights?partner=soarchsoarch&depart_after=" + dateString + "&direct_flights=true&depart_before=" + departBeforeString + "&curr=USD&sort=price&asc=1&limit=200&fly_from=airport:" + srcList + "&fly_to=airport:" + midList1
		}
		fmt.Println(APIqueryString)
		response, err = http.Get(APIqueryString)

		responseData, err = ioutil.ReadAll(response.Body)
		json.Unmarshal(responseData, &responseObject)

		for option := range responseObject.FlightData {

			dtime := time.Unix(int64(responseObject.FlightData[option].DTime), 0).UTC().Format(time.UnixDate)[4:16]

			atime := time.Unix(int64(responseObject.FlightData[option].ATime), 0).UTC().Format(time.UnixDate)[4:16]

			insertStmt := "insert into price_list_m2_" + table + "(iata_from,iata_to,city_from,city_to,duration,price,dTime,dTimeUTC,aTime,aTimeUTC,deepLink,airline,flightno) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11, $12, $13);"
			_, err := db.Exec(insertStmt, responseObject.FlightData[option].FlyFrom, responseObject.FlightData[option].FlyTo, responseObject.FlightData[option].CityFrom, responseObject.FlightData[option].CityTo, responseObject.FlightData[option].FlyDuration, responseObject.FlightData[option].Price, dtime, responseObject.FlightData[option].DTimeUTC, atime, responseObject.FlightData[option].ATimeUTC, responseObject.FlightData[option].DeepLink, responseObject.FlightData[option].Route[0].Airline, responseObject.FlightData[option].Route[0].Flightno)

			if err != nil {
				fmt.Print(err.Error())
				return
			}

		}

	}

}

func flightsThreeStops(w http.ResponseWriter, srcList string, dstList string, table string, dateString string) {

	rows, err := db.Query("select iata_code from src_3;")
	var toSrc3 string
	for rows.Next() {

		var iata string
		err = rows.Scan(&iata)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)

		}

		toSrc3 = toSrc3 + "," + iata
	}
	if len(toSrc3) > 5 {
		toSrc3 = toSrc3[5:]
	}

	flightsOneStop(w, srcList, toSrc3, "three_nearSrc", false, false, dateString, 0, 1)

	rows1_3, err := db.Query("SELECT distinct IATA_to from price_list_three_nearSrc;")
	toSrc3 = ""
	for rows1_3.Next() {

		var iata string
		err = rows1_3.Scan(&iata)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)

		}

		toSrc3 = toSrc3 + "," + iata
	}
	if len(toSrc3) > 1 {
		toSrc3 = toSrc3[1:]
	}

	fromDst3 := ""
	rows2_3, err := db.Query("SELECT distinct IATA_CODE from dst_3;")

	for rows2_3.Next() {

		var iata string
		err = rows2_3.Scan(&iata)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)

		}

		fromDst3 = fromDst3 + "," + iata
	}

	if len(fromDst3) > 1 {
		fromDst3 = fromDst3[1:]
	}

	fmt.Println("Mid Flights:-")
	flightsOneStop(w, toSrc3, fromDst3, "mid_flights", false, false, dateString, 0, 2)
	//var fromDst3 string;
	for rows.Next() {

		var iata string
		err = rows.Scan(&iata)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)

		}

		fromDst3 = fromDst3 + "," + iata
	}

	tableSizeRows, err := db.Query("SELECT COUNT(*) FROM price_list_mid_flights;")
	defer tableSizeRows.Close()
	var tableSize int
	for tableSizeRows.Next() {

		err = tableSizeRows.Scan(&tableSize)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(500), 500)

		}
	}

	for i := 0; i < tableSize; i += 450 {
		s := fmt.Sprint(i)
		rows, err = db.Query("select distinct iata_to from price_list_mid_flights limit 450 offset" + s)

		flightsOneStop(w, fromDst3, dstList, "three_nearDst", false, false, dateString, 1, 3)
	}

	_, err = db.Exec(" INSERT INTO  final_prices_three(src, src_iata, mid1, mid1_iata, mid2, mid2_iata, dst, dst_iata, total,d1, d2, d3, atime1,  dtime1, atime1utc,  dtime1utc, deeplink1,airline1, flightno1,atime2,  dtime2, atime2utc,  dtime2utc, deeplink2,airline2, flightno2, atime3,  dtime3, atime3utc,  dtime3utc, deeplink3,airline3, flightno3) SELECT price_list_three_nearSrc.city_from as src,price_list_three_nearSrc.iata_from as src_iata, price_list_three_nearSrc.city_to as mid1,price_list_three_nearSrc.iata_to as mid1_iata,price_list_three_nearDst.city_from as mid2,price_list_three_nearDst.iata_from as mid2_iata, price_list_three_nearDst.city_to as dst,price_list_three_nearDst.iata_to as dst_iata, price_list_three_nearSrc.price+price_list_mid_flights.price+price_list_three_nearDst.price as total, price_list_three_nearSrc.duration as d1, price_list_mid_flights.duration as d2, price_list_three_nearDst.duration as d3,price_list_three_nearSrc.aTime as atime1,price_list_three_nearSrc.dTime as dtime1,price_list_three_nearSrc.aTimeUTC as atime1UTC,price_list_three_nearSrc.dTimeUTC as dtime1UTC,price_list_three_nearSrc.deepLink as deeplink1, price_list_three_nearSrc.airline as airline1,price_list_three_nearSrc.flightno as flightno1, price_list_mid_flights.aTime as atime2,price_list_mid_flights.dTime as dtime2,price_list_mid_flights.aTimeUTC as atime2UTC,price_list_mid_flights.dTimeUTC as dtime2UTC,price_list_mid_flights.deepLink as deeplink2, price_list_mid_flights.airline as airline2, price_list_mid_flights.flightno as flightno2, price_list_three_nearDst.aTime as atime3,price_list_three_nearDst.dTime as dtime3,price_list_three_nearDst.aTimeUTC as atime3UTC,price_list_three_nearDst.dTimeUTC as dtime3UTC,price_list_three_nearDst.deepLink as deeplink3, price_list_three_nearDst.airline as airline3, price_list_three_nearDst.flightno as flightno3 from price_list_three_nearSrc, price_list_mid_flights, price_list_three_nearDst where price_list_three_nearSrc.iata_to=price_list_mid_flights.iata_from  AND price_list_mid_flights.iata_to = price_list_three_nearDst.iata_from AND price_list_three_nearSrc.aTimeUTC<price_list_mid_flights.dTimeUTC AND price_list_mid_flights.aTimeUTC<price_list_three_nearDst.dTimeUTC order by total;")

}

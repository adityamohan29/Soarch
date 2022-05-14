function flightSearchAPI(fromFlight, toFlight, date) {

    if (document.body.children.length > 0)
        document.body.removeChild(document.body.lastChild);

    fromFlight = fromFlight.toUpperCase();
    toFlight = toFlight.toUpperCase();
    var constructURL = "http://localhost:8090/find/" + fromFlight + "&" + toFlight + "&1&1&" + date;

    callAPI(constructURL);

}

async function callAPI(url) {




    var flightDivLoad = document.createElement("p");
    flightDivLoad.setAttribute("id", "flightDivLoad")
    flightDivLoad.innerHTML = "<div class='lds-ripple'><div></div><div></div></div>"
    document.body.appendChild(flightDivLoad);
    var response = '';
    var resp = '';
    try {
        response = await fetch(url);
        resp = await response.json();
    } catch (e) {
        console.log('error:', e.message);
    }


    console.log(resp);
    readJson(resp);


}


function readJson(resp) {

    if (document.body.children.length > 0)
        document.body.removeChild(document.body.lastChild);


    var flightDivParent = document.createElement("p");
    flightDivParent.setAttribute("id", "flightDiv")

    var flightDivHead1 = document.createElement("p");
    flightDivHead1.setAttribute("id", "flightDivHead1");
    flightDivHead1.innerHTML = " &nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; &nbsp;&nbsp;&nbsp;&nbsp; Route ";
    flightDivParent.appendChild(flightDivHead1);
    var flightDivHead2 = document.createElement("p");
    flightDivHead2.setAttribute("id", "flightDivHead2");
    flightDivHead2.innerHTML = " Duration ";
    flightDivParent.appendChild(flightDivHead2);
    var flightDivHead3 = document.createElement("p");
    flightDivHead3.setAttribute("id", "flightDivHead3");
    flightDivHead3.innerHTML = "  Price";
    flightDivParent.appendChild(flightDivHead3);


    for (var i = 0; i < resp.length; i++) {

        var flightDivChild = document.createElement("p");
        flightDivChild.setAttribute("id", "flightDivChildID")
        var flightDivGrandChild = document.createElement("p");
        flightDivGrandChild.setAttribute("id", "flightDivGrandChildID")
        //const flightDivChild = document.createElement("flightDiv")
        flightDivGrandChild.innerHTML = resp[i]["Srciata"] + "&nbsp;&nbsp;&nbsp;&nbsp;  <hr id='line1'></hr>  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;";
        flightDivChild.appendChild(flightDivGrandChild)

        var flightDivGrandChild_time = document.createElement("p");
        flightDivGrandChild_time.setAttribute("id", "flightDivGrandChildID_time")
        flightDivGrandChild_time.innerHTML = resp[i]["Dtime1"];
        flightDivChild.appendChild(flightDivGrandChild_time)

        var flightDivGrandChild_city = document.createElement("p");
        flightDivGrandChild_city.setAttribute("id", "flightDivGrandChildID_city")
        flightDivGrandChild_city.innerHTML = resp[i]["A1"];
        flightDivGrandChild.appendChild(flightDivGrandChild_city)

        if (resp[i]["Mid1iata"] != "invalid") {
            var flightDivGrandChild1 = document.createElement("p");
            flightDivGrandChild1.setAttribute("id", "flightDivGrandChildID1")
            flightDivGrandChild1.innerHTML = resp[i]["Mid1iata"] + "&nbsp;&nbsp;&nbsp;&nbsp; <hr id='line2'></hr> &nbsp;&nbsp;&nbsp    &nbsp;&nbsp;&nbsp;&nbsp;&nbsp &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;";



            var flightDivGrandChild1_city = document.createElement("p");
            flightDivGrandChild1_city.setAttribute("id", "flightDivGrandChildID1_city")
            flightDivGrandChild1_city.innerHTML = resp[i]["A2"];
            flightDivGrandChild1.appendChild(flightDivGrandChild1_city);

            flightDivChild.appendChild(flightDivGrandChild1);

            var flightDivGrandChild1_airline = document.createElement("p");
            flightDivGrandChild1_airline.setAttribute("id", "flightDivGrandChildID1_airline")
            flightDivGrandChild1_airline.innerHTML = "<a class='dl'  href=" + resp[i]["Deeplink1"] + " target='_blank'> " + resp[i]["Airline1"] + " " + resp[i]["Flightno1"] + "</a>";
            flightDivChild.appendChild(flightDivGrandChild1_airline)

            var flightDivGrandChild1_time = document.createElement("p");
            flightDivGrandChild1_time.setAttribute("id", "flightDivGrandChildID1_time")
            flightDivGrandChild1_time.innerHTML = resp[i]["Atime1"] + "&nbsp; &nbsp;  " + resp[i]["Dtime2"];
            flightDivChild.appendChild(flightDivGrandChild1_time)



            flightDivChild.appendChild(flightDivGrandChild1)


        }
        if (resp[i]["Mid2iata"] != "invalid") {
            var flightDivGrandChild2 = document.createElement("p");
            flightDivGrandChild2.setAttribute("id", "flightDivGrandChildID2")
            flightDivGrandChild2.innerHTML = resp[i]["Mid2iata"] + "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;  <hr id='line3'></hr>   &nbsp;&nbsp;&nbsp  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;";
            flightDivChild.appendChild(flightDivGrandChild2)


            var flightDivGrandChild2_city = document.createElement("p");
            flightDivGrandChild2_city.setAttribute("id", "flightDivGrandChildID2_city")
            flightDivGrandChild2_city.innerHTML = resp[i]["A3"];
            flightDivGrandChild2.appendChild(flightDivGrandChild2_city)

            flightDivChild.appendChild(flightDivGrandChild2);

            var flightDivGrandChild2_airline = document.createElement("p");
            flightDivGrandChild2_airline.setAttribute("id", "flightDivGrandChildID2_airline")
            flightDivGrandChild2_airline.innerHTML = "<a class='dl' href=" + resp[i]["Deeplink2"] + " target='_blank'>" + resp[i]["Airline2"] + " " + resp[i]["Flightno2"] + "</a>";
            flightDivChild.appendChild(flightDivGrandChild2_airline)

            var flightDivGrandChild2_time = document.createElement("p");
            flightDivGrandChild2_time.setAttribute("id", "flightDivGrandChildID2_time")
            flightDivGrandChild2_time.innerHTML = resp[i]["Atime2"] + "&nbsp;&nbsp;   " + resp[i]["Dtime3"];
            flightDivChild.appendChild(flightDivGrandChild2_time)


        }

        var flightDivGrandChild3 = document.createElement("p");

        if (resp[i]["Mid1iata"] == "invalid")
            flightDivGrandChild3.setAttribute("id", "flightDivGrandChildID3_1")
        else if (resp[i]["Mid2iata"] == "invalid")
            flightDivGrandChild3.setAttribute("id", "flightDivGrandChildID3_2")
        else
            flightDivGrandChild3.setAttribute("id", "flightDivGrandChildID3_3")

        flightDivGrandChild3.innerHTML = resp[i]["Dstiata"];
        flightDivChild.appendChild(flightDivGrandChild3)


        var flightDivGrandChild3_city = document.createElement("p");
        flightDivGrandChild3_city.setAttribute("id", "flightDivGrandChildID3_city")
        flightDivGrandChild3_city.innerHTML = resp[i]["A4"];
        flightDivGrandChild3.appendChild(flightDivGrandChild3_city)

        flightDivChild.appendChild(flightDivGrandChild3);

        var flightDivGrandChild3_airline = document.createElement("p");


        if (resp[i]["Mid1iata"] == "invalid") {
            flightDivGrandChild3_airline.setAttribute("id", "flightDivGrandChildID3_airline1")
            flightDivGrandChild3_airline.innerHTML = "<a class='dl' href=" + resp[i]["Deeplink1"] + " target='_blank'>" + resp[i]["Airline1"] + " " + resp[i]["Flightno1"] + "</a>";
        } else if (resp[i]["Mid2iata"] == "invalid") {
            flightDivGrandChild3_airline.setAttribute("id", "flightDivGrandChildID3_airline2")
            flightDivGrandChild3_airline.innerHTML = "<a class='dl' href=" + resp[i]["Deeplink2"] + " target='_blank'>" + resp[i]["Airline2"] + " " + resp[i]["Flightno2"] + "</a>";
        } else {
            flightDivGrandChild3_airline.setAttribute("id", "flightDivGrandChildID3_airline3")
            flightDivGrandChild3_airline.innerHTML = "<a class='dl' href=" + resp[i]["Deeplink3"] + " target='_blank'>" + resp[i]["Airline3"] + " " + resp[i]["Flightno3"] + "</a>";
        }
        flightDivChild.appendChild(flightDivGrandChild3_airline)

        var flightDivGrandChild3_time = document.createElement("p");

        if (resp[i]["Mid1iata"] == "invalid") {
            flightDivGrandChild3_time.setAttribute("id", "flightDivGrandChildID3_time1")
            flightDivGrandChild3_time.innerHTML = resp[i]["Atime1"];
        } else if (resp[i]["Mid2iata"] == "invalid") {
            flightDivGrandChild3_time.setAttribute("id", "flightDivGrandChildID3_time2")
            flightDivGrandChild3_time.innerHTML = resp[i]["Atime2"];
        } else {
            flightDivGrandChild3_time.setAttribute("id", "flightDivGrandChildID3_time3")
            flightDivGrandChild3_time.innerHTML = resp[i]["Atime3"];
        }
        flightDivChild.appendChild(flightDivGrandChild3_time)


        var flightDivGrandChildPrice = document.createElement("p");
        flightDivGrandChildPrice.setAttribute("id", "flightDivGrandChildPriceID")
        flightDivGrandChildPrice.innerHTML = "&nbsp; " + resp[i]["Price"];

        var flightDivGrandChildCurr = document.createElement("p");
        flightDivGrandChildCurr.setAttribute("id", "flightDivGrandChildCurrID")
        flightDivGrandChildCurr.innerHTML = "&nbsp;&nbsp; USD";
        flightDivGrandChildPrice.appendChild(flightDivGrandChildCurr)

        var flightDivGrandChildDuration = document.createElement("p");
        flightDivGrandChildDuration.setAttribute("id", "flightDivGrandChildDurationID")
        flightDivGrandChildDuration.innerHTML = "&nbsp;&nbsp;<br><br><br><br>&nbsp;&nbsp;" + resp[i]["TotalDuration"];


        flightDivChild.appendChild(flightDivGrandChildPrice)
        flightDivChild.appendChild(flightDivGrandChildDuration)
        flightDivParent.appendChild(flightDivChild);

    }
    document.body.appendChild(flightDivParent);

}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	gtfs "github.com/artonge/go-gtfs"
	"github.com/rs/cors"
)

var gtfsStops []gtfs.Stop
var gtfsRoutes []gtfs.Route
var routes []route
var gtfsStoptimes []gtfs.StopTime
var gtfsTrips []gtfs.Trip
var gtfsCalendars []gtfs.Calendar
var gtfsCalendarDates []gtfs.CalendarDate
var gtfsTransfers []gtfs.Transfer

func writeJson(w http.ResponseWriter, data interface{}) {
	json, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func getTrips(w http.ResponseWriter, r *http.Request) {
	writeJson(w, gtfsTrips)
}

func getTransfers(w http.ResponseWriter, r *http.Request) {
	writeJson(w, gtfsTransfers)
}

func getStops(w http.ResponseWriter, r *http.Request) {
	writeJson(w, gtfsStops)
}

type trip struct {
	ID        string `json:id`
	Direction string `json:direction`
	Shape     string `json:shape`
}

type shape struct {
	ID        string `json:id`
	Longitude string `json:longitude`
	Latitude  string `json:latitude`
}

type route struct {
	ID    string `json:id`
	Name  string `json:name`
	Trips []trip `json:trips`
}

func getRoutes(w http.ResponseWriter, r *http.Request) {
	writeJson(w, routes)
}

type info struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func getApplicationInfo(w http.ResponseWriter, r *http.Request) {
	info := info{"regionalmap", "0.1.0"}
	writeJson(w, info)
}

func initWebServer() {
	fmt.Println("Initializing Webserver")
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/stops", getStops)
	mux.HandleFunc("/v1/routes", getRoutes)
	mux.HandleFunc("/v1/transfers", getTransfers)
	mux.HandleFunc("/v1/trips", getTrips)
	mux.HandleFunc("/info", getApplicationInfo)
	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":8080", handler)
}

func getTripsForRoute(routeID string) []trip {
	var tripsForRoute []trip
	for _, triped := range gtfsTrips {
		if triped.RouteID == routeID {
			tripsForRoute = append(tripsForRoute, trip{triped.ID, triped.DirectionID, triped.ShapeID})
		}
	}

	return tripsForRoute
}

func main() {
	fmt.Println("Loading GTFS data")
	gs, err := gtfs.LoadSplitted("gtfs", nil)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}

	fmt.Printf("Loaded %d GTFS datasets\n", len(gs))

	loadedGtfs := gs[0]
	fmt.Println("Loading GTFS data for " + loadedGtfs.Agency.Name)
	gtfsRoutes = loadedGtfs.Routes
	gtfsStops = loadedGtfs.Stops
	gtfsStoptimes = loadedGtfs.StopsTimes
	gtfsTrips = loadedGtfs.Trips
	gtfsCalendars = loadedGtfs.Calendars
	gtfsCalendarDates = loadedGtfs.CalendarDates
	gtfsTransfers = loadedGtfs.Transfers

	for _, routed := range gtfsRoutes {
		tripsForRoute := getTripsForRoute(routed.ID)
		routes = append(routes, route{routed.ID, routed.ShortName, tripsForRoute})
	}

	initWebServer()
}

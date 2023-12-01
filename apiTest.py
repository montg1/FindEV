import requests
import json

here_api_key = "g-gI_EzetmmNF8WDJNLXRu-hYIqtoRj8OiGtGZADXeM"
la_long = []  # List to store latitude and longitude pairs

def get_ev_charge_points(lat, lng, limit="1"):
    base_url = "https://discover.search.hereapi.com/v1/discover"
    params = {
        "apiKey": here_api_key,
        "q": "EV Charging Station",
        "at": f"{lat},{lng}",
        "limit": limit
    }

    response = requests.get(base_url, params=params)

    if response.status_code == 200:
        ev_charge_points = response.json()
        return ev_charge_points
    else:
        print(f"API request failed with status code: {response.status_code}")
        return None

def get_route(start, end, charging_station):
    base_url = "https://route.ls.hereapi.com/routing/7.2/calculateroute.json"
    params = {
        "apiKey": here_api_key,
        "waypoint0": f"{start[0]},{start[1]}",
        "waypoint1": f"{charging_station['position']['lat']},{charging_station['position']['lng']}",
        "waypoint2": f"{end[0]},{end[1]}",
        "mode": "fastest;car;traffic:disabled"
    }

    response = requests.get(base_url, params=params)

    if response.status_code == 200:
        route_data = response.json()
        return route_data
    else:
        print(f"Route calculation failed with status code: {response.status_code}")
        return None

def read_csv():
    with open("Wales.csv", "r") as file:
        for line in file:
            item = line.strip().split(",")
            list_item = item[9:11]
            la_long.append(list_item)

def main():
    read_csv()

    # Input for maximum capacitor, current location, and goal location
    max_capacitor = float(input("Enter the maximum capacitor: "))
    current_location = input("Enter the current location (latitude,longitude): ").split(',')
    goal_location = input("Enter the goal location (latitude,longitude): ").split(',')

    best_route = None
    best_charging_station = None

    for lat, lng in la_long:
        ev_charge_points = get_ev_charge_points(lat, lng)

        if ev_charge_points:
            print("EV Charging Stations:")
            for station in ev_charge_points["items"]:
                print(f"Name: {station['title']}, Latitude: {station['position']['lat']}, Longitude: {station['position']['lng']}")
                if "connectors" in station and "connector" in station["connectors"]:
                    for connector in station["connectors"]["connector"]:
                        print(f"  Supplier: {connector['supplierName']}, Charge Capacity: {connector['chargeCapacity']}")

                    # Check if the charging station can cover the trip
                    if float(connector['chargeCapacity']) >= max_capacitor:
                        route_data = get_route(current_location, goal_location, station)
                        if route_data:
                            if best_route is None or route_data['response']['route'][0]['summary']['travelTime'] < best_route['response']['route'][0]['summary']['travelTime']:
                                best_route = route_data
                                best_charging_station = station

            print("-------------------------------")

    if best_route is not None:
        print("Best Route:")
        for leg in best_route['response']['route'][0]['leg']:
            for maneuver in leg['maneuver']:
                print(f"  {maneuver['instruction']}")
        print(f"\nBest Charging Station: {best_charging_station['title']}, Latitude: {best_charging_station['position']['lat']}, Longitude: {best_charging_station['position']['lng']}")
    else:
        print("No suitable charging station found.")

if __name__ == "__main__":
    main()

import csv

charging_stations = []  # List to store charging station information

def read_csv():
    with open("locations.csv", newline='') as csvfile:
        reader = csv.DictReader(csvfile)
        for row in reader:
            charging_stations.append({
                'title': row['Location of Charging Point'],
                'position': {
                    'lat': float(row['Latitude']),
                    'lng': float(row['Longitude'])
                },
                'connectors': {
                    'connector': [
                        {
                            'supplierName': '',  # You can add supplier information if available
                            'chargeCapacity': 0  # You can add charge capacity information if available
                        }
                    ]
                }
            })

def main():
    read_csv()

    # Input for maximum capacitor, current location, and goal location
    max_capacitor = float(input("Enter the maximum capacitor: "))
    current_location = input("Enter the current location (latitude,longitude): ").split(',')
    goal_location = input("Enter the goal location (latitude,longitude): ").split(',')

    best_route = None
    best_charging_station = None

    for station in charging_stations:
        #print(f"Name: {station['title']}, Latitude: {station['position']['lat']}, Longitude: {station['position']['lng']}")

        # Check if the charging station can cover the trip
        if float(station['connectors']['connector'][0]['chargeCapacity']) >= max_capacitor:
            route_data = get_route(current_location, goal_location, station)
            if route_data:
                if best_route is None or route_data['response']['route'][0]['summary']['travelTime'] < best_route['response']['route'][0]['summary']['travelTime']:
                    best_route = route_data
                    best_charging_station = station

        #print("-------------------------------")

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

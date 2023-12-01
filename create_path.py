import pandas as pd
import folium
from geopy.distance import geodesic

# Load CSV data into a Pandas DataFrame
df = pd.read_csv('Wales.csv', header=None, names=['ID', 'Location', 'City', 'Unused', 'Postal Code', 'County', 'Unused2', 'Country', 'Latitude', 'Longitude'])

# Extract name, latitude, and longitude columns
coordinates_df = df[['Location', 'Latitude', 'Longitude']]

# Save coordinates with names as a new CSV file
coordinates_df.to_csv('coordinates.csv', index=False)

# Create a Folium map centered around Wales
map_wales = folium.Map(location=[52.1307, -3.7837], zoom_start=8)

# Add markers for each location
for index, row in df.iterrows():
    folium.Marker([row['Latitude'], row['Longitude']], popup=row['Location']).add_to(map_wales)

# Initialize variables for path calculation
current_path = []
total_distance = 0

# Helper function to calculate distance between two coordinates
def calculate_distance(coord1, coord2):
    return geodesic(coord1, coord2).kilometers

# Add a path using the latitude and longitude coordinates, limiting to 300 km
for index, row in df.iterrows():
    if current_path:
        # Calculate distance between the last point in the path and the current point
        distance = calculate_distance((current_path[-1]['Latitude'], current_path[-1]['Longitude']),
                                       (row['Latitude'], row['Longitude']))

        if total_distance + distance <= 300:
            # Add the current point to the path
            current_path.append({'Location': row['Location'], 'Latitude': row['Latitude'], 'Longitude': row['Longitude']})
            total_distance += distance
        else:
            # If adding the current point exceeds the limit, draw the current path and start a new one
            folium.PolyLine(locations=[[p['Latitude'], p['Longitude']] for p in current_path], color='blue').add_to(map_wales)
            current_path = [{'Location': row['Location'], 'Latitude': row['Latitude'], 'Longitude': row['Longitude']}]
            total_distance = 0
    else:
        # If the path is empty, add the current point
        current_path.append({'Location': row['Location'], 'Latitude': row['Latitude'], 'Longitude': row['Longitude']})

# Draw the last path
folium.PolyLine(locations=[[p['Latitude'], p['Longitude']] for p in current_path], color='blue').add_to(map_wales)

# Save the map as an HTML file
map_wales.save('path_map.html')

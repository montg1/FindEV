import pandas as pd
import folium
from geopy.distance import geodesic
import networkx as nx
from folium import plugins

# Load CSV data into a Pandas DataFrame
df = pd.read_csv('Wales.csv', header=None, names=['ID', 'Location', 'City', 'Unused', 'Postal Code', 'County', 'Unused2', 'Country', 'Latitude', 'Longitude'])

# Extract name, latitude, and longitude columns
coordinates_df = df[['Location', 'Latitude', 'Longitude']]

# Save coordinates with names as a new CSV file
coordinates_df.to_csv('coordinates.csv', index=False)

# Create a Folium map centered around Wales
map_wales = folium.Map(location=[52.1307, -3.7837], zoom_start=8)

# Add markers for each location with latitude and longitude
for index, row in df.iterrows():
    if not pd.isna(row['Latitude']) and not pd.isna(row['Longitude']):
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

# Now, let's implement Dijkstra's algorithm to find the shortest path between current and goal coordinates
def dijkstra(graph, start, goal):
    shortest_path = nx.shortest_path(graph, source=start, target=goal, weight='weight')
    return shortest_path

# Create a graph using NetworkX
G = nx.Graph()

# Add nodes and edges to the graph based on the coordinates
for index, row in df.iterrows():
    G.add_node(row['Location'], pos=(row['Latitude'], row['Longitude']))
    if index > 0:
        distance = calculate_distance((df.at[index-1, 'Latitude'], df.at[index-1, 'Longitude']),
                                       (row['Latitude'], row['Longitude']))
        G.add_edge(df.at[index-1, 'Location'], row['Location'], weight=distance)

# User input for current and goal locations
current_location = input("Enter the current location (latitude,longitude): ").split(',')
goal_location = input("Enter the goal location (latitude,longitude): ").split(',')

# Convert input to float
current_coordinates = (float(current_location[0]), float(current_location[1]))
goal_coordinates = (float(goal_location[0]), float(goal_location[1]))

# Find the closest nodes to the current and goal coordinates
start_node = min(G.nodes, key=lambda n: geodesic(G.nodes[n]['pos'], current_coordinates).kilometers)
goal_node = min(G.nodes, key=lambda n: geodesic(G.nodes[n]['pos'], goal_coordinates).kilometers)

# Find the shortest path using Dijkstra's algorithm
shortest_path = dijkstra(G, start_node, goal_node)

# Create a new Folium map for the path from current to goal
map_path = folium.Map(location=current_coordinates, zoom_start=12)

# Add markers for the current and goal coordinates
folium.Marker(current_coordinates, popup='Current Location').add_to(map_path)
folium.Marker(goal_coordinates, popup='Goal Location').add_to(map_path)

# Add the shortest path to the map
path_coordinates = [G.nodes[node]['pos'] for node in shortest_path]
plugins.AntPath(locations=path_coordinates, color='red', delay=1000).add_to(map_path)

# Save the map as an HTML file
map_path.save('shortest_path_map.html')

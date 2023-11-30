import pandas as pd
import folium
import networkx as nx
import heapq
from bs4 import BeautifulSoup
import re
import sys

def dijkstra(graph, start, goal):
    queue = [(0, start, [])]
    visited = set()

    while queue:
        (cost, node, path) = heapq.heappop(queue)
        if node not in visited:
            visited.add(node)
            path = path + [node]

            if node == goal:
                return path

            for next_node in graph[node]:
                if next_node not in visited:
                    heapq.heappush(queue, (cost + graph[node][next_node]['weight'], next_node, path))

    return None

# Load CSV data into a Pandas DataFrame
df = pd.read_csv('Wales.csv', header=None, names=['ID', 'Location', 'City', 'Unused', 'Postal Code', 'County', 'Unused2', 'Country', 'Latitude', 'Longitude'])

# Create a NetworkX graph using latitude and longitude as nodes and distance as weights
G = nx.Graph()

# Load the path_map.html and extract the edges
path_map = folium.Map(location=[52.1307, -3.7837], zoom_start=8)
path_map.save('path_map.html')  # Make sure path_map.html is saved before running this script

with open('path_map.html', 'r') as f:
    path_map_html = f.read()

# Print the HTML content for inspection
print(path_map_html)

soup = BeautifulSoup(path_map_html, 'html.parser')

# Find the JavaScript code that adds the polyline
script_pattern = re.compile(r'addLatLng\((.*?), (.*?)\);')
matches = script_pattern.findall(path_map_html)

if matches:
    polyline = matches
else:
    print("Error: Unable to find polyline information in JavaScript code.")
    sys.exit(1)


# Add nodes and edges only for the paths in path_map.html
for match in re.finditer(r'addLatLng\((.*?), (.*?)\);', polyline):
    loc = tuple(map(float, match.groups()))
    G.add_node(loc, location=df[(df['Latitude'] == loc[0]) & (df['Longitude'] == loc[1])]['Location'].values[0])

for i in range(len(df)):
    for j in range(i + 1, len(df)):
        loc1 = (df.at[i, 'Latitude'], df.at[i, 'Longitude'])
        loc2 = (df.at[j, 'Latitude'], df.at[j, 'Longitude'])

        if G.has_node(loc1) and G.has_node(loc2):
            distance = ((loc1[0] - loc2[0]) ** 2 + (loc1[1] - loc2[1]) ** 2) ** 0.5
            G.add_edge(loc1, loc2, weight=distance)

# Get user input for current and goal locations
current_location = tuple(map(float, input("Enter the current location (latitude,longitude): ").split(',')))
goal_location = tuple(map(float, input("Enter the goal location (latitude,longitude): ").split(',')))

# Find the path using Dijkstra's algorithm
path = dijkstra(G, current_location, goal_location)

# Create a Folium map centered around Wales
map_wales_dijkstra = folium.Map(location=[52.1307, -3.7837], zoom_start=8)

# Add markers for each location
for index, row in df.iterrows():
    folium.Marker([row['Latitude'], row['Longitude']], popup=row['Location']).add_to(map_wales_dijkstra)

# Add a path using the Dijkstra path
if path:
    folium.PolyLine(locations=path, color='red').add_to(map_wales_dijkstra)

# Save the map as an HTML file
map_wales_dijkstra.save('path_map_dijkstra.html')

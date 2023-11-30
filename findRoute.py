import pandas as pd
from geopy.distance import geodesic
import networkx as nx
import matplotlib.pyplot as plt

def calculate_distance(coord1, coord2):
    return geodesic(coord1, coord2).kilometers

def create_graph(df):
    G = nx.Graph()

    for i, row in df.iterrows():
        G.add_node(i, pos=(row['Latitude'], row['Longitude']))

    for i in range(len(df)):
        for j in range(i + 1, len(df)):
            distance = calculate_distance((df.iloc[i]['Latitude'], df.iloc[i]['Longitude']),
                                          (df.iloc[j]['Latitude'], df.iloc[j]['Longitude']))
            G.add_edge(i, j, weight=distance)

    return G

def find_best_route(graph, start, goal):
    path = nx.shortest_path(graph, source=start, target=goal, weight='weight')
    return path

def main():
    # Load CSV file
    csv_file = 'Wales.csv'
    df = pd.read_csv(csv_file)

    # Get user input for start and goal coordinates
    start_latitude = float(input("Enter current latitude: "))
    start_longitude = float(input("Enter current longitude: "))
    goal_latitude = float(input("Enter goal latitude: "))
    goal_longitude = float(input("Enter goal longitude: "))

    # Find nearest nodes in the graph to the entered coordinates
    start_node = ((df['Latitude'] - start_latitude)**2 + (df['Longitude'] - start_longitude)**2).idxmin()
    goal_node = ((df['Latitude'] - goal_latitude)**2 + (df['Longitude'] - goal_longitude)**2).idxmin()

    # Create graph and find the best route
    graph = create_graph(df)
    best_route = find_best_route(graph, start=start_node, goal=goal_node)

    # Print the best route
    print(f"Best route: {best_route}")

    # Plot the graph with the best route
    pos = nx.get_node_attributes(graph, 'pos')
    nx.draw(graph, pos, with_labels=True, font_weight='bold')
    nx.draw_networkx_edges(graph, pos, edgelist=[(best_route[i], best_route[i + 1]) for i in range(len(best_route) - 1)], edge_color='r', width=2)
    plt.show()

if __name__ == "__main__":
    main()

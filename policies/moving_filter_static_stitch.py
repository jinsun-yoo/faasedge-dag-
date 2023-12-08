from collections import defaultdict
from sklearn.neighbors import KDTree
from sklearn.cluster import KMeans
import numpy as np
from tqdm import tqdm
from utils import post_process_links_bw_consumption

def policy(time_index, cell_tower_ids, cell_tower_locations, edge_node_tree, meta):
    # given car position at time t, find nearest edge node for filter function, find stitch 1, find stitch 2. 
    # make trees for filter, s1, and s2. 
    FIRST_HOP_DATA, SECOND_HOP_DATA, THIRD_HOP_DATA = meta
    
    times = list(time_index.keys())
    min_t = min(times)
    max_t = max(times)
    
    time_index_ = defaultdict(dict, time_index)
    
    edge_node_bandwidth_consumption = [0] * len(cell_tower_ids)
    links_bandwidth_consumption = defaultdict(int)
    
    filter_function_vehicle_node_map = {}   # { car_id : edge_node_id }
    # populate this as a new car is seen. 
    
    # add filter locations for cars at t=0 to determine S1, S2 locations
    time_zero_cars = time_index_[min_t]
    filter_locations = []
    
    for car, position in time_zero_cars.items():
        _, closest_node = edge_node_tree.query(np.array([list(position)]))
        filter_node_id = closest_node[0][0]
        filter_function_vehicle_node_map[car] = filter_node_id
        filter_node_position = cell_tower_locations[filter_node_id]
        filter_locations.append(filter_node_position)
        
    filter_locations = np.concatenate((cell_tower_locations, np.array(filter_locations)), axis=0)
        
    # calculate S1 placement
    first_stitch_functions = KMeans(n_clusters=5, random_state=0).fit(filter_locations)
    first_stitch_functions_centroids_raw = first_stitch_functions.cluster_centers_
    
    first_stitch_functions_centroids = []
    for centroid in first_stitch_functions_centroids_raw:
        dist, closest_node = edge_node_tree.query(centroid[np.newaxis, :])
        node_idx = closest_node[0][0]
        first_stitch_functions_centroids.append((node_idx, cell_tower_locations[node_idx]))
        
    first_stitch_tree_index_to_node_id = [id_ for id_, _ in first_stitch_functions_centroids]
    
    # calculate S2 placement
    second_stitch_functions = KMeans(n_clusters=1, random_state=0).fit(filter_locations)
    second_stitch_functions_centroids_raw = second_stitch_functions.cluster_centers_
    
    second_stitch_functions_centroids = []
    for centroid in second_stitch_functions_centroids_raw:
        dist, closest_node = edge_node_tree.query(centroid[np.newaxis, :])
        node_idx = closest_node[0][0]
        second_stitch_functions_centroids.append((node_idx, cell_tower_locations[node_idx]))
        
    second_stitch_tree_index_to_node_id = [id_ for id_, _ in second_stitch_functions_centroids]
    
    # filter function will find closest S1
    first_stitch_node_tree = KDTree(np.array([point for _, point in first_stitch_functions_centroids]), leaf_size=2)
    
    # S1 function will find closest S2 (only one in this case)
    second_stitch_node_tree = KDTree(np.array([point for _, point in second_stitch_functions_centroids]), leaf_size=2)

    for current_time in tqdm(range(min_t, max_t + 1)):
        current_car_positions = time_index_[current_time]
        for car_id, position in current_car_positions.items():                
            # add data to links
            _, node = edge_node_tree.query(np.array([list(position)]))
            filter_node = node[0][0]
            
            _, node = first_stitch_node_tree.query(cell_tower_locations[filter_node][np.newaxis, :])
            first_stitch = first_stitch_tree_index_to_node_id[node[0][0]]
 
            _, node = second_stitch_node_tree.query(cell_tower_locations[first_stitch][np.newaxis, :])
            second_stitch = second_stitch_tree_index_to_node_id[node[0][0]]
            
            edge_node_bandwidth_consumption[filter_node] += FIRST_HOP_DATA
            links_bandwidth_consumption[(filter_node, first_stitch)] += SECOND_HOP_DATA
            links_bandwidth_consumption[(first_stitch, second_stitch)] += THIRD_HOP_DATA

    return edge_node_bandwidth_consumption, post_process_links_bw_consumption(links_bandwidth_consumption)
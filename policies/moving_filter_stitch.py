from collections import defaultdict
from sklearn.neighbors import KDTree
from sklearn.cluster import KMeans
import numpy as np
from tqdm import tqdm
from utils import post_process_links_bw_consumption

def policy(time_index, cell_tower_ids, cell_tower_locations, edge_node_tree, meta):
    
    def filter_map_to_filter_locations(filter_function_vehicle_node_map):
        return np.array([ cell_tower_locations[filter_node] for _, filter_node in filter_function_vehicle_node_map.items() ])
    
    # given car position at time t, find nearest edge node for filter function, recompute stitch 1, recompute stitch 2. 
    # make trees for filter, s1, and s2. 

    times = list(time_index.keys())
    min_t = min(times)
    max_t = max(times)
    
    FIRST_HOP_DATA, SECOND_HOP_DATA, THIRD_HOP_DATA = meta
    
    time_index_ = defaultdict(dict, time_index)
    
    edge_node_bandwidth_consumption = [0] * len(cell_tower_ids)
    links_bandwidth_consumption = defaultdict(int)
    
    filter_function_vehicle_node_map = {}   # { car_id : edge_node_id }
    filter_locations = None
    first_stitch_node_tree = None
    first_stitch_tree_index_to_node_id = None
    second_stitch_node_tree = None
    second_stitch_tree_index_to_node_id = None
    updated = 0
    
    for current_time in tqdm(range(min_t, max_t + 1)):
        current_car_positions = time_index_[current_time]
    
        update_placement = False
        # filter_locations = np.concatenate((cell_tower_locations, np.array(filter_locations)), axis=0)
        for car_id, position in current_car_positions.items():
            _, node = edge_node_tree.query(np.array([list(position)]))
            filter_node = node[0][0]
            
            if car_id not in filter_function_vehicle_node_map or filter_function_vehicle_node_map[car_id] != filter_node:
                filter_function_vehicle_node_map[car_id] = filter_node
                update_placement = True
        
        if update_placement:
            updated += 1
            filter_locations = filter_map_to_filter_locations(filter_function_vehicle_node_map=filter_function_vehicle_node_map)
        
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
        
        for car_id, position in current_car_positions.items():
            filter_node = filter_function_vehicle_node_map[car_id]
            
            _, node = first_stitch_node_tree.query(cell_tower_locations[filter_node][np.newaxis, :])
            first_stitch = first_stitch_tree_index_to_node_id[node[0][0]]
 
            _, node = second_stitch_node_tree.query(cell_tower_locations[first_stitch][np.newaxis, :])
            second_stitch = second_stitch_tree_index_to_node_id[node[0][0]]
            
            edge_node_bandwidth_consumption[filter_node] += FIRST_HOP_DATA
            links_bandwidth_consumption[(filter_node, first_stitch)] += SECOND_HOP_DATA
            links_bandwidth_consumption[(first_stitch, second_stitch)] += THIRD_HOP_DATA

    print(f"Placement updated {updated} times.")
    return edge_node_bandwidth_consumption, post_process_links_bw_consumption(links_bandwidth_consumption)
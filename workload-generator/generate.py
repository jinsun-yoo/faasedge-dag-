import json
import threading
import requests
import random
import time
import sys

with open('map.json', 'r') as file:
    map_data = json.load(file)

def generate_coordinate(lowerLeft, upperRight):
    x = random.uniform(lowerLeft['x'], upperRight['x'])
    y = random.uniform(lowerLeft['y'], upperRight['y'])
    return {'x': x, 'y': y}

def client_request(client_id, lowerLeft, upperRight):
    while True:
        coord = generate_coordinate(lowerLeft, upperRight)
        headers = {'Content-Type': 'application/json', 'Client-ID': str(client_id)}
        data = {'coordinate': coord}
        # response = requests.post('http://localhost:8000/dag/invoke/DagName', headers=headers, json=data)
        # print(f"Client {client_id} Response: {response.status_code}")
        print(f"client {client_id} calling endpoint with coords : {data}")
        time.sleep(random.randint(1, 5))  # simulate periodic requests

def generate_clients(n, lowerLeft, upperRight):
    for i in range(n):
        client_thread = threading.Thread(target=client_request, args=(i, lowerLeft, upperRight))
        client_thread.start()

n_clients = int(sys.argv[1])
generate_clients(n_clients, map_data['lowerLeft'], map_data['upperRight'])

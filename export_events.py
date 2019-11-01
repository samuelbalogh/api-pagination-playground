import time
import random

import requests

HOST = "localhost"
PORT = 8000
API_URI = f"http://{HOST}:{PORT}"

LIMIT = 10

def offset_paginate():
    events = []
    finished = False
    offset = 0
    while not finished:
        response = requests.get(f"{API_URI}/events?limit={LIMIT}&offset={offset}").json()
        offset += LIMIT
        events.extend(response.get("Value"))
        if len(response.get("Value")) < limit:
            finished = True
            total_count = requests.get(f"{API_URI}/events_count").json()
        time.sleep(.3)

    print(f"total count: {total_count}")
    print(f"exported {len(events)} events")
    print(f"exported {len(set(event['ID'] for event in events))} unique events")
    return events

def keyset_paginate():
    events = []
    response = requests.get(f"{API_URI}/events?limit={LIMIT}").json()
    events.extend(response.get("Value"))
    cursor = response.get("Cursor")
    finished = False
    print(response)
    while not finished:
        response = requests.get(f"{API_URI}/events?limit={LIMIT}&cursor={cursor}").json()
        cursor = response.get("Cursor")
        events.extend(response.get("Value"))
        if len(response.get("Value")) < LIMIT:
            finished = True
            total_count = requests.get(f"{API_URI}/events_count").json()
    print(f"total count: {total_count}")
    print(f"exported {len(events)} events")
    print(f"exported {len(set(event['ID'] for event in events))} unique events")




if __name__ == "__main__":
    keyset_paginate()

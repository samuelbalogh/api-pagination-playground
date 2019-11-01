import time
import random

import requests

HOST = "localhost"
PORT = 8000
ENDPOINT = "events"
API_URI = f"http://{HOST}:{PORT}"

def get_events(limit=10):
    response = requests.get(f"{API_URI}/events?limit={limit}")
    return response.json().get("Value")

def update_event(event):
    endpoint = f"{API_URI}/events/{event['ID']}"
    payload = {"description": f"description {random.randint(0, 100) * 'long string'}"}
    res = requests.put(endpoint, data=payload)
    return

def create_events(count=10):
    payload = {"count": count}
    requests.post(f"{API_URI}/fakedata", data=payload)

def run():
    while True:
        try:
            if random.randint(0, 10) == 100:
                create_events(count=random.randint(1, 10))
            events = get_events(limit=100)
            if not events:
                continue
            event = random.choice(events)
            update_event(event)
            time.sleep(.2)
        except:
            time.sleep(2)

if __name__ == "__main__":
    run()

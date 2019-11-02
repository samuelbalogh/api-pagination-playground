import time
import random

import requests

HOST = "localhost"
PORT = 8000
API_URI = f"http://{HOST}:{PORT}"

LIMIT = 50

def offset_paginate(orderby=None):
    events = []
    finished = False
    offset = 0
    if orderby is None:
        orderby = ''
    while not finished:
        response = requests.get(f"{API_URI}/events?limit={LIMIT}&offset={offset}&orderby={orderby}").json()
        offset += LIMIT
        events.extend(response.get("Value"))
        time.sleep(.1)
        if len(response.get("Value")) < LIMIT:
            finished = True
            total_count = requests.get(f"{API_URI}/events_count").json()
    _print_results(total_count, events, "offset", orderby)
    return events

def keyset_paginate(orderby=None):
    events = []
    response = requests.get(f"{API_URI}/events?limit={LIMIT}").json()
    events.extend(response.get("Value"))
    cursor = response.get("Cursor")
    finished = False
    while not finished:
        response = requests.get(f"{API_URI}/events?limit={LIMIT}&cursor={cursor}").json()
        cursor = response.get("Cursor")
        events.extend(response.get("Value"))
        time.sleep(.1)
        if len(response.get("Value")) < LIMIT:
            finished = True
            total_count = requests.get(f"{API_URI}/events_count").json()
    _print_results(total_count, events, "keyset", orderby)


def _print_results(total_count, events, pagination_method, orderby="default"):
    print("---")
    print(f"Ordering used: {orderby}")
    print(f"Pagination method used: {pagination_method}")
    print(f"total count: {total_count}")
    print(f"exported {len(events)} events")
    print(f"exported {len(set(event['ID'] for event in events))} unique events")
    print(f"{total_count - len(set(event['ID'] for event in events))} events missed")


if __name__ == "__main__":
    offset_paginate()
    offset_paginate(orderby="updated_at desc")
    offset_paginate(orderby="updated_at asc")
    offset_paginate(orderby="created_at asc")
    offset_paginate(orderby="updated_at desc")
    keyset_paginate()

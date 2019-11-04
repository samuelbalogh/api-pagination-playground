import time
import random

import requests

HOST = "localhost"
PORT = 8000
API_URI = f"http://{HOST}:{PORT}"

LIMIT = 500

def timeit(method):
    def timed(*args, **kw):
        start = time.time()
        result = method(*args, **kw)
        end = time.time()
        print(f"{method.__name__, round(end - start, 2)} seconds")
        return result
    return timed

@timeit
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
        if len(response.get("Value")) < LIMIT:
            finished = True
            total_count = int(response.get("Count"))
    _print_results(total_count, events, "offset", orderby)
    return events

@timeit
def keyset_paginate():
    events = []
    response = requests.get(f"{API_URI}/events?limit={LIMIT}").json()
    events.extend(response.get("Value"))
    cursor = response.get("Cursor")
    finished = False
    while not finished:
        response = requests.get(f"{API_URI}/events?limit={LIMIT}&cursor={cursor}").json()
        cursor = response.get("Cursor")
        events.extend(response.get("Value"))
        if len(response.get("Value")) < LIMIT:
            finished = True
            total_count = int(response.get("Count"))
    _print_results(total_count, events, "keyset", orderby="default as per implementation")


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
    offset_paginate(orderby="updated_at asc")
    offset_paginate(orderby="updated_at desc")
    offset_paginate(orderby="created_at asc")
    offset_paginate(orderby="created_at desc")
    keyset_paginate()

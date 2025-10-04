import asyncio
import datetime
import json
import random
import uuid
from fastapi import FastAPI, Request
from sse_starlette.sse import EventSourceResponse

app = FastAPI()

STREAM_DELAY = 1  # second
RETRY_TIMEOUT = 15000  # millisecond


async def generate_random_numbers(request: Request):
    """Generate random numbers continuously with proper SSE format."""
    count = 0
    while True:
        # Check if client disconnected
        if await request.is_disconnected():
            print("Client disconnected")
            break

        random_number = random.randint(1, 100)
        print(f"Generated random number: {random_number}")

        # Yield properly formatted SSE message
        yield {
            "event": "random_number",
            "retry": RETRY_TIMEOUT,
            "data": json.dumps({
                "value": random_number,
                "timestamp": datetime.datetime.now().isoformat(sep="T", timespec="auto")
            }),
            "id": str(uuid.uuid4())
        }

        await asyncio.sleep(STREAM_DELAY)

        # Optional: break after some time to allow clean shutdown
        count += 1
        if count >= 60:  # Stop after 60 seconds
            break


@app.get("/poll")
async def poll(request: Request):
    """SSE endpoint that streams random numbers."""
    return EventSourceResponse(generate_random_numbers(request))


@app.get("/")
async def root():
    """Root endpoint."""
    return {"message": "SSE Server running. Visit /poll for the event stream."}

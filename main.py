import requests
import json
import time

# URL of the Go API server
url = "http://localhost:8080/query"

# Sample payload for generating traffic
payload = {
    "metric_name": "up",
    "labels": {
        # "mode": "idle",
        "instance": "prometheus:9090"
    },
    "start_time": "2024-09-15T03:30:00Z",
    "end_time": "2024-09-15T03:32:00Z"
}


# Function to generate traffic by sending requests
def generate_traffic():
    for i in range(10):  # Change the range value to generate more traffic
        try:
            response = requests.post(url, data=json.dumps(payload), headers={'Content-Type': 'application/json'})
            if response.status_code == 200:
                print(f"Request {i + 1}: Success - Response received")
                print(f"Response: {response.json()}")
            else:
                print(f"Request {i + 1}: Failed - Status Code: {response.status_code}, Message: {response.text}")
        except Exception as e:
            print(f"Request {i + 1}: Error - {e}")

        time.sleep(1)  # Adjust the sleep time as needed to control the request rate


# Run the traffic generator
if __name__ == "__main__":
    generate_traffic()

import os
import requests

log_directory = '/home/sftp/uploads'

central_server_url = 'http://10.0.0.204:5000/receive_data'

sftp_server = 'sftp-2 10.0.0.202'

def parse_logs(log_directory):
    parsed_data = []
    for filename in os.listdir(log_directory):
        if filename.endswith(".txt"):
            file_path = os.path.join(log_directory, filename)
            with open(file_path, 'r') as file:
                content = file.read()
                numbers = [int(s) for s in content.split() if s.isdigit()]
                parsed_data.append((filename, sum(numbers)))
    return parsed_data

def send_data(parsed_data):
    for filename, result in parsed_data:
        response = requests.post(central_server_url, json={
            'filename': filename,
            'result': result,
            'sftp_server': sftp_server
        })
        print(f"File {filename} with result {result} has been sent: {response.status_code}, {response.text}")

if __name__ == "__main__":
    parsed_data = parse_logs(log_directory)
    send_data(parsed_data)


from flask import Flask, request, jsonify, render_template
import sqlite3
import os
from datetime import datetime

app = Flask(__name__)
DATABASE = '/home/vagrant/app/logs.db'

def get_db():
    conn = sqlite3.connect(DATABASE)
    return conn

@app.route('/receive_data', methods=['POST'])
def receive_data():
    data = request.get_json()
    filename = data.get('filename')
    result = data.get('result')
    timestamp = datetime.now().strftime('%Y-%m-%d %H:%M:%S')

    conn = get_db()
    c = conn.cursor()
    
    # Check if a record with the same filename already exists
    c.execute("SELECT COUNT(*) FROM logs WHERE filename = ?", (filename,))
    exists = c.fetchone()[0] > 0

    if not exists:
        # Insert the new record if it does not already exist
        c.execute("INSERT INTO logs (filename, result, timestamp) VALUES (?, ?, ?)",
                  (filename, result, timestamp))
        conn.commit()

    return jsonify({"message": "Data successfully received"}), 200

@app.route('/report', methods=['GET'])
def generate_report():
    conn = get_db()
    c = conn.cursor()
    c.execute("SELECT filename, result, timestamp FROM logs")
    logs = c.fetchall()

    report = {}
    for log in logs:
        filename, result, timestamp = log
        server_name = filename.split('_')[0]
        if server_name not in report:
            report[server_name] = []
        report[server_name].append({"result": result, "timestamp": timestamp})

    return jsonify(report)

@app.route('/report/html', methods=['GET'])
def report_html():
    conn = get_db()
    c = conn.cursor()
    c.execute("SELECT filename, result, timestamp FROM logs")
    logs = c.fetchall()
    
    report = {}
    for log in logs:
        filename, result, timestamp = log
        server_name = filename.split('_')[0]
        if server_name not in report:
            report[server_name] = []
        report[server_name].append({"result": result, "timestamp": timestamp})

    return render_template('report.html', report=report)

if __name__ == '__main__':
    if not os.path.exists('/home/vagrant/app/uploads'):
        os.makedirs('/home/vagrant/app/uploads')
    if not os.path.exists(DATABASE):
        conn = get_db()
        c = conn.cursor()
        c.execute('''CREATE TABLE logs (id INTEGER PRIMARY KEY AUTOINCREMENT, filename TEXT, result INTEGER, timestamp TEXT)''')
        conn.commit()
    app.run(host='0.0.0.0', port=5000)
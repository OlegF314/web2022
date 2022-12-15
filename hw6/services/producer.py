from http.server import BaseHTTPRequestHandler, HTTPServer
from kafka import KafkaProducer
import redis

from parameters import REDIS_PORT, KAFKA_PORT, SERVER_PORT


producer = KafkaProducer(bootstrap_servers=f'localhost:{KAFKA_PORT}')
pool = redis.ConnectionPool(host='localhost', port=REDIS_PORT, db=0)
redis = redis.Redis(connection_pool=pool)

class HttpProcessor(BaseHTTPRequestHandler):
    def do_GET(self):
        request = self.requestline
        task_no = request.split()[1]
        redis.set(f'task{task_no} in_queue')
        producer.send('task_queue', b'task_no')  
        self.send_response(200)
        self.send_header('content-type','text/html')
        self.end_headers()
        self.wfile.write(b"text")


def main():
    serv = HTTPServer(("localhost", SERVER_PORT), HttpProcessor)
    serv.serve_forever()


if __name__ == '__main__':
    main()

from kafka import KafkaConsumer
import time
import redis
from circuitbreaker import circuit
from healthcheck import HealthCheck
from parameters import REDIS_PORT, KAFKA_PORT, SERVER_PORT, WAIT_TIME

def redis_correct():
    redis_ = redis.Redis(connection_pool=pool)


pool = redis.ConnectionPool(host='localhost', port=REDIS_PORT, db=0)
@circuit
def main():
    consumer = KafkaConsumer('task_queue')
    redis_ = redis.Redis(connection_pool=pool)
    health.add_check(redis_correct)
    health.run()
    for message in consumer:
        print(message)
        redis_.set(f'task "{message}"', 'accepted')
        time.sleep(WAIT_TIME)
        redis_.set(f'task "{message}"', 'done')

main()

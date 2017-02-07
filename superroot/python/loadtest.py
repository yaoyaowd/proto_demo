import grpc
import time
from multiprocessing import Pool

import superroot_pb2
# python -m grpc_tools.protoc -I=../protos/ --python_out=. --grpc_python_out=. ../protos/superroot.proto
import superroot_pb2_grpc


def run():
    # Generate n requests / sec, total m requests
    channel = grpc.insecure_channel('localhost:8999')
    stub = superroot_pb2_grpc.SuperRootStub(channel)
    requests = 0
    start_time = time.time()
    while True:
        response = stub.Search(superroot_pb2.SearchRequest(
            query="shoes", offset=0, limit=25))
        requests += 1
        if requests % 1000 == 0:
            time_spend = time.time() - start_time
            print 'qps:', 1000 / float(time_spend)
            start_time = time.time()

        if requests == 10000:
            break


def parallel(n):
    pool = Pool(processes=n)
    results = []
    for i in range(n):
        results.append(pool.apply_async(run))
    [x.get() for x in results]


if __name__ == '__main__':
    parallel(10)

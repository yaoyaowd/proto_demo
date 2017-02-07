import grpc

import superroot_pb2
# python -m grpc_tools.protoc -I=../protos/ --python_out=. --grpc_python_out=. ../protos/superroot.proto
import superroot_pb2_grpc


def run():
    channel = grpc.insecure_channel('localhost:8999')
    stub = superroot_pb2_grpc.SuperRootStub(channel)
    response = stub.Search(superroot_pb2.SearchRequest(
        query="shoes", offset=0, limit=25))
    print response


if __name__ == '__main__':
    run()

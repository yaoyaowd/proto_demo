from concurrent import futures
import requests
import time
import ujson
import urllib

import grpc
import gevent

from superroot_pb2 import SearchResponse, SearchDoc
import superroot_pb2_grpc


def search(host, query, offset, limit):
    q_str = urllib.urlencode({'q': '{!cache=false}' + query})
    url = 'http://%s:8999/solr/wishsolrcluster/select?' \
          'defType=edismax&qf=product_description%%20text&' \
          'wt=json&start=%d&count=%d&docsAllowed=50000&omitHeader=true' \
          '&fl=id%%20score' % (host, offset, limit)
    url = url + '&' + q_str
    try:
        r = requests.get(url, timeout=2)
    except Exception:
        return []
    response = ujson.loads(r.content)
    return response.get('response', {}).get('docs', [])


class SuperRoot(superroot_pb2_grpc.SuperRootServicer):

    hosts = ['10.10.32.23',
             '10.10.32.89',
             '10.10.32.185',
             '10.10.32.42',
             '10.10.32.165']

    def Search(self, request, context):
        greenlets = [gevent.spawn(search, host, request.query, request.offset, request.limit)
                     for host in self.hosts]
        gevent.joinall(greenlets)
        return SearchResponse()


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=500))
    superroot_pb2_grpc.add_SuperRootServicer_to_server(SuperRoot(), server)
    server.add_insecure_port('[::]:8999')
    server.start()
    try:
        while True:
            time.sleep(10)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == '__main__':
    serve()
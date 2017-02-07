from concurrent import futures
from multiprocessing import Pool
import requests
import time
import ujson
import urllib

import grpc

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

worker_pool = Pool(processes=200)

class SuperRoot(superroot_pb2_grpc.SuperRootServicer):

    hosts = ['10.10.32.23',
             '10.10.32.89',
             '10.10.32.185',
             '10.10.32.42',
             '10.10.32.165']

    def Search(self, request, context):
        future_results = []
        for host in self.hosts:
            future_results.append(worker_pool.apply_async(
                search,
                (host, request.query, request.offset, request.limit, )))
        docs_dict = {}
        for fr in future_results:
            for doc in fr.get():
                docs_dict[doc.get('id')] = doc.get('score')

        return SearchResponse()
        # sorted_docs = sorted(docs_dict, key=docs_dict.get, reverse=True)
        #
        # docs = []
        # for index in range(request.offset, request.offset + request.limit):
        #     if index < len(sorted_docs):
        #         docs.append(SearchDoc(id=sorted_docs[index],
        #                               score=docs_dict.get(sorted_docs[index])))
        # return SearchResponse(docs=docs)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=200))
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
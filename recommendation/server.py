#!/usr/bin/python

from concurrent import futures
import grpc
import os
import random
import time

import app_pb2
import app_pb2_grpc
from grpc_health.v1 import health_pb2
from grpc_health.v1 import health_pb2_grpc

from logger import getJSONLogger
logger = getJSONLogger('server')

class RecommendationService(app_pb2_grpc.RecommendationServiceServicer):
    def ListRecommendations(self, request, context):
        max_responses = 5
        response = catalog_stub.ListProducts(app_pb2.Empty())
        product_ids = [x.id for x in response.products]
        filtered_products = list(set(product_ids)-set(request.product_ids))
        num_products = len(filtered_products)
        num_return = min(max_responses, num_products)
        indices = random.sample(range(num_products), num_return)
        prod_list = [filtered_products[i] for i in indices]
        logger.info('[Recv ListRecommendaitons] product_ids={}'.format(prod_list))
        resp = app_pb2.ListRecommendationsResponse()
        resp.product_ids.extend(prod_list)
        return resp

    def Check(self, request, context):
        return health_pb2.HealthCheckResponse(
            status=health_pb2.HealthCheckResponse.SERVING)

    def Watch(self, request, context, send_respone_callback=None):
        return None


if __name__ == "__main__":
    logger.info('initializing recommendation server')

    port = os.environ.get('PORT', '14000')
    catalog_addr = os.environ.get('CATALOG_ADDR', '')
    if catalog_addr == '':
        raise Exception('CATALOG_ADDR env variable not set')

    logger.info('catalog address: {}'.format(catalog_addr))
    channel = grpc.insecure_channel(catalog_addr)
    catalog_stub = app_pb2_grpc.ProductCatalogServiceStub(channel)

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    service = RecommendationService()
    app_pb2_grpc.add_RecommendationServiceServicer_to_server(service, server)
    health_pb2_grpc.add_HealthServicer_to_server(service, server)

    logger.info('listening on port: {}'.format(port))
    server.add_insecure_port('[::]:{}'.format(port))
    server.start()

    try:
        while True:
            time.sleep(1000)
    except KeyboardInterrupt:
        server.stop(0)

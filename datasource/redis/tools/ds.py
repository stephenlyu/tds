import redis
from entity_pb2 import *

class DataSource:
    def __init__(self, proto_class, ds_type="tick", redis_url="localhost", redis_port=6379, redis_password=""):
        self.proto_class = proto_class
        self.ds_type = ds_type
        self.redis = redis.StrictRedis(host=redis_url, port=redis_port, password=redis_password)

    def get_data(self, code):
        l = self.redis.zrange('%s.%s' % (self.ds_type, code), 0, 0xFFFFFFFFFFFFFFF)
        ret = []
        for bs in l:
            o = self.proto_class()
            o.ParseFromString(bs)
            ret.append(o)
        return ret


class TickDataSource(DataSource):
    def __init__(self, redis_url="localhost", redis_port=6379, redis_password=""):
        super(TickDataSource, self).__init__(ProtoTick, "tick", redis_url, redis_port, redis_password)


class M1DataSource(DataSource):
    def __init__(self, redis_url="localhost", redis_port=6379, redis_password=""):
        super(M1DataSource, self).__init__(ProtoRecord, "M1", redis_url, redis_port, redis_password)

import json,time,queue,threading,logging,sys
from cassandra.cluster import Cluster
from cassandra.policies import RoundRobinPolicy
import asyncio

from nats.aio.client import Client as NATS
from stan.aio.client import Client as STAN


UPDATE_QUERY = """UPDATE counterks.page_view_counts
SET counter_value = counter_value + 1
WHERE url_name=? AND page_name=?;"""


class DB:
    def __init__(self):
        self.cluster = None
        self.session = None
        self.query_stmt = None

    def ensure_db(self):
        if self.session is not None:
            return

        logging.info('DB is not cached, connecting.')


        while True:
            try:
                self.cluster = Cluster(['cassandra'], load_balancing_policy = RoundRobinPolicy())
                self.session = self.cluster.connect('system')
                self.session.execute("CREATE KEYSPACE IF NOT EXISTS counterks WITH replication = {'class':'SimpleStrategy','replication_factor':'1'};")
                self.session.execute("CREATE TABLE IF NOT EXISTS counterks.page_view_counts (counter_value counter, url_name varchar, page_name varchar, PRIMARY KEY (url_name, page_name) );")

                self.session.set_keyspace('counterks')
            except Exception as e:
                self.close()
                logging.info('%s failed to init cassandra. retrying', e)
                time.sleep(5)
            else:
                logging.info('connected to cassandra')
                break
        logging.info("Connected to DB, preparing statement")
        self.query_stmt = self.session.prepare(UPDATE_QUERY)

    def close(self):
        if self.cluster is not None:
            self.cluster.shutdown()
        self.cluster = None
        self.session = None
        self.query_stmt = None

def rundb(db, q):
    db.ensure_db()
    while True:
        data = q.get()
        logging.info("Received a message from Q %s", data)
        deserdata = json.loads(data)
        db.session.execute(db.query_stmt, [deserdata["Url"], deserdata["Page"]])
     
async def run(loop, q):
    # Use borrowed connection for NATS then mount NATS Streaming
    # client on top.
    nc = NATS()
    await nc.connect(servers=["nats://nats-streaming:4222"], io_loop=loop)

    # Start session with NATS Streaming cluster.
    sc = STAN()
    await sc.connect("test-cluster", "analytics-client-1", nats=nc)

    async def cb(msg):
        q.put(msg.data)
        logging.info("Received a message (seq={}): {}".format(msg.seq, msg.data))

    # Subscribe to get all messages since beginning.
    sub = await sc.subscribe("analytics", start_at='first', cb=cb)

    while True:
        await asyncio.sleep(1, loop=loop)

    await sub.unsubscribe()

    # Close NATS Streaming session
    await sc.close()

    # We are using a NATS borrowed connection so we need to close manually.
    await nc.close()

if __name__ == '__main__':
    logging.basicConfig(stream = sys.stderr, level=logging.DEBUG)
    logging.info("hello!")
    q = queue.Queue()
    db = DB()
    def rundbwithq():
        while True:
            try:
                logging.info("Starting DB event loop")
                rundb(db, q)
            except Exception as e:
                db.close()
                logging.info("Run DB: exception %s", e)
                time.sleep(5)
    t = threading.Thread(target=rundbwithq)
    t.start()


    while True:
        try:
            logging.info("Starting NATS event loop")
            loop = asyncio.get_event_loop()
            loop.run_until_complete(run(loop, q))
            loop.close()
        except Exception as e:
            logging.info("NATS: exception %s", e)
            time.sleep(5)
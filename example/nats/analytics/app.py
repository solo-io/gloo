import json,time
from cassandra.cluster import Cluster
from flask import Flask, request, abort, g
app = Flask(__name__)

UPDATE_QUERY = """UPDATE counterks.page_view_counts
SET counter_value = counter_value + 1
WHERE url_name=? AND page_name=?;"""


SELECT_QUERY = "SELECT * FROM counterks.page_view_counts;"



def get_db():
    db = getattr(g, '_database', None)
    if db is not None:
        return db

    app.logger.info('DB is not cached, connecting.')

    while True:
        try:
            cluster = Cluster(['cassandra'])
            session = cluster.connect('system')
            session.execute("CREATE KEYSPACE IF NOT EXISTS counterks WITH replication = {'class':'SimpleStrategy','replication_factor':'1'};")
            session.execute("CREATE TABLE IF NOT EXISTS counterks.page_view_counts (counter_value counter, url_name varchar, page_name varchar, PRIMARY KEY (url_name, page_name) );")

            session.set_keyspace('counterks')
        except Exception as e:
            cluster.shutdown()
            app.logger.info('%s failed to init cassandra. retrying', e)
            time.sleep(5)
        else:
            app.logger.info('connected to cassandra')
            break
    query_stmt = session.prepare(UPDATE_QUERY)
    g._database = (cluster, session, query_stmt)
    return g._database

@app.teardown_appcontext
def teardown_db(exception):
    db = getattr(g, '_database', None)
    if db is not None:
        cluster,_, _ = db
        cluster.shutdown()
        g._database = None

@app.route("/analytics", methods=['GET', 'POST'])
def analytics():
    app.logger.info('Analytics %s request!', request.method)
    _, session, query = get_db()
    if request.method == 'POST':
        data = request.get_json(force=True)
        session.execute(query, [data["Url"], data["Page"]])
        return ''
    counters = session.execute(SELECT_QUERY)
    countersdict = {}
    for counter in counters:
        key = counter.url_name + ";" + counter.page_name
        countersdict[key] = counter.counter_value
    return json.dumps(countersdict)

if __name__ == '__main__':
    app.run(host='0.0.0.0')
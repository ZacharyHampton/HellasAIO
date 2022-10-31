import time
from buzz.utils import GetPIDs
from buzz.pid import BuzzPIDFlow
from buzz.msku import BuzzMSKUFlow
from buzz.image import BuzzImageFlow
from buzz.backendlink import BackendLinkFlow
from buzz.keyword import KeywordFlow
from buzz.types import Thread
from db import db
from buzz.threads import RunningThreads


def Start():
    t = Thread(BackendLinkFlow, None)
    t.start()

    cur = db.cursor()

    keywordRows = cur.execute('SELECT keyword FROM keywords').fetchall()
    keywords = [keyword[0] for keyword in keywordRows]
    keywordThread = Thread(KeywordFlow, keywords)
    keywordThread.start()

    while True:
        time.sleep(1)

        pidRows = cur.execute('SELECT id FROM products').fetchall()
        for row in pidRows:
            if row[0] not in RunningThreads:
                RunningThreads[row[0]] = Thread(BuzzPIDFlow, row[0])
                RunningThreads[row[0]].start()

        mskuRows = cur.execute('SELECT msku FROM mskus').fetchall()
        for row in mskuRows:
            if row[0] not in RunningThreads:
                RunningThreads[row[0]] = Thread(BuzzMSKUFlow, row[0])
                RunningThreads[row[0]].start()

        imageRows = cur.execute('SELECT msku_partial FROM msku_partials_image').fetchall()
        for row in imageRows:
            if row[0] not in RunningThreads:
                RunningThreads[row[0]] = Thread(BuzzImageFlow, row[0])
                RunningThreads[row[0]].start()

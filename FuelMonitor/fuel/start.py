import threading
import time
import fuel.flow
import fuel.keywords
import fuel.mskus
from fuel.utils import CurrentRunningKeywords
import fuel.loading


def Start():
    time.sleep(10)  #: wait for discord to start
    print("Starting monitor.")
    threading.Thread(target=fuel.mskus.MSKUFlowRewrite).start()

    while True:
        time.sleep(1)
        keywords = set(fuel.loading.GetKeywords())
        keywords = list(keywords)
        for keyword in keywords:
            if keyword not in CurrentRunningKeywords:
                CurrentRunningKeywords.append(keyword)
                print('[{}] Starting keyword monitor.'.format(keyword))
                threading.Thread(target=fuel.flow.KeywordFlow, args=(keyword,)).start()

from fastapi import FastAPI
import uvicorn
import tensorflow as tf

from schemas import Traffic
from typing import List
import random

from traffic_detect import traffic_detect

# test_traffic = Traffic(
#     UserAgent="Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0",
#     Method="GET",
#     Proto="HTTP/1.1",
#     ProtoMajor=1,
#     ProtoMinor=1,
#     ContentLength=0,
#     TransferEncoding=[],
#     Close=False,
#     RemoteAddr="172.20.48.1:50743",
#     RequestURI="http://www.yiidian.com/questions/30562",
#     Scheme="http",
#     Host="www.yiidian.com",
#     Path="/images/ic8.gif"
# )

app = FastAPI()


@app.post("/detect")
def detect(traffic_list: List[Traffic]):
    urls = [traffic.Method + ' ' + traffic.Path for traffic in traffic_list]
    result = traffic_detect(urls)

    for i in range(len(traffic_list)):
        traffic_list[i].IsMalicious = result[i]

    print(result)
    return traffic_list


@app.post("/detect2")
def detect2(traffic: Traffic):
    traffic.IsMalicious = traffic_detect([traffic.Method + ' ' + traffic.Path])[0]
    print(traffic.IsMalicious)
    return traffic


IS_MALICIOUS_ENUM = [0, 1, 2]
WEIGHTS = [0.85, 0.1, 0.05]


@app.post("/test")
def test_detect(traffic_list: List[Traffic]):
    for traffic in traffic_list:
        traffic.IsMalicious = random.choices(IS_MALICIOUS_ENUM, weights=WEIGHTS)[0]

    result = [traffic.IsMalicious for traffic in traffic_list]
    print(result)
    return traffic_list


@app.post("/test2")
def test_detect2(traffic: Traffic):
    traffic.IsMalicious = random.choices(IS_MALICIOUS_ENUM, weights=WEIGHTS)[0]
    print(traffic.IsMalicious)
    return traffic


if __name__ == '__main__':
    uvicorn.run("main:app", host="0.0.0.0", port=8000, reload=True)
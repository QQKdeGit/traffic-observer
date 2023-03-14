FROM python:3.7

COPY requirements.txt requirements.txt
COPY main.py main.py
COPY schemas.py schemas.py
COPY data_loader.py data_loader.py
COPY traffic_detect.py traffic_detect.py

COPY model/ model/
COPY params/ params/

RUN pip install -r requirements.txt \
    -i https://pypi.tuna.tsinghua.edu.cn/simple

CMD ["python3", "main.py"]

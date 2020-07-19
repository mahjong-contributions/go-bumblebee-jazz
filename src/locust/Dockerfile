FROM locustio/locust:1.1.1

WORKDIR /workspace
COPY ./locustfile.py /workspace/locustfile.py

CMD ["-f", "/workspace/locustfile.py"]

ENV PYTHONUNBUFFERED=1
FROM python:3.8-slim
RUN apt-get update -qy \
 && apt-get -qy install g++ wget \
 && rm -rf /var/lib/apt/lists/*
ENV PYTHONUNBUFFERED=0

RUN GRPC_HEALTH_PROBE_VERSION=v0.3.1 \
 && wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 \
 && chmod +x /bin/grpc_health_probe
WORKDIR /recommendation
COPY requirements.txt requirements.txt
RUN pip install -r requirements.txt

COPY . .

ENV PORT "14000"
EXPOSE 14000

ENTRYPOINT ["python", "/recommendation/server.py"]

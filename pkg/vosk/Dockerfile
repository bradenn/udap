ARG TAG=0.3.37-m1

FROM sskorol/vosk-api:$TAG

WORKDIR /opt/vosk-server

COPY asr_server_m1.py asr_server_m1.py

# Download common server

RUN python3 -m pip install --upgrade pip
RUN pip3 install websockets asyncio

EXPOSE 2700

CMD ["python3", "./asr_server_m1.py"]

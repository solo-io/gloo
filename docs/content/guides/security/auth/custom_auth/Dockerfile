FROM python:3-slim

RUN apt update && apt -y upgrade

ADD server.py /

CMD [ "python", "./server.py" ]

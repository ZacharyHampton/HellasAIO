#FROM golang:1.16
#RUN go install mvdan.cc/garble@v0.4.0

FROM python:3.10

WORKDIR /opt/app

COPY . .

RUN pip install --no-cache-dir -r requirements.txt

RUN pip install helheim-0.9.5-py310-linux.x86_64.tar.gz

EXPOSE 80

CMD ["uvicorn", "app:App", "--host=0.0.0.0", "--port=80", "--reload"]
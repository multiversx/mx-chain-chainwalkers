FROM python:3.8

WORKDIR /chainwalkers_app

COPY . .

RUN pip3 install -r requirements.txt

ENV PYTHONPATH "${PYTHONPATH}:/chainwalkers_app"

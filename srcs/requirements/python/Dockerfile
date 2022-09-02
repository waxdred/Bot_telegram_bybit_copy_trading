FROM python:3

WORKDIR /usr/src/app
COPY ./app/telegram.py ./
COPY ./config/* ./config/

RUN pip3.10  install python-dotenv
RUN pip3.10  install Telethon

WORKDIR /usr/src/app/config

CMD [ "python3.10", "../telegram.py" ]

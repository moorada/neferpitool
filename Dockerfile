FROM debian
WORKDIR /app
RUN echo "deb http://ftp.debian.org/debian sid main" >> /etc/apt/sources.list
RUN apt-get update
RUN apt-get install -y libc6 libc-bin
RUN apt-get install -y wget
RUN wget https://XXXX_NEFERPITOOL_
RUN chmod +x ./cmd/cmd
ENTRYPOINT ["./neferpitool"]
CMD ["-h"]
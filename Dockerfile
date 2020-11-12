FROM alpine
ADD exam-setter /go/exam-setter
ADD res/ /go/res
WORKDIR /go
RUN chmod +x exam-setter
EXPOSE 8080
ENTRYPOINT ["./exam-setter"]
CMD []
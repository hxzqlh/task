FROM alpine
ADD task-service /task-service
ENTRYPOINT [ "/task-service" ]

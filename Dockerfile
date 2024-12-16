FROM busybox

WORKDIR /opt
ADD app app
CMD ["./app", "serve"]
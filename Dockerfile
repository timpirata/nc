#
# A multi-stage build is used to obtain a slim service image
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
# ( This is stolen from my SAGA execution coordinator, a bit )

# First stage - Build the app binary:

FROM golang
COPY . /go/src/github.com/schnoddelbotz/nc
WORKDIR /go/src/github.com/schnoddelbotz/nc
RUN make nc

# Second stage - Build the final service image, containing the app:

FROM alpine:latest
WORKDIR /root/
COPY --from=0 /go/src/github.com/schnoddelbotz/nc/build/nc .
CMD ["./nc"]
EXPOSE 2001

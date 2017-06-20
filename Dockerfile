FROM scratch
MAINTAINER Kelsey Hightower <kelsey.hightower@gmail.com>
ADD istio-webhook /istio-webhook
ENTRYPOINT ["/istio-webhook"]

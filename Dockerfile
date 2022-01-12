FROM scratch
MAINTAINER Rid <rid@cylo.io>
ADD dist/hapesay_linux_amd64/hapettpsay hapettpsay
CMD ["/hapettpsay"]
EXPOSE 8000
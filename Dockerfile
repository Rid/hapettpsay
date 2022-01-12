FROM scratch
MAINTAINER Rid <rid@cylo.io>
ADD dist/hapettpsay_linux_amd64/hapettpsay hapettpsay
ADD static static
ADD templates templates
CMD ["/hapettpsay"]
EXPOSE 8000
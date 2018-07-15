FROM scratch
ADD main /
ADD config /config
CMD ["/main"]

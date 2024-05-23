FROM alpine:3.20.0@sha256:77726ef6b57ddf65bb551896826ec38bc3e53f75cdde31354fbffb4f25238ebd

# bubblewrap is for sandboxing, and git permits pulling modules via
# the git protocol
RUN apk add --no-cache bubblewrap git

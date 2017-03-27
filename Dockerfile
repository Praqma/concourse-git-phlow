FROM concourse/buildroot:base

COPY assets/ /opt/resource/
COPY test.json /opt/resource/

ENV TMPDIR=/gitphlow
RUN mkdir $TMPDIR

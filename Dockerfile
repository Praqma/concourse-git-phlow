FROM concourse/buildroot:git

COPY assets/ /opt/resource/

ENV TMPDIR=/gitphlow
RUN mkdir $TMPDIR

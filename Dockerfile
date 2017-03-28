FROM concourse/buildroot:git

COPY assets/ /opt/resource/
COPY git-phlow /bin

ENV TMPDIR=/gitphlow
RUN mkdir $TMPDIR

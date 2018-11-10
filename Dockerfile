FROM node:latest
FROM python:2
FROM ruby:latest

# specific libs for python
RUN pip install \
    cachetools \
    nose \
    python-dateutil \
    pytz \
    six

# add user
RUN groupadd cb-code-runner
RUN useradd -m -d /home/cb-code-runner -g cb-code-runner -s /bin/bash cb-code-runner

# install code runner
ADD https://github.com/danielborowski/cb-code-runner/releases/download/v0.1/runner /home/cb-code-runner/runner
RUN chmod +x /home/cb-code-runner/runner

USER cb-code-runner
WORKDIR /home/cb-code-runner
CMD ["/home/cb-code-runner/runner"]
ENTRYPOINT "/home/cb-code-runner/runner"
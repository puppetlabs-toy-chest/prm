# based upon the work in https://github.com/da-ar/pcv-test/
# this file is a mock up of what a generated file based on the prm-config.yml
# would have to contain in order to work.

# build with:
# docker build -t pdk:puppet-7.11.0-epp -f examples/tools/epp/Dockerfile .

# run with:
# docker run -v ${PWD}:/module -w /module pdk:puppet-7.11.0-epp

FROM puppet/puppet-agent:5.5.0
VOLUME [ "/code" ]
# The working directory must be where the code lives
WORKDIR /code

COPY ./content/* /tmp/

# use_entrypoint_script: true
ENTRYPOINT [ "/tmp/entrypoint.sh" ]
# default_args is empty, do not write CMD
# CMD []

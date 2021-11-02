# based upon the work in https://github.com/da-ar/pcv-test/
# this file is a mock up of what a generated file based on the prm-config.yml
# would have to contain in order to work.

# build with:
# docker build -t pdk:puppet-7.11.0-puppet_spec -f examples/tools/spec_puppet/Dockerfile .

# run with:
# needs_write_access: true
# docker run -v ${PWD}:/module -w /module pdk:puppet-7.11.0-puppet_spec

# gem - gems need ruby - so get puppet...
FROM puppet/puppet-agent:7.11.0
# requires_git: true
RUN apt update
RUN apt install git -y
# name: [puppetlabs_spec_helper, rspec-puppet-facts]
RUN /opt/puppetlabs/puppet/bin/gem install puppetlabs_spec_helper --no-document
RUN /opt/puppetlabs/puppet/bin/gem install rspec-puppet-facts --no-document
COPY examples/tools/spec_puppet/content/Rakefile /tmp/Rakefile
VOLUME [ "/module" ]

# executable: rake
ENTRYPOINT [ "/opt/puppetlabs/puppet/bin/rake" ,"-f", "/tmp/Rakefile"]
# default_args: [spec]
CMD [ "spec" ]

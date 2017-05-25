from setuptools import setup
from subprocess import check_output

with open('VERSION.txt', 'r') as content_file:
    version = content_file.read()

    setup(
        name='charon-client',
        version=version[1:],
        description='charon service grpc client library',
        url='http://github.com/piotrkowalczuk/charon',
        author='Piotr Kowalczuk',
        author_email='p.kowalczuk.priv@gmail.com',
        license='MIT',
        packages=['charonrpc'],
        install_requires=[
            'protobuf',
            'grpcio',
            'protobuf-ntypes',
            'protobuf-qtypes',
        ],
        zip_safe=False,
        keywords=['charon', 'grpc', 'authentication', 'authorization', 'service', 'client'],
        download_url='https://github.com/piotrkowalczuk/charon/archive/%s.tar.gz' % version.rstrip(),
      )
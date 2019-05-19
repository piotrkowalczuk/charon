#################
Charon Client
#################

TODO

*************
Example
*************

::

  $ pip install charon-client


.. code-block:: python

    from  github.com.piotrkowalczuk.charon.pb.rpc.charond.v1 import auth_pb2, auth_pb2_grpc
    import grpc

    charonChannel = grpc.insecure_channel('ADDRESS')
    auth = auth_pb2_grpc.AuthStub(charonChannel)
    try:
        res = auth.Login(auth_pb2.LoginRequest(
            username="USERNAME",
            password="PASSWORD",
        ))

        print "access token: %s" % res.value
    except grpc.RpcError as e:
        if e.code() == grpc.StatusCode.UNAUTHENTICATED:
            print "login failure, username and/or password do not match"
        else:
            print "grpc error: %s" % e
    except Exception as e:
        print "unexpected error: %s" % e



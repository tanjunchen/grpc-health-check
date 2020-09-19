package hello



func (s *server) Hello(helloReq *proto.HelloRequest, srv proto.GreetService_HelloServer) error {
	logrus.Infof("Server received an rpc request with the following parameter %v", helloReq.Hello)

	for i := 0; i<=10 ; i++ {
		resp := &proto.HelloResponse{
			Greet: fmt.Sprintf("Hello %s for %d time",helloReq.Hello, i),
		}

		srv.SendMsg(resp)

	}
	return nil
}
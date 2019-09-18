package main

func main() {

	bouncer_uri := URI{
		Address: "0.0.0.0",
		Port:    6667,
	}

	carousel := Server{
		URI: bouncer_uri,
	}

	carousel.Serve()
}

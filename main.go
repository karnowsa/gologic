package main

func main() {
	var admin = AdminServer{ipAdress: "127.0.0.1", port: 7001, username: "weblogic", password: "password123"}
	admin.init()
	admin.status(nil)
	admin.startAll()
}

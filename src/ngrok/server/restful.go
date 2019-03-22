package server

import (
	"encoding/json"
	"github.com/emicklei/go-restful"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"ngrok/server/ports_db"
)

// PortResource is the model for holding port information
type PortResource struct {
	Token string
	Port  string
}

// GET http://localhost:8000/v1/tokens
func (t PortResource) getTokens(request *restful.Request, response *restful.Response) {

	ports,err := ports_db.ReadAllFromDB()
	if err != nil {
		log.Println(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "port could not be found.")
	} else {
		response.WriteHeaderAndEntity(http.StatusOK, ports)
	}
}

// DELETE http://localhost:8000/v1/ports/1
func (t PortResource) removeToken(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("token-id")

	// jimmy: todo: disconnect the live connection and free the port resource
	//remove port id from port pool
	//{
	//	v, err := readFromDB(id)
	//	if err != nil {
	//		response.AddHeader("Content-Type", "text/plain")
	//		response.WriteErrorString(http.StatusInternalServerError,
	//			err.Error())
	//		return
	//	}
	//
	//	parts := strings.Split(v, ":")
	//	portPart := parts[len(parts)-1]
	//	if port, err := strconv.Atoi(portPart); err == nil {
	//		RemovePort(port)
	//	}
	//}


	err := ports_db.DeleteFromDB(id)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError,
			err.Error())
		return
	}

	response.WriteHeader(http.StatusOK)
}



// POST http://localhost:8000/v1/tokens/
func (t PortResource) createToken(request *restful.Request, response *restful.Response) {
	//log.Println(request.Request.Body)

	decoder := json.NewDecoder(request.Request.Body)
	var b PortResource
	err := decoder.Decode(&b)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError,
			err.Error())
		return
	}

	b.Port = "client-id-tcp:" + b.Port
	err = ports_db.WriteToDB(b.Token, b.Port)
	if err == nil {
		response.WriteHeaderAndEntity(http.StatusCreated, b)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError,
			err.Error())
	}
}

// Register adds paths and routes to container
func (t *PortResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/v1/tokens").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("").To(t.getTokens))
	ws.Route(ws.POST("").To(t.createToken))
	ws.Route(ws.DELETE("/{token-id}").To(t.removeToken))
	container.Add(ws)
}

func startRestful(port string){

	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	t := PortResource{}
	t.Register(wsContainer)
	log.Printf("start listening on localhost:8000")
	server := &http.Server{Addr: ":"+ port, Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
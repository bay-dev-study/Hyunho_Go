package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"nomad_coin/blockchain"

	"github.com/gorilla/mux"
)

const (
	port        string = ":4000"
	templateDir string = "templates/"
)

var templates *template.Template

var rootUrlWithPort string

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	data := homeData{"Home", blockchain.GetBlockchain().AllBlocks()}
	templates.ExecuteTemplate(rw, "home", data)
}

func Start(port int) {
	router := mux.NewRouter()

	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	router.HandleFunc("/", home)

	portInString := fmt.Sprintf(":%d", port)
	rootUrlWithPort = fmt.Sprintf("http://localhost%s", portInString)

	fmt.Printf("Explorer server listening on %s\n", rootUrlWithPort)
	log.Fatal(http.ListenAndServe(portInString, router))
}

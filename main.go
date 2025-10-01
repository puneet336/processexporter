package main
import (
       "github.com/puneets336/customexporter1/pkg/restapi"
       "github.com/puneets336/customexporter1/pkg/prom"
       "github.com/puneets336/customexporter1/pkg/aparser"
       "github.com/puneets336/customexporter1/pkg/yamlreader"
       "github.com/puneets336/customexporter1/pkg/commandrunner"
       "net/http"
       "log"
       )


func main() {
    _functionname:="main.main"
    _parsercfg:=aparser.SetupParser()
    //fmt.Printf("%+v\n",_parsercfg)
    //log.Println(_parsercfg.Configfilepath)

    _yamlcfg:=yamlreader.ReadYaml(_parsercfg.Configfilepath)
    //log.Println(_yamlcfg)


    _cmr:=commandrunner.CreateAndStartCommandRunner(_parsercfg,_yamlcfg)

    _mux:=http.NewServeMux()

    restapi.SetupRestHandler(_mux,_parsercfg,_yamlcfg,_cmr)

    prom.SetupPromHandler(_mux,_parsercfg,_yamlcfg,_cmr)

    log.Println(_functionname," started Listening on ",_parsercfg.ListenAddress," ....")
    log.Fatal(http.ListenAndServe(_parsercfg.ListenAddress,_mux) )
    log.Println(_functionname," Exiting  ",_parsercfg.ListenAddress," ....")

}

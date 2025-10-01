package yamlreader
import (
         "io/ioutil"
         "gopkg.in/yaml.v3"
         "log"
         "os"
       )


type YamlConfig struct {
    ListenAddress string
    Queryinterval int
    Commonservertags_cfg map[string]string
    Serverspecifictags_cfg  map[string]map[string]string
    Commonservertags_list map[string]bool
    Serverspecifictags_list map[string]bool
}



func ReadYaml(configfilepath string) *YamlConfig {
    _functionname := "yamlreader.ReadYaml"

    var _yamlconfig YamlConfig

    //If you have provided incorrect path
    if _,_err := os.Stat(configfilepath); _err != nil {
        log.Println(_functionname," WARN: configfile(",configfilepath , ")does not exist")
        _yamlconfig.ListenAddress=":41000"
        _yamlconfig.Queryinterval = 15
        return &_yamlconfig
    }


    //Read the entire file
    _yfile, _err1 := ioutil.ReadFile(configfilepath)
    if _err1 != nil {
          log.Fatal(_functionname," ERROR: while reading ",configfilepath," file")
          log.Fatal(_err1)
     }

    //make storage area and populate the data via Unmarshal
    // key for the data will always be  string type, value can be nested structure
    _data := make(map[string]interface{})
    _err2 := yaml.Unmarshal(_yfile, &_data)
    if _err2 != nil {
          log.Fatal(_functionname," ERROR: ",_err2)
     }



    //this variable holds list of all the tags provided under commonservertags section
    _commonservertags_list := make(map[string]bool)
    //this variable holds list of all the tags + values provided under commonservertags section
    var _commonservertags_cfg = map[string]string{}


    //this variable holds list of all the server specific tags provided under serverspecific section
    _serverspecifictags_list := make(map[string]bool)
    //this variable holds list of all the server specific tags + respective value provided under serverspecific section
    var _serverspecifictags_cfg = map[string]map[string]string{}



    for _k, _v := range _data {
          //obtaining the query interval value from the Yaml
          if _k == "queryinterval" {
              _yamlconfig.Queryinterval = _v.(int)
              continue
          }

          //obtain value of listen address if Provided
          if _k == "listenaddress" {
               _yamlconfig.ListenAddress = _v.(string)
           }

          //obtaining the common server tags
          //
          if _k == "commonservertags" {

              //format of the value variable v is [tag1:val1 tag2:val2]
              // need to further iterate this to ket hold of tag1,tag2 and val1,val2 , as key is a string so create another map with key as string type
              _v_convert:=_v.(map[string]interface{})
              for _k1,_v1 := range _v_convert{

                  _commonservertags_cfg[_k1]=_v1.(string)
                  _commonservertags_list[_k1] = true

              }//for loop ends

          }//if condition for commonservertags

          // if we are dealing with serverspecific tags section
          // we have to deal with 2 more nested levels -  first level will have compute server names as key and a dictionary ( source of level 2 data)
          //                                              second level will have tag name as key and tag value .
          if _k == "serverspecifictags" {
              //v=[server1:val1 server2:val2]
              _v_convert:=_v.(map[string]interface{})

              for _k1,_v1 := range _v_convert {
                  //den-c-001 map[tag1:val10]
                  _v1_convert:=_v1.(map[string]interface{})
                  _serverspecifictags_cfg[_k1]=map[string]string{}

                  for _k2,_v2 := range _v1_convert {
                      _serverspecifictags_cfg[_k1][_k2]=_v2.(string)
                      _serverspecifictags_list[_k2] = true
                  } //for loop for v1_convert

              } // for loop for v_convert
          } //if condition for serverspecifictags
     } // for loop ends

    // were you able to get query interval?
    if  _yamlconfig.ListenAddress == "" {
        _yamlconfig.ListenAddress=":41000"
    }
    // were you able to get ListenAddress?
    if _yamlconfig.Queryinterval == 0 {
        _yamlconfig.Queryinterval = 15
    }
    // is common server tags mentioned?

    // is server specific tags mentioned?


     // lets set all the values collected as property of struct and return that value.

     _yamlconfig.Commonservertags_cfg = _commonservertags_cfg
     _yamlconfig.Commonservertags_list = _commonservertags_list
     _yamlconfig.Serverspecifictags_cfg=_serverspecifictags_cfg
     _yamlconfig.Serverspecifictags_list=_serverspecifictags_list

     return &_yamlconfig
}

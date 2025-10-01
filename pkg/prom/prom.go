package prom


import  (
        "fmt"
        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promhttp"
        "github.com/puneets336/customexporter1/pkg/aparser"
        "github.com/puneets336/customexporter1/pkg/yamlreader"
        "github.com/puneets336/customexporter1/pkg/commandrunner"
        "net/http"
        "time"
        "os"
        "log"
        "regexp"
        "strings"
        "strconv"
        )


type LSFService struct {
    sbatchd bool;
    lim bool;
    pim bool;

}

var _cmr *commandrunner.CommandRunner
var _yamlcfg *yamlreader.YamlConfig
var _parsercfg *aparser.ParserCfg



func updateData(_parsercfg_local *aparser.ParserCfg,_yamlcfg_local *yamlreader.YamlConfig,_gaugeMap map[string]*prometheus.GaugeVec,_all_labels map[string][]string,_thisserver_map map[string]string) {
       _cmr.Lock()
       defer _cmr.Unlock()
       _gaugeMap["ntp_drift"].Reset()
       _gaugeMap["lsf_service"].Reset()
       _gaugeMap["localdisk_storage"].Reset()

       lim_status(_gaugeMap["lsf_service"],_all_labels["lsf_service"],_thisserver_map)
       res_status(_gaugeMap["lsf_service"],_all_labels["lsf_service"],_thisserver_map)
       sbatchd_status(_gaugeMap["lsf_service"],_all_labels["lsf_service"],_thisserver_map)
       pim_status(_gaugeMap["lsf_service"],_all_labels["lsf_service"],_thisserver_map)
       bld_status(_gaugeMap["lsf_service"],_all_labels["lsf_service"],_thisserver_map)
       mbatchd_status(_gaugeMap["lsf_service"],_all_labels["lsf_service"],_thisserver_map)
       mbschd_status(_gaugeMap["lsf_service"],_all_labels["lsf_service"],_thisserver_map)
       blcollect_status(_gaugeMap["lsf_service"],_all_labels["lsf_service"],_thisserver_map)
       licpollerd_status(_gaugeMap["lsf_service"],_all_labels["lsf_service"],_thisserver_map)
       lsfpollerd_status(_gaugeMap["lsf_service"],_all_labels["lsf_service"],_thisserver_map)
       crond_status(_gaugeMap["lsf_service"],_all_labels["lsf_service"],_thisserver_map)

       ntp_drift_status( _parsercfg_local,_yamlcfg_local,_gaugeMap["ntp_drift"],_all_labels["ntp_drift"],_thisserver_map   )


       //localmount_status(_gaugeMap["localdisk_storage"],_all_labels["localdisk_storage"],_thisserver_map)
}


func recordMetrics(_parsercfg_local *aparser.ParserCfg,_yamlcfg_local *yamlreader.YamlConfig ,_gaugeMap map[string]*prometheus.GaugeVec,_all_labels map[string][]string,_thisserver_map map[string]string) {

    var _query_interval int
    if _yamlcfg.Queryinterval == 0 {
        _query_interval=_yamlcfg.Queryinterval
    } else {
        _query_interval=15
    }


    for {
       //gauge..Set( float64(rand.Intn(100)) )
       //fmt.Println("Running..")
       updateData(_parsercfg_local,_yamlcfg_local,_gaugeMap,_all_labels,_thisserver_map)
       time.Sleep(   time.Duration(_query_interval) * time.Second)
    }
}


func appendThisServerMap(_newmap map[string]string, _thisserver_map map[string]string)  {
    for _k,_v := range _thisserver_map {
        _newmap[_k] = _v
    }

}


func lim_status(gauge *prometheus.GaugeVec,_all_labels []string,_thisserver_map map[string]string)   {
    _stdout:=_cmr.GetLim_status_stdout()
    _users:=[]string{"root","lsfadmin"}

    _labels_map:=make(map[string]string)
    _labels_map[  _all_labels[0]    ] = "lim"

    //check if root owns lim


    appendThisServerMap(_labels_map,_thisserver_map)

    for _,_user := range _users {
        //_pattern:=\s+_user+`\s+.*/lim `
        _pattern:=\s+_user+\s+/tools/.*/etc/lim(?:\s|\n|$)
        _re:=regexp.MustCompile(_pattern)
        _matches := _re.FindAllString( _stdout , -1)
        _labels_map[  _all_labels[1]    ] = _user
        gauge.With(_labels_map).Set(   float64(len(_matches))    )
    }

}


func res_status(gauge *prometheus.GaugeVec,_all_labels []string,_thisserver_map map[string]string)   {
     _functionname:="prom.res_status"
     _users:=[]string{"root","lsfadmin"}
     if _functionname == "" {
         log.Println("promt.res_status ","_functionname variable is blank")
     }

     _stdout:=_cmr.GetRes_status_stdout()

     _labels_map:=make(map[string]string)
     _labels_map[  _all_labels[0]   ] = "res"
     appendThisServerMap(_labels_map,_thisserver_map)

    for _,_user := range _users {
        _pattern:=\s+_user+\s+/tools/.*/etc/res(?:\s|\n|$)
        _re:=regexp.MustCompile(_pattern)
        _matches := _re.FindAllString( _stdout , -1)
        _labels_map[  _all_labels[1]    ] = _user
         gauge.With(_labels_map).Set(   float64(len(_matches))    )
    }

}
func sbatchd_status(gauge *prometheus.GaugeVec,_all_labels []string,_thisserver_map map[string]string)   {
     _stdout:=_cmr.GetSbatchd_status_stdout()
     _users:=[]string{"root","lsfadmin"}
     _labels_map:=make(map[string]string)
     _labels_map[  _all_labels[0]   ] = "sbatchd"
     appendThisServerMap(_labels_map,_thisserver_map)

    for _,_user := range _users {
        _pattern:=\s+_user+\s+/tools/.*/etc/sbatchd(?:\s|\n|$)
        _re:=regexp.MustCompile(_pattern)

        _matches := _re.FindAllString( _stdout , -1)
        _labels_map[  _all_labels[1]    ] = _user
         gauge.With(_labels_map).Set(   float64(len(_matches))    )
    }

}

func pim_status(gauge *prometheus.GaugeVec,_all_labels []string,_thisserver_map map[string]string)   {
    _functionname:="prom.pim_status"
    _users:=[]string{"root","lsfadmin"}
    _stdout:=_cmr.GetPim_status_stdout()

     if _functionname == "" {
         log.Println("promt.pim_status ","_functionname variable is blank")
     }

    _labels_map:=make(map[string]string)
    _labels_map[    _all_labels[0]    ] = "pim"
    appendThisServerMap(_labels_map,_thisserver_map)

    for _,_user := range _users {
        _pattern:=\s+_user+\s+/tools/.*/etc/pim(?:\s|\n|$)
        _re:=regexp.MustCompile(_pattern)
        _matches := _re.FindAllString( _stdout , -1)
        _labels_map[  _all_labels[1]    ] = _user
        gauge.With(_labels_map).Set(   float64(len(_matches))    )
    }
}

func bld_status(gauge *prometheus.GaugeVec,_all_labels []string,_thisserver_map map[string]string)   {
    _functionname:="prom.bld_status"
    _users:=[]string{"root","lsfadmin"}
    _stdout:=_cmr.GetBld_status_stdout()
     if _functionname == "" {
         log.Println("prom.bld_status ","_functionname variable is blank")
     }

    _labels_map:=make(map[string]string)
    _labels_map[ _all_labels[0]   ]   = "bld"
    appendThisServerMap(_labels_map,_thisserver_map)

    for _,_user := range _users {
        _pattern:=\s+_user+\s+/tools/.*/etc/bld(?:\s|\n|$)
        _re:=regexp.MustCompile(_pattern)
        _matches := _re.FindAllString( _stdout , -1)
        _labels_map[  _all_labels[1]    ] = _user
        gauge.With(_labels_map).Set(   float64(len(_matches))    )
    }

}


func mbschd_status(gauge *prometheus.GaugeVec,_all_labels []string,_thisserver_map map[string]string)   {
    _functionname:="prom.mbschd_status"
    _users:=[]string{"root","lsfadmin"}
    _stdout:=_cmr.GetMbschd_status_stdout()
     if _functionname == "" {
         log.Println("prom.mbschd_status ","_functionname variable is blank")
     }

    _labels_map:=make(map[string]string)
    _labels_map[ _all_labels[0]   ]   = "mbschd"
    appendThisServerMap(_labels_map,_thisserver_map)
    for _,_user := range _users {
        _pattern:=\s+_user+\s+/tools/.*/etc/mbschd(?:\s|\n|$)
        _re:=regexp.MustCompile(_pattern)
        _matches := _re.FindAllString( _stdout , -1)
        _labels_map[  _all_labels[1]    ] = _user
        gauge.With(_labels_map).Set(   float64(len(_matches))    )
    }

}

func mbatchd_status(gauge *prometheus.GaugeVec,_all_labels []string,_thisserver_map map[string]string)   {
    _functionname:="prom.mbatchd_status"
    _users:=[]string{"root","lsfadmin"}
    _stdout:=_cmr.GetMbatchd_status_stdout()
     if _functionname == "" {
         log.Println("prom.mbatchd_status ","_functionname variable is blank")
     }

    _labels_map:=make(map[string]string)
    _labels_map[ _all_labels[0]   ]   = "mbatchd"
    appendThisServerMap(_labels_map,_thisserver_map)

    for _,_user := range _users {
        _pattern:=\s+_user+\s+/tools/.*/etc/mbatchd(?:\s|\n|$)
        _re:=regexp.MustCompile(_pattern)
        _matches := _re.FindAllString( _stdout , -1)
        _labels_map[  _all_labels[1]    ] = _user
        gauge.With(_labels_map).Set(   float64(len(_matches))    )
    }

}

func blcollect_status(gauge *prometheus.GaugeVec,_all_labels []string,_thisserver_map map[string]string)   {
    _functionname:="prom.blcollect_status"
    _users:=[]string{"root","lsfadmin"}
    _stdout:=_cmr.GetBlcollect_status_stdout()
     if _functionname == "" {
         log.Println("prom.blcollect_status ","_functionname variable is blank")
     }

    _labels_map:=make(map[string]string)
    _labels_map[ _all_labels[0]   ]   = "blcollect"
    appendThisServerMap(_labels_map,_thisserver_map)

    for _,_user := range _users {
        _pattern:=\s+_user+\s+/tools/.*/etc/blcollect(?:\s|\n|$)
        _re:=regexp.MustCompile(_pattern)
        _matches := _re.FindAllString( _stdout , -1)
        _labels_map[  _all_labels[1]    ] = _user
        gauge.With(_labels_map).Set(   float64(len(_matches))    )
    }

}

func lsfpollerd_status(gauge *prometheus.GaugeVec,_all_labels []string,_thisserver_map map[string]string)   {
    _functionname:="prom.lsfpollerd_status"
    _users:=[]string{"root","lsfadmin","apache"}
    _stdout:=_cmr.GetLsfpollerd_status_stdout()
     if _functionname == "" {
         log.Println("prom.lsfpollerd_status ","_functionname variable is blank")
     }

    _labels_map:=make(map[string]string)
    _labels_map[ _all_labels[0]   ]   = "lsfpollerd"
    appendThisServerMap(_labels_map,_thisserver_map)

    for _,_user := range _users {
        _pattern:=\s+_user+\s+/opt/.*/bin/lsfpollerd(?:\s|\n|$)
        _re:=regexp.MustCompile(_pattern)
        _matches := _re.FindAllString( _stdout , -1)
        _labels_map[  _all_labels[1]    ] = _user
        gauge.With(_labels_map).Set(   float64(len(_matches))    )
    }

}
func licpollerd_status(gauge *prometheus.GaugeVec,_all_labels []string,_thisserver_map map[string]string)   {
    _functionname:="prom.licpollerd_status"
    _users:=[]string{"root","lsfadmin","apache"}
    _stdout:=_cmr.GetLicpollerd_status_stdout()
     if _functionname == "" {
         log.Println("prom.licpollerd_status ","_functionname variable is blank")
     }

    _labels_map:=make(map[string]string)
    _labels_map[ _all_labels[0]   ]   = "licpollerd"
    appendThisServerMap(_labels_map,_thisserver_map)

    for _,_user := range _users {
        _pattern:=\s+_user+\s+/opt/.*/bin/licpollerd(?:\s|\n|$)
        _re:=regexp.MustCompile(_pattern)
        _matches := _re.FindAllString( _stdout , -1)
        _labels_map[  _all_labels[1]    ] = _user
        gauge.With(_labels_map).Set(   float64(len(_matches))    )
    }

}
func crond_status(gauge *prometheus.GaugeVec,_all_labels []string,_thisserver_map map[string]string)   {
    _functionname:="prom.crond_status"
    _users:=[]string{"root"}
    _stdout:=_cmr.GetCrond_status_stdout()
     if _functionname == "" {
         log.Println("prom.crond_status ","_functionname variable is blank")
     }

    _labels_map:=make(map[string]string)
    _labels_map[ _all_labels[0]   ]   = "crond"
    appendThisServerMap(_labels_map,_thisserver_map)

    for _,_user := range _users {
        _labels_map[  _all_labels[1]    ] = _user
        _isrunning, _ := regexp.MatchString(\sroot\s*.*\scron, _stdout)
        if _isrunning {
             gauge.With(_labels_map).Set( 1.0 )
         } else {
             gauge.With(_labels_map).Set( 0.0 )
         }
    }
}




////////////////////////CHRONYD


func ntp_chronyd_gettimediff() float64 {
    _stdout:=_cmr.GetNtp_Chronyd_status_stdout()
    _re:=regexp.MustCompile(System time     : \d+\.\d+)
    _data1 := _re.FindStringSubmatch(_stdout)

    if len(_data1) == 0 {
        return -1.11223344
    }

    _data2:=strings.Split( _data1[0],":")[1]
    data3,:=strconv.ParseFloat(strings.TrimSpace(_data2), 64)
    isrunning, := regexp.MatchString(" slow", _stdout)
    if _isrunning {
        _data3=-_data3
        return _data3
    } else {
        return _data3
    }

}


func ntp_chronyd_getserver() string {
    //Reference ID    : C6D367D1 (mailhost.netads.com)
    _stdout:=_cmr.GetNtp_Chronyd_status_stdout()
    _pattern := regexp.MustCompile(".Reference ID.\\((\\S+)\\)\\nStratum.*")
    _subStr := _pattern.FindStringSubmatch( _stdout )

    if len(_subStr) > 1 {
        return _subStr[1]
    } else {
        return "notfound"
    }

}


func ntp_chronyd_getsync() string {
    _stdout:=_cmr.GetNtp_Chronyd_status_stdout()
    isrunning, :=regexp.MatchString("NTP synchronized: yes|System clock synchronized: yes",_stdout)
    if _isrunning {
        return "yes"
    } else {
        return "no"
    }

}



func ntp_chronyd_drift_status(gauge *prometheus.GaugeVec,_all_labels []string,_thisserver_map map[string]string) {
    _functionname:="prom.ntp_chronyd_drift_status"
     if _functionname == "" {
         log.Println("prom.ntp_chronyd_drift_status ","_functionname variable is blank")
     }

    _labels_map:=make(map[string]string)

    //_all_labels_ntpdrift:=[]string{"NTPSERVER","ISSYNCED"}
    _labels_map[ _all_labels[0]   ]   = ntp_chronyd_getserver()
    _labels_map[ _all_labels[1]   ]   = ntp_chronyd_getsync()
    _value:=ntp_chronyd_gettimediff()


    appendThisServerMap(_labels_map,_thisserver_map)

    gauge.With(_labels_map).Set( _value )


}

///////////////////////NTPD

func ntp_ntpd_gettimediff() float64 {
    //synchronised to NTP server (198.211.103.209) at stratum 3
    //time correct to within 49 ms
    //polling server every 512 s

    _stdout:=_cmr.GetNtp_Ntpd_status_stdout()
    _re:=regexp.MustCompile(time correct to within (\d+) m.*)
    _data1 := _re.FindStringSubmatch(_stdout)


    if len(_data1) == 0 {
        return -0.11223344
    }


    data2,:=strconv.ParseFloat(strings.TrimSpace(_data1[1]), 64)
    _data2=_data2/1000

   return _data2
}


func ntp_ntpd_getserver() string {
    //synchronised to NTP server (198.211.103.209) at stratum 3
    _stdout:=_cmr.GetNtp_Ntpd_status_stdout()
    _pattern := regexp.MustCompile(".ynchronised to NTP server.\\((\\S+)\\) at stratum.*")
    _subStr := _pattern.FindStringSubmatch( _stdout )


    if len(_subStr) > 1 {
        return _subStr[1]
    } else {
        return "notfound"
    }

}


func ntp_ntpd_getsync() string {
    _stdout:=_cmr.GetNtp_Ntpd_status_stdout()
    isrunning, :=regexp.MatchString(".ynchronised to NTP server.",_stdout)
    if _isrunning {
        return "yes"
    } else {
        return "no"
    }
}



func ntp_drift_status(_parsercfg_local *aparser.ParserCfg,_yamlcfg_local *yamlreader.YamlConfig, gauge *prometheus.GaugeVec,_all_labels []string,_thisserver_map map[string]string) {
    _functionname:="prom.ntp_drift_status"
     if _functionname == "" {
         log.Println("prom.ntp_drift_status ","_functionname variable is blank")
     }

    _labels_map:=make(map[string]string)

    //_all_labels_ntpdrift:=[]string{"NTPSERVER","ISSYNCED"}

    //check if ntpd command returns a server name, If not, lets fallback to chronyd
    _servername:=ntp_ntpd_getserver()
    if _servername != "notfound" {

        _labels_map[ _all_labels[0]   ]   = ntp_ntpd_getserver()
        _labels_map[ _all_labels[1]   ]   = ntp_ntpd_getsync()
        _value:=ntp_ntpd_gettimediff()

        appendThisServerMap(_labels_map,_thisserver_map)

        gauge.With(_labels_map).Set( _value )
        if _parsercfg_local.GetDebugMode() > 0 {
            log.Println("NTPD VALUE:",_value)
        }


    } else {

        _labels_map[ _all_labels[0]   ]   = ntp_chronyd_getserver()
        _labels_map[ _all_labels[1]   ]   = ntp_chronyd_getsync()
        _value:=ntp_chronyd_gettimediff()

        appendThisServerMap(_labels_map,_thisserver_map)
        gauge.With(_labels_map).Set( _value )

        if _parsercfg_local.GetDebugMode() > 0 {
            log.Println("Chronyd VALUE:",_value)
        }

    }



}





func localmount_status(gauge *prometheus.GaugeVec,_all_labels []string,_thisserver_map map[string]string)   {
     _functionname:="prom.localmount_status"
     _stdout:=_cmr.GetLocalmountpoint_status_stdout()
     //Filesystem                     1K-blocks     Used Available Use% Mounted on
     //devtmpfs                       495271232        0 495271232   0% /dev
     //tmpfs                          495284404        0 495284404   0% /dev/shm
     // "FILESYSTEM","TYPE","MOUNTED_ON"

     if _functionname == "" {
         log.Println("prom.localmount_status ","_functionname variable is blank")
     }
     _stdout_list:=strings.Split(_stdout, "\n")

     for _,_line := range _stdout_list {
         _ismatch, _ := regexp.MatchString("/run/user|Filesystem",_line)
         if _ismatch {

             continue
         } else {

             _labels_map:=make(map[string]string)
             for _k,_v := range _thisserver_map {
                 _labels_map[_k]=_v
             }


             _words:=strings.Fields(_line)
             //log.Println("-->",_index,_words,_words[4])
             //[devtmpfs 495271232 0 495271232 0% /dev]
             //_all_labels[0] = FILESYSTEM = devtmpfs
             _labels_map[ _all_labels[0] ] = _words[0]

             //_all_labels[2] = MOUNTED_ON = /dev/shm
             _labels_map[ _all_labels[2] ] = _words[5]

             _labels_map[ "TYPE" ] = "TOTAL"
             total_bytes,:=strconv.ParseFloat(_words[1],64)
             gauge.With(_labels_map).Set(  _total_bytes  )

             _labels_map[ "TYPE" ] = "AVAIL"
             used_bytes,:=strconv.ParseFloat(_words[3],64)
             gauge.With(_labels_map).Set(  _used_bytes  )

         }

     }



}

func getHostname() string{
     _hostnamefqdn, err := os.Hostname()
    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }

    _hostname:=strings.Split(_hostnamefqdn,".")[0]
    return _hostname
}

//returns 2 arguments
//1 - list of unique tags which are mentioned in the yaml file
//2 - map with - all the tags as keys
//               all the values applicable to this server as value, if some values are not applicable for this server, then value is "NA"

func getServertagsfromyamlfile() ([]string,map[string]string) {
    //first get the common labels which is applicable for all servers
    //Commonservertags_list
    _hostname:=getHostname()

    _yamlfile_labels_map1 := make(map[string]bool)

    for k1 , v1 := range _yamlcfg.Serverspecifictags_list {
        _yamlfile_labels_map1[k1]=v1
    }



    for k2 , v2 := range _yamlcfg.Commonservertags_list {
        _yamlfile_labels_map1[k2]=v2
    }

    // we may have repeated tags in common or server specific section, since we dont have concept of set here, we use  _label(servername,bool) to get uniqueservernames
    //get list of keys from the total k-v pair (global + server specific)
    //_yamlfile_labels_map1 = dict_common_server_tags + dict_specific_server_tag
    //_yamlfile_labels_map1 is a map  with bool type as value, so value not required.
    i1:=0
    _all_keys_yamlfile := make([]string, len(_yamlfile_labels_map1))
    for k,_ := range _yamlfile_labels_map1 {
        _all_keys_yamlfile[i1] = k
        i1++
    }

    //Objective 1 finished, we have obtained list of keys mentioned in yaml fil - _all_keys_yamlfile - this will be used to create a gauge with available properties/labels
    //Objective 2 - generate a map of tags which are valid for current server.
    //i.e. yaml file has -
    //commonservertags:
    //    key1:val1
    //serverspecifictags :
    //  den-c-003:
    //    key2:val2
    //  den-c-001:
    //    key3:val3
    //in this scenario,
    //    [key1:val1  key2:val2 key3:NA] is the expected data if this code is running on den-c-003,
    //this will be used when we will set value of metrics on this server.


    //dictionary of key value pairs
    _yamlfile_labels_map2:=make(map[string]string)

    //loop over common server tags , and populate map/dictionary with the data discovered from yaml file
    for k3 , v3 := range _yamlcfg.Commonservertags_cfg {
        _yamlfile_labels_map2[k3]=v3
    }


    //next up, we will add more k-v pairs in the dictionary if we have more tags for the server on thic this exporter is running. i e. if current server is den-c-001 and yaml file has -
    //den-c-003:
    //    tag1: value1
    //in this case, we will add tag1:value1 also in the _yamlfile_labels_map2

    //den-c-003 map[tag1:val11 tag2:val12] <--- select this element and merge elements with _yamlfile_labels_map2 created above.
    //den-c-001 map[tag1:val10] den-l-013
    thisserverspecifictags_cfg,:=_yamlcfg.Serverspecifictags_cfg[_hostname]
    for k4 , v4 := range _thisserverspecifictags_cfg {
        _yamlfile_labels_map2[k4]=v4
    }




   // in _yamlfile_labels_map2 we have -
   // a) keys/value pairs from  commonservertags section - done
   // b) keys/value pairs from the serverspecifictags section for this server - done
   // c) keys/value pairs from the serverspecifictags section for other servers - to be done
   // we need all the labels (_all_keys_yamlfile) , only then we can set the gauge values for this server correctliy.value will be "NA".

   for _k5 := range _all_keys_yamlfile {
       //_k5 is index
       // first value is value, second is the status of the lookup
       _, ok5 := _yamlfile_labels_map2[  _all_keys_yamlfile[_k5] ]
       if  ok5 {
           continue
       } else {
           _yamlfile_labels_map2[ _all_keys_yamlfile[_k5] ]="NA"
       }

   }


    _all_keys_yamlfile=append(_all_keys_yamlfile,"HOSTNAME")
    _yamlfile_labels_map2["HOSTNAME"] = _hostname

    return _all_keys_yamlfile , _yamlfile_labels_map2
}



func setupAndRegisterGauge() (map[string]*prometheus.GaugeVec,*prometheus.Registry,map[string][]string,map[string]string) {

    //_all_labels_yamlfile = list of labels mentioned in yaml file
    //_labels_thisserver_map = list of kv pairs applicable for this server
    _all_labels_yamlfile,_labels_thisserver_map:=getServertagsfromyamlfile()



    _all_labels_lsfsvc:=[]string{"SERVICETYPE","USER"}
    _all_labels_lsfsvc=append(_all_labels_lsfsvc,_all_labels_yamlfile...)

    _all_labels_localdisk:=[]string{"FILESYSTEM","TYPE","MOUNTED_ON"}
    _all_labels_localdisk=append(_all_labels_localdisk,_all_labels_yamlfile...)

    _all_labels_ntpdrift:=[]string{"NTPSERVER","ISSYNCED"}
    _all_labels_ntpdrift=append(_all_labels_ntpdrift,_all_labels_yamlfile...)

    // dont need to worry about setting the  to value for _all_labels_yamlfile , adn you can use map _labels_thisserver_map. just worry about the labels you are going to apply other than the ones you have mentioned in yamlfile.
    _allLabelsMap:=make( map[string][]string   )
    _allLabelsMap["lsf_service"]=_all_labels_lsfsvc
    _allLabelsMap["localdisk_storage"]=_all_labels_localdisk
    _allLabelsMap["ntp_drift"]=_all_labels_ntpdrift




    _gaugeMap := make( map[string]*prometheus.GaugeVec )
    //gaugeVecs := make(map[string]*prometheus.GaugeVec)

    promReg := prometheus.NewRegistry()


    var  gauge1 = prometheus.NewGaugeVec(
                prometheus.GaugeOpts{
                        Namespace: "go",
                        Name:      "lsf_service_status",
                        Help:      "Gauge for the current status of lim pim sbatchd",
                },_all_labels_lsfsvc)


    promReg.MustRegister(gauge1)

    var  gauge2 = prometheus.NewGaugeVec(
                prometheus.GaugeOpts{
                        Namespace: "go",
                        Name:      "localdisk_storage_status",
                        Help:      "Gauge for the current status of local disk",
                },_all_labels_localdisk)

    promReg.MustRegister(gauge2)

    var  gauge3 = prometheus.NewGaugeVec(
                prometheus.GaugeOpts{
                        Namespace: "go",
                        Name:      "ntp_drift_status",
                        Help:      "Gauge for the current status of ntp time sync via chronyd",
                },_all_labels_ntpdrift)

    promReg.MustRegister(gauge3)




    _gaugeMap["lsf_service"]=gauge1
    _gaugeMap["localdisk_storage"]=gauge2
    _gaugeMap["ntp_drift"]=gauge3



    return _gaugeMap,promReg,_allLabelsMap,_labels_thisserver_map

}


func SetupPromHandler(mux *http.ServeMux,_parsercfg_local *aparser.ParserCfg ,_yamlcfg_local *yamlreader.YamlConfig,_cmr_local *commandrunner.CommandRunner) {

     _functionname:="prom.SetupPromHandler"
     _cmr=_cmr_local
     _yamlcfg=_yamlcfg_local
     _parsercfg=_parsercfg_local

     if _cmr == nil || _yamlcfg == nil || _parsercfg == nil || _cmr_local == nil || _yamlcfg_local == nil || _parsercfg_local == nil {
         log.Println(_functionname,"one of the supplied parameter is blank")
     }

    //prometheus.MustRegister(gauge)
    _gaugeMap,promReg,_allLabelsMap,_labels_thisserver_map:=setupAndRegisterGauge()
    handler := promhttp.HandlerFor(promReg, promhttp.HandlerOpts{})





    go recordMetrics(_parsercfg_local,_yamlcfg_local,_gaugeMap,_allLabelsMap,_labels_thisserver_map)

    //mux.Handle("/v1/metrics", promhttp.Handler())
    mux.Handle("/v1/processes/prom", handler)
    fmt.Println(_functionname," registered handler at /v1/processes/prom")

}

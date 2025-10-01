package restapi

import (
        "net/http"
        "log"
        "os"
        "time"
        "encoding/json"
        "regexp"
        "strings"
        "strconv"
        "github.com/puneets336/customexporter1/pkg/aparser"
        "github.com/puneets336/customexporter1/pkg/yamlreader"
        "github.com/puneets336/customexporter1/pkg/commandrunner"
      )




type ServerProcessesData struct {
    SBATCHD int json:"SBATCHD"
    RES int json:"RES"
    BLD int json:"BLD"
    LIM int json:"LIM"
    PIM int json:"PIM"
    MBATCHD int json:"MBATCHD"
    BLCOLLECT int json:"BLCOLLECT"
    MBSCHD int json:"MBSCHD"
    LICPOLLERD int json:"LICPOLLERD"
    LSFPOLLERD int json:"LSFPOLLERD"
    CROND int json:"CROND"
    CORES int json:"CORES"
    MEMORY int json:"MEMORY"


    TIMEDIFF_NTPD float64 json:"TIMEDIFF_NTPD"
    SERVER_NTPD string json:"SERVER_NTPD"
    SYNCED_NTPD int json:"SYNCED_NTPD"

    TIMEDIFF_CHRONYD float64 json:"TIMEDIFF_CHRONYD"
    SERVER_CHRONYD string json:"SERVER_CHRONYD"
    SYNCED_CHRONYD int json:"SYNCED_CHRONYD"

    PROCESSES *[]ProcessInfo json:"PROCESSES"
    LOCALMOUNTPOINTS *[]LocalmountpointInfo json:"LOCALMOUNTPOINTS"
    TIMESTAMP int64 json:"TIMESTAMP"
    HOSTNAME string  json:"HOSTNAME"
    DISTRIBUTION string json:"DISTRIBUTION"
}


type ProcessInfo struct {
    PID string json:"PID"
    PPID string json:"PPID"
    USER string json:"USER"
    NI string json:"NI"
    RSS string json:"RSS"
    SZ string json:"SZ"
    STAT string json:"STAT"
    MEM string json:"MEM"
    CPU string json:"CPU"
    TIME string json:"TIME"
    COMM string json:"COMM"
    CMD string json:"CMD"
}

type LocalmountpointInfo struct {
    FILESYSTEM string json:"FILESYSTEM"
    MOUNTED_ON string json:"MOUNTED_ON"
    AVAIL float64  json:"AVAIL"
    TOTAL float64  json:"TOTAL"
}



var _cmr *commandrunner.CommandRunner
var _yamlcfg *yamlreader.YamlConfig
var _parsercfg *aparser.ParserCfg


func parseRunningProcessData() ([]ProcessInfo,error)  {
    _functionname:="restapi.parseRunningProcessData"
    //ps -e -o 'pid ppid user:30 ni rss sz stat %mem %cpu time comm cmd'
    if _functionname != "" {
        log.Println(_functionname)
    }


    var _processes []ProcessInfo
            //TIMESTAMP: _time,

    _lines := strings.Split( _cmr.GetProcesses_status_stdout() ,"\n" )



    for _,_line := range _lines {
        _fields:=strings.Fields(_line)
        _process:= ProcessInfo{
            PID: _fields[0],
            PPID: _fields[1],
            USER: _fields[2],
            NI: _fields[3],
            RSS: _fields[4],
            SZ: _fields[5],
            STAT: _fields[6],
            MEM: _fields[7],
            CPU: _fields[8],
            TIME: _fields[9],
            COMM: _fields[10],
            CMD: strings.Join(_fields[11:]," "),
        }
        _processes=append(_processes,_process)
    }//forloop
    return _processes,nil
}


func parseLocalmountpointData() ([]LocalmountpointInfo,error)  {
    _functionname:="restapi.gparseLocalmountpointData"
    //ps -e -o 'pid ppid user:30 ni rss sz stat %mem %cpu time comm cmd'
    if _functionname != "" {
        log.Println(_functionname)
    }

    _lines := strings.Split( _cmr.GetLocalmountpoint_status_stdout() ,"\n" )

    var _localmountpointsinfo []LocalmountpointInfo
            //TIMESTAMP: _time,




    for _,_line := range _lines {
        _fields:=strings.Fields(_line)
        match,:=regexp.MatchString("Filesystem ",_line)
        if _match {
            continue
        }

        _avail,_err1:=strconv.ParseFloat(_fields[2] , 64)
        if _err1 != nil {
            log.Println("WARN: unable to parse the available storage field from ",_line,_err1)
            _avail=0.0
        }

        _total,_err2:=strconv.ParseFloat(_fields[3] , 64)
        if _err2 != nil {
            log.Println("WARN: unable to parse the total storage field",_line,_err2)
            _total=0.0

        }

        //Filesystem                     1K-blocks     Used Available Use% Mounted on
       //devtmpfs                       495271232        0 495271232   0% /dev
       //tmpfs                          495284404        0 495284404   0% /dev/shm
        _localmountpointinfo:= LocalmountpointInfo{
            FILESYSTEM: _fields[0],
            MOUNTED_ON: _fields[5],
            AVAIL: _avail,
            TOTAL: _total,
        }
        _localmountpointsinfo=append(_localmountpointsinfo,_localmountpointinfo)
    }//forloop
    return _localmountpointsinfo,nil
}





func lim_status() int   {
    _stdout:=_cmr.GetLim_status_stdout()
    //_pattern:="/tools/ibm/lsf/10.1/install/10.1/linux3.10-glibc2.17-x86_64/etc/lim"
     _pattern:=\s(root|lsfadmin)\s+/tools/ibm/lsf/10.1/install/10.1/linux3.10-glibc2.17-x86_64/etc/lim(?:\s|\n|$)

    _re:=regexp.MustCompile(_pattern)
    _matches := _re.FindAllString( _stdout , -1)
    return len(_matches)

}

func res_status() int   {
    _stdout:=_cmr.GetRes_status_stdout()
    //_pattern:="/tools/ibm/lsf/10.1/install/10.1/linux3.10-glibc2.17-x86_64/etc/res"
    _pattern:=\s(root|lsfadmin)\s+/tools/ibm/lsf/10.1/install/10.1/linux3.10-glibc2.17-x86_64/etc/res(?:\s|\n|$)
    _re:=regexp.MustCompile(_pattern)
    _matches := _re.FindAllString( _stdout , -1)
    return len(_matches)


}

func sbatchd_status() int   {
    _stdout:=_cmr.GetSbatchd_status_stdout()
    //_pattern:="/tools/ibm/lsf/10.1/install/10.1/linux3.10-glibc2.17-x86_64/etc/sbatchd"
    _pattern:=\s(root|lsfadmin)\s+/tools/ibm/lsf/10.1/install/10.1/linux3.10-glibc2.17-x86_64/etc/sbatchd(?:\s|\n|$)
    _re:=regexp.MustCompile(_pattern)
    _matches := _re.FindAllString( _stdout , -1)
    return len(_matches)
}

func pim_status() int   {
    _stdout:=_cmr.GetPim_status_stdout()
    //_pattern:="\s(root|lsfadmin)\s+.*/pim"
    _pattern:=\s(root|lsfadmin)\s+/tools/ibm/lsf/10.1/install/10.1/linux3.10-glibc2.17-x86_64/etc/pim(?:\s|\n|$)
    _re:=regexp.MustCompile(_pattern)
    _matches := _re.FindAllString( _stdout , -1)
    return len(_matches)

}


func bld_status() int   {
    _stdout:=_cmr.GetBld_status_stdout()
    _pattern:=\s(root|lsfadmin)\s+/tools/ibm/lsf/10.1/install/10.1/linux3.10-glibc2.17-x86_64/etc/bld(?:\s|\n|$)
    _re:=regexp.MustCompile(_pattern)
    _matches := _re.FindAllString( _stdout , -1)
    return len(_matches)
}

func mbatchd_status() int   {
    _stdout:=_cmr.GetMbatchd_status_stdout()
    _pattern:=\s(root|lsfadmin)\s+/tools/ibm/lsf/10.1/install/10.1/linux3.10-glibc2.17-x86_64/etc/mbatchd(?:\s|\n|$)
    _re:=regexp.MustCompile(_pattern)
    _matches := _re.FindAllString( _stdout , -1)
    return len(_matches)
}

func mbschd_status() int   {
    _stdout:=_cmr.GetMbschd_status_stdout()
    _pattern:=\s(root|lsfadmin)\s+/tools/ibm/lsf/10.1/install/10.1/linux3.10-glibc2.17-x86_64/etc/mbschd(?:\s|\n|$)
    _re:=regexp.MustCompile(_pattern)
    _matches := _re.FindAllString( _stdout , -1)
    return len(_matches)

}

func blcollect_status() int   {
    _stdout:=_cmr.GetBlcollect_status_stdout()
    _pattern:=\s(root|lsfadmin)\s+/tools/ibm/lsf/10.1/install/10.1/linux3.10-glibc2.17-x86_64/etc/blcollect(?:\s|\n|$)
    _re:=regexp.MustCompile(_pattern)
    _matches := _re.FindAllString( _stdout , -1)
    return len(_matches)
}

func licpollerd_status() int   {
    _stdout:=_cmr.GetLicpollerd_status_stdout()
    _pattern:=\s(root|lsfadmin|apache)\s+/opt/IBM/rtm/bin/licpollerd(?:\s|\n|$)
    _re:=regexp.MustCompile(_pattern)
    _matches := _re.FindAllString( _stdout , -1)
    return len(_matches)
}

func lsfpollerd_status() int   {
    _stdout:=_cmr.GetLsfpollerd_status_stdout()
    _pattern:=\s(root|lsfadmin|apache)\s+/opt/IBM/rtm/bin/lsfpollerd(?:\s|\n|$)
    _re:=regexp.MustCompile(_pattern)
    _matches := _re.FindAllString( _stdout , -1)
    return len(_matches)
}



func crond_status() int   {
    _stdout:=_cmr.GetCrond_status_stdout()
    isrunning, := regexp.MatchString(\sroot\s*.*\scrond, _stdout)
    if _isrunning {
        return 1
    } else {
        return 0
    }
}


func cores_status() int   {
    _stdout:=_cmr.GetCore_status_stdout()
    _pattern:=(\d+)
    _re:=regexp.MustCompile(_pattern)
    _matches:=_re.FindStringSubmatch(_stdout)

    if len(_matches)  >= 2  {
       _num,_err:= strconv.Atoi(_matches[1])
        if _err != nil {
        return 2
        }
        return _num


    } else {
        return 1
    }
}

func memory_status() int   {
    _stdout:=_cmr.GetMemory_status_stdout()
    _patt:=`Mem:\s+(\d+) `
    _re := regexp.MustCompile(_patt)
    _matches:=_re.FindStringSubmatch(_stdout)

    if len(_matches)  >= 2  {
        _num,_err:= strconv.Atoi(_matches[1])
        if _err != nil {
        return 2
        }
        return _num

    } else {
        return 1
    }


}



func timediff_chronyd_status() float64 {
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


func ntpserver_chronyd_status() string {
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


func ntpsynced_chronyd_status() int {
    _stdout:=_cmr.GetNtp_Chronyd_status_stdout()
    isrunning, :=regexp.MatchString("NTP synchronized: yes",_stdout)
    if _isrunning {
        return 1
    } else {
        return 0
    }

}



func timediff_ntpd_status() float64 {
    //synchronised to NTP server (198.211.103.209) at stratum 3
    //time correct to within 49 ms
    //polling server every 512 s

    _stdout:=_cmr.GetNtp_Ntpd_status_stdout()
    _re:=regexp.MustCompile(time correct to within (\d+) m)
    _data1 := _re.FindStringSubmatch(_stdout)

    if len(_data1) == 0 {
        return -1.11223344
    }

    data2,:=strconv.ParseFloat(strings.TrimSpace(_data1[1]), 64)

   return _data2

}


func ntpserver_ntpd_status() string {
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


func ntpsynced_ntpd_status() int {
    _stdout:=_cmr.GetNtp_Ntpd_status_stdout()
    isrunning, :=regexp.MatchString(".ynchronised to NTP server.",_stdout)
    if _isrunning {
        return 1
    } else {
        return 0
    }

}



func getHostName() string {
    _functionname:="restapi.getHostname"
    _hostname,_err:=os.Hostname()
    if _err != nil {
        log.Fatal(_functionname," unable to get hostname")
        log.Fatal(_functionname,_err)
        return ""
    }
    return _hostname
}

func getDistribution() string {
    _stdout:=_cmr.GetDistribution_status_stdout()
    _pattern := regexp.MustCompile("Description:\\s+(.*)")
    _subStr := _pattern.FindStringSubmatch( _stdout )

    _distribution:="notfound"
    if len(_subStr) > 1 {
        _distribution=_subStr[1]
    }

    return _distribution
}




//we need to pass additional argument here,
func runningProcessDataHandler(w http.ResponseWriter, r *http.Request) {
    _functionname:="restapi.runningProcessDataHandler"
    //_processes:=[]ProcessInfo {           ProcessInfo{"123","PPID","USER","NI","RSS","SZ","STAT","MEM","CPU","TIME","COMM","CMD",time.Now().Unix()} }
    //fmt.Println(r.URL.Path)
    var _spd ServerProcessesData

    log.Println(_functionname," restendpoint was queried")

    _cmr.Lock()
    defer _cmr.Unlock()
    _processes,_err:=parseRunningProcessData()
    _localmountpointsinfo,_err:=parseLocalmountpointData()
    _spd.LIM=lim_status()
    _spd.RES=res_status()
    _spd.SBATCHD=sbatchd_status()
    _spd.PIM=pim_status()
    _spd.BLD=bld_status()
    _spd.MBSCHD=mbschd_status()
    _spd.BLCOLLECT=blcollect_status()
    _spd.MBATCHD=mbatchd_status()
    _spd.LICPOLLERD=licpollerd_status()
    _spd.LSFPOLLERD=lsfpollerd_status()
    _spd.CROND=crond_status()

    _spd.CORES=cores_status()
    _spd.MEMORY=memory_status()


    _spd.TIMEDIFF_CHRONYD=timediff_chronyd_status()
    _spd.SERVER_CHRONYD=ntpserver_chronyd_status()
    _spd.SYNCED_CHRONYD=ntpsynced_chronyd_status()

    _spd.TIMEDIFF_NTPD=timediff_ntpd_status()
    _spd.SERVER_NTPD=ntpserver_ntpd_status()
    _spd.SYNCED_CHRONYD=ntpsynced_ntpd_status()

    _spd.PROCESSES=&_processes
    _spd.LOCALMOUNTPOINTS=&_localmountpointsinfo

    _time:=time.Now().Unix()
    _spd.TIMESTAMP=_time

    _hostname:=getHostName()
    _spd.HOSTNAME=_hostname

    _distribution:=getDistribution()
    _spd.DISTRIBUTION=_distribution

    if len(_processes) == 0{
        log.Print(_functionname," No running process founf on this host!! total processes=",len(_processes))
    }


    _data, _err := json.Marshal(_spd)
    w.Header().Set("Content-Type", "application/json")
    if _err != nil {
        log.Println(_functionname," ERROR occured while converting data to json: ",_err)
         w.Write([]byte("{}"))
    } else {
        w.Write(_data)
    }
}





func SetupRestHandler(mux *http.ServeMux,_parsercfg_local *aparser.ParserCfg ,_yamlcfg_local *yamlreader.YamlConfig,_cmr_local *commandrunner.CommandRunner) {
     _functionname:="restapi.SetupRestHandler"
     _cmr=_cmr_local
     _yamlcfg=_yamlcfg_local
     _parsercfg=_parsercfg_local

     if _cmr == nil || _yamlcfg == nil || _parsercfg == nil || _cmr_local == nil || _yamlcfg_local == nil || _parsercfg_local == nil {
         log.Println(_functionname,"one of the supplied parameter is blank")
     }


     mux.HandleFunc("/v1/processes/rest", runningProcessDataHandler)
     log.Println(_functionname," registered /v1/processes/rest handler")

}

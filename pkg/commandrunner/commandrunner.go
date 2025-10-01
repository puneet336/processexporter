package commandrunner

import (
        //"fmt"
        "os/exec"
        "sync"
        "log"
        //"regexp"
        "strings"
       "time"
       "github.com/puneets336/customexporter1/pkg/aparser"
       "github.com/puneets336/customexporter1/pkg/yamlreader"
       )


var _yamlcfg *yamlreader.YamlConfig
var _parsercfg *aparser.ParserCfg



type CommandRunner struct{
    Pim_status_stdout string
    Lim_status_stdout string
    Res_status_stdout string
    Sbatchd_status_stdout string
    Bld_status_stdout string
    Mbatchd_status_stdout string
    Blcollect_status_stdout string
    Mbschd_status_stdout string
    Lsfpollerd_status_stdout string
    Licpollerd_status_stdout string
    Crond_status_stdout string
    Ntp_Chronyd_status_stdout string
    Ntp_Ntpd_status_stdout string
    Processes_status_stdout  string
    Localmountpoint_status_stdout string
    Core_status_stdout string
    Memory_status_stdout string
    Distribution_status_stdout string


    Mutex sync.Mutex
}



func NewCommandRunner () *CommandRunner {
    return &CommandRunner {
        Pim_status_stdout: "",
        Lim_status_stdout: "",
        Res_status_stdout: "",
        Sbatchd_status_stdout: "",
        Localmountpoint_status_stdout: "",
    }
}

// for lsf status command :if command fails , you get output,error
// for other command : if command fails you get blankstring and error
func (c *CommandRunner) RunShellCommand(command []string) (string,error) {
    var _stdout string
    _cmd:=exec.Command(command[0],command[1:]...)
    stdout_bytes, _err:=_cmd.Output()


    if stdout_bytes != nil {
        _stdout = strings.TrimSpace(string(stdout_bytes))
    } else {
        _stdout = ""
    }




    return _stdout,_err
}



func (c *CommandRunner) Lock(){
    //log.Println(c.Mutex)
    c.Mutex.Lock()
}

func (c *CommandRunner) Unlock(){
    //log.Println(c.Mutex)
    c.Mutex.Unlock()
}


func (c *CommandRunner) Update() {
    _functionname:="commandrunner.Update"
    if _functionname == ""{
        log.Println("Please set _functionname variable to correct value")
    }
    c.Mutex.Lock()
    defer c.Mutex.Unlock()

    //log.Println(_functionname,"Updating data")
    _output:=c.setSbatchd_status_stdout()
    c.setPim_status_stdout()
    c.setLim_status_stdout(_output)
    c.setRes_status_stdout(_output)
    c.setBld_status_stdout()
    c.setMbatchd_status_stdout()
    c.setBlcollect_status_stdout()
    c.setMbschd_status_stdout()
    c.setLsfpollerd_status_stdout()
    c.setLicpollerd_status_stdout()
    c.setCrond_status_stdout()
    c.setNtp_Ntpd_status_stdout()
    c.setNtp_Chronyd_status_stdout()

    c.setProcesses_status_stdout()


    c.setLocalmountpoint_status_stdout()

    c.setCore_status_stdout()
    c.setMemory_status_stdout()
    c.setDistribution_status_stdout()


    c.Print()
    //log.Println(_functionname," Updated data")
}

func (c *CommandRunner) Print() {
    _functionname:="commandrunner.Print"
    log.Println(_functionname," Displaying data")
    log.Println(_functionname," PIM:",  strings.Replace( c.GetPim_status_stdout() ,"\n"," ",-1) )
    log.Println(_functionname," LIM:",  strings.Replace( c.GetLim_status_stdout() ,"\n"," ",-1) )
    log.Println(_functionname," RES:",  strings.Replace( c.GetRes_status_stdout() ,"\n"," ",-1) )
    log.Println(_functionname," SBATCHD:", strings.Replace( c.GetSbatchd_status_stdout(),"\n"," ",-1)  )
    log.Println(_functionname," BLD:", strings.Replace( c.GetBld_status_stdout(),"\n"," ",-1)  )
    log.Println(_functionname," MBATCHD:", strings.Replace( c.GetMbatchd_status_stdout(),"\n"," ",-1)  )
    log.Println(_functionname," BLCOLLECT:", strings.Replace( c.GetBlcollect_status_stdout(),"\n"," ",-1)  )
    log.Println(_functionname," MBSCHD:", strings.Replace( c.GetMbschd_status_stdout(),"\n"," ",-1)  )
    log.Println(_functionname," LSFPOLLERD:", strings.Replace( c.GetLsfpollerd_status_stdout(),"\n"," ",-1)  )
    log.Println(_functionname," LICPOLLERD:", strings.Replace( c.GetLicpollerd_status_stdout(),"\n"," ",-1)  )
    log.Println(_functionname," CROND:", strings.Replace( c.GetCrond_status_stdout(),"\n"," ",-1)  )
    log.Println(_functionname," CORES:", strings.Replace( c.GetCore_status_stdout(),"\n"," ",-1)  )
    log.Println(_functionname," MEM:", strings.Replace( c.GetMemory_status_stdout(),"\n"," ",-1)  )




    //process data we can't print full string , lets print first 11 words

    _proc_data:=c.GetProcesses_status_stdout()
    _words1:=strings.Replace(_proc_data,"\n"," ",-1)

    if len(_words1) > 11 {
        log.Println(_functionname," PROC:", strings.Fields(_words1)[:11]  )
    } else {
        log.Println(_functionname," PROC:", _words1  )
    }

    _localmountpoint_data:=c.GetLocalmountpoint_status_stdout()
    _words2:=strings.Replace(_localmountpoint_data,"\n"," ",-1)

    if len(_words2) > 11 {
        log.Println(_functionname," MOUNT:", strings.Fields(_words2)[:11]  )
    } else {
        log.Println(_functionname," MOUNT:", _words2  )
    }


    _ntp_chronyd_data:=c.GetNtp_Chronyd_status_stdout()
    _words3:=strings.Replace( _ntp_chronyd_data,"\n"," ",-1)
    if len(_words3) > 11 {
        log.Println(_functionname," NTP CHRONYD:", strings.Fields(_words3)[:11]  )
    } else {
        log.Println(_functionname," NTP CHRONYD:", _words3  )
    }


    _ntp_ntpd_data:=c.GetNtp_Ntpd_status_stdout()
    _words4:=strings.Replace( _ntp_ntpd_data,"\n"," ",-1)
    if len(_words4) > 11 {
        log.Println(_functionname," NTP NTPD:", strings.Fields(_words4)[:11]  )
    } else {
        log.Println(_functionname," NTP NTPD:", _words4  )
    }



}


func (c *CommandRunner) setPim_status_stdout() {
    _functionname:="commandrunner.setPim_status_stdout"
    _cmd:=[]string{"ps","--noheader","-C","pim","-o","pid,ppid,user:20,cmd"}
    _output,_err:=c.RunShellCommand(_cmd)
    if _err != nil {
        log.Println(_functionname," Unable to get data from ",_cmd,". no running pim process found")
        _output=""
    }
    c.Pim_status_stdout=_output

}


func (c *CommandRunner) setMbatchd_status_stdout() {
    _functionname:="commandrunner.setMbatchd_status_stdout"
    _cmd:=[]string{"ps","--noheader","-C","mbatchd","-o","pid,ppid,user:20,cmd"}
    _output,_err:=c.RunShellCommand(_cmd)
    if _err != nil {
        log.Println(_functionname," Unable to get data from ",_cmd,". no running mbatchd process found")
        _output=""
    }
    c.Mbatchd_status_stdout=_output

}

func (c *CommandRunner) setBlcollect_status_stdout() {
    _functionname:="commandrunner.setBlcollect_status_stdout"
    _cmd:=[]string{"ps","--noheader","-C","blcollect","-o","pid,ppid,user:20,cmd"}
    _output,_err:=c.RunShellCommand(_cmd)
    if _err != nil {
        log.Println(_functionname," Unable to get data from ",_cmd,". no running blcollect process found")
        _output=""
    }
    c.Blcollect_status_stdout=_output

}

func (c *CommandRunner) setMbschd_status_stdout() {
    _functionname:="commandrunner.setMbschd_status_stdout"
    _cmd:=[]string{"ps","--noheader","-C","mbschd","-o","pid,ppid,user:20,cmd"}
    _output,_err:=c.RunShellCommand(_cmd)
    if _err != nil {
        log.Println(_functionname," Unable to get data from ",_cmd,". no running mbschd process found")
        _output=""
    }
    c.Mbschd_status_stdout=_output

}
func (c *CommandRunner) setLsfpollerd_status_stdout() {
    _functionname:="commandrunner.setLsfpollerd_status_stdout"
    _cmd:=[]string{"ps","--noheader","-C","lsfpollerd","-o","pid,ppid,user:20,cmd"}
    _output,_err:=c.RunShellCommand(_cmd)
    if _err != nil {
        log.Println(_functionname," Unable to get data from ",_cmd,". no running lsfpollerd process found")
        _output=""
    }
    c.Lsfpollerd_status_stdout=_output
}
func (c *CommandRunner) setLicpollerd_status_stdout() {
    _functionname:="commandrunner.setLicpollerd_status_stdout"
    _cmd:=[]string{"ps","--noheader","-C","licpollerd","-o","pid,ppid,user:20,cmd"}
    _output,_err:=c.RunShellCommand(_cmd)
    if _err != nil {
        log.Println(_functionname," Unable to get data from ",_cmd,". no running licpollerd process found")
        _output=""
    }
    c.Licpollerd_status_stdout=_output
}
func (c *CommandRunner) setCrond_status_stdout() {
    _functionname:="commandrunner.setCrond_status_stdout"
    _cmd1:=[]string{"ps","--noheader","-C","crond","-o","pid,ppid,user:20,comm"}
    _output1,_err1:=c.RunShellCommand(_cmd1)

    _cmd2:=[]string{"ps","--noheader","-C","cron","-o","pid,ppid,user:20,comm"}
    _output2,_err2:=c.RunShellCommand(_cmd2)

    if _err1 == nil {
        c.Crond_status_stdout=_output1
    } else if _err2 == nil  {
        c.Crond_status_stdout=_output2
    } else {
        c.Crond_status_stdout=""
    }

    if _err1 != nil {
        log.Println(_functionname," Unable to get data from ",_cmd1,". no running crond process found")
    }
    if  _err2 != nil {
        log.Println(_functionname," Unable to get data from ",_cmd2,". no running crond process found")
    }
}


func (c *CommandRunner) setNtp_Chronyd_status_stdout() {
    _functionname:="commandrunner.setNtp_Chronyd_status_stdout"
    _cmd1:=[]string{"chronyc","tracking"}
    _output1,_err1:=c.RunShellCommand(_cmd1)
    if _err1 != nil {
        log.Println(_functionname," Unable to get data from ",_cmd1)
        _output1=""
    }

    _cmd2:=[]string{"timedatectl"}
    _output2,_err2:=c.RunShellCommand(_cmd2)
    if _err2 != nil {
        log.Println(_functionname," Unable to get data from ",_cmd2)
        _output2=""
    }

    c.Ntp_Chronyd_status_stdout=_output1 + _output2
    //fmt.Println(c.Ntp_Chronyd_status_stdout)
}


func (c *CommandRunner) setNtp_Ntpd_status_stdout() {
    _functionname:="commandrunner.setNtp_Ntpd_status_stdout"
    _cmd:=[]string{"ntpstat"}
    _output,_err:=c.RunShellCommand(_cmd)
    if _err != nil {
        log.Println(_functionname," Unable to get data from ",_cmd)
        _output=""
    }

    c.Ntp_Ntpd_status_stdout=_output
}



func (c *CommandRunner)  setLocalmountpoint_status_stdout() {
    _functionname:="commandrunner.setLocalmountpoint_status_stdout"
    _cmd:=[]string{"df","-l"}
    _output,_err:=c.RunShellCommand(_cmd)
    if _err != nil {
        log.Println(_functionname," Unable to get data from ",_cmd,". is the command missing? or command is not behaving as expected")
        _output=""
    }

    c.Localmountpoint_status_stdout=_output
}


func (c *CommandRunner) setLim_status_stdout(_output string) {
    _functionname:="commandrunner.setLim_status_stdout"
    _cmd:=[]string{"ps","--noheader","-C","lim","-o","pid,ppid,user:20,cmd"}
    _output,_err:=c.RunShellCommand(_cmd)

    if _err != nil {
        log.Println(_functionname," Unable to get data from ",_cmd,". no running bld process found")
        _output=""
    }
    c.Lim_status_stdout=_output
}

func (c *CommandRunner) setRes_status_stdout(_output string) {
     _functionname:="commandrunner.setRes_status_stdout"
    _cmd:=[]string{"ps","--noheader","-C","res","-o","pid,ppid,user:20,cmd"}
    _output,_err:=c.RunShellCommand(_cmd)

    if _err != nil {
        log.Println(_functionname," Unable to get data from ",_cmd,". no running bld process found")
        _output=""
    }

    c.Res_status_stdout=_output

}


func (c *CommandRunner) setBld_status_stdout() {
    _functionname:="commandrunner.setBld_status_stdout"
    //c.Res_status_stdout=""

    _cmd:=[]string{"ps","--noheader","-C","bld","-o","pid,ppid,user:20,cmd"}
    _output,_err:=c.RunShellCommand(_cmd)

    if _err != nil {
        log.Println(_functionname," Unable to get data from ",_cmd,". no running bld process found")
        _output=""
    }

    c.Bld_status_stdout=_output
}

func (c *CommandRunner) setSbatchd_status_stdout() string {
    //from output of this script we will check status of sbatchd,lim,pim
    _functionname:="commandrunner.setSbatchd_status_stdout"
    _cmd:=[]string{"ps","--noheader","-C","sbatchd","-o","pid,ppid,user:20,cmd"}
    _output,_err:=c.RunShellCommand(_cmd)

    if _err != nil {
        log.Println(_functionname," Unable to get data from ",_cmd,". no running bld process found")
        _output=""
    }

    //var command [3]string=[3]string{_lsffile,"status"}
    //output,:=c.RunShellCommand([]string{_lsffile,"status"})
    c.Sbatchd_status_stdout=_output


    return _output

}

func (c *CommandRunner) setProcesses_status_stdout() {
    _functionname:="commandrunner.setProcesses_status_stdout"
    _cmd:=[]string{"ps","--no-headers", "-e", "-o", "pid,ppid,user:30,ni,rss,sz,stat,%mem,%cpu,time,comm,cmd"}
    _output,_err:=c.RunShellCommand(_cmd)

    if _err != nil {
         _output=""
         log.Print(_functionname," Error while running command ",_cmd," - ",_err)
    }

    c.Processes_status_stdout=_output
}


func (c *CommandRunner) setCore_status_stdout() {
    _functionname:="commandrunner.setProcesses_status_stdout"
    _cmd:=[]string{"getconf" ,"_NPROCESSORS_ONLN"}
    _output,_err:=c.RunShellCommand(_cmd)

    if _err != nil {
         _output=""
         log.Print(_functionname," Error while running command ",_cmd," - ",_err)
    }
    c.Core_status_stdout=_output

}

func (c *CommandRunner) setMemory_status_stdout() {
    _functionname:="commandrunner.setProcesses_status_stdout"
    _cmd:=[]string{"free","-g"}
    _output,_err:=c.RunShellCommand(_cmd)

    if _err != nil {
         _output=""
         log.Print(_functionname," Error while running command ",_cmd," - ",_err)
    }

    c.Memory_status_stdout=_output
}

func (c *CommandRunner) setDistribution_status_stdout() {
    _functionname:="commandrunner.setDistribution_status_stdout"
    _cmd:=[]string{"lsb_release","-d"}
    _output,_err:=c.RunShellCommand(_cmd)

    if _err != nil {
         _output=""
         log.Print(_functionname," Error while running command ",_cmd," - ",_err)
    }

    c.Distribution_status_stdout=_output
}




func (c *CommandRunner) GetPim_status_stdout() string{
    //c.Mutex.Lock()
    //defer c.Mutex.Unlock()
    return c.Pim_status_stdout
}

func (c *CommandRunner) GetMbatchd_status_stdout() string{
    //c.Mutex.Lock()
    //defer c.Mutex.Unlock()
    return c.Mbatchd_status_stdout
}

func (c *CommandRunner) GetBlcollect_status_stdout() string{
    //c.Mutex.Lock()
    //defer c.Mutex.Unlock()
    return c.Blcollect_status_stdout
}

func (c *CommandRunner) GetMbschd_status_stdout() string{
    //c.Mutex.Lock()
    //defer c.Mutex.Unlock()
    return c.Mbschd_status_stdout
}


func (c *CommandRunner) GetLim_status_stdout() string{
    //c.Mutex.Lock()
    //defer c.Mutex.Unlock()
    return c.Lim_status_stdout
}

func (c *CommandRunner) GetRes_status_stdout() string{
    //c.Mutex.Lock()
    //defer c.Mutex.Unlock()
    return c.Res_status_stdout
}

func (c *CommandRunner) GetSbatchd_status_stdout() string{
    //c.Mutex.Lock()
    //defer c.Mutex.Unlock()
    return c.Sbatchd_status_stdout
}

func (c *CommandRunner) GetBld_status_stdout() string{
    //c.Mutex.Lock()
    //defer c.Mutex.Unlock()
    return c.Bld_status_stdout
}

func (c *CommandRunner) GetLicpollerd_status_stdout() string{
    //c.Mutex.Lock()
    //defer c.Mutex.Unlock()
    return c.Licpollerd_status_stdout
}

func (c *CommandRunner) GetLsfpollerd_status_stdout() string{
    //c.Mutex.Lock()
    //defer c.Mutex.Unlock()
    return c.Lsfpollerd_status_stdout
}

func (c *CommandRunner) GetCrond_status_stdout() string{
    //c.Mutex.Lock()
    //defer c.Mutex.Unlock()
    return c.Crond_status_stdout
}

func (c *CommandRunner) GetNtp_Chronyd_status_stdout() string{
    //c.Mutex.Lock()
    //defer c.Mutex.Unlock()
    return c.Ntp_Chronyd_status_stdout
}

func (c *CommandRunner) GetNtp_Ntpd_status_stdout() string{
    //c.Mutex.Lock()
    //defer c.Mutex.Unlock()
    return c.Ntp_Ntpd_status_stdout
}


func (c *CommandRunner) GetProcesses_status_stdout() string{
   // c.Mutex.Lock()
   // defer c.Mutex.Unlock()
    return c.Processes_status_stdout
}

func (c *CommandRunner) GetLocalmountpoint_status_stdout() string{
    return c.Localmountpoint_status_stdout
}

func (c *CommandRunner) GetCore_status_stdout() string{
    return c.Core_status_stdout
}
func (c *CommandRunner) GetMemory_status_stdout() string{
    return c.Memory_status_stdout
}

func (c *CommandRunner) GetDistribution_status_stdout() string{
    return c.Distribution_status_stdout
}



func Update(_cmr *CommandRunner) {
    for {
    _cmr.Update()
    time.Sleep( time.Duration(_yamlcfg.Queryinterval)   * time.Second )
    }
}

func CreateAndStartCommandRunner( _parsercfg_local *aparser.ParserCfg ,_yamlcfg_local *yamlreader.YamlConfig )  *CommandRunner{
    _parsercfg=_parsercfg_local
    _yamlcfg=_yamlcfg_local
    // reference type
    _cmr:=NewCommandRunner()
    go Update(_cmr)

    time.Sleep( 30  * time.Second )
    return _cmr
}

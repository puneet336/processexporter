package aparser

import (
       "fmt"
       "log"
       "github.com/akamensky/argparse"
       "os"
       )

type ParserCfg struct {
    Configfilepath string
    ListenAddress string
    DebugMode int
}

func (pc *ParserCfg) GetDebugMode() int {
    return pc.DebugMode

}

func SetupParser() *ParserCfg  {
    _functionname:="aparser.SetupParser"
    _parser := argparse.NewParser("servicestatus", "Prints provided string to stdout")

    //var myString *servertype = parser.String("t", "server-type", []string{"MASTER","CLIENT"},&argparse.Options{Required: true, Help: "running on LSF masger of LSF client"})

    //if you are not supplying ath of config file
    _configfilepath := _parser.String("f", "config" ,&argparse.Options{Help: "configfilepath",Default: "/tools/ampere/common/lib/utils/go_process_monitor/config.yaml",})
    _version := _parser.Flag("v", "version" ,&argparse.Options{Help: "print version",})
    _port := _parser.String( "p" ,"port",&argparse.Options{Help: "listen port for application",Default: "41000",}  )
    _debug :=  _parser.Int( "d" ,"debug",&argparse.Options{Help: "debug mode",Default: 0,})

    _err:=_parser.Parse(os.Args)


    if _err != nil {
                // In case of error print error and print usage
                // This can also be done by passing -h or --help flags
                fmt.Print(_functionname,_parser.Usage(_err))
                os.Exit(1)
    }

    _parsercfg:=ParserCfg{}
    _parsercfg.Configfilepath=*_configfilepath
    _parsercfg.ListenAddress=":" + *_port

    if *_version {
         log.Println("version: v0.3.5")
         os.Exit(0)
    }
    _parsercfg.DebugMode=*_debug

    //log.Println(_functionname,"INFO: parserconfig-",_parsercfg.Configfilepath,"-",*_configfilepath)
    return &_parsercfg

}

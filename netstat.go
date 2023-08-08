package main

import (
	"bufio"
	"fmt"
	"strings"
	"errors"
	"os"
	"strconv"
	"git.zabbix.com/ap/plugin-support/plugin/container"
	"git.zabbix.com/ap/plugin-support/plugin"
)

type Plugin struct {
           plugin.Base
       }
       
var impl Plugin

func (p *Plugin) getProcNetstat(protocol string, stat string) (result uint64, err error) {	

	file, err := os.Open("/proc/net/netstat")
	
	if err != nil {
		return 0, errors.New("Cannot open /proc/net/netstat")
	}
	defer file.Close()
	
	protocol += ":"
	var fields []string
	var values []string
	for scanner := bufio.NewScanner(file); scanner.Scan(); {
		fields = strings.Split(scanner.Text(), " ")
		scanner.Scan()
		if fields[0] == protocol {
			values = strings.Split(scanner.Text(), " ")
			for i := 1; i < len(fields); i++ {
				if fields[i] == stat {
					result, err =  strconv.ParseUint(values[i], 10, 64)
				}
			}
			return result, nil
		}
	}
	return 0, errors.New("Cannot find protocol or stat")
}

func (p *Plugin) Export(key string, params []string, ctx plugin.ContextProvider) (result interface{}, err error) {
	var protocol, stat string
	
	if len(params) == 2 {
		protocol = params[0]
		stat = params[1]
	} else {
		return nil, errors.New("Wrong number of parameters")
	}
	
	return p.getProcNetstat(protocol, stat)
}

func init() {
	plugin.RegisterMetrics(&impl, "netstat", "net.proc.netstat", "Get stat in /proc/net/netstat.")
}

func main() {
           h, err := container.NewHandler(impl.Name())
           if err != nil {
               panic(fmt.Sprintf("failed to create plugin handler %s", err.Error()))
           }
           impl.Logger = &h

           err = h.Execute()
           if err != nil {
               panic(fmt.Sprintf("failed to execute plugin handler %s", err.Error()))
           }
       }

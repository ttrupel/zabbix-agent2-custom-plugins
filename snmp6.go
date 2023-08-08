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

func (p *Plugin) getProcSnmp6(stat string) (result uint64, err error) {
	
	file, err := os.Open("/proc/net/snmp6")
	
	if err != nil {
		return 0, errors.New("Cannot open /proc/net/snmp6")
	}
	defer file.Close()
	
	var stats []string
	for scanner := bufio.NewScanner(file); scanner.Scan(); {
		
		stats = strings.Fields(scanner.Text())
		
		if stats[0] == stat {
			result, err =  strconv.ParseUint(stats[1], 10, 64)
			if err != nil {
				return 0, errors.New("Cannot parse stat")
			}
			return result, nil
		}
	}
	return 0, errors.New("Cannot find stat")
}

func (p *Plugin) Export(key string, params []string, ctx plugin.ContextProvider) (result interface{}, err error) {	
	var stat string

	if len(params) == 1 {
		stat = params[0]
	} else {
		return nil, errors.New("Wrong number of parameters")
	}
	
	return p.getProcSnmp6(stat)
}

func init() {
	plugin.RegisterMetrics(&impl, "snmp6", "net.proc.snmp6", "Get stat in /proc/net/snmp6.")
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

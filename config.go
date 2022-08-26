package main
import (
	"strconv"
	"fmt"
	"os"
	"errors"
	"encoding/json"
	"io/ioutil"
        )
var conf_fn = "conf.json"

type Arg_func func(*configuration, []string) (error)

type Secret struct {
        Secret string   `json:"secret"`
        Fix string      `json:"fix"`
        Counter uint64  `json:"counter"`
        Len int         `json:"len"`
        }
type cmd_line_items struct {
	id		int
	Switch		string
	Help_srt	string
	Has_arg		bool
	Func		Arg_func
}

type configuration struct {
	LinuxWDebug	string
	LinuxWODebug	string
	StripBin	string
	DBURL		string
	DBPort		int
	DBUser		string
	DBPassword	string
	DBTargetDB	string
	Maintainers_fn	string
	KConfig_fn	string
	KMakefile	string
	Mode		int
	Note		string
}
var	Default_config  configuration = configuration{
	LinuxWDebug:	"vmlinux",
	LinuxWODebug:	"vmlinux.work",
	StripBin:	"/usr/bin/strip",
	DBURL:		"dbs.hqhome163.com",
	DBPort:		5432,
	DBUser:		"alessandro",
	DBPassword:	"<password>",
	DBTargetDB:	"kernel_bin",
	Maintainers_fn:	"MAINTAINERS",
	KConfig_fn:	"include/generated/autoconf.h",
	KMakefile:	"Makefile",
	Mode:		ENABLE_SYBOLSNFILES|ENABLE_XREFS|ENABLE_MAINTAINERS|ENABLE_VERSION_CONFIG,
	Note:		"upstream",
	}

func push_cmd_line_item(Switch string, Help_str string, Has_arg bool, Func Arg_func, cmd_line *[]cmd_line_items){
	*cmd_line = append(*cmd_line, cmd_line_items{id: len(*cmd_line)+1, Switch: Switch, Help_srt: Help_str, Has_arg: Has_arg, Func: Func})
}

func cmd_line_item_init() ([]cmd_line_items){
	var res	[]cmd_line_items

	push_cmd_line_item("-f", "specifies json configuration file",			true,  func_jconf,	&res)
	push_cmd_line_item("-s", "Forces use specified strip binary",			true,  func_forceStrip,	&res)
	push_cmd_line_item("-u", "Forces use specified database userid",		true,  func_DBUser,	&res)
	push_cmd_line_item("-p", "Forecs use specified password",			true,  func_DBPass,	&res)
	push_cmd_line_item("-d", "Forecs use specified DBhost",				true,  func_DBHost,	&res)
	push_cmd_line_item("-p", "Forecs use specified DBPort",				true,  func_DBPort,	&res)
	push_cmd_line_item("-n", "Forecs use specified note (default 'upstream')",	true,  func_Note,	&res)
	push_cmd_line_item("-c", "Checks dependencies",					false, func_check,	&res)
	push_cmd_line_item("-h", "This Help",						false, func_help,	&res)

	return res
}
func func_help          (conf *configuration,fn []string)               (error){
	return errors.New("Dummy")
}
func func_jconf		(conf *configuration,fn []string)		(error){
	jsonFile, err := os.Open(fn[0])
	if err != nil {
                return err
                }
        byteValue, _ := ioutil.ReadAll(jsonFile)
        jsonFile.Close()
	err=json.Unmarshal(byteValue, conf)
	if err != nil {
		return err
		}
	return nil
}

func func_forceStrip	(conf *configuration, fn []string)	(error){
	(*conf).StripBin=fn[0]
	return nil
}

func func_DBUser	(conf *configuration, user []string)	(error){
	(*conf).DBUser=user[0]
	return nil
}

func func_DBPass	(conf *configuration, pass []string)	(error){
	(*conf).DBPassword=pass[0]
	return nil
}

func func_DBHost	(conf *configuration, host []string)	(error){
	(*conf).DBURL=host[0]
	return nil
}

func func_DBPort	(conf *configuration, port []string)	(error){
	s, err := strconv.Atoi(port[0])
	if err!=nil {
		return err
		}
	(*conf).DBPort=s
	return nil
}

func func_Note		(conf *configuration, note []string)	(error){
	(*conf).Note=note[0]
	return nil
}

func func_check		(conf *configuration, dummy []string)			(error){
	return nil
}


func print_help(lines []cmd_line_items){

	for _,item := range lines{
		fmt.Printf(
			"\t%s\t%s\t%s\n",
			item.Switch,
			func (a bool)(string){
				if a {
					return "<v>"
					}
				return ""
			}(item.Has_arg),
			item.Help_srt,
			)
		}
}

func args_parse(lines []cmd_line_items)(configuration, error){
	var	skip		bool=false;
	var	conf		configuration=Default_config
	var 	f		Arg_func

	for _, os_arg := range os.Args[1:] {
//		fmt.Printf("osarg=%s\n", os_arg)
		if !skip {
//			fmt.Printf("check if I have it\n")
			for _, arg := range lines{
//				fmt.Printf("consider configured arg=%s\n", arg.Switch)
				if arg.Switch==os_arg {
//					fmt.Printf("have it\n")
					if arg.Has_arg{
//						fmt.Printf("it needs more arguments\n")
						f=arg.Func
						skip=true
						break
						}
					err := arg.Func(&conf, []string{})
					if err != nil {
//						fmt.Println("-----------------------------_")
						return Default_config, err
						}
					}
				}
			continue
			}
		if skip{
//			fmt.Printf("Fetch extra arg (%s)and use it\n", os_arg)
			err := f(&conf,[]string{os_arg})
			if err != nil {
				return Default_config, err
				}
			skip=false
			}

		}
	return	conf, nil
}

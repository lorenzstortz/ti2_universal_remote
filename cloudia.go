package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
)

var (
	IRDevices []IRDevice
	RFDevices []RFDevice
)

const IRDEVICE="IRDevice"
const RFDEVICE="RFDevice"

type IRCommand struct {
	DeviceName string `json:"DeviceName"`
	DeviceKey  string `json:"DeviceKey"`
}

type RFCommand struct{
	DeviceName string `json:"DeviceName"`
	Status string `json:"Status"`
}

type RFDevice struct {
	Name string
	OnCode int
	OffCode int
}

type IRDevice struct {
	Name string
}


func handleIR(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		respBytes,err:=ioutil.ReadAll(r.Body)
		if err!=nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		cmd := IRCommand{}
		err=json.Unmarshal(respBytes, &cmd)
		if err!=nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if sort.Search(len(IRDevices), func(i int) bool {
			return IRDevices[i].Name== cmd.DeviceName
		}) >= len(IRDevices) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if sendIRCommand(cmd) {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)

	default:
		// Give an error message.
		fmt.Println("default")
		w.WriteHeader(http.StatusNotFound)
	}
}


func sendIRCommand(cmd IRCommand) bool{
	// do stuff
}

func sendRFCommand(cmd RFCommand) bool{
	//do stuff
	var code string
	if
        cmd := exec.Command("./controller", cmd.)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func handleRF(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		respBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		cmd := RFCommand{}
		err = json.Unmarshal(respBytes, &cmd)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if sort.Search(len(RFDevices), func(i int) bool {
			return RFDevices[i].Name == cmd.DeviceName
		}) >= len(RFDevices) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if sendRFCommand(cmd) {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)

	default:
		// Give an error message.
		fmt.Println("default")
		w.WriteHeader(http.StatusNotFound)
	}
}
func parseDevices(path string){

	IRDevices=make([]IRDevice,0)
	RFDevices=make([]RFDevice,0)

	file,err:=os.Open(path)
	if err!=nil {
		fmt.Println(err)
		os.Exit(1)
	}

	csvReader := csv.NewReader(file)

	for line,err:=csvReader.Read();line!=nil ; line,err=csvReader.Read() {
		if err!=nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if line[1] == IRDEVICE{
			IRDevices= append(IRDevices, IRDevice{Name:line[0]})
		} else if line[1] == RFDEVICE {
			onCode,err:=strconv.Atoi(line[2])
			if err!=nil {
				fmt.Println(err)
				os.Exit(1)
			}
			offCode,err:=strconv.Atoi(line[3])
			if err!=nil {
				fmt.Println(err)
				os.Exit(1)
			}
			RFDevices=append(RFDevices,RFDevice{
				Name:    line[0],
				OnCode:  onCode,
				OffCode: offCode,
			})
		}
	}

}

func main(){
	if len(os.Args)!=2 {
		fmt.Println("Need device file!")
		return
	}
	parseDevices(os.Args[1])

	http.HandleFunc("/ir", handleIR)
	http.HandleFunc("/rf", handleRF)
	log.Fatal(http.ListenAndServe(":80", nil))
}

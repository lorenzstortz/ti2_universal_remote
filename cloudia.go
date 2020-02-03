package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var (
	IRDevices map[string]IRDevice
	RFDevices map[string]RFDevice
)

const IRDEVICE = "IRDevice"
const RFDEVICE = "RFDevice"

type IRCommand struct {
	DeviceName string `json:"DeviceName"`
	DeviceKey  string `json:"DeviceKey"`
}

type RFCommand struct {
	DeviceName string `json:"DeviceName"`
	Status     string `json:"Status"`
}

type RFDevice struct {
	Name    string
	OnCode  string
	OffCode string
}

type IRDevice struct {
	Name string
}

func handleIR(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		respBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		cmd := IRCommand{}
		err = json.Unmarshal(respBytes, &cmd)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, exists := IRDevices[cmd.DeviceName]; !exists {
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

func sendIRCommand(cmd IRCommand) bool {
	// do stuff
	ex := exec.Command("irsend", "-d", "/dev/lirc0", "SEND_ONCE", cmd.DeviceName, cmd.DeviceKey)
	err := ex.Run()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func sendRFCommand(cmd RFCommand) bool {
	//do stuff
	device := RFDevices[cmd.DeviceName]
	var code string
	if cmd.Status == "on" {
		code = device.OnCode
	} else if cmd.Status == "off" {
		code = device.OffCode
	} else {
		return false
	}
	ex := exec.Command("./controller", code)
	err := ex.Run()
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
		if _, exists := RFDevices[cmd.DeviceName]; !exists {
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
func parseDevices(path string) {

	IRDevices = make(map[string]IRDevice, 0)
	RFDevices = make(map[string]RFDevice, 0)

	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	csvReader := csv.NewReader(file)

	for line, err := csvReader.Read(); line != nil; line, err = csvReader.Read() {
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if line[1] == IRDEVICE {
			IRDevices[line[0]] = IRDevice{Name: line[0]}
		} else if line[1] == RFDEVICE {
			RFDevices[line[0]] = RFDevice{
				Name:    line[0],
				OnCode:  line[2],
				OffCode: line[3],
			}
		}
	}

}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Need device file!")
		return
	}
	parseDevices(os.Args[1])

	http.HandleFunc("/ir", handleIR)
	http.HandleFunc("/rf", handleRF)
	log.Fatal(http.ListenAndServe(":80", nil))
}

package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

//Envia una petición al endpoint indicado y extrae la respuesta del campo Data para retornarla

func GetRequestNew(endpoint string, route string, target interface{}) error {
	url := beego.AppConfig.String("ProtocolAdmin") + "://" + beego.AppConfig.String(endpoint) + route
	fmt.Println("URL ", url)
	var response map[string]interface{}
	var err error
	err = GetJson(url, &response)
	err = ExtractData(response, &target)
	return err
}

func GetJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			beego.Error(err)
		}
	}()
	return json.NewDecoder(r.Body).Decode(target)
}

//Esta función extrae la información cuando se recibe encapsulada en una estructura
// y da manejo a las respuestas que contienen arreglos de objetos vacios

func ExtractData(respuesta map[string]interface{}, v interface{}) error {
	var err error
	if respuesta["Success"] == false {
		err = errors.New(fmt.Sprint(respuesta["Data"], respuesta["Message"]))
		panic(err)
	}
	datatype := fmt.Sprint("%v", respuesta["Data"])
	switch datatype {
	case "map[]", "[map[]]": //response vacio
		break
	default:
		err = FillStruct(respuesta["Data"], &v)
		respuesta = nil
	}
	return err
}

//Convertir estructura a JSON

func FillStruct(in interface{}, out interface{}) (err error) {
	var str []byte
	if str, err = json.Marshal(in); err != nil {
		return
	}
	err = json.Unmarshal(str, &out)
	return
}

//Manejo de errores
func ErrorController(c beego.Controller, controller string) {
	if err := recover(); err != nil {
		logs.Error(err)
		localError := err.(map[string]interface{})
		c.Data["mesaage"] = (beego.AppConfig.String("appname") + "/" + controller + "/" + (localError["funcion"]).(string))
		c.Data["data"] = (localError["err"])
		if status, ok := localError["status"]; ok {
			c.Abort(status.(string))
		} else {
			c.Abort("500")
		}
	}
}

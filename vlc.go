//Control of VLC through http interface
package vlc

import(
	"net/url"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"errors"
)

type Server struct{
	addr string
}

type Item struct{
	Name string
	URI string
	
	//Internal things
	Children []Item
	ID string
}

const statusPath = "/requests/status.json?"

func New(address string)(*Server,error){
	s := new(Server)
	s.addr = fmt.Sprintf("http://%s",address)
	//TODO connection testing
	
	return s,nil
}

//
// Playback Control
//

func (s *Server) Play()(error){
	v := url.Values{}
	v.Set("command", "pl_play")
	
	http.Get(s.addr+statusPath+v.Encode())
	//TODO error testing
	
	return nil	
}

func (s *Server) Pause()(error){
	v := url.Values{}
	v.Set("command", "pl_pause")
	
	http.Get(s.addr+statusPath+v.Encode())
	//TODO error testing
	
	return nil	
}

func (s *Server) Stop()(error){
	v := url.Values{}
	v.Set("command", "pl_stop")
	
	http.Get(s.addr+statusPath+v.Encode())
	//TODO error testing
	
	return nil
}

//Removes item from playlist
func (s *Server) Next()(error){	
	v := url.Values{}
	v.Set("command", "pl_next")
	
	http.Get(s.addr+statusPath+v.Encode())
	//TODO error testing
	
	//Delete first
	s.Delete(s.Playlist()[0].URI)
	
	return nil
}
func (s *Server) Previous()(error){
	v := url.Values{}
	v.Set("command", "pl_previous")
	
	http.Get(s.addr+statusPath+v.Encode())
	//TODO error testing
	
	return nil
}

//Recieves volume as a percentage 0-100
func (s *Server) SetVolume(vol uint8)(error){
	v := url.Values{}
	v.Set("command", "volume")
	v.Set("val",fmt.Sprintf("%d%%",vol))
	
	http.Get(s.addr+statusPath+v.Encode())
	//TODO error testing
	
	return nil
}

//
// Playlist
//
func (s *Server) NowPlaying()(string){
	resp,_ := http.Get(s.addr+statusPath)
	data,_ := ioutil.ReadAll(resp.Body);
	
	var val1 map[string](interface{})
	json.Unmarshal(data, &val1)
	
	info,ok := val1["information"].(map[string]interface{})
	if ok{
		category := info["category"].(map[string]interface{})
		meta := category["meta"].(map[string]interface{})
		file := meta["filename"].(string)
	
		return file //URI of current song
	}else{
		//No song playing
		return ""
	}
}

func (s *Server) Playlist()([]Item){
	resp,err:= http.Get(s.addr+"/requests/playlist.json")
	if err != nil{
		fmt.Println(err)
		return nil
	}
	
	data,_ := ioutil.ReadAll(resp.Body);
		
	var val Item
	err = json.Unmarshal(data, &val)
	if err != nil{
		fmt.Println(err)
	}
	
	if len(val.Children) == 0{
		return nil
	}

	playlist := val.Children[0]
	return playlist.Children
}

func (s *Server) Add(uri string)(error){
	v := url.Values{}
	v.Set("command", "in_enqueue")
	v.Set("input",uri)
	
	http.Get(s.addr+statusPath+v.Encode())
	//TODO error testing
	
	return nil
	
}
func (s *Server) AddAndPlay(uri string)(error){
	v := url.Values{}
	v.Set("command", "in_play")
	v.Set("input",uri)
	
	http.Get(s.addr+statusPath+v.Encode())
	//TODO error testing

	return nil
}

func (s *Server) Delete(uri string)(error){
	//Find item in playlist
	for _,item := range s.Playlist(){
		if item.URI == uri{
			v := url.Values{}
			v.Set("command", "pl_delete")
			v.Set("id",item.ID)
	
			http.Get(s.addr+statusPath+v.Encode())
			return nil
		}
	}
	
	return errors.New("Not Found")
}

func (s *Server) EmptyPlaylist()(error){
	v := url.Values{}
	v.Set("command", "pl_empty")
	
	http.Get(s.addr+statusPath+v.Encode())
	//TODO error testing
	
	return nil
}
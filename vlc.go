//Control of VLC through http interface
package vlc

import(
	"net/url"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"fmt"
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

const path = "/requests/status.json?"

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
	
	http.Get(s.addr+"/requests/status.json?"+v.Encode())
	//TODO error testing
	
	return nil	
}
func (s *Server) Stop()(error){
	v := url.Values{}
	v.Set("command", "pl_stop")
	
	http.Get(s.addr+"/requests/status.json?"+v.Encode())
	//TODO error testing
	
	return nil
}
func (s *Server) Next()(error){
	v := url.Values{}
	v.Set("command", "pl_next")
	
	http.Get(s.addr+"/requests/status.json?"+v.Encode())
	//TODO error testing
	
	return nil
}
func (s *Server) Previous()(error){
	v := url.Values{}
	v.Set("command", "pl_previous")
	
	http.Get(s.addr+"/requests/status.json?"+v.Encode())
	//TODO error testing
	
	return nil
}

//
// Playlist
//

func (s *Server) Playlist()([]Item){
	playlist,err:= http.Get(s.addr+"/requests/playlist.json")
	if err != nil{
		return nil
	}
	
	data,_ := ioutil.ReadAll(playlist.Body);
	
	var val *Item
	json.Unmarshal(data, val)
	
	if val == nil{
		return nil
	}else{
		fmt.Println()
		return val.Children[0].Children
	}		
}

func (s *Server) Add(uri string)(error){
	v := url.Values{}
	v.Set("command", "in_enqueue")
	v.Set("input",uri)
	
	http.Get(s.addr+"/requests/status.json?"+v.Encode())
	//TODO error testing
	
	return nil
	
}
func (s *Server) AddAndPlay(uri string)(error){
	v := url.Values{}
	v.Set("command", "in_play")
	v.Set("input",uri)
	
	http.Get(s.addr+"/requests/status.json?"+v.Encode())
	//TODO error testing

	return nil
}
func (s *Server) EmptyPlaylist()(error){
	v := url.Values{}
	v.Set("command", "pl_empty")
	
	http.Get(s.addr+"/requests/status.json?"+v.Encode())
	//TODO error testing
	
	return nil
}
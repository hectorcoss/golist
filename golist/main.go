package main

import (
	"appengine"
	"appengine/datastore"
	"html/template"
	"lib/validate"
	"net/http"
)

const GUESTBOOKTEMPLATEHTML = `
<html>
    <title>Anime List</title>
        <body>
        </br>
         <form action="/add" method="post">
            <table>
                <tr>
                    <td><label>Name: </label></td>
                    <td><input type="text" name="name"></td>
                </tr>
                <tr>
                    <td><label>Fansub: </label></td>
                    <td><input type="text" name="fansub"></td>
                </tr>
                <tr>
                    <td><label>Episode #: </label></td>
                    <td><input type="text" name="episode"></td>
                </tr>
                <tr>
                    <td><label>Filename: </label></td>
                    <td><input type="text" name="filename"></td>
                </tr>
            </table>
            <div><input type="submit" value="Save"></div>
         </form>
         </br>
         <hr>
         </br>
         <ul>
          <table>
         {{range .}}
             <tr>
                <td><li>{{.Name}}</td>
                <td>{{.Fansub}}</td>
                <td>{{.FileName}}</td>
                <td>{{.EpisodeNum}}</td>
                <td><input type="checkbox" name="cbxWatched"></li>
                </td></tr>
         {{end}}
          </table>
         </ul>
        </body>
</html>
`

var guestbookTemplate = template.Must(template.New("serie").Parse(GUESTBOOKTEMPLATEHTML))

type Serie struct {
	Name       string
	Fansub     string
	FileName   string
	EpisodeNum string
}

func init() {
	http.HandleFunc("/anime", root)
	http.HandleFunc("/add", add)
}

func root(w http.ResponseWriter, r *http.Request) {
	//    fmt.Fprintf(w, guestbookForm)
	c := appengine.NewContext(r)
	q := datastore.NewQuery("Serie").Order("Name").Limit(10)
	series := make([]Serie, 0, 10)
	if _, err := q.GetAll(c, &series); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := guestbookTemplate.Execute(w, series); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func add(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	x := "S"
	validate.ValidateString(&x)
	s := Serie{
		Name:     r.FormValue("name"),
		Fansub:   r.FormValue("fansub"),
		FileName: r.FormValue("filename"),
		//FileName:   x,
		EpisodeNum: r.FormValue("episode"),
	}
	_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Serie", nil), &s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

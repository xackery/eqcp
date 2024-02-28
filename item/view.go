package item

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/xackery/eqcp/config"
	"github.com/xackery/eqcp/tlog"

	"github.com/xackery/eqcp/content"
	"github.com/xackery/eqcp/db"
)

var (
	viewTemplate *template.Template
)

func init() {
	var err error
	viewTemplate = template.New("view")
	viewTemplate, err = viewTemplate.ParseFS(content.FS(), "template/item/view.tpl")
	if err != nil {
		tlog.Fatalf("template.ParseFS: %v", err)
		return
	}
}

// View handles item view requests
func View(w http.ResponseWriter, r *http.Request) {
	var err error

	var id int
	strID := r.URL.Query().Get("id")
	if len(strID) > 0 {
		id, err = strconv.Atoi(strID)
		if err != nil {
			tlog.Errorf("strconv.Atoi: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = viewRender(ctx, id, r.URL.Query().Get("name"), w)
	if err != nil {
		tlog.Errorf("viewRender: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func viewRender(ctx context.Context, id int, name string, w http.ResponseWriter) error {

	name = "%" + name + "%"
	query := "SELECT * FROM items WHERE name like :name LIMIT 1"
	if config.IsEnabled(config.IsDiscoveredOnly) {
		query += "SELECT * FROM items, discovered items WHERE items.name like :name AND discovered_items.item_id=items.id LIMIT 1"
	}

	if id > 0 {
		query = "SELECT * FROM items WHERE id=:id LIMIT 1"
		if config.IsEnabled(config.IsDiscoveredOnly) {
			query += "SELECT * FROM items, discovered items WHERE items.id=:id AND discovered_items.item_id=:id LIMIT 1"
		}
	}

	rows, err := db.Query(ctx,
		query,
		map[string]interface{}{
			"id":   id,
			"name": name,
		})
	if err != nil {
		return fmt.Errorf("query items: %w", err)
	}
	defer rows.Close()

	type TemplateData struct {
		Item *Table
	}

	data := TemplateData{}
	item := &Table{}

	if !rows.Next() {
		http.Error(w, "Not Found", http.StatusNotFound)
		return nil
	}

	err = rows.StructScan(item)
	if err != nil {
		return fmt.Errorf("rows.StructScan: %w", err)
	}
	data.Item = item

	buf := &bytes.Buffer{}

	err = viewTemplate.ExecuteTemplate(buf, "view.tpl", data)
	if err != nil {
		return fmt.Errorf("viewTemplate.Execute: %w", err)
	}

	_, err = w.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("w.Write: %w", err)
	}

	return nil
}

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	_ "reflect"
	_ "strings"

	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func updateSenderProfiles(s sdr, w http.ResponseWriter) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading file")
	}

	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}

	_, inserr := db.Exec("INSERT into sender VALUES($1,$2,$3,$4,$5) ON CONFLICT (sender_id) DO UPDATE SET sender_address = EXCLUDED.sender_address, sender_phone = EXCLUDED.sender_phone,  sender_email = EXCLUDED.sender_email, sender_name = EXCLUDED.sender_name, sender_id = EXCLUDED.sender_id; ", s.Address, s.Number, s.Email, s.Name, s.Id)

	if inserr != nil {
		panic(inserr)
	}
	defer db.Close()
}

func deletePreset(preset string, w http.ResponseWriter) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading file")
	}

	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}

	_, delerr := db.Exec("DELETE FROM sender WHERE sender_id=$1", preset)

	if delerr != nil {
		panic(delerr)
	}

}

// Struct definitions to match the database
type sender struct {
	Sender_name string
	Address     string
	Phone       string
	Email       string
	Sender_id   int
}
type receipient struct {
	Receipient_name string
	Address         string
	Phone           string
}

type invoice struct {
	Invoice_date string
	Subtotal     float32
	Tax          float32
	Shipping     float32
	Total        float32
}
type items struct {
	Description string
	Ppu         float32
	Qty         int
	Total       float32
}
type collection struct {
	S        sender
	R        receipient
	Inv      invoice
	Items    []items
	Pre      int
	Invn     string
	Notes    string
	Currency string
}
type pcollection struct {
	Sender_name string
	Sender_id   int
	info        []data
}
type data struct {
	Invoice_date string     `json:"invoice_date"`
	R            receipient `json:"receipient"`
	Invn         string     `json:"invn"`
	Total        float32    `json:"total"`
}

func toDb(clr collection, w http.ResponseWriter) error {

	err := godotenv.Load()
	if err != nil {

		return err
	}

	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return err
	}

	var inserr error

	// invoice_string := fmt.Sprintf("INSERT INTO invoice VALUES(%d, %s, %s, %f, %f,%f,%f)", clr.Pre, clr.Invn, clr.Inv.Invoice_date, clr.Inv.Total, clr.Inv.Subtotal, clr.Inv.Tax, clr.Inv.Shipping)

	fmt.Print("Here")
	_, inserr = db.Exec("INSERT INTO invoice VALUES($1, $2, $3, $4, $5, $6, $7);", clr.Inv.Invoice_date, clr.Inv.Total, clr.Inv.Subtotal, clr.Inv.Tax, clr.Inv.Shipping, clr.Pre, clr.Invn)

	if inserr != nil {
		return inserr
	}
	_, inserr = db.Exec("INSERT INTO recipient VALUES($1, $2, $3, $4, $5);", clr.R.Receipient_name, clr.R.Address, clr.R.Phone, clr.Pre, clr.Invn)
	if inserr != nil {
		return inserr
	}

	//Optimize this, but doing one insert statement
	for _, item := range clr.Items {
		_, inserr = db.Exec("INSERT INTO items VALUES($1, $2, $3, $4, $5, $6);", item.Description, item.Qty, item.Ppu, item.Total, clr.Pre, clr.Invn)
		if inserr != nil {
			return inserr
		}
	}

	defer db.Close()

	return nil

}

// var templates = template.Must(template.ParseFiles("templates/index.html", "templates/analytics.html"))
var templates = template.Must(template.ParseGlob("templates/*.html"))

func renderTemplate(w http.ResponseWriter, dir string, data any) {
	err := templates.ExecuteTemplate(w, dir, data)
	//Handling errors associated with ParseFiles method
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type sdr struct {
	Address string `json:"address"`
	Number  string `json:"number"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Id      string `json:"id"`
}
type q struct {
	Name string
	Id   int
}

func initt(w http.ResponseWriter) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading file")
	}

	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}

	defer db.Close()
	rows, e := db.Query(`SELECT * FROM sender`)
	if e != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}

	var data sdr
	//Create a "slice" with the make function, to form a growable array of sender structs
	a := make([]sdr, 0)

	for rows.Next() {
		read_err := rows.Scan(&data.Address, &data.Number, &data.Email, &data.Name, &data.Id)
		//Error handling for issue with database reads
		if read_err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}

		defer rows.Close()
		a = append(a, data)

	}

	str, re := json.Marshal(a)
	if re != nil {
		fmt.Print(re)
	}

	//Setup response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, writerr := w.Write(str)

	//Error handling for issues with writing the response
	if writerr != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
}

func main() {
	//Registering route handlers
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	//Setup as the root directory for the app, rejects any request that is not defined in the server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		if r.URL.Path != "/" {
			http.Error(w, "404 not found", http.StatusInternalServerError)
			return
		}
		renderTemplate(w, "index.html", nil)
	})

	//http package provides the HandleFunc method, which accepts a path and a handler containing response and request information.
	http.HandleFunc("/analytics", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "analytics.html", nil)
	})

	//make sure this route only accepts get requests
	http.HandleFunc("/presets", func(w http.ResponseWriter, r *http.Request) {
		initt(w)
	})

	http.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		var rec receipient
		var inv invoice
		inv.Subtotal = 0
		var ite items
		var sen sender
		var ds collection

		var err error
		var pre int
		var invn string

		//Storing receipient information
		rec.Receipient_name = r.FormValue("receiver_Name")
		rec.Address = r.FormValue("receiver_address")
		rec.Phone = r.FormValue("receiver_number")
		_, err = fmt.Sscan(r.Form.Get("preset"), &pre)
		if err != nil {
			panic(err)
		}
		_, err = fmt.Sscan(r.Form.Get("invoice_num"), &invn)
		if err != nil {
			panic(err)
		}

		//Storing the invoice items in the struct
		a := make([]items, 0)
		lent := len(r.Form["desc"])

		for i := 0; i < lent; i++ {
			ite.Description = r.Form["desc"][i]

			_, err = fmt.Sscan(r.Form["qty"][i], &ite.Qty)
			if err != nil {
				http.Error(w, err.Error(), 406)
			}
			_, err = fmt.Sscan(r.Form["cpu"][i], &ite.Ppu)

			if err != nil {
				http.Error(w, err.Error(), 406)
			}

			ite.Total = float32(ite.Qty) * ite.Ppu

			inv.Subtotal = inv.Subtotal + ite.Total
			a = append(a, ite)
		}

		//Storing invoice informtion
		inv.Invoice_date = r.FormValue("invoice_date")
		_, err = fmt.Sscan(r.Form.Get("tax"), &inv.Tax)
		if err != nil {
			panic(err)
		}
		_, err = fmt.Sscan(r.Form.Get("shipping"), &inv.Shipping)
		if err != nil {
			panic(err)
		}
		inv.Invoice_date = r.FormValue("invoice_date")
		inv.Total = (inv.Subtotal * inv.Tax) + inv.Subtotal + inv.Shipping

		//Storing sender information
		sen.Sender_name = r.FormValue("sender_name")
		sen.Address = r.FormValue("sender_address")
		sen.Email = r.FormValue("sender_email")
		sen.Phone = r.FormValue("sender_number")
		sen.Sender_id = pre

		//Storing collection to send back to templating
		ds.S = sen
		ds.Items = a
		ds.R = rec
		ds.Inv = inv
		ds.Pre = pre
		ds.Invn = invn
		ds.Notes = r.Form.Get("extra_notes")
		ds.Currency = r.Form.Get("currency")

		//Issue, refreshing the page resubmits the form. Submitting the same data causes the database to raise the constaint key error. Simply capturing the error, but database is still opened in the process
		errs := toDb(ds, w)

		//Consider using fetch request to mitigate having to open the database unecessarily.. Optimize
		fmt.Print(errs)

		switch r.Method {
		case "POST":
			renderTemplate(w, "printout.html", ds)

		default:
			http.Error(w, "404 not found", http.StatusInternalServerError)
		}

	})

	http.HandleFunc("/loadhistory", func(w http.ResponseWriter, r *http.Request) {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading file")
		}

		connStr := os.Getenv("DATABASE_URL")
		db, err := sql.Open("postgres", connStr)

		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}

		defer db.Close()

		presetData, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		preset := string(presetData)

		b := make([]data, 0)
		var c data

		//Join based on the invoice number and the sender id to retrieve all invoices (and recipients) associated with a sender.
		rows, err := db.Query(`SELECT invoice.invoice_num, invoice_date, invoice_total, recipient_name, recipient_address, recipient_phone FROM invoice JOIN recipient on invoice.invoice_num = recipient.invoice_num AND invoice.sender_id = recipient.sender_id AND invoice.sender_id = $1;`, preset)

		//This will store all invoices and recipient data in an array of structs. Each array instance is an invoice.
		for rows.Next() {
			read_err := rows.Scan(&c.Invn, &c.Invoice_date, &c.Total, &c.R.Receipient_name, &c.R.Address, &c.R.Phone)
			//Error handling for issue with database reads
			if read_err != nil {
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
			}

			defer rows.Close()
			b = append(b, c)
			if err != nil {
				panic(err)
			}
		}

		str, re := json.Marshal(b)

		if re != nil {
			fmt.Print(re)
		}
		//Setup response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, writerr := w.Write(str)

		//Error handling for issues with writing the response
		if writerr != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}

	})

	//make sure this route only accepts get requests
	http.HandleFunc("/sdrUpdate", func(w http.ResponseWriter, r *http.Request) {

		var data sdr
		d := json.NewDecoder(r.Body)
		d.DisallowUnknownFields()

		err := d.Decode(&data)

		if err != nil {
			panic(err)
		}
		updateSenderProfiles(data, w)
	})

	http.HandleFunc("/deletePreset", func(w http.ResponseWriter, r *http.Request) {

		presetData, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		preset := string(presetData)
		deletePreset(preset, w)

	})
	http.HandleFunc("/history", func(w http.ResponseWriter, r *http.Request) {

		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading file")
		}

		connStr := os.Getenv("DATABASE_URL")
		db, err := sql.Open("postgres", connStr)

		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}

		defer db.Close()

		//Retrieve all senders from the database with simple query
		rows, e := db.Query(`SELECT sender_name, sender_id FROM sender`)
		if e != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}

		a := make([]q, 0)
		var d q

		for rows.Next() {
			//Only store the name and id/preset
			read_err := rows.Scan(&d.Name, &d.Id)
			//Error handling for issue with database reads
			if read_err != nil {
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
			}
			defer rows.Close()
			a = append(a, d)

		}

		renderTemplate(w, "history.html", a)
	})

	//Sets up port for listening

	fmt.Println("Server Started")
	log.Fatal(http.ListenAndServe("127.0.0.1:3000", nil))
}

// http.HandleFunc("/archive", func(w http.ResponseWriter, r *http.Request) {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading file")
// 	}

// 	connStr := os.Getenv("DATABASE_URL")
// 	db, err := sql.Open("postgres", connStr)

// 	if err != nil {
// 		http.Error(w, "Something went wrong", http.StatusInternalServerError)
// 	}

// 	defer db.Close()

// 	//Retrieve all senders from the database with simple query
// 	rows, e := db.Query(`SELECT sender_name, sender_id FROM sender`)
// 	if e != nil {
// 		http.Error(w, "Something went wrong", http.StatusInternalServerError)
// 	}

// 	var name string
// 	var id int
// 	//Create a "slice" with the make function, to form a growable array of sender structs

// 	//Create a slice of partial collections. Each partial collection represents the information needed for a specific sender profile (to be displayed on the summary cards)
// 	a := make([]pcollection, 0)

// 	//For each row in the sender tables
// 	for rows.Next() {
// 		//Only store the name and id/preset
// 		read_err := rows.Scan(&name, &id)
// 		//Error handling for issue with database reads
// 		if read_err != nil {
// 			http.Error(w, "Something went wrong", http.StatusInternalServerError)
// 		}
// 		defer rows.Close()

// 		//Append a new partial collection for each sender
// 		a = append(a, pcollection{name, id, nil})
// 	}

// 	b := make([]data, 0)
// 	var c data

// 	//Loop through partial collections
// 	for i := 0; i < len(a); i++ {

// 		//Join based on the invoice number and the sender id to retrieve all invoices (and recipients) associated with a sender.
// 		rows, err := db.Query(`SELECT invoice.invoice_num, invoice_date, invoice_total, recipient_name, recipient_address, recipient_phone FROM invoice JOIN recipient on invoice.invoice_num = recipient.invoice_num AND invoice.sender_id = recipient.sender_id AND invoice.sender_id = $1;`, a[i].Sender_id)

// 		//This will store all invoices and recipient data in an array of structs. Each array instance is an invoice.
// 		for rows.Next() {
// 			read_err := rows.Scan(&c.Invn, &c.Invoice_date, &c.total, &c.receipient.Receipient_name, &c.receipient.Address, &c.receipient.Phone)
// 			//Error handling for issue with database reads
// 			if read_err != nil {
// 				http.Error(w, "Something went wrong", http.StatusInternalServerError)
// 			}

// 			defer rows.Close()
// 			b = append(b, c)

// 		}

// 		//Store all invoices in each partial collection. A.K.A, store every invoice for every user.
// 		a[i].info = b
// 		if err != nil {
// 			panic(err)
// 		}
// 	}

// 	str, re := json.Marshal(a[0].info)
// 	if re != nil {
// 		fmt.Print(re)
// 	}

// 	//Setup response
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	_, writerr := w.Write(str)

// 	//Error handling for issues with writing the response
// 	if writerr != nil {
// 		http.Error(w, "Something went wrong", http.StatusInternalServerError)
// 	}

// })

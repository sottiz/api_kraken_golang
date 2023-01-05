package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func createXMLFile(myAssets []AssetPair) {
	err := os.Mkdir("Archive", 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	// Create a file
	file, err := os.Create("Archive/asset_data.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	systemStat, err := getSystemStatus()
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.WriteString(fmt.Sprintf("<Data>\n"))
	_, err = file.WriteString(fmt.Sprintf("\t<Time>\n\t\t%s\n\t</Time>\n", systemStat.Timestamp))
	_, err = file.WriteString(fmt.Sprintf("\t<SystemStatus>\n\t\t%s\n\t</SystemStatus>\n", systemStat.Status))
	_, err = file.WriteString(fmt.Sprintf("\t<AssetCount>\n\t\t%d\n\t</AssetCount>\n", len(myAssets)))
	_, err = file.WriteString(fmt.Sprintf("\t<Assets>\n"))
	for _, asset := range myAssets {
		ticker := getTickerInfos(asset.Altname)
		_, err = file.WriteString(fmt.Sprintf("\t\t<Asset>\n"))
		_, err = file.WriteString(fmt.Sprintf("\t\t\t<Base>%s</Base>\n", asset.Base))
		_, err = file.WriteString(fmt.Sprintf("\t\t\t<Quote>%s</Quote>\n", asset.Quote))
		_, err = file.WriteString(fmt.Sprintf("\t\t\t<Altname>%s</Altname>\n", asset.Altname))
		_, err = file.WriteString(fmt.Sprintf("\t\t\t<Wsname>%s</Wsname>\n", asset.Wsname))
		_, err = file.WriteString(fmt.Sprintf("\t\t\t<Price>%f</Price>\n", ticker.Price))
		_, err = file.WriteString(fmt.Sprintf("\t\t\t<Volume>%f</Volume>\n", ticker.Volume))
		_, err = file.WriteString(fmt.Sprintf("\t\t\t<NbOfTrades>%d</NbOfTrades>\n", ticker.NbOfTrades))
		_, err = file.WriteString(fmt.Sprintf("\t\t\t<LowPrice>%f</LowPrice>\n", ticker.LowPrice))
		_, err = file.WriteString(fmt.Sprintf("\t\t\t<HighPrice>%f</HighPrice>\n", ticker.HighPrice))
		_, err = file.WriteString(fmt.Sprintf("\t\t</Asset>\n"))
	}
	_, err = file.WriteString(fmt.Sprintf("</Data>"))

	if err != nil {
		log.Fatal(err)
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// Ouvrez le fichier à télécharger
	file, err := os.Open("./Archive/asset_data.xml")
	if err != nil {
		http.Error(w, "Cannot open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Définissez les en-têtes HTTP pour indiquer au client que c'est un téléchargement
	w.Header().Set("Content-Disposition", "attachment; filename=asset_data.xml")
	w.Header().Set("Content-Type", "text/plain")

	// Copiez le contenu du fichier dans le corps de la réponse HTTP
	io.Copy(w, file)
}

func main() {
	// Connexion à la base de données
	connectionString := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	myAssets := getAssets()

	createXMLFile(myAssets)

	// dropTables(db)
	createTables(db)

	writeDb(db, myAssets)

	http.HandleFunc("/assets", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM ASSETS")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		assets := []AssetPair{}
		for rows.Next() {
			asset := AssetPair{}
			err := rows.Scan(&asset.Base, &asset.Quote, &asset.Altname, &asset.Wsname)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			assets = append(assets, asset)
		}

		json.NewEncoder(w).Encode(assets)
	})

	http.HandleFunc("/download", downloadHandler)
	http.ListenAndServe(":8080", nil)

	fmt.Println("Listening on port 8080")

}

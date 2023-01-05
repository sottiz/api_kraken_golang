package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "toto"
	password = "mysecretpassword"
	dbname   = "mydatabase"
	schema   = "public"
)

func createTables(db *sql.DB) {
	sqlStat := "CREATE TABLE IF NOT EXISTS public.SYSTEM ( id SERIAL NOT NULL, TIMESTAMP date NOT NULL, STATUS character varying NOT NULL, PRIMARY KEY (id, TIMESTAMP) ); ALTER TABLE IF EXISTS public.SYSTEM OWNER to toto;"
	_, errors := db.Exec(sqlStat)
	if errors != nil {
		fmt.Println(errors)
	}

	sqlStat = "CREATE TABLE IF NOT EXISTS public.ASSETS ( id SERIAL NOT NULL, TIMESTAMP_ID integer NOT NULL, BASE character varying NOT NULL, QUOTE character varying NOT NULL, ALTNAME character varying NOT NULL, WSNAME character varying NOT NULL, PRICE real NOT NULL, VOLUME real NOT NULL, NB_OF_TRADES integer NOT NULL, LOW_PRICE real NOT NULL, HIGH_PRICE real NOT NULL, PRIMARY KEY (id) ); ALTER TABLE IF EXISTS public.ASSETS OWNER to toto;"
	_, errors = db.Exec(sqlStat)
	if errors != nil {
		fmt.Println(errors)
	}
}

func dropTables(db *sql.DB) {
	sqlStat := "DROP TABLE IF EXISTS public.SYSTEM; DROP TABLE IF EXISTS public.ASSETS;"
	_, errors := db.Exec(sqlStat)
	if errors != nil {
		fmt.Println(errors)
	}
}

func writeDb(db *sql.DB, myAssets []AssetPair) {
	systemStat, err := getSystemStatus()
	if err != nil {
		log.Fatal(err)
	}

	// Insertion de la date dans la table SYSTEM
	sqlStat := fmt.Sprintf("INSERT INTO public.SYSTEM (TIMESTAMP, STATUS) VALUES ('%s', '%s');", systemStat.Timestamp, systemStat.Status)
	_, errors := db.Exec(sqlStat)
	if errors != nil {
		fmt.Println(errors)
	}

	// Récupération de l'id de la date insérée
	var id int
	sqlStat = fmt.Sprintf("SELECT id FROM public.SYSTEM WHERE TIMESTAMP = '%s';", systemStat.Timestamp)
	err = db.QueryRow(sqlStat).Scan(&id)
	if err != nil {
		fmt.Println(err)
	}

	// Insertion des données dans la table ASSETS
	for _, asset := range myAssets {
		ticker := getTickerInfos(asset.Altname)
		sqlStat = fmt.Sprintf("INSERT INTO public.ASSETS (TIMESTAMP_ID, BASE, QUOTE, ALTNAME, WSNAME, PRICE, VOLUME, NB_OF_TRADES, LOW_PRICE, HIGH_PRICE) VALUES (%d, '%s', '%s', '%s', '%s', %f, %f, %d, %f, %f);", id, asset.Base, asset.Quote, asset.Altname, asset.Wsname, ticker.Price, ticker.Volume, ticker.NbOfTrades, ticker.LowPrice, ticker.HighPrice)
		_, errors = db.Exec(sqlStat)
		if errors != nil {
			fmt.Println(errors)
		}
	}
}

package main

import (
	"context"
	"log"
	"my_finance/internal/database"
	loggerHelper "my_finance/internal/logger"
	r "my_finance/internal/routes"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func createDefaultCategory(db *pgxpool.Pool) {
	var existis bool

	if err := db.QueryRow(
		context.Background(),
		`
			SELECT EXISTS
				(
					SELECT	
						id
					FROM
						categories
				)
		`,
	).Scan(&existis); err != nil {
		loggerHelper.InfoLogger.Fatal("Erro ao criar categoria padrão:", err)
	}

	if !existis {
		loggerHelper.InfoLogger.Println("Categoria não cadastrada.")

		if _, err := db.Exec(
			context.Background(),
			`
				INSERT INTO categories
					(name, type)
				VALUES
					('Recebimentos', 'Entrada')
			`,
		); err != nil {
			loggerHelper.InfoLogger.Fatal("Erro ao criar categoria padrão:", err)
		}
	}

	loggerHelper.InfoLogger.Println("Categoria já cadastrada.")
}

func main() {
	loggerHelper.Logger()

	wd, _ := os.Getwd()

	log.Println("Rodando em:", wd)

	dbConn, err := database.Init()

	if err != nil {
		loggerHelper.InfoLogger.Fatal("Erro ao conectar ao banco de dados:", err)
	}

	defer dbConn.Close()

	createDefaultCategory(dbConn)

	loggerHelper.InfoLogger.Println("Banco de dados conectado com sucesso!")
	r.StartServer(dbConn) // Start server
}

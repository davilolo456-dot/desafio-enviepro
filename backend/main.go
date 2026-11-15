package main

import (
	"database/sql"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// ---------- Formatos de resposta (viram JSON) ----------

type Cliente struct {
	ID             int    `json:"id"`
	Nome           string `json:"nome"`
	Empresa        string `json:"empresa"`
	Email          string `json:"email"`
	Telefone       string `json:"telefone"`
	Status         string `json:"status"`
	Responsavel    string `json:"responsavel"`
	UltimoContato  string `json:"ultimo_contato"`
	DiasSemContato int    `json:"dias_sem_contato"`
}

type EtapaFunil struct {
	Etapa      string  `json:"etapa"`
	Quantidade int     `json:"quantidade"`
	Valor      float64 `json:"valor"`
}

type ResumoVendas struct {
	Mes            string  `json:"mes"`
	Total          float64 `json:"total"`
	Quantidade     int     `json:"quantidade"`
	TicketMedio    float64 `json:"ticket_medio"`
	ClientesAtivos int     `json:"clientes_ativos"`
}

type Venda struct {
	ID          int     `json:"id"`
	ClienteID   int     `json:"cliente_id"`
	ClienteNome string  `json:"cliente_nome"`
	Valor       float64 `json:"valor"`
	Data        string  `json:"data"`
	Produto     string  `json:"produto"`
}

func main() {
	db, err := abrirBanco()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Modo seed:  go run . seed        (insere dados se o banco estiver vazio)
	//             go run . seed reset  (apaga tudo e insere de novo)
	if len(os.Args) > 1 && os.Args[1] == "seed" {
		reset := len(os.Args) > 2 && os.Args[2] == "reset"
		if err := executarSeed(db, reset); err != nil {
			log.Fatal(err)
		}
		return
	}

	// Modo servidor: go run .
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// CORS: o navegador só deixa o frontend (localhost:3000) chamar esta API
	// (localhost:8080) se a API avisar que permite. É isso que esta linha faz.
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods: []string{"GET", "OPTIONS"},
		AllowHeaders: []string{"Content-Type"},
	}))

	r.GET("/health", health(db))
	r.GET("/clientes", listarClientes(db))
	r.GET("/funil", funilPorEtapa(db))
	r.GET("/vendas/resumo", resumoVendas(db))
	r.GET("/vendas", listarVendas(db))

	log.Println("API rodando em http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func erroInterno(c *gin.Context, err error) {
	log.Println("erro:", err)
	c.JSON(http.StatusInternalServerError, gin.H{"erro": "erro interno do servidor"})
}

// GET /health: confirma que a API e o banco estão no ar.
func health(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "erro", "banco": "indisponivel"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "banco": "conectado"})
	}
}

// GET /clientes: lista de clientes. Filtros opcionais: ?busca=nome&status=Ativo
func listarClientes(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		busca := c.Query("busca")
		status := c.Query("status")

		rows, err := db.Query(`
			SELECT id, nome, empresa, email, COALESCE(telefone, ''), status, responsavel, ultimo_contato,
			       CAST(julianday(date('now','localtime')) - julianday(ultimo_contato) AS INTEGER)
			FROM clientes
			WHERE (? = '' OR LOWER(nome) LIKE '%' || LOWER(?) || '%')
			  AND (? = '' OR status = ?)
			ORDER BY nome`, busca, busca, status, status)
		if err != nil {
			erroInterno(c, err)
			return
		}
		defer rows.Close()

		clientes := make([]Cliente, 0)
		for rows.Next() {
			var cl Cliente
			if err := rows.Scan(&cl.ID, &cl.Nome, &cl.Empresa, &cl.Email, &cl.Telefone,
				&cl.Status, &cl.Responsavel, &cl.UltimoContato, &cl.DiasSemContato); err != nil {
				erroInterno(c, err)
				return
			}
			clientes = append(clientes, cl)
		}
		if err := rows.Err(); err != nil {
			erroInterno(c, err)
			return
		}
		c.JSON(http.StatusOK, clientes)
	}
}

// GET /funil: quantidade e valor por etapa, sempre na ordem do funil.
func funilPorEtapa(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(`
			SELECT etapa_funil, COUNT(*), COALESCE(SUM(valor), 0)
			FROM oportunidades
			GROUP BY etapa_funil`)
		if err != nil {
			erroInterno(c, err)
			return
		}
		defer rows.Close()

		porEtapa := map[string]EtapaFunil{}
		for rows.Next() {
			var e EtapaFunil
			if err := rows.Scan(&e.Etapa, &e.Quantidade, &e.Valor); err != nil {
				erroInterno(c, err)
				return
			}
			porEtapa[e.Etapa] = e
		}
		if err := rows.Err(); err != nil {
			erroInterno(c, err)
			return
		}

		// Monta a lista na ordem Lead -> Contato -> Proposta -> Negociação -> Fechado.
		funil := make([]EtapaFunil, 0, len(etapasFunil))
		for _, nome := range etapasFunil {
			e, ok := porEtapa[nome]
			if !ok {
				e = EtapaFunil{Etapa: nome}
			}
			funil = append(funil, e)
		}
		c.JSON(http.StatusOK, funil)
	}
}

// GET /vendas/resumo: total, quantidade e ticket médio do mês atual + clientes ativos.
func resumoVendas(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var r ResumoVendas
		err := db.QueryRow(`
			SELECT strftime('%Y-%m', 'now', 'localtime'),
			       COALESCE(SUM(valor), 0),
			       COUNT(*)
			FROM vendas
			WHERE strftime('%Y-%m', data) = strftime('%Y-%m', 'now', 'localtime')`).
			Scan(&r.Mes, &r.Total, &r.Quantidade)
		if err != nil {
			erroInterno(c, err)
			return
		}
		if r.Quantidade > 0 {
			r.TicketMedio = math.Round(r.Total/float64(r.Quantidade)*100) / 100
		}

		if err := db.QueryRow(`SELECT COUNT(*) FROM clientes WHERE status = 'Ativo'`).
			Scan(&r.ClientesAtivos); err != nil {
			erroInterno(c, err)
			return
		}
		c.JSON(http.StatusOK, r)
	}
}

// GET /vendas: todas as vendas, da mais recente para a mais antiga.
func listarVendas(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(`
			SELECT v.id, v.cliente_id, cl.nome, v.valor, v.data, v.produto
			FROM vendas v
			JOIN clientes cl ON cl.id = v.cliente_id
			ORDER BY v.data DESC, v.id DESC`)
		if err != nil {
			erroInterno(c, err)
			return
		}
		defer rows.Close()

		vendas := make([]Venda, 0)
		for rows.Next() {
			var v Venda
			if err := rows.Scan(&v.ID, &v.ClienteID, &v.ClienteNome, &v.Valor, &v.Data, &v.Produto); err != nil {
				erroInterno(c, err)
				return
			}
			vendas = append(vendas, v)
		}
		if err := rows.Err(); err != nil {
			erroInterno(c, err)
			return
		}
		c.JSON(http.StatusOK, vendas)
	}
}

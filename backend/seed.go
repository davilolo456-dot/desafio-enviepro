package main

import (
	"database/sql"
	"fmt"
	"time"
)

type clienteSeed struct {
	nome, empresa, email, telefone, status, responsavel string
	diasSemContato                                      int
}

var clientesSeed = []clienteSeed{
	{"Mariana Alves", "Padaria Pão Dourado", "mariana@paodourado.com.br", "(31) 99111-0001", "Ativo", "Ana Souza", 2},
	{"Carlos Pereira", "Oficina Pereira", "carlos@oficinapereira.com.br", "(31) 99111-0002", "Ativo", "Bruno Lima", 5},
	{"Juliana Costa", "Clínica Sorriso", "juliana@clinicasorriso.com.br", "(31) 99111-0003", "Prospecto", "Carla Mendes", 9},
	{"Roberto Nunes", "Nunes Contabilidade", "roberto@nunescontabil.com.br", "(31) 99111-0004", "Ativo", "Ana Souza", 20},
	{"Fernanda Rocha", "Studio Rocha", "fernanda@studiorocha.com.br", "(31) 99111-0005", "Em negociação", "Bruno Lima", 3},
	{"Paulo Menezes", "Menezes Transportes", "paulo@menezestransp.com.br", "(31) 99111-0006", "Inativo", "Carla Mendes", 45},
	{"Aline Barros", "Barros Moda", "aline@barrosmoda.com.br", "(31) 99111-0007", "Ativo", "Ana Souza", 7},
	{"Eduardo Lopes", "Lopes Engenharia", "eduardo@lopeseng.com.br", "(31) 99111-0008", "Em negociação", "Bruno Lima", 1},
	{"Camila Duarte", "Duarte Café", "camila@duartecafe.com.br", "(31) 99111-0009", "Prospecto", "Carla Mendes", 16},
	{"Ricardo Tavares", "Tavares Móveis", "ricardo@tavaresmoveis.com.br", "(31) 99111-0010", "Ativo", "Ana Souza", 12},
	{"Patrícia Gomes", "Gomes Advocacia", "patricia@gomesadv.com.br", "(31) 99111-0011", "Ativo", "Bruno Lima", 4},
	{"Marcelo Ribeiro", "Ribeiro Auto Peças", "marcelo@ribeiroautopecas.com.br", "(31) 99111-0012", "Inativo", "Carla Mendes", 60},
	{"Larissa Martins", "Martins Pet Shop", "larissa@martinspet.com.br", "(31) 99111-0013", "Prospecto", "Ana Souza", 18},
	{"Gustavo Ferraz", "Ferraz Tecnologia", "gustavo@ferraztec.com.br", "(31) 99111-0014", "Em negociação", "Bruno Lima", 6},
	{"Beatriz Cardoso", "Cardoso Eventos", "beatriz@cardosoeventos.com.br", "(31) 99111-0015", "Ativo", "Carla Mendes", 10},
	{"Thiago Moreira", "Moreira Construções", "thiago@moreiraconst.com.br", "(31) 99111-0016", "Ativo", "Ana Souza", 25},
	{"Renata Freitas", "Freitas Farmácia", "renata@freitasfarma.com.br", "(31) 99111-0017", "Prospecto", "Bruno Lima", 8},
	{"Diego Azevedo", "Azevedo Academia", "diego@azevedofit.com.br", "(31) 99111-0018", "Em negociação", "Carla Mendes", 2},
	{"Vanessa Teixeira", "Teixeira Papelaria", "vanessa@teixeirapapel.com.br", "(31) 99111-0019", "Inativo", "Ana Souza", 38},
	{"Leandro Castro", "Castro Mercado", "leandro@castromercado.com.br", "(31) 99111-0020", "Ativo", "Bruno Lima", 14},
	{"Silvia Barbosa", "Barbosa Educação", "silvia@barbosaedu.com.br", "(31) 99111-0021", "Prospecto", "Carla Mendes", 17},
	{"André Siqueira", "Siqueira Imóveis", "andre@siqueiraimoveis.com.br", "(31) 99111-0022", "Ativo", "Ana Souza", 3},
}

type oportunidadeSeed struct {
	clienteIdx int // posição do cliente em clientesSeed (0 = primeiro)
	titulo     string
	etapa      string
	valor      float64
	diasAtras  int // há quantos dias a oportunidade foi criada
}

// 30 oportunidades: 6 em cada etapa do funil.
var oportunidadesSeed = []oportunidadeSeed{
	// Lead
	{2, "Gestão de agendamentos da clínica", "Lead", 4500, 2},
	{8, "Cardápio digital e delivery", "Lead", 3200, 4},
	{12, "Loja virtual para pet shop", "Lead", 6800, 1},
	{16, "Sistema de estoque da farmácia", "Lead", 5400, 6},
	{20, "Plataforma de matrículas", "Lead", 9800, 3},
	{5, "Rastreamento de frota", "Lead", 7200, 8},
	// Contato
	{0, "Programa de fidelidade", "Contato", 3800, 5},
	{6, "Catálogo online de coleção", "Contato", 4100, 7},
	{9, "Showroom virtual", "Contato", 8600, 9},
	{14, "Gestão de eventos corporativos", "Contato", 6300, 10},
	{18, "Reativação de contrato de papelaria", "Contato", 2900, 12},
	{21, "Portal de captação de imóveis", "Contato", 11500, 4},
	// Proposta
	{1, "Sistema de ordens de serviço", "Proposta", 5200, 14},
	{3, "Automação de relatórios contábeis", "Proposta", 7900, 15},
	{10, "Gestão de processos jurídicos", "Proposta", 12400, 11},
	{15, "Controle de obras e medições", "Proposta", 15800, 16},
	{19, "PDV e controle de caixa", "Proposta", 6700, 13},
	{11, "Reativação: integração de estoque", "Proposta", 4300, 20},
	// Negociação
	{4, "Site institucional e agenda", "Negociação", 5900, 18},
	{7, "Painel de projetos de engenharia", "Negociação", 18900, 17},
	{13, "Suporte Premium anual", "Negociação", 9400, 19},
	{17, "App de treinos e check-in", "Negociação", 8200, 21},
	{0, "Plano Empresarial", "Negociação", 14500, 22},
	{9, "Treinamento da equipe de vendas", "Negociação", 3600, 24},
	// Fechado
	{1, "Plano Profissional", "Fechado", 6400, 26},
	{6, "Consultoria de e-commerce", "Fechado", 7800, 28},
	{10, "Suporte Premium", "Fechado", 5500, 30},
	{14, "Plano Básico", "Fechado", 2400, 33},
	{19, "Treinamento de operadores", "Fechado", 3300, 36},
	{21, "Consultoria de lançamento", "Fechado", 10200, 40},
}

var etapasFunil = []string{"Lead", "Contato", "Proposta", "Negociação", "Fechado"}

var produtosSeed = []string{
	"Plano Básico", "Plano Profissional", "Plano Empresarial",
	"Consultoria", "Treinamento", "Suporte Premium",
}

const dataFmt = "2006-01-02"

// executarSeed popula o banco com dados fictícios.
// Se o banco já tiver clientes, não faz nada (evita duplicar).
// Com reset=true, apaga tudo antes de inserir de novo.
func executarSeed(db *sql.DB, reset bool) error {
	if reset {
		for _, tabela := range []string{"vendas", "oportunidades", "clientes"} {
			if _, err := db.Exec("DELETE FROM " + tabela); err != nil {
				return fmt.Errorf("limpar %s: %w", tabela, err)
			}
		}
		if _, err := db.Exec("DELETE FROM sqlite_sequence"); err != nil {
			return fmt.Errorf("zerar contadores: %w", err)
		}
		fmt.Println("dados antigos apagados")
	}

	var total int
	if err := db.QueryRow("SELECT COUNT(*) FROM clientes").Scan(&total); err != nil {
		return fmt.Errorf("contar clientes: %w", err)
	}
	if total > 0 {
		fmt.Printf("banco já tem %d clientes; nada foi inserido (use \"seed reset\" para recriar)\n", total)
		return nil
	}

	agora := time.Now()
	inicioMes := time.Date(agora.Year(), agora.Month(), 1, 0, 0, 0, 0, agora.Location())
	diaAtual := agora.Day()

	// Tudo dentro de uma transação: ou entra tudo, ou nada.
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("iniciar transação: %w", err)
	}
	defer tx.Rollback()

	ids := make([]int64, 0, len(clientesSeed))
	for _, c := range clientesSeed {
		ultimo := agora.AddDate(0, 0, -c.diasSemContato).Format(dataFmt)
		res, err := tx.Exec(
			`INSERT OR IGNORE INTO clientes (nome, empresa, email, telefone, status, responsavel, ultimo_contato)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			c.nome, c.empresa, c.email, c.telefone, c.status, c.responsavel, ultimo,
		)
		if err != nil {
			return fmt.Errorf("inserir cliente %s: %w", c.nome, err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("id do cliente %s: %w", c.nome, err)
		}
		ids = append(ids, id)
	}

	// Oportunidades: cada uma ligada a um cliente (cliente_id) pela lista acima.
	for i, o := range oportunidadesSeed {
		criacao := agora.AddDate(0, 0, -o.diasAtras).Format(dataFmt)
		if _, err := tx.Exec(
			`INSERT OR IGNORE INTO oportunidades (cliente_id, titulo, etapa_funil, valor, data_criacao)
			 VALUES (?, ?, ?, ?, ?)`,
			ids[o.clienteIdx], o.titulo, o.etapa, o.valor, criacao,
		); err != nil {
			return fmt.Errorf("inserir oportunidade %d: %w", i+1, err)
		}
	}

	// 25 vendas: 15 no mês atual e 10 no mês anterior.
	for i := 0; i < 25; i++ {
		clienteID := ids[(i*3)%len(ids)]
		produto := produtosSeed[i%len(produtosSeed)]
		valor := float64(1500 + (i*937)%9000)

		var data time.Time
		if i < 15 {
			data = inicioMes.AddDate(0, 0, i%diaAtual)
		} else {
			data = inicioMes.AddDate(0, 0, -(i - 14))
		}

		if _, err := tx.Exec(
			`INSERT OR IGNORE INTO vendas (cliente_id, valor, data, produto)
			 VALUES (?, ?, ?, ?)`,
			clienteID, valor, data.Format(dataFmt), produto,
		); err != nil {
			return fmt.Errorf("inserir venda %d: %w", i+1, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("confirmar transação: %w", err)
	}

	fmt.Printf("seed concluído: %d clientes, %d oportunidades, 25 vendas\n", len(clientesSeed), len(oportunidadesSeed))
	return nil
}

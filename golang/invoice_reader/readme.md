## Setup

1. Crie um arquivo `invoices.csv` na raiz do projeto com o seguinte conteúdo:

```csv
Data;Estabelecimento;Portador;Valor;Parcela
2023-01-01;Ifd A;Cartão de Crédito;R$ 100,00;1
2023-01-02;Steam B;Cartão de Débito;R$ 50,00;1
2023-01-03;Uber C;Cartão de Crédito;R$ 30,00;1
```


2. Use o comando abaixo para executar o programa:

```bash
go run main.go
```

3. O programa irá gerar um arquivo `categorized_invoice.csv` com o seguinte conteúdo:

```csv
Category,Total
Food,100.00
Transport,30.00
Entertainment,50.00
Total,180.00
```


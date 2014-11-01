library(xlsx)
dataset <- read.xlsx("karina/dataset.xlsx", sheetIndex=1)
options(width=10000) 

get <- function(geneName) {
  a = dataset[dataset$Genes == geneName,];
  return (a) 
}

fc <- function(genes) {
    a = dataset[match(genes, dataset$Genes),];
    b = as.numeric(as.character(a$dm))
    return (b) 
}

pvalues <- function(genes) { 
    a = dataset[match(genes, dataset$Genes),];
    b = as.numeric(as.character(a$BH_adj_pval))
    return(b) 
}

add <- function(a,b){
  return (a+b)
}

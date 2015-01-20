library(xlsx)
dataset <- read.xlsx("karina/dataset.xlsx", sheetIndex=1)
options(width=10000) 

get <- function(geneName) {
  a = dataset[dataset$Genes == geneName,];
  return (a) 
}

genes <- function() {
    return (as.character(dataset$David.Input))
} 

fc <- function(genes) {
#    a = dataset[match(genes, dataset$Genes),];
#    b = as.numeric(as.character(a$dm))

  r <- rnorm2(length(genes), 0, 1)
  if(length(genes) == 1) {
    r <- rnorm2(2, 0, 1)
    return (c(r[1]))
  }
  return (c(r))
#   return (b) 
}

pvalues <- function(genes) { 
    a = dataset[match(genes, dataset$Genes),];
    b = as.numeric(as.character(a$BH_adj_pval))
    return(b) 
}

exprs <- function(gene) { 
    r <- rnorm2(100,0,1)
  return 
}

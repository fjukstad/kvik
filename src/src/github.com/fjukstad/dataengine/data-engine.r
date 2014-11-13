library(xlsx)
dataset <- read.xlsx("karina/dataset.xlsx", sheetIndex=1)
options(width=10000) 

#setwd("/Users/bjorn/Dropbox/repos/mine/kvik/src/src/github.com/fjukstad/dataengine")

get <- function(geneName) {
  a = dataset[dataset$Genes == geneName,];
  return (a) 
}

genes <- function() {
    return (as.character(dataset$David.Input))
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

exprs <- function(gene) { 
  mean = fc(gene)
  r <- rnorm2(100,mean,1)
  return (c(r))
}

rnorm2 <- function(n,mean,sd) { 
  return (mean+sd*scale(rnorm(n)))
}



add <- function(a,b){
  return (a+b)
}

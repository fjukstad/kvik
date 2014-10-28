library(xlsx)
dataset <- read.xlsx("karina/dataset.xlsx", sheetIndex=1)

get <- function(geneName) {
  a = dataset[dataset$Genes == geneName,];
  return (a) 
}

add <- function(a,b){
  return (a+b)
}

library(XLConnect)
setwd("/Users/bjorn/Dropbox/repos/mine/kvik/src/src/github.com/fjukstad/dataengine/")
library(xlsx)
dataset <- read.xlsx("karina/dataset.xlsx", sheetIndex=1)

get <- function(geneName) {
  a = dataset[dataset$Genes == geneName,]
  return (a) 
}
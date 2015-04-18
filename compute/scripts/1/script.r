#load("/project/data1/tice2/StressTest/genlister.rda")
load("/Users/bjorn/stallo/genlister.rda")

cc0 = data.frame(c_v_c_0, row.names=NULL, stringsAsFactors=FALSE)
cc1 = data.frame(c_v_c_1, row.names=NULL, stringsAsFactors=FALSE)

columns <- names(cc0)
for (column in columns) {
  row.names(cc0[[column]]) <- NULL 
  row.names(cc1[[column]]) <- NULL
}

### Fetches available data for the given gene(s). and returns
### the filename of a csv file where the results are written 
getGenes <- function(group, geneSymbols) {
    if(group == 0){
        dataset = cc0
    }
    else {
        dataset = cc1
    }
  fn = paste0("cc",group)

  for (gene in geneSymbols) { 
    
    res = dataset[dataset$Gene.symbol==gene,]
    
    if(dim(res)[1] == 0 ){
      return(paste0("ERROR: COULD NOT FIND GENE ",gene))
    }
    
    fn = paste0(fn, "-", gene)
  }
  
  filename = paste0(fn, ".csv")
  write.table(res, file=filename, sep=",", row.names=FALSE)
  return(filename)
} 

## Get avail info from the top n genes (sorted by p-values)
getTopN <- function(group, n=50) {
    if(group == 0){
        dataset = cc0
    }
    else {
        dataset = cc1
    }
   filename = paste0("top",n,"-cc",group,".csv") 

   res = data.frame(head(dataset,n))
   
   write.table(res, filename, sep=",", row.names=FALSE)
   return(filename)
}

### Get all available information for ALL genes 
getAll <- function(group){
  return(getTopN(group, dim(cc0)[1]))
}


fc <- function(group=0, gene) { 
    if(group == 0){
        dataset = cc0
    }
    else {
        dataset = cc1
    }
    res = dataset[dataset$Gene.symbol==gene,]
    return(res$log.fold.change)
} 

pvalue <- function(group=0, gene) { 
    if(group == 0){
        dataset = cc0
    }
    else {
        dataset = cc1
    }
    res = dataset[dataset$Gene.symbol==gene,]
    return(res$FDR.q.value)
} 

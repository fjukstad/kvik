packages <- c("ggplot2", "Hmisc", "Rcpp", "roxygen2", "jsonlite", "igraph", "dplyr", "multicore", "colorspace", "ic10", "igraph", "network","GGally","sna") 
install.packages(packages, repos='http://cran.rstudio.com/')

source("https://bioconductor.org/biocLite.R")
biocLite(c("genefu", "WGCNA", "impute", "preprocessCore", "GO.db"))
